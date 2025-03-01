package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
)

// Repository интерфейс общего репозитория
type Repository interface {
	GetAuthRepo() AuthRepository
	GetAccountRepo() AccountRepository
	Close()
}

// NewRepository конструктор репозитория
func NewRepository(cnf config.DataBaseConfigurator) (Repository, error) {
	pool, err := pgxpool.New(context.TODO(), cnf.GetDSN())
	if err != nil {
		return nil, errors.Wrap(err, "failed connect to database")
	}

	r := &AppRepository{
		auth:    NewAuthRepository(pool),
		account: NewAccountRepo(pool),
	}

	return r, nil
}

// AppRepository репозиторий приложения
type AppRepository struct {
	auth    AuthRepository
	account AccountRepository
	pool    *pgxpool.Pool
}

// GetAuthRepo возвращает репозиторий авторизации
func (r *AppRepository) GetAuthRepo() AuthRepository {
	return r.auth
}

// GetAccountRepo возвращает репозиторий аккаунта
func (r *AppRepository) GetAccountRepo() AccountRepository {
	return r.account
}

// Close закрывает соединение с базой данных
func (r *AppRepository) Close() {
	r.pool.Close()
}
