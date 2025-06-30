package main

import (
	"golang-workshop/src/api"
	"golang-workshop/src/config"
	"golang-workshop/src/infra/persistence/database"
	"golang-workshop/src/pkg/logging"
	"log"
)

func main() {
	cfg := config.GetConfig()
	cfg.Env.AppName = "api"
	logger, err := logging.NewLogger(cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	err = database.InitDb(cfg)
	defer database.CloseDb()
	if err != nil {
		logger.Fatal(logging.Postgres, logging.Startup, err.Error(), nil)
	}

	api.InitServer(cfg, logger)
}
