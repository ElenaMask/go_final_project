package main

import (
	"log"
	"os"
	"strconv"

	"github.com/ElenaMask/go_final_project/pkg/db"
	"github.com/ElenaMask/go_final_project/pkg/server"
)

const WEB_DIR string = "web"
const TODO_PORT string = "TODO_PORT"
const TODO_DBFILE string = "TODO_DBFILE"

var defaultPort = 7540
var defaultDBFile = "scheduler.db"

func main() {
	logger := log.New(os.Stdout, "server: ", log.LstdFlags|log.Lshortfile)

	port := defaultPort
	portStr := os.Getenv(TODO_PORT)
	if portStr != "" {
		portInt, err := strconv.Atoi(portStr)
		if err != nil {
			logger.Fatalln("could not recognize port:", portStr, err)
		}
		port = portInt
	}

	dbFile := os.Getenv(TODO_DBFILE)
	if dbFile == "" {
		dbFile = defaultDBFile
	}
	err := db.Init(dbFile)
	if err != nil {
		logger.Fatalln("error when initializing database:", err)
	}

	srv := server.NewServer(port, logger, WEB_DIR)
	if err := srv.HttpServer.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed: %v", err)
	}
}
