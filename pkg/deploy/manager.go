package deploy

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"dario.lol/gotils/pkg/slice"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/servling/servling/pkg/constants"
	"github.com/servling/servling/pkg/deploy/runtime"
	"github.com/servling/servling/pkg/model"
	"github.com/servling/servling/pkg/util"
)

//goland:noinspection GoNameStartsWithPackageName
type DeployManager struct {
	runtime runtime.Runtime
	pubSub  *gochannel.GoChannel
}

func NewDeployManager(runtime runtime.Runtime, pubSub *gochannel.GoChannel) *DeployManager {
	return &DeployManager{
		runtime: runtime,
		pubSub:  pubSub,
	}
}

const pollingInterval = 10 * time.Second

func (d *DeployManager) WatchForServiceStatusInfoUpdates(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(pollingInterval)
		defer ticker.Stop()

		pollAndPublish := func() {
			serviceIDs, err := d.runtime.GetAllServiceIDs(ctx)
			if err != nil {
				log.Error().Err(err).Msg("Polling failed: could not get service IDs")
				_ = runtime.PublishServiceError(
					d.pubSub,
					"*",
					err,
					"error connecting to docker daemon",
				)
				return
			}
			for _, serviceID := range slice.FilterNotNil(serviceIDs) {
				log.Debug().Str("serviceId", serviceID).Msg("Getting service info")
				statusInfo, err := d.GetServiceStatusInfo(ctx, serviceID)
				if err != nil {
					log.Error().Err(err).Str("serviceId", serviceID).Msg("Polling failed: could not get service status")
					_ = runtime.PublishServiceError(
						d.pubSub,
						serviceID,
						err,
						"error getting service status",
					)
					return
				}

				// 3. Publish the update
				err = util.Publish(d.pubSub, constants.TopicServiceStatusChanged, &model.ServiceStatusChangedMessage{
					ID:     serviceID,
					Status: statusInfo.Status,
					Error:  statusInfo.Error,
				})
				if err != nil {
					log.Error().Err(err).Str("serviceId", serviceID).Msg("Polling failed: could not publish status update")
				}
			}
		}

		pollAndPublish()
		for {
			select {
			case <-ticker.C:
				pollAndPublish()
			case <-ctx.Done():
				log.Info().Msg("Stopping periodic status poller.")
				return
			}
		}
	}()

	err := d.runtime.WatchForChanges(ctx, func(statusInfo *model.ServiceStatusInfoUpdate) {
		err := util.Publish(d.pubSub, constants.TopicServiceStatusChanged, &model.ServiceStatusChangedMessage{
			ID:     statusInfo.ID,
			Status: statusInfo.Status,
			Error:  statusInfo.Error,
		})
		if err != nil {
			log.Error().Err(err).Msg("Failed to publish service status info update")
		}
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to start WatchForChanges")
	}
}

func (d *DeployManager) StartService(ctx context.Context, service *model.Service) error {
	return d.runtime.StartService(ctx, service)
}

func (d *DeployManager) StopService(ctx context.Context, serviceID string) error {
	return d.runtime.StopService(ctx, serviceID)
}

func (d *DeployManager) GetServiceStatusInfo(ctx context.Context, serviceID string) (*model.ServiceStatusInfo, error) {
	return d.runtime.GetServiceStatusInfo(ctx, serviceID)
}

func (d *DeployManager) Deploy(ctx context.Context, application *model.Application) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(application.Services))
	log.Debug().Str("applicationId", application.ID).Int("serviceCount", len(application.Services)).Msg("Starting services for application...")

	for _, service := range application.Services {
		service := service
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := d.StartService(ctx, service); err != nil {
				errorChannel <- fmt.Errorf("service '%s' failed: %w", service.Name, err)
			}
		}()
	}

	wg.Wait()
	close(errorChannel)

	var allErrors []string
	for err := range errorChannel {
		allErrors = append(allErrors, err.Error())
	}

	if len(allErrors) > 0 {
		consolidatedError := strings.Join(allErrors, "; ")
		log.Error().Str("applicationId", application.ID).Err(errors.New(consolidatedError)).Msg("Application failed to start.")
	} else {
		log.Debug().Str("applicationId", application.ID).Msg("All services for application started successfully.")
	}
}
