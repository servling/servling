package deploy

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/servling/servling/pkg/deploy/runtime"
	"github.com/servling/servling/pkg/types"
)

//goland:noinspection GoNameStartsWithPackageName
type DeployManager struct {
	runtime runtime.Runtime
}

func NewDeployManager(runtime runtime.Runtime) *DeployManager {
	return &DeployManager{
		runtime: runtime,
	}
}

func (d *DeployManager) StartService(ctx context.Context, service *types.Service) error {
	return d.runtime.StartService(ctx, service)
}

func (d *DeployManager) StopService(ctx context.Context, service *types.Service) error {
	return d.runtime.StopService(ctx, service)
}

func (d *DeployManager) GetServiceStatusInfo(ctx context.Context, service *types.Service) (*types.ServiceStatusInfo, error) {
	return d.runtime.GetServiceStatusInfo(ctx, service)
}

func (d *DeployManager) Deploy(ctx context.Context, application *types.Application) {
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(application.Services))
	log.Debug().Str("applicationId", application.ID).Int("serviceCount", len(application.Services)).Msg("Starting services for application...")

	for _, service := range application.Services {
		service := service
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := d.StartService(ctx, &service); err != nil {
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
