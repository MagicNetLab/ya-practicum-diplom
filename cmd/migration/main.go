package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

func main() {
	err := logger.Init()
	if err != nil {
		log.Fatal(err)
	}

	cnf, err := conf.GetCnf()
	if err != nil {
		logger.Fatal("failed to load configuration", logger.StrArg("error", err.Error()))
	}

	connString, err := getDbConnectString(cnf)
	if err != nil {
		logger.Fatal("failed migration: connect params error", logger.StrArg("error", err.Error()))
	}

	m, err := migrate.New(
		"file://migrations",
		connString,
	)
	if err != nil {
		logger.Fatal("failed migration: error with init migration", logger.StrArg("error", err.Error()))
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal("failed migration", logger.StrArg("error", err.Error()))
	}
}

// получение строки подключения из конфигурации
func getDbConnectString(cnf conf.Configurator) (string, error) {
	if cnf.DBHost() == "" || cnf.DBPort() == "" || cnf.DBUser() == "" || cnf.DBPassword() == "" || cnf.DBName() == "" {
		return "", errors.New("db connection string is incorrect")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cnf.DBUser(),
		cnf.DBPassword(),
		cnf.DBHost(),
		cnf.DBPort(),
		cnf.DBName())

	return connStr, nil
}
