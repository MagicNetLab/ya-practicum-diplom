package repo

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/models"
	"time"
)

const TokenLifeTime = time.Hour * 1
const RefreshTokenLifeTime = time.Hour * 24

type Storage struct {
	data  DataStorage
	files FileStorage
}

// GetUserByLoginAndPassword получение пользователя из хранилища по логину и паролю
func (s *Storage) GetUserByLoginAndPassword(ctx context.Context, login, password string) (models.UserEntity, error) {
	user, err := s.data.GetUserByLoginPassword(ctx, login, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByUID получение пользователя по UID
func (s *Storage) GetUserByUID(ctx context.Context, uid string) (models.UserEntity, error) {
	user, err := s.data.GetUserByUID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// HasLogin проверка существования логина
func (s *Storage) HasLogin(ctx context.Context, login string) (bool, error) {
	_, err := s.data.GetUserByLogin(ctx, login)
	if err != nil {
		return true, nil
	}

	return false, nil
}

// CreateUser создание пользователя
func (s *Storage) CreateUser(ctx context.Context, login, password string) (models.UserEntity, error) {
	user, err := s.data.CreateUser(ctx, login, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateToken создание токена
func (s *Storage) CreateToken(ctx context.Context, token, uid string, isRefresh bool) error {
	var lifeTime time.Duration
	if isRefresh {
		lifeTime = RefreshTokenLifeTime
	} else {
		lifeTime = TokenLifeTime
	}

	err := s.data.AddToken(ctx, token, uid, isRefresh, time.Now().Add(lifeTime))
	if err != nil {
		return err
	}

	return nil
}

// HasToken проверка существования действующего токена
func (s *Storage) HasToken(ctx context.Context, token string, uid string, isRefresh bool) (bool, error) {
	hasToken, err := s.data.HasToken(ctx, token, uid, isRefresh)
	if err != nil {
		return false, err
	}

	return hasToken, nil
}
