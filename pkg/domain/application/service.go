package application

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"dario.lol/gotils/pkg/encoding"
	"dario.lol/gotils/pkg/pointer"
	"dario.lol/gotils/pkg/slice"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
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
	log.Printf("--> Subscribed to topic: %s", constants.TopicServiceStatusChanged)

	for msg := range channel {
		log.Printf("Received message: UUID=%s", msg.UUID)

		receivedMsg, err := encoding.UnmarshalJSON[types.ServiceStatusChangedMessage](msg.Payload)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
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

		if err := s.ChangeServiceStatus(msg.Context(), update); err != nil {
			log.Printf("Failed to process status change for service %s: %v", receivedMsg.ID, err)
			msg.Nack()
			continue
		}

		log.Printf("Successfully processed status change for service %s to %s", receivedMsg.ID, receivedMsg.Status)
		msg.Ack()
	}

	return nil
}

func (s *ApplicationService) ChangeServiceStatus(ctx context.Context, update types.ServiceStatusInfoUpdate) error {
	err := s.repository.UpdateServiceStatus(ctx, update.ID, update.ServiceStatusInfo)
	if err != nil {
		return err
	}
	service, err := s.repository.GetServiceWithApplicationServices(ctx, update.ID)
	if err != nil {
		return err
	}
	application, err := service.Edges.ApplicationOrErr()
	if err != nil {
		return err
	}
	services, err := application.Edges.ServicesOrErr()
	if err != nil {
		return err
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
		log.Printf("CRITICAL: Failed to publish started status for service %s: %v", service.ID, pubErr)
	}

	return s.ChangeApplicationStatus(ctx, types.ApplicationStatusInfoUpdate{
		ID:                application.ID,
		ServiceStatusInfo: applicationUpdate.ServiceStatusInfo,
	})
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
	log.Printf("Application %v created", databaseApplication)
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
	log.Printf("[SERVICE] Starting individual service '%s'...", service.Name)
	if err := s.deployManager.StartService(ctx, service); err != nil {
		log.Printf("[SERVICE] Individual service '%s' failed to start. Error: %v", service.Name, err)
		msg := &types.ApplicationStatusChangedMessage{
			ID:     service.ID,
			Status: types.ServiceStatusError,
			Error:  pointer.Of(err.Error()),
		}
		if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, msg); pubErr != nil {
			log.Printf("CRITICAL: Failed to publish error status for service %s: %v", service.ID, pubErr)
		}
	} else {
		log.Printf("[SERVICE] Individual service '%s' started successfully.", service.Name)
		msg := &types.ApplicationStatusChangedMessage{
			ID:     service.ID,
			Status: types.ServiceStatusRunning,
		}
		if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, msg); pubErr != nil {
			log.Printf("CRITICAL: Failed to publish started status for service %s: %v", service.ID, pubErr)
		}
	}
}

func (s *ApplicationService) StopService(ctx context.Context, service *types.Service) {
	log.Printf("[SERVICE] Stopping individual service '%s'...", service.Name)
	if err := s.deployManager.StopService(ctx, service); err != nil {
		log.Printf("[SERVICE] Individual service '%s' failed to stop Error: %v", service.Name, err)
		msg := &types.ApplicationStatusChangedMessage{
			ID:     service.ID,
			Status: types.ServiceStatusError,
			Error:  pointer.Of(err.Error()),
		}
		if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, msg); pubErr != nil {
			log.Printf("CRITICAL: Failed to publish error status for service %s: %v", service.ID, pubErr)
		}
	} else {
		log.Printf("[SERVICE] Individual service '%s' started successfully.", service.Name)
		msg := &types.ApplicationStatusChangedMessage{
			ID:     service.ID,
			Status: types.ServiceStatusStopped,
		}
		if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, msg); pubErr != nil {
			log.Printf("CRITICAL: Failed to publish started status for service %s: %v", service.ID, pubErr)
		}
	}
}

func (s *ApplicationService) Start(ctx context.Context, application *types.Application) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(application.Services))
	log.Printf("[SERVICE] Starting %d services for application %s...", len(application.Services), application.Name)

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
		log.Printf("[SERVICE] Application '%s' failed to start. Errors: %s", application.Name, consolidatedError)
		msg := &types.ApplicationStatusChangedMessage{
			ID:     application.ID,
			Status: types.ServiceStatusError,
			Error:  pointer.Of(consolidatedError),
		}
		if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, msg); pubErr != nil {
			log.Printf("CRITICAL: Failed to publish error status for app %s: %v", application.ID, pubErr)
		}
	} else {
		log.Printf("[SERVICE] All services for application '%s' started successfully.", application.Name)
		msg := &types.ApplicationStatusChangedMessage{
			ID:     application.ID,
			Status: types.ServiceStatusRunning,
		}
		if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, msg); pubErr != nil {
			log.Printf("CRITICAL: Failed to publish started status for app %s: %v", application.ID, pubErr)
		}
	}
}

func (s *ApplicationService) Stop(ctx context.Context, application *types.Application) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(application.Services))
	log.Printf("[SERVICE] Stopping all services for application %s...", application.Name)

	for _, service := range application.Services {
		wg.Add(1)
		go func(srv *types.Service) {
			defer wg.Done()
			if err := s.deployManager.StopService(ctx, srv); err != nil {
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
		log.Printf("[SERVICE] Application '%s' failed to stop. Errors: %s", application.Name, consolidatedError)
		msg := &types.ApplicationStatusChangedMessage{
			ID:     application.ID,
			Status: types.ServiceStatusError,
			Error:  pointer.Of(consolidatedError),
		}
		if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, msg); pubErr != nil {
			log.Printf("CRITICAL: Failed to publish error status for app %s: %v", application.ID, pubErr)
		}
	} else {
		log.Printf("[SERVICE] All services for application '%s' stopped successfully.", application.Name)
		msg := &types.ApplicationStatusChangedMessage{
			ID:     application.ID,
			Status: types.ServiceStatusStopped,
		}
		if pubErr := util.Publish(s.pubSub, constants.TopicApplicationStatusChanged, msg); pubErr != nil {
			log.Printf("CRITICAL: Failed to publish started status for app %s: %v", application.ID, pubErr)
		}
	}
}
