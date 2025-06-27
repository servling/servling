package runtime

import (
	"context"

	"github.com/servling/servling/pkg/types"
)

type Runtime interface {
	StartService(ctx context.Context, service *types.Service) error
	StopService(ctx context.Context, service *types.Service) error
	GetServiceStatusInfo(ctx context.Context, service *types.Service) (*types.ServiceStatusInfo, error)
	PrepareStack(ctx context.Context, service *types.Application) error
}
