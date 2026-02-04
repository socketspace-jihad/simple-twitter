package logger

import (
	"errors"
	"testing"
)

func Test_Logger(t *testing.T) {
	t.Run("initialize default Logger", func(t *testing.T) {
		Logger = NewLogger()
		if Logger == nil {
			t.Error(errors.New("cannot initialize Logger with default config"))
		}
		Logger.Close()
	})

	t.Run("write to a close Logger", func(t *testing.T) {
		Logger = NewLogger()
		Logger.Close()

		Logger.Info().Msg("this should not be written inside the log")

		if Logger.File != nil {
			t.Error(errors.New("logger is not closed yet"))
		}
	})
}
