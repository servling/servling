package runtime

import (
	"context"

	"github.com/servling/servling/pkg/types"
)

type Runtime interface {
	StartService(ctx context.Context, service *types.Service) error
	StopService(ctx context.Context, serviceID string) error
	GetServiceStatusInfo(ctx context.Context, serviceID string) (*types.ServiceStatusInfo, error)
	PrepareStack(ctx context.Context, service *types.Application) error
	WatchForChanges(ctx context.Context, onUpdate func(statusInfo *types.ServiceStatusInfoUpdate)) error
	GetAllServiceIDs(ctx context.Context) ([]*string, error)
}
