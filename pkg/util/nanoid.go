package util

import (
	"log"

	"github.com/matoous/go-nanoid/v2"
)

func NewNanoID() string {
	id, err := gonanoid.New()
	if err != nil {
		// This is highly unlikely to happen with the standard alphabet
		// and is considered a fatal error for application startup.
		log.Fatalf("failed to generate nanoid: %v", err)
	}
	return id
}
