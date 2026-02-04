package main

import (
	_ "simple_twitter/internal/db/postgresql"
	"simple_twitter/internal/monitoring"
	"simple_twitter/internal/server"
)

func main() {
	p := monitoring.NewPrometheus(monitoring.WithPromHTTP())
	go p.Monitor()

	srv := server.NewHTTPServer()
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
