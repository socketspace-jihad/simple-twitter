package server

import (
	"net/http"
	"os"
	"simple_twitter/internal/logger"
)

type HTTPConfig func(*http.Server)

func WithHostEnv() HTTPConfig {
	return func(s *http.Server) {
		val := os.Getenv("APP_HOST")
		if val != "" {
			s.Addr = val
		}
	}
}

func NewHTTPServer(configs ...HTTPConfig) *http.Server {
	router := StartRouter()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	for _, config := range configs {
		if config != nil {
			config(srv)
		}
	}
	logger.Logger.Info().Str("addr", srv.Addr).Msg("http server is listening")
	return srv
}
