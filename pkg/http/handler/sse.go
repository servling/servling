package handler

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-fuego/fuego"
)

// FIX: Define a struct to hold both the message and its original topic.
// This is necessary because the received message from the channel doesn't
// inherently know which topic it came from.
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

	log.Printf("SSE: Client %s connected. Subscribing to topics: %v", r.RemoteAddr, topics)

	// FIX: The channel now carries our new struct type.
	mergedMessages := make(chan messageWithTopic)
	var wg sync.WaitGroup

	for _, topic := range topics {
		// Capture the topic variable for the goroutine.
		currentTopic := topic
		messages, err := pubSub.Subscribe(r.Context(), currentTopic)
		if err != nil {
			return nil, fuego.InternalServerError{Detail: fmt.Sprintf("Failed to subscribe to topic '%s': %v", currentTopic, err)}
		}

		wg.Add(1)
		// FIX: The goroutine now accepts the topic name.
		go func(topicName string, messages <-chan *message.Message) {
			defer wg.Done()
			for msg := range messages {
				select {
				// FIX: Send the message *and* its topic down the channel.
				case mergedMessages <- messageWithTopic{msg: msg, topic: topicName}:
				case <-r.Context().Done():
					return
				}
			}
		}(currentTopic, messages) // Pass the captured topic name here.
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

	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			log.Printf("SSE: Client %s disconnected.", r.RemoteAddr)
			return nil, nil
		// FIX: Receive our new struct from the channel.
		case msgWithTopic, ok := <-mergedMessages:
			if !ok {
				log.Printf("SSE: All topic streams ended for client %s.", r.RemoteAddr)
				return nil, nil
			}

			// Extract the original message and the event name.
			msg := msgWithTopic.msg
			eventName := msgWithTopic.topic // This is now guaranteed to be correct.

			msg.Ack()

			log.Printf("SSE: Sending event '%s' with UUID '%s' to client %s", eventName, msg.UUID, r.RemoteAddr)

			sseMessage := fmt.Sprintf("id: %s\nevent: %s\ndata: %s\n\n", msg.UUID, eventName, string(msg.Payload))

			if _, err := fmt.Fprint(w, sseMessage); err != nil {
				log.Printf("SSE: Error writing to client %s: %v", r.RemoteAddr, err)
				return nil, nil
			}
			flusher.Flush()
		}
	}
}
