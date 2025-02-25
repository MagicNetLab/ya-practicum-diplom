package inmemory

import (
	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

func NewRepository(cnf conf.Configurator) (*Repository, error) {
	r := Repository{
		dumpPath: cnf.InMemoryDumpPath(),
		users:    make(map[string]UserModel),
		tokens:   make(map[string]TokenModel),
	}

	err := r.Import()
	if err != nil {
		logger.Error("failed import data to in-memory storage", logger.StrArg("error", err.Error()))
	}

	return &r, nil
}
