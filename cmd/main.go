package main

import (
	"log"
	"os"

	"github.com/ElenaMask/go_final_project/pkg/config"
	"github.com/ElenaMask/go_final_project/pkg/db"
	"github.com/ElenaMask/go_final_project/pkg/server"
)

func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags|log.Lshortfile)

	err := db.Init(config.DBFile)
	if err != nil {
		logger.Fatalln("error when initializing database:", err)
	}

	srv := server.NewServer(config.Port, logger, config.WebDir)
	if err := srv.HttpServer.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
