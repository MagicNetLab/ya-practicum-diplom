package main

import (
	"errors"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/s3"
	"log"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
)

func main() {
	cnf, repo, s3Client, err := appInit()
	if err != nil {
		log.Fatal(err)
	}

	_ = cnf
	_ = s3Client

	defer appExit(repo)

	logger.Info("Application started")
}

// appInit - Инициализация приложения
func appInit() (config.AppConfigurator, repository.Repository, s3.S3Client, error) {
	err := logger.Init()
	if err != nil {
		return nil, nil, nil, err
	}

	err = config.InitConfiguration()
	if err != nil {
		return nil, nil, nil, err
	}

	appCnf := config.GetAppConfig()
	if !appCnf.IsValid() {
		return nil, nil, nil, errors.New("app config is not valid")
	}

	s3Client, err := s3.New(appCnf.GetS3Conf())
	if err != nil {
		return nil, nil, nil, err
	}

	repo, err := repository.NewRepository(appCnf.GetDBConf())

	if err != nil {
		return nil, nil, nil, errors.New("failed to create repository")
	}

	return appCnf, repo, s3Client, nil
}

func appExit(repo repository.Repository) {
	// TODO close repository
	//repo.Close()
}
