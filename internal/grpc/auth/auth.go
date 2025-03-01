package auth

import (
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
)

// MakeService возвращает настроенный сервис авторизации
func MakeService(cnf config.JWTConfigurator, repo repository.AuthRepository) (Service, error) {
	return Service{jwt: cnf, store: repo}, nil
}
