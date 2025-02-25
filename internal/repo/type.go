package repo

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/models"
	"time"
)

type storage struct {
	data  DataStorage
	files FileStorage
}

// DataStorage интерфейс хранилища данных
type DataStorage interface {
	// GetUserByUID получение пользователя по UID
	GetUserByUID(ctx context.Context, uid string) (models.UserEntity, error)
	// GetUserByLogin получение пользователя по логину
	GetUserByLogin(ctx context.Context, login string) (models.UserEntity, error)
	// GetUserByLoginPassword получение пользователя по логину и паролю
	GetUserByLoginPassword(ctx context.Context, login, password string) (models.UserEntity, error)
	// CreateUser создание пользователя
	CreateUser(ctx context.Context, login, password string) (models.UserEntity, error)
	// RemoveUserByUID удаление пользователя по UID
	RemoveUserByUID(ctx context.Context, uid string) error
	// RemoveUserByLogin удаление пользователя по логину
	RemoveUserByLogin(ctx context.Context, login string) error
	// AddToken добавление токена
	AddToken(ctx context.Context, token, uid string, isRefresh bool, expired time.Time) error
	// RemoveToken удаление токена
	RemoveToken(ctx context.Context, value string) error
	// RemoveUserTokens удаление всех токенов пользователя
	RemoveUserTokens(ctx context.Context, uid string) error
	// HasToken проверка токена
	HasToken(ctx context.Context, value, uid string, isRefresh bool) (bool, error)
	// Close закрытие хранилища
	Close(ctx context.Context) error
}

// FileStorage интерфейс хранилища файлов
type FileStorage interface {
	Close() error
}
