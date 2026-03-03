package main

import (
	"simple_twitter/internal/cache"
	_ "simple_twitter/internal/cache/redis_st"
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

	cacheEngine, err := cache.UseCache("redis")
	if err != nil {
		panic(err)
	}
	if err := cacheEngine.Connect(); err != nil {
		panic(err)
	}
	defer cacheEngine.Disconnect()

	srv := server.NewHTTPServer(
		server.WithHostEnv(),
	)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

	logger.Logger.Close()
}
