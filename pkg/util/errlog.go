package util

import (
	"io"

	"github.com/rs/zerolog/log"
)

func CloserOrLog(closer io.Closer, message string) {
	err := closer.Close()
	if err != nil {
		log.Error().Err(err).Msg(message)
	}
}

func MustOrLog[T any](value T, err error) func(message string) *T {
	return func(message string) *T {
		if err != nil {
			log.Error().Err(err).Msg(message)
			return &value
		}
		return nil
	}
}
