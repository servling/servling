package runtime

import (
	"context"
	"fmt"
	"strings"

	"dario.lol/gotils/pkg/maps"
	"dario.lol/gotils/pkg/pointer"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/rs/zerolog/log"
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

func (d DockerRuntime) GetServiceStatusInfo(ctx context.Context, service *types.Service) (*types.ServiceStatusInfo, error) {
	summary, err := d.GetContainerByService(ctx, service)
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
		log.Error().Str("scope", "docker").Str("serviceId", service.ID).Msg("Individual service failed to start.")
	}
	out, err := d.client.ImagePull(ctx, service.Image, image.PullOptions{})
	if err != nil {
		return d.publishServiceError(
			service.ID,
			err,
			"failed to pull image %s", service.Image,
		)
	}
	defer out.Close()

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
	createdContainer, err := d.client.ContainerCreate(ctx, &container.Config{
		Image:        service.Image,
		Labels:       service.Labels,
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

func (d DockerRuntime) StopService(ctx context.Context, service *types.Service) error {
	err := util.Publish(d.pubSub, constants.TopicServiceStatusChanged, types.ServiceStatusChangedMessage{
		ID:     service.ID,
		Status: types.ServiceStatusStopping,
	})
	if err != nil {
		log.Error().Str("serviceId", service.ID).Err(err).Msg("Individual service failed to stop.")
	}
	containerSummary, err := d.GetContainerByService(ctx, service)
	if err != nil {
		return d.publishServiceError(
			service.ID,
			err,
			"failed to pull image %s", service.Image,
		)
	}
	err = d.client.ContainerStop(ctx, containerSummary.ID, container.StopOptions{})
	if err != nil {
		return d.publishServiceError(
			service.ID,
			err,
			"failed to stop container %s", service.Image,
		)
	}
	err = d.client.ContainerRemove(ctx, containerSummary.ID, container.RemoveOptions{})
	if err != nil {
		return d.publishServiceError(
			service.ID,
			err,
			"failed to remove container %s", service.Image,
		)
	}
	return util.Publish(d.pubSub, constants.TopicServiceStatusChanged, types.ServiceStatusChangedMessage{
		ID:     service.ID,
		Status: types.ServiceStatusStopped,
	})
}

func (d DockerRuntime) GetContainerByService(ctx context.Context, service *types.Service) (*container.Summary, error) {
	containerFilters := filters.NewArgs()
	containerFilters.Add("name", "^/"+service.ServiceName+"$")

	containers, err := d.client.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: containerFilters,
	})
	if err != nil {
		return nil, d.publishServiceError(
			service.ID,
			err,
			"failed to list containers %s", service.Image,
		)
	}

	if len(containers) == 0 {
		return nil, fmt.Errorf("no container found with name: %s", service.ServiceName)
	}

	return &containers[0], nil
}
