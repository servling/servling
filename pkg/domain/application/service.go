package application

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"dario.lol/gotils/pkg/encoding"
	"dario.lol/gotils/pkg/slice"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/servling/servling/ent"
	"github.com/servling/servling/pkg/config"
	"github.com/servling/servling/pkg/constants"
	"github.com/servling/servling/pkg/deploy"
	"github.com/servling/servling/pkg/types"
	"github.com/servling/servling/pkg/util"
)

//goland:noinspection GoNameStartsWithPackageName
type ApplicationService struct {
	config        *config.Config
	repository    *ApplicationRepository
	pubSub        *gochannel.GoChannel
	deployManager *deploy.DeployManager
}

func NewApplicationService(config *config.Config, client *ent.Client, pubSub *gochannel.GoChannel, deployManager *deploy.DeployManager) *ApplicationService {
	return &ApplicationService{
		config:        config,
		repository:    NewApplicationRepository(client),
		pubSub:        pubSub,
		deployManager: deployManager,
	}
}

func (s *ApplicationService) GetPubSub() *gochannel.GoChannel {
	return s.pubSub
}

func (s *ApplicationService) GetAll(ctx context.Context) ([]types.Application, error) {
	apps, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return slice.FromPtr(slice.Map(apps, types.ApplicationFromEnt)), nil
}

func (s *ApplicationService) SubscribeToServiceEvents() error {
	channel, err := s.pubSub.Subscribe(context.Background(), constants.TopicServiceStatusChanged)
	if err != nil {
		return err
	}
	log.Debug().Str("topic", constants.TopicServiceStatusChanged).Msg("Subscribed to topic.")

	for msg := range channel {
		log.Debug().Str("UUID", msg.UUID).Msg("Received message.")

		receivedMsg, err := encoding.UnmarshalJSON[types.ServiceStatusChangedMessage](msg.Payload)
		if err != nil {
			log.Debug().Err(err).Msg("Error unmarshalling message.")
			msg.Nack()
			continue
		}

		update := types.ServiceStatusInfoUpdate{
			ID: receivedMsg.ID,
			ServiceStatusInfo: types.ServiceStatusInfo{
				Status: receivedMsg.Status,
				Error:  receivedMsg.Error,
			},
		}

		s.ChangeServiceStatus(msg.Context(), update)
		log.Debug().Str("serviceId", receivedMsg.ID).Str("status", string(receivedMsg.Status)).Msg("Successfully processed status change for service.")

		msg.Ack()
	}

	return nil
}

func (s *ApplicationService) ChangeServiceStatus(ctx context.Context, update types.ServiceStatusInfoUpdate) {
	err := s.repository.UpdateServiceStatus(ctx, update.ID, update.ServiceStatusInfo)
	if err != nil {
		log.Error().Str("serviceId", update.ID).Err(err).Msg("Service status could not be updated.")
		return
	}
	service, err := s.repository.GetServiceWithApplicationServices(ctx, update.ID)
	if err != nil {
		log.Error().Str("serviceId", update.ID).Err(err).Msg("Service not found.")
		return
	}
	application, err := service.Edges.ApplicationOrErr()
	if err != nil {
		log.Error().Str("serviceId", service.ID).Err(err).Msg("Service's Application edge not found.")
		return
	}
	services, err := application.Edges.ServicesOrErr()
	if err != nil {
		log.Error().Str("serviceId", service.ID).Err(err).Msg("Application's service edges not found.")
		return
	}

	var overallStatus types.ServiceStatus
	var overallError string

	if len(services) == 0 {
		overallStatus = types.ServiceStatusStopped
	} else {
		statusPriority := map[types.ServiceStatus]int{
			types.ServiceStatusError:    5,
			types.ServiceStatusStopping: 4,
			types.ServiceStatusStarting: 3,
			types.ServiceStatusRunning:  2,
			types.ServiceStatusStopped:  1,
		}

		overallStatus = types.ServiceStatusStopped
		for _, svc := range services {
			serviceStatus := types.ServiceStatus(svc.Status)
			if statusPriority[serviceStatus] > statusPriority[overallStatus] {
				overallStatus = serviceStatus
			}
		}

		if overallStatus == types.ServiceStatusError {
			var errorMessages []string
			for _, svc := range services {
				if types.ServiceStatus(svc.Status) == types.ServiceStatusError && svc.Error != nil {
					errorMessages = append(errorMessages, *svc.Error)
				}
			}
			overallError = strings.Join(errorMessages, "; ")
		}
	}

	applicationUpdate := types.ServiceStatusInfoUpdate{
		ID: application.ID,
		ServiceStatusInfo: types.ServiceStatusInfo{
			Status: overallStatus,
			Error:  &overallError,
		},
	}

	if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, types.ApplicationStatusChangedMessage{
		ID:     application.ID,
		Status: applicationUpdate.Status,
		Error:  applicationUpdate.Error,
	}); pubErr != nil {
		log.Error().Err(pubErr).Str("serviceId", service.ID).Msg("Failed to publish started status for service.")
	}

	err = s.ChangeApplicationStatus(ctx, types.ApplicationStatusInfoUpdate{
		ID:                application.ID,
		ServiceStatusInfo: applicationUpdate.ServiceStatusInfo,
	})
	if err != nil {
		log.Error().Str("serviceId", service.ID).Err(err).Msg("Failed to change application status for service.")
	}
}

func (s *ApplicationService) ChangeApplicationStatus(ctx context.Context, update types.ApplicationStatusInfoUpdate) error {
	return s.repository.UpdateApplicationStatus(ctx, update.ID, update.ServiceStatusInfo)
}

func (s *ApplicationService) GetByID(ctx context.Context, id string) (*types.Application, error) {
	app, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return types.ApplicationFromEnt(app), nil
}

func (s *ApplicationService) Delete(ctx context.Context, application *types.Application) (*types.Application, error) {
	go s.Stop(ctx, application)
	err := s.repository.Delete(ctx, application.ID)
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (s *ApplicationService) Create(ctx context.Context, input types.CreateApplicationInput) (*types.Application, error) {
	databaseApplication, err := s.repository.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	resultApplication := types.ApplicationFromEnt(databaseApplication)
	if input.Start {
		go s.Start(context.Background(), resultApplication)
	}
	createdApp, err := s.repository.GetByID(ctx, resultApplication.ID)
	return types.ApplicationFromEnt(createdApp), err
}

func (s *ApplicationService) StartService(ctx context.Context, service *types.Service) {
	log.Debug().Str("serviceId", service.ID).Msg("Starting individual service...")
	if err := s.deployManager.StartService(ctx, service); err != nil {
		log.Error().Err(err).Str("serviceId", service.ID).Msg("Individual service failed to start.")
	} else {
		log.Debug().Str("serviceId", service.ID).Msg("Individual service started successfully.")
	}
}

func (s *ApplicationService) StopService(ctx context.Context, serviceID string) {
	log.Debug().Str("serviceId", serviceID).Msg("Stopping individual service...")
	if err := s.deployManager.StopService(ctx, serviceID); err != nil {
		log.Error().Err(err).Str("serviceId", serviceID).Msg("Individual service failed to stop.")
	} else {
		log.Debug().Str("serviceId", serviceID).Msg("Individual service stopped successfully.")
	}
}

func (s *ApplicationService) Start(ctx context.Context, application *types.Application) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(application.Services))
	log.Debug().Str("applicationId", application.ID).Int("serviceCount", len(application.Services)).Msg("Starting services for application...")

	for _, service := range application.Services {
		wg.Add(1)
		go func(srv *types.Service) {
			defer wg.Done()
			if err := s.deployManager.StartService(ctx, srv); err != nil {
				errorChannel <- fmt.Errorf("service '%s' failed: %w", srv.Name, err)
			}
		}(&service)
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

func (s *ApplicationService) Stop(ctx context.Context, application *types.Application) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(application.Services))
	log.Debug().Str("applicationId", application.ID).Msg("Stopping all services for application...")

	for _, service := range application.Services {
		service := service
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.deployManager.StopService(ctx, service.ID); err != nil {
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
		log.Error().Str("applicationId", application.ID).Err(errors.New(consolidatedError)).Msg("Application failed to stop.")
	} else {
		log.Debug().Str("applicationId", application.ID).Msg("All services for application stopped successfully.")
	}
}
