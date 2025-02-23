package local

import (
	"os"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

// New инициализация локального хранилища файлов
func New(cnf conf.Configurator) (Repository, error) {
	r := Repository{}
	r.baseFolder = cnf.FileStoragePath()

	// проверка наличия папки и ее создание в случае отсутствия
	fInfo, err := os.Stat(r.baseFolder)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Error("failed check local file storage folder", logger.StrArg("error", err.Error()))
		}

		err = createStorageFolder(r.baseFolder)
		if err != nil {
			return r, err
		}

		return r, nil
	}

	if !fInfo.IsDir() {
		err = createStorageFolder(r.baseFolder)
		if err != nil {
			return r, err
		}
	}

	return r, nil
}

func createStorageFolder(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		logger.Error("failed create local file storage folder", logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}
