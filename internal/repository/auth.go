package repository

import (
	"context"
	"errors"
	"time"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/encryptor"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

// NewAuthRepository возвращает новый репозиторий для работы с базой данных.
func NewAuthRepository(pool *pgxpool.Pool) AuthRepository {
	return &AuthRepo{pool: pool}
}

// AuthRepository интерфейс репозитория авторизации
type AuthRepository interface {
	GetUserByUID(ctx context.Context, uid string) (models.UserModel, error)
	GetUserByLogin(ctx context.Context, login string) (models.UserModel, error)
	GetUserByLoginAndPassword(ctx context.Context, login, password string) (models.UserModel, error)
	CreateUser(ctx context.Context, login, password string) (models.UserModel, error)
	RemoveUser(ctx context.Context, uid string) error
	CreateToken(ctx context.Context, uid, token string, isRefresh bool, expired time.Time) error
	RemoveToken(ctx context.Context, Token string) error
	RemoveUserTokens(ctx context.Context, uid string) error
	HasToken(ctx context.Context, Token string) (bool, error)
	HasLogin(ctx context.Context, login string) (bool, error)
}

// AuthRepo репозиторий авторизации
type AuthRepo struct {
	pool *pgxpool.Pool
}

// GetUserByUID получение пользователя по UID
func (r *AuthRepo) GetUserByUID(ctx context.Context, uid string) (models.UserModel, error) {
	user := models.User{}
	rows := r.pool.QueryRow(ctx, "SELECT uid, login, password, created_at, updated_at FROM users WHERE uid=$1", uid)
	err := rows.Scan(&user.UID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.Error("failed get user by uid", logger.StrArg("uid", uid), logger.StrArg("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

// GetUserByLogin получение пользователя по логину
func (r *AuthRepo) GetUserByLogin(ctx context.Context, login string) (models.UserModel, error) {
	user := models.User{}
	rows := r.pool.QueryRow(ctx, "SELECT uid, login, password, created_at, updated_at FROM users WHERE login=$1", login)
	err := rows.Scan(&user.UID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.Error("failed get user by login", logger.StrArg("login", login), logger.StrArg("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

// GetUserByLoginAndPassword получение пользователя по логину и паролю
func (r *AuthRepo) GetUserByLoginAndPassword(ctx context.Context, login, password string) (models.UserModel, error) {
	user := models.User{}
	password, err := encryptor.EncryptPassword(password)
	if err != nil {
		return nil, err
	}

	rows := r.pool.QueryRow(ctx, "SELECT uid, login, password, created_at, updated_at FROM users WHERE login=$1 AND password=$2", login, password)
	err = rows.Scan(&user.UID, &user.Login, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.Error("failed get user by login and password", logger.StrArg("login", login), logger.StrArg("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

// CreateUser создание пользователя
func (r *AuthRepo) CreateUser(ctx context.Context, login, password string) (models.UserModel, error) {
	if login == "" || password == "" {
		return nil, errors.New("empty login or password")
	}

	uid := uuid.New().String()
	encryptPassword, err := encryptor.EncryptPassword(password)
	if err != nil {
		return nil, err
	}
	user := models.User{
		UID:       uid,
		Login:     login,
		Password:  encryptPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	sql := "INSERT INTO users (uid, login, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err = r.pool.Exec(ctx, sql, user.UID, user.Login, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		logger.Error(
			"failed create user",
			logger.StrArg("uid", uid),
			logger.StrArg("login", login),
			logger.StrArg("password", user.Password),
			logger.StrArg("created_at", user.CreatedAt.String()),
			logger.StrArg("updated_at", user.UpdatedAt.String()),
			logger.StrArg("error", err.Error()),
		)
		return nil, err
	}
	return &user, nil
}

// RemoveUser удаление пользователя
func (r *AuthRepo) RemoveUser(ctx context.Context, uid string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM users WHERE uid=$1", uid)
	if err != nil {
		logger.Error("failed remove user", logger.StrArg("uid", uid), logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}

// CreateToken создание токена
func (r *AuthRepo) CreateToken(ctx context.Context, uid, token string, isRefresh bool, expired time.Time) error {
	id := uuid.New().String()
	sql := "INSERT INTO tokens (id, token, uid, is_refresh, expired) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.pool.Exec(ctx, sql, id, token, uid, isRefresh, expired)
	if err != nil {
		logger.Error("failed create token", logger.StrArg("uid", uid), logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}

// RemoveToken удаление токена
func (r *AuthRepo) RemoveToken(ctx context.Context, Token string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM tokens WHERE token=$1", Token)
	if err != nil {
		logger.Error("failed remove token", logger.StrArg("token", Token), logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}

// RemoveUserTokens удаление токенов пользователя
func (r *AuthRepo) RemoveUserTokens(ctx context.Context, uid string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM tokens WHERE uid=$1", uid)
	if err != nil {
		logger.Error("failed remove user tokens", logger.StrArg("uid", uid), logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}

// HasToken проверка наличия действующего токена
func (r *AuthRepo) HasToken(ctx context.Context, Token string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM tokens WHERE token=$1 AND expired > $2", Token, time.Now()).Scan(&count)
	if err != nil {
		logger.Error("failed has token", logger.StrArg("token", Token), logger.StrArg("error", err.Error()))
		return false, err
	}
	return count > 0, nil
}

// HasLogin проверка наличия пользователя с таким логином
func (r *AuthRepo) HasLogin(ctx context.Context, login string) (bool, error) {
	var count int
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE login=$1", login).Scan(&count)
	if err != nil {
		logger.Error("failed has login", logger.StrArg("login", login), logger.StrArg("error", err.Error()))
		return false, err
	}
	return count > 0, nil
}
