package util

import (
	"github.com/matoous/go-nanoid/v2"
	"github.com/rs/zerolog/log"
)

func NewNanoID() string {
	id, err := gonanoid.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate nanoid")
	}
	return id
}
