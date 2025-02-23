package repo

import (
	"errors"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/files/local"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/files/s3"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/data/inmemory"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/data/postgres"
)

var repo storage

// InitRepo инициализация репозитория
func InitRepo(cnf conf.Configurator) error {

	// data storage
	if cnf.DataStorageType() == "inmemory" {
		s, err := inmemory.NewRepository(cnf)
		if err != nil {
			logger.Error("failed to init in-memory data storage", logger.StrArg("error", err.Error()))
			return err
		}
		repo.data = s
	} else if cnf.DataStorageType() == "postgres" {
		s, err := postgres.NewRepository(cnf)
		if err != nil {
			logger.Error("failed to init postgres data storage", logger.StrArg("error", err.Error()))
			return err
		}
		repo.data = s
	} else {
		logger.Error("data storage type not correct. support only 'inmemory', 'postgres", logger.StrArg("dataStorageType", cnf.DataStorageType()))
		return errors.New("data storage type not support")
	}

	// file storage
	if cnf.FileStorageType() == "local" {
		f, err := local.New(cnf)
		if err != nil {
			logger.Error("failed init local file storage error", logger.StrArg("error", err.Error()))
			return err
		}
		repo.files = f
	} else if cnf.FileStorageType() == "s3" {
		f, err := s3.New(cnf)
		if err != nil {
			logger.Error("failed init s3 storage error", logger.StrArg("error", err.Error()))
			return err
		}
		repo.files = f
	} else {
		logger.Error("file storage type not correct. support only 'local', 's3", logger.StrArg("fileStorageType", cnf.FileStorageType()))
		return errors.New("file storage type not support")
	}

	return nil
}

// Close закрытие репозитория
func Close() error {
	return nil
}
