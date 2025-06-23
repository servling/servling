package deploy

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/servling/servling/pkg/deploy/runtime"
	"github.com/servling/servling/pkg/types"
)

//goland:noinspection GoNameStartsWithPackageName
type DeployManager struct {
	pubSub  *gochannel.GoChannel
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
