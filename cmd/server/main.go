package main

import (
	"log"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/app"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}

	err = config.InitConfiguration()
	if err != nil {
		logger.Fatal("failed to init configuration", logger.StrArg("error", err.Error()))

	}

	application, err := app.New(config.GetAppConfig())
	if err != nil {
		log.Fatalf("appInit err: %v", err)
	}

	err = application.InitServer()
	if err != nil {
		application.Stop()
		log.Fatalf("appInit err: %v", err)
	}
	defer application.Stop()

	err = application.Start()
	if err != nil {
		application.Stop()
		log.Fatalf("appStart error: %v", err)
	}

	logger.Info("Application started")
}
