package main

import (
	"errors"
	"log"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
)

func main() {
	cnf, repo, err := appInit()
	if err != nil {
		log.Fatal(err)
	}

	_ = cnf

	defer appExit(repo)

	logger.Info("Application started")
}

// appInit - Инициализация приложения
func appInit() (config.AppConfigurator, repository.Repository, error) {
	err := logger.Init()
	if err != nil {
		return nil, nil, err
	}

	err = config.InitConfiguration()
	if err != nil {
		return nil, nil, err
	}

	appCnf := config.GetAppConfig()
	if !appCnf.IsValid() {
		return nil, nil, errors.New("app config is not valid")
	}

	repo, err := repository.NewRepository(appCnf.GetDBConf())

	if err != nil {
		return nil, nil, errors.New("failed to create repository")
	}

	return appCnf, repo, nil
}

func appExit(repo repository.Repository) {
	repo.Close()
}
