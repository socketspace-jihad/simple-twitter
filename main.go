package main

import (
	_ "simple_twitter/internal/db/postgresql"
	"simple_twitter/internal/logger"
	"simple_twitter/internal/monitoring"
	"simple_twitter/internal/server"
)

func main() {
	p := monitoring.NewPrometheus(
		monitoring.WithPromHTTP(),
		monitoring.WithHostEnv(),
	)
	go p.Monitor()

	srv := server.NewHTTPServer(
		server.WithHostEnv(),
	)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

	logger.Logger.Close()
}
