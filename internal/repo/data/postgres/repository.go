package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/models"
)

type Repository struct {
	pool *pgxpool.Pool
}

// GetUserByUID получение пользователя по UID
func (r Repository) GetUserByUID(ctx context.Context, uid string) (models.UserEntity, error) {
	user := UserModel{}
	err := r.pool.QueryRow(ctx, "SELECT * FROM users WHERE uid=$1", uid).Scan(&user)
	if err != nil {
		logger.Error("failed get user by uid", logger.StrArg("uid", uid), logger.StrArg("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

// GetUserByLogin получение пользователя по логину
func (r Repository) GetUserByLogin(ctx context.Context, login string) (models.UserEntity, error) {
	user := UserModel{}
	err := r.pool.QueryRow(ctx, "SELECT * FROM users WHERE login=$1", login).Scan(&user)
	if err != nil {
		logger.Error("failed get user by login", logger.StrArg("login", login), logger.StrArg("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

// GetUserByLoginPassword получение пользователя по логину и паролю
func (r Repository) GetUserByLoginPassword(ctx context.Context, login, password string) (models.UserEntity, error) {
	user := UserModel{}
	err := r.pool.QueryRow(ctx, "SELECT * FROM users WHERE login=$1 AND password=$2", login, password).Scan(&user)
	if err != nil {
		logger.Error("failed get user by login and password", logger.StrArg("login", login), logger.StrArg("error", err.Error()))
		return nil, err
	}
	return &user, nil
}

// CreateUser создание пользователя
func (r Repository) CreateUser(ctx context.Context, login, password string) (models.UserEntity, error) {
	uid := uuid.New().String()
	sql := "INSERT INTO users (uid, login, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.pool.Exec(ctx, sql, uid, login, password, time.Now(), time.Now())
	if err != nil {
		return nil, err
	}

	user := UserModel{}
	err = r.pool.QueryRow(ctx, "SELECT * FROM users WHERE uid=$1", uid).Scan(&user)
	if err != nil {
		logger.Error("failed get user by uid after create", logger.StrArg("uid", uid), logger.StrArg("login", login), logger.StrArg("error", err.Error()))
		return nil, err
	}

	return &user, nil
}

// RemoveUserByUID удаление пользователя по UID
func (r Repository) RemoveUserByUID(ctx context.Context, uid string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM users WHERE uid=$1", uid)
	if err != nil {
		logger.Error("failed remove user by uid", logger.StrArg("uid", uid), logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}

// RemoveUserByLogin удаление пользователя по логину
func (r Repository) RemoveUserByLogin(ctx context.Context, login string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM users WHERE login=$1", login)
	if err != nil {
		logger.Error("failed remove user by login", logger.StrArg("login", login), logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}

// AddToken добавление токена в бд
func (r Repository) AddToken(ctx context.Context, value, uid string, isRefresh bool, expired time.Time) error {
	_, err := r.pool.Exec(ctx, "INSERT INTO tokens (token, uid, is_refresh, expired) VALUES ($1, $2, $3, $4)", value, uid, isRefresh, expired.Format(time.RFC3339))
	if err != nil {
		logger.Error("failed add token", logger.StrArg("token", value), logger.StrArg("uid", uid), logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}

// RemoveToken удаление токена из бд
func (r Repository) RemoveToken(ctx context.Context, value string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM tokens WHERE token=$1", value)
	if err != nil {
		logger.Error("failed remove token", logger.StrArg("token", value), logger.StrArg("error", err.Error()))
		return err
	}
	return nil
}

// RemoveUserTokens удаление всех токенов пользователя
func (r Repository) RemoveUserTokens(ctx context.Context, uid string) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM tokens WHERE uid=$1", uid)
	if err != nil {
		return err
	}
	return nil
}

// HasToken проверка токена в бд
func (r Repository) HasToken(ctx context.Context, value, uid string, isRefresh bool) (bool, error) {
	count := 0
	err := r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM tokens WHERE token=$1 AND uid=$2 AND is_refresh=$3 AND expired > NOW()", value, uid, isRefresh).Scan(&count)
	if err != nil {
		logger.Error("failed check token", logger.StrArg("token", value), logger.StrArg("uid", uid), logger.StrArg("error", err.Error()))
		return false, err
	}
	return count > 0, nil
}

func (r Repository) Close(ctx context.Context) error {
	r.pool.Close()
	return nil
}
