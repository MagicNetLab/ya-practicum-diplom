package main

import (
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatal(err)
	}

	cnf, err := config.MakeConfig()
	if err != nil || !cnf.IsValid() {
		log.Fatal(err)
	}

	if !cnf.IsValid() {
		logger.Fatal("db config is not valid")
	}

	m, err := migrate.New(
		"file://migrations",
		cnf.DBConnectionString(),
	)
	if err != nil {
		logger.Fatal("failed migration: error with init migration", logger.StrArg("error", err.Error()))
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal("failed migration", logger.StrArg("error", err.Error()))
	}
}
