package util

import (
	"dario.lol/gotils/pkg/encoding"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

func Publish[T any](pubSub *gochannel.GoChannel, topic string, data T) error {
	encodedData, err := encoding.MarshalJSON(data)
	if err != nil {
		return err
	}
	return pubSub.Publish(topic, message.NewMessage(watermill.NewUUID(), encodedData))
}
