package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ElenaMask/go_final_project/pkg/server"
)

const WEB_DIR string = "web"
const TODO_PORT string = "TODO_PORT"

var defaultPort = 7540

func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags|log.Lshortfile)
	port := defaultPort
	portStr := os.Getenv(TODO_PORT)
	if portStr != "" {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			logger.Println("could not recognize port:", portStr, err)
		} else {
			port = portInt
		}
	}
	srv := server.NewServer(port, logger, WEB_DIR)
	if err := srv.HttpServer.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
