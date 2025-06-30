package main

import (
	"golang-workshop/src/api"
	"golang-workshop/src/config"
	"golang-workshop/src/infra/persistence/database"
	"golang-workshop/src/pkg/logging"
	"log"
)

func main() {
	log.Println("main called.")
	cfg := config.GetConfig()
	logger := logging.NewLogger(cfg)

	err := database.InitDb(cfg)
	defer database.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	api.InitServer(cfg)
}
