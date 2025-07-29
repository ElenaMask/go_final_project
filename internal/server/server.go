package server

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	Logger     *log.Logger
	HttpServer *http.Server
}

func NewServer(logger *log.Logger, webDir string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	httpServer := &http.Server{
		Addr:         ":7540",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	return &Server{
		Logger:     logger,
		HttpServer: httpServer,
	}
}
