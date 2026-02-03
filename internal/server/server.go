package server

import (
	"log"
	"net/http"
)

func NewHTTPServer() *http.Server {
	router := StartRouter()

	addr := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	log.Println("http server is listening on port 8080..")

	return addr
}
