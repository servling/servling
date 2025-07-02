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
	"github.com/servling/servling/pkg/model"
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

func (s *ApplicationService) GetAll(ctx context.Context) ([]model.Application, error) {
	apps, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return slice.FromPtr(slice.Map(apps, model.ApplicationFromEnt)), nil
}

func (s *ApplicationService) SubscribeToServiceEvents() error {
	channel, err := s.pubSub.Subscribe(context.Background(), constants.TopicServiceStatusChanged)
	if err != nil {
		return err
	}
	log.Debug().Str("topic", constants.TopicServiceStatusChanged).Msg("Subscribed to topic.")

	for msg := range channel {
		log.Debug().Str("UUID", msg.UUID).Msg("Received message.")

		receivedMsg, err := encoding.UnmarshalJSON[model.ServiceStatusChangedMessage](msg.Payload)
		if err != nil {
			log.Debug().Err(err).Msg("Error unmarshalling message.")
			msg.Nack()
			continue
		}

		if receivedMsg.ID == "*" {
			allServices, err := s.repository.GetAllServices(context.Background())
			if err != nil {
				log.Debug().Err(err).Msg("Error getting all services from repository.")
			}
			for _, service := range allServices {
				s.ChangeServiceStatus(msg.Context(), model.ServiceStatusInfoUpdate{
					ID: service.ID,
					ServiceStatusInfo: model.ServiceStatusInfo{
						Status: receivedMsg.Status,
						Error:  receivedMsg.Error,
					},
				})
			}
			log.Debug().Str("status", string(receivedMsg.Status)).Msg("Successfully changed status for all services.")
		} else {
			s.ChangeServiceStatus(msg.Context(), model.ServiceStatusInfoUpdate{
				ID: receivedMsg.ID,
				ServiceStatusInfo: model.ServiceStatusInfo{
					Status: receivedMsg.Status,
					Error:  receivedMsg.Error,
				},
			})
			log.Debug().Str("serviceId", receivedMsg.ID).Str("status", string(receivedMsg.Status)).Msg("Successfully processed status change for service.")
		}

		msg.Ack()
	}

	return nil
}

func (s *ApplicationService) ChangeServiceStatus(ctx context.Context, update model.ServiceStatusInfoUpdate) {
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

	var overallStatus model.ServiceStatus
	var overallError string

	if len(services) == 0 {
		overallStatus = model.ServiceStatusStopped
	} else {
		statusPriority := map[model.ServiceStatus]int{
			model.ServiceStatusError:    5,
			model.ServiceStatusStopping: 4,
			model.ServiceStatusStarting: 3,
			model.ServiceStatusRunning:  2,
			model.ServiceStatusStopped:  1,
		}

		overallStatus = model.ServiceStatusStopped
		for _, svc := range services {
			serviceStatus := model.ServiceStatus(svc.Status)
			if statusPriority[serviceStatus] > statusPriority[overallStatus] {
				overallStatus = serviceStatus
			}
		}

		if overallStatus == model.ServiceStatusError {
			var errorMessages []string
			for _, svc := range services {
				if model.ServiceStatus(svc.Status) == model.ServiceStatusError && svc.Error != nil {
					errorMessages = append(errorMessages, *svc.Error)
				}
			}
			overallError = strings.Join(errorMessages, "; ")
		}
	}

	applicationUpdate := model.ServiceStatusInfoUpdate{
		ID: application.ID,
		ServiceStatusInfo: model.ServiceStatusInfo{
			Status: overallStatus,
			Error:  &overallError,
		},
	}

	if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, model.ApplicationStatusChangedMessage{
		ID:     application.ID,
		Status: applicationUpdate.Status,
		Error:  applicationUpdate.Error,
	}); pubErr != nil {
		log.Error().Err(pubErr).Str("serviceId", service.ID).Msg("Failed to publish started status for service.")
	}

	err = s.ChangeApplicationStatus(ctx, model.ApplicationStatusInfoUpdate{
		ID:                application.ID,
		ServiceStatusInfo: applicationUpdate.ServiceStatusInfo,
	})
	if err != nil {
		log.Error().Str("serviceId", service.ID).Err(err).Msg("Failed to change application status for service.")
	}
}

func (s *ApplicationService) ChangeApplicationStatus(ctx context.Context, update model.ApplicationStatusInfoUpdate) error {
	return s.repository.UpdateApplicationStatus(ctx, update.ID, update.ServiceStatusInfo)
}

func (s *ApplicationService) GetByID(ctx context.Context, id string) (*model.Application, error) {
	app, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ApplicationFromEnt(app), nil
}

func (s *ApplicationService) GetByIDWithIngresses(ctx context.Context, id string) (*model.Application, error) {
	app, err := s.repository.GetByIDWithIngresses(ctx, id)
	if err != nil {
		return nil, err
	}
	return model.ApplicationFromEnt(app), nil
}

func (s *ApplicationService) Delete(ctx context.Context, application *model.Application) (*model.Application, error) {
	go s.Stop(ctx, application)
	err := s.repository.Delete(ctx, application.ID)
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (s *ApplicationService) Create(ctx context.Context, input model.CreateApplicationInput) (*model.Application, error) {
	databaseApplication, err := s.repository.Create(ctx, input)
	if err != nil {
		return nil, err
	}
	resultApplication := model.ApplicationFromEnt(databaseApplication)
	if input.Start {
		go s.Start(context.Background(), resultApplication)
	}
	createdApp, err := s.repository.GetByID(ctx, resultApplication.ID)
	return model.ApplicationFromEnt(createdApp), err
}

func (s *ApplicationService) StartService(ctx context.Context, service *model.Service) {
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

func (s *ApplicationService) Start(ctx context.Context, application *model.Application) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(application.Services))
	log.Debug().Str("applicationId", application.ID).Int("serviceCount", len(application.Services)).Msg("Starting services for application...")

	for _, service := range application.Services {
		wg.Add(1)
		go func(srv *model.Service) {
			defer wg.Done()
			if err := s.deployManager.StartService(ctx, srv); err != nil {
				errorChannel <- fmt.Errorf("service '%s' failed: %w", srv.Name, err)
			}
		}(service)
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

func (s *ApplicationService) Stop(ctx context.Context, application *model.Application) {
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
