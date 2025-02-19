package local

import "github.com/MagicNetLab/ya-practicum-diplom/internal/conf"

// New инициализация локального хранилища файлов
func New(cnf conf.Configurator) (Repository, error) {
	r := Repository{}

	return r, nil
}
