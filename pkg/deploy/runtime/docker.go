package runtime

import (
	"context"
	"fmt"
	"strings"

	"dario.lol/gotils/pkg/maps"
	"dario.lol/gotils/pkg/pointer"
	"dario.lol/gotils/pkg/slice"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"
	"github.com/servling/servling/ent/service"
	"github.com/servling/servling/pkg/constants"
	"github.com/servling/servling/pkg/types"
	"github.com/servling/servling/pkg/util"
)

type DockerRuntime struct {
	client *client.Client
	pubSub *gochannel.GoChannel
}

var _ Runtime = (*DockerRuntime)(nil)

func NewDockerRuntime(client *client.Client, pubSub *gochannel.GoChannel) *DockerRuntime {
	return &DockerRuntime{
		client: client,
		pubSub: pubSub,
	}
}

func GetServiceStatusInfo(summary *container.Summary) types.ServiceStatusInfo {
	switch summary.State {
	case "running":
		if strings.Contains(summary.Status, "(unhealthy)") {
			return types.ServiceStatusInfo{
				Status: types.ServiceStatusError,
				Error:  pointer.Of(fmt.Sprintf("container is unhealthy: %s", summary.Status)),
			}
		}
		if strings.Contains(summary.Status, "(health: starting)") {
			return types.ServiceStatusInfo{Status: types.ServiceStatusStarting}
		}
		return types.ServiceStatusInfo{Status: types.ServiceStatusRunning}

	case "created", "restarting":
		return types.ServiceStatusInfo{Status: types.ServiceStatusStarting}

	case "removing", "paused":
		return types.ServiceStatusInfo{Status: types.ServiceStatusStopping}

	case "exited", "dead":
		var exitCode int
		if n, _ := fmt.Sscanf(summary.Status, "Exited (%d)", &exitCode); n == 1 && exitCode != 0 {
			return types.ServiceStatusInfo{
				Status: types.ServiceStatusError,
				Error:  pointer.Of(fmt.Sprintf("container exited with non-zero code: %s", summary.Status)),
			}
		}
		if summary.State == "dead" {
			return types.ServiceStatusInfo{
				Status: types.ServiceStatusError,
				Error:  pointer.Of(fmt.Sprintf("container is dead: %s", summary.Status)),
			}
		}
		return types.ServiceStatusInfo{Status: types.ServiceStatusStopped}

	default:
		return types.ServiceStatusInfo{
			Status: types.ServiceStatusError,
			Error:  pointer.Of(fmt.Sprintf("unknown container state '%s'", summary.State)),
		}
	}
}

func (d DockerRuntime) publishServiceError(
	serviceID string,
	originalErr error,
	format string,
	args ...any,
) error {
	detailedErr := fmt.Errorf(format+": %w", append(args, originalErr)...)

	msg := types.ServiceStatusChangedMessage{
		ID:     serviceID,
		Status: types.ServiceStatusError,
		Error:  pointer.Of(detailedErr.Error()),
	}

	if pubErr := util.Publish(d.pubSub, constants.TopicServiceStatusChanged, msg); pubErr != nil {
		return fmt.Errorf("failed to publish error status for service %s: %w", serviceID, pubErr)
	}

	return detailedErr
}

func (d DockerRuntime) GetServiceStatusInfo(ctx context.Context, serviceID string) (*types.ServiceStatusInfo, error) {
	summary, err := d.GetContainerByServiceID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	return pointer.Of(GetServiceStatusInfo(summary)), nil
}

func (d DockerRuntime) StartService(ctx context.Context, service *types.Service) error {
	err := util.Publish(d.pubSub, constants.TopicServiceStatusChanged, types.ServiceStatusChangedMessage{
		ID:     service.ID,
		Status: types.ServiceStatusStarting,
	})
	if err != nil {
		log.Error().Str("scope", "docker").Str("serviceId", service.ID).Msg("Failed to publish status change message.")
	}
	out, err := d.client.ImagePull(ctx, service.Image, image.PullOptions{})
	if err != nil {
		return d.publishServiceError(
			service.ID,
			err,
			"failed to pull image %s", service.Image,
		)
	}
	defer util.CloserOrLog(out, "Error closing image pull response")

	existingContainers, err := d.client.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("label", "servling.serviceId="+service.ID)),
		All:     true,
	})
	if err != nil {
		return d.publishServiceError(
			service.ID,
			err,
			"failed to find existing container %s", service.ServiceName,
		)
	}

	if existingContainers == nil || len(existingContainers) == 0 {

		exposedPorts := make(nat.PortSet)
		portBindings := make(nat.PortMap)

		for containerPort, hostPort := range service.Ports {
			port := nat.Port(containerPort)
			exposedPorts[port] = struct{}{}
			portBindings[port] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: hostPort,
				},
			}
		}

		labels := map[string]string{
			"servling.managed":   "true",
			"servling.serviceId": service.ID,
		}

		for key, value := range service.Labels {
			labels[key] = value
		}

		createdContainer, err := d.client.ContainerCreate(ctx, &container.Config{
			Image:        service.Image,
			Labels:       labels,
			ExposedPorts: exposedPorts,
			Env: maps.MapEntries(service.Environment, func(e maps.Entry[string, string]) string {
				return e.Key + "=" + e.Value
			}),
		}, &container.HostConfig{
			PortBindings: portBindings,
		}, nil, nil, service.ServiceName)
		if err != nil {
			return d.publishServiceError(
				service.ID,
				err,
				"failed to create container %s", service.ServiceName,
			)
		}
		err = d.client.ContainerStart(ctx, createdContainer.ID, container.StartOptions{})
	} else {
		err = d.client.ContainerStart(ctx, existingContainers[0].ID, container.StartOptions{})
	}
	if err != nil {
		return d.publishServiceError(
			service.ID,
			err,
			"failed start container %s", service.ServiceName,
		)
	}
	return util.Publish(d.pubSub, constants.TopicServiceStatusChanged, types.ServiceStatusChangedMessage{
		ID:     service.ID,
		Status: types.ServiceStatusRunning,
	})
}

func (d DockerRuntime) StopService(ctx context.Context, serviceID string) error {
	err := util.Publish(d.pubSub, constants.TopicServiceStatusChanged, types.ServiceStatusChangedMessage{
		ID:     serviceID,
		Status: types.ServiceStatusStopping,
	})
	if err != nil {
		log.Error().Str("serviceId", serviceID).Err(err).Msg("Individual service failed to stop.")
	}
	containerSummary, err := d.GetContainerByServiceID(ctx, serviceID)
	if err != nil {
		return d.publishServiceError(
			serviceID,
			err,
			"failed to pull image %s", service.Image,
		)
	}
	err = d.client.ContainerStop(ctx, containerSummary.ID, container.StopOptions{})
	if err != nil {
		return d.publishServiceError(
			serviceID,
			err,
			"failed to stop container %s", service.Image,
		)
	}
	err = d.client.ContainerRemove(ctx, containerSummary.ID, container.RemoveOptions{})
	if err != nil {
		return d.publishServiceError(
			serviceID,
			err,
			"failed to remove container %s", service.Image,
		)
	}
	return util.Publish(d.pubSub, constants.TopicServiceStatusChanged, types.ServiceStatusChangedMessage{
		ID:     serviceID,
		Status: types.ServiceStatusStopped,
	})
}

func (d DockerRuntime) GetContainerByServiceID(ctx context.Context, serviceID string) (*container.Summary, error) {
	containerFilters := filters.NewArgs()
	containerFilters.Add("label", "servling.serviceId="+serviceID)

	containers, err := d.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: containerFilters,
	})
	if err != nil {
		return nil, d.publishServiceError(
			serviceID,
			err,
			"failed to list containers %s", service.Image,
		)
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("no container found with name: %s", service.ServiceName)
	}

	return &containers[0], nil
}

func (d DockerRuntime) PrepareStack(ctx context.Context, service *types.Application) error {
	return nil
}

func (d DockerRuntime) WatchForChanges(ctx context.Context, onUpdate func(statusInfo *types.ServiceStatusInfoUpdate)) error {
	eventFilters := filters.NewArgs(
		filters.Arg("type", "container"),
		filters.Arg("event", "start"),
		filters.Arg("event", "die"),
		filters.Arg("event", "stop"),
		filters.Arg("event", "pause"),
		filters.Arg("event", "unpause"),
		filters.Arg("event", "health_status"),
		filters.Arg("event", "oom"), // Out of Memory
		filters.Arg("event", "destroy"),
	)

	messages, errs := d.client.Events(ctx, events.ListOptions{
		Filters: eventFilters,
	})

	go func() {
		log.Info().Msg("Listening for container events...")
		for {
			select {
			case <-ctx.Done():
				return // Stop listening when context is cancelled
			case err := <-errs:
				if err != nil && !strings.Contains(err.Error(), "context canceled") {
					log.Error().Err(err).Msg("Error receiving Docker event")
				}
				return
			case msg := <-messages:
				d.processEvent(ctx, msg, onUpdate)
			}
		}
	}()

	return nil
}

// processEvent fetches container details and calls the callback.
func (d DockerRuntime) processEvent(ctx context.Context, msg events.Message, onUpdate func(statusInfo *types.ServiceStatusInfoUpdate)) {
	serviceID, ok := msg.Actor.Attributes["servling.serviceId"]
	if !ok || serviceID == "" {
		// This should not happen due to the event filter, but it's a good safeguard.
		return
	}

	if msg.Action == "destroy" {
		onUpdate(&types.ServiceStatusInfoUpdate{
			ID: serviceID,
			ServiceStatusInfo: types.ServiceStatusInfo{
				Status: types.ServiceStatusStopped,
			},
		})
		return
	}

	containerID := msg.Actor.ID
	containerFilters := filters.NewArgs(filters.Arg("id", containerID))
	summaries, err := d.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: containerFilters,
	})
	if err != nil {
		log.Printf("Error inspecting container %s after event '%s': %v", containerID[:12], msg.Action, err)
		return
	}

	if len(summaries) == 0 {
		onUpdate(&types.ServiceStatusInfoUpdate{
			ID: serviceID,
			ServiceStatusInfo: types.ServiceStatusInfo{
				Status: types.ServiceStatusStopped,
			},
		})
		return
	}

	// Use your provided function to get the status
	statusInfo := GetServiceStatusInfo(&summaries[0])
	statusInfoUpdate := types.ServiceStatusInfoUpdate{
		ID:                serviceID,
		ServiceStatusInfo: statusInfo,
	}

	onUpdate(&statusInfoUpdate)
}

func (d DockerRuntime) GetAllServiceIDs(ctx context.Context) ([]*string, error) {
	containers, err := d.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("label", "serving.managed=true")),
	})
	if err != nil {
		return nil, err
	}
	return slice.Map(containers, func(container container.Summary) *string {
		id, ok := container.Labels["servling.serviceId"]
		if !ok {
			return nil
		}
		return &id
	}), nil
}
