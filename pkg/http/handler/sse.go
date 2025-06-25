package handler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-fuego/fuego"
	"github.com/rs/zerolog/log"
)

type messageWithTopic struct {
	msg   *message.Message
	topic string
}

// SSEEventsController is a Fuego controller for the SSE stream.
func SSEEventsController[T any](c fuego.Context[any, any], pubSub *gochannel.GoChannel, topics ...string) (*T, error) {
	w := c.Response()
	r := c.Request()

	if len(topics) == 0 {
		return nil, fuego.BadRequestError{Detail: "At least one topic must be provided."}
	}

	log.Debug().Str("scope", "SSE").Str("remoteAddress", r.RemoteAddr).Strs("topics", topics).Msg("Client connected. Subscribing to topics.")

	mergedMessages := make(chan messageWithTopic)
	var wg sync.WaitGroup

	for _, topic := range topics {
		currentTopic := topic
		messages, err := pubSub.Subscribe(r.Context(), currentTopic)
		if err != nil {
			return nil, fuego.InternalServerError{Detail: fmt.Sprintf("Failed to subscribe to topic '%s': %v", currentTopic, err)}
		}

		wg.Add(1)
		go func(topicName string, messages <-chan *message.Message) {
			defer wg.Done()
			for msg := range messages {
				select {
				case mergedMessages <- messageWithTopic{msg: msg, topic: topicName}:
				case <-r.Context().Done():
					return
				}
			}
		}(currentTopic, messages)
	}

	go func() {
		wg.Wait()
		close(mergedMessages)
	}()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		return nil, fuego.InternalServerError{Detail: "Streaming unsupported!"}
	}
	_, err := w.Write([]byte{'\n'})
	if err != nil {
		return nil, err
	}

	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			log.Debug().Str("scope", "SSE").Str("remoteAddress", r.RemoteAddr).Msg("Client disconnected.")
			return nil, nil
		case msgWithTopic, ok := <-mergedMessages:
			if !ok {
				log.Debug().Str("scope", "SSE").Str("remoteAddress", r.RemoteAddr).Msg("All topic streams ended for client.")
				return nil, nil
			}

			msg := msgWithTopic.msg
			eventName := msgWithTopic.topic

			msg.Ack()

			log.Debug().Str("scope", "SSE").Str("remoteAddress", r.RemoteAddr).Str("UUID", msg.UUID).Str("eventName", eventName).Msg("Sending event to client.")

			sseMessage := fmt.Sprintf("id: %s\nevent: %s\ndata: %s\n\n", msg.UUID, eventName, string(msg.Payload))

			if _, err := fmt.Fprint(w, sseMessage); err != nil {
				log.Debug().Str("scope", "SSE").Str("remoteAddress", r.RemoteAddr).Err(err).Msg("Error writing to client.")
				return nil, nil
			}
			flusher.Flush()
		}
	}
}
