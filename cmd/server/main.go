package main

import (
	"log"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo"
)

func main() {
	err := appInit()
	if err != nil {
		log.Fatal(err)
	}

	cnf, err := conf.GetCnf()
	if err != nil {
		logger.Fatal("failed to load configuration", logger.StrArg("error", err.Error()))
	}

	if err = repo.InitRepo(cnf); err != nil {
		logger.Fatal("failed to init repository", logger.StrArg("error", err.Error()))
	}

	defer appExit()

	logger.Info("Application started")
}

func appInit() error {
	err := logger.Init()
	if err != nil {
		return err
	}

	return nil
}

func appExit() {
	if err := repo.Close(); err != nil {
		logger.Error("failed close repo with app exit", logger.StrArg("error", err.Error()))
	}
}
