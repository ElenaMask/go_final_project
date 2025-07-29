package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Logger     *log.Logger
	HttpServer *http.Server
}

func NewServer(port int, logger *log.Logger, webDir string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	addr := fmt.Sprintf(":%d", port)
	httpServer := &http.Server{
		Addr:         addr,
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
