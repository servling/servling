package runtime

import (
	"context"
	"fmt"

	"dario.lol/gotils/pkg/pointer"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/servling/servling/pkg/constants"
	"github.com/servling/servling/pkg/model"
	"github.com/servling/servling/pkg/util"
)

type Runtime interface {
	StartService(ctx context.Context, service *model.Service) error
	StopService(ctx context.Context, serviceID string) error
	GetServiceStatusInfo(ctx context.Context, serviceID string) (*model.ServiceStatusInfo, error)
	PrepareStack(ctx context.Context, service *model.Application) error
	WatchForChanges(ctx context.Context, onUpdate func(statusInfo *model.ServiceStatusInfoUpdate)) error
	GetAllServiceIDs(ctx context.Context) ([]*string, error)
}

func PublishServiceError(
	pubSub *gochannel.GoChannel,
	serviceID string,
	originalErr error,
	format string,
	args ...any,
) error {
	detailedErr := fmt.Errorf(format+": %w", append(args, originalErr)...)

	msg := model.ServiceStatusChangedMessage{
		ID:     serviceID,
		Status: model.ServiceStatusError,
		Error:  pointer.Of(detailedErr.Error()),
	}

	if pubErr := util.Publish(pubSub, constants.TopicServiceStatusChanged, msg); pubErr != nil {
		return fmt.Errorf("failed to publish error status for service %s: %w", serviceID, pubErr)
	}

	return detailedErr
}
