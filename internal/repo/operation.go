package repo

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/models"
)

// GetUserByLoginAndPassword получение пользователя из хранилища по логину и паролю
func GetUserByLoginAndPassword(ctx context.Context, login, password string) (models.UserEntity, error) {
	user, err := repository.data.GetUserByLoginPassword(ctx, login, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUID получение пользователя по UID
func GetUserByUID(ctx context.Context, uid string) (models.UserEntity, error) {
	user, err := repository.data.GetUserByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// HasLogin проверка существования логина
func HasLogin(ctx context.Context, login string) (bool, error) {
	_, err := repository.data.GetUserByLogin(ctx, login)
	if err != nil {
		return true, nil
	}

	return false, nil
}

// CreateUser создание пользователя
func CreateUser(ctx context.Context, login, password string) (models.UserEntity, error) {
	user, err := repository.data.CreateUser(ctx, login, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// HasToken проверка существования действующего токена
func HasToken(ctx context.Context, token string, uid string, isRefresh bool) (bool, error) {
	hasToken, err := repository.data.HasToken(ctx, token, uid, isRefresh)
	if err != nil {
		return false, err
	}

	return hasToken, nil
}
