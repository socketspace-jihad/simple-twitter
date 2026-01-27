package main

import (
	_ "simple_twitter/internal/db/postgresql"
	"simple_twitter/internal/server"
)

func main() {
	srv := server.NewHTTPServer()
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
