package main

import (
	"os"
	"simple_twitter/internal/cache"
	_ "simple_twitter/internal/cache/redis_st"
	_ "simple_twitter/internal/db/postgresql"
	"simple_twitter/internal/logger"
	"simple_twitter/internal/monitoring"
	_ "simple_twitter/internal/nats_st"
	"simple_twitter/internal/server"
	"simple_twitter/internal/worker"
	_ "simple_twitter/internal/worker/mail_sender"
)

func main() {
	p := monitoring.NewPrometheus(
		monitoring.WithPromHTTP(),
		monitoring.WithHostEnv(),
	)
	go p.Monitor()

	//cache
	cacheEngine, err := cache.UseCache(os.Getenv("CACHE_ENGINE"))
	if err != nil {
		panic(err)
	}
	if err := cacheEngine.Connect(); err != nil {
		panic(err)
	}
	defer cacheEngine.Disconnect()

	//message queue
	worker.RunAll()

	srv := server.NewHTTPServer(
		server.WithHostEnv(),
	)
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}

	logger.Logger.Close()
}
