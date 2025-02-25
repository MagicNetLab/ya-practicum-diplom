package auth

import "github.com/MagicNetLab/ya-practicum-diplom/internal/conf"

// New возвращает настроенный сервис авторизации
func New(cnf conf.Configurator) (Service, error) {

	return Service{}, nil
}
