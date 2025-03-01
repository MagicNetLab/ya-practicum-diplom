package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

func NewAccountRepo(pool *pgxpool.Pool) AccountRepository {
	return &AccountRepo{pool: pool}
}

type AccountRepository interface {
	GetAccount(ctx context.Context, id string) (models.AccountModel, error)
	CreateAccount(ctx context.Context, uid, login, password, url, description string) (models.AccountModel, error)
	RemoveAccount(ctx context.Context, id string) error
	SearchAccounts(ctx context.Context, search models.AccountSearchModel) ([]models.AccountModel, error)
	UpdateAccount(ctx context.Context, account models.AccountModel) error
}

type AccountRepo struct {
	pool *pgxpool.Pool
}

func (a *AccountRepo) GetAccount(ctx context.Context, id string) (models.AccountModel, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepo) CreateAccount(ctx context.Context, uid, login, password, url, description string) (models.AccountModel, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepo) RemoveAccount(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepo) SearchAccounts(ctx context.Context, search models.AccountSearchModel) ([]models.AccountModel, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AccountRepo) UpdateAccount(ctx context.Context, account models.AccountModel) error {
	//TODO implement me
	panic("implement me")
}
