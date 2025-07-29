package main

import (
	"log"
	"os"

	"github.com/ElenaMask/go_final_project/internal/server"
)

var webDir string = "web"

func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags|log.Lshortfile)
	srv := server.NewServer(logger, webDir)
	if err := srv.HttpServer.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
