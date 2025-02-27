package auth

import (
	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo"
)

// MakeService возвращает настроенный сервис авторизации
func MakeService(cnf conf.Configurator) (Service, error) {
	store := repo.GetStorage()
	cnf, err := conf.GetCnf()
	if err != nil {
		return Service{}, err
	}

	return Service{cnf: cnf, store: store}, nil
}
