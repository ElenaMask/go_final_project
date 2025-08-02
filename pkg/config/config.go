package config

import (
	"log"
	"os"
	"strconv"
)

const (
	TODO_PORT   = "TODO_PORT"
	TODO_DBFILE = "TODO_DBFILE"
)

var (
	WebDir = "web"
	Port   = 7540
	DBFile = "scheduler.db"
)

func init() {
	Port = getIntEnvOrDefault(TODO_PORT, Port)
	DBFile = getEnvOrDefault(TODO_DBFILE, DBFile)
}

func getIntEnvOrDefault(key string, def int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return def
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Variable %s: can not parse int, use default value: %d", key, def)
		return def
	}

	return val
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
