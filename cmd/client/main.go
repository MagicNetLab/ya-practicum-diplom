package main

import (
	"context"
	cli "github.com/MagicNetLab/ya-practicum-diplom/internal/app/client"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	appClianet "github.com/MagicNetLab/ya-practicum-diplom/internal/services/client"
	"log"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatalf("failed to start application. Logger init error: %v", err)
	}

	cnf, err := config.MakeConfig()
	if err != nil {
		log.Fatalf("failed to start application. Configuration init error: %v", err)
	}

	client, err := appClianet.NewAppClient(cnf)
	if err != nil {
		log.Fatalf("failed to start application. Application client error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = cli.Run(ctx, client)
	if err != nil {
		log.Fatalf("failed to run application")
	}

}
