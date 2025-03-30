package repository

import (
	"context"
	"errors"
	"time"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/encryptor"
)

// NewAccountRepo возвращает новый репозиторий аккаунтов
func NewAccountRepo(pool *pgxpool.Pool) AccountRepository {
	return &AccountRepo{pool: pool}
}

// AccountRepository интерфейс репозитория аккаунтов
type AccountRepository interface {
	GetAccount(ctx context.Context, id string, uid string) (models.AccountModel, error)
	CreateAccount(ctx context.Context, uid, login, password, url, description string) (models.AccountModel, error)
	RemoveAccount(ctx context.Context, id string, uid string) error
	SearchAccounts(ctx context.Context, search models.AccountSearchModel) ([]models.AccountModel, error)
}

// AccountRepo репозиторий аккаунтов
type AccountRepo struct {
	pool *pgxpool.Pool
}

// GetAccount возвращает аккаунт по идентификатору
func (a *AccountRepo) GetAccount(ctx context.Context, id string, uid string) (models.AccountModel, error) {
	account := &models.Account{}
	sql := "SELECT id, uid, login, password, url, description, created_at, updated_at FROM accounts WHERE id=$1 AND uid=$2"
	err := a.pool.QueryRow(ctx, sql, id, uid).Scan(&account.ID, &account.UID, &account.Login, &account.Password, &account.URL, &account.Description, &account.CreatedAt, &account.UpdatedAt)
	if err != nil {
		return nil, err
	}

	decryptLogin, err := encryptor.DecryptData(account.Login)
	if err != nil {
		return nil, err
	}
	account.Login = decryptLogin

	decryptPass, err := encryptor.DecryptData(account.Password)
	if err != nil {
		return nil, err
	}
	account.Password = decryptPass

	decryptDesc, err := encryptor.DecryptData(account.Description)
	if err != nil {
		return nil, err
	}
	account.Description = decryptDesc

	return account, nil
}

// CreateAccount создает аккаунт
func (a *AccountRepo) CreateAccount(ctx context.Context, uid, login, password, url, description string) (models.AccountModel, error) {
	account := &models.Account{}
	err := account.SetID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	err = account.SetUID(uid)
	if err != nil {
		return nil, err
	}

	if login == "" {
		return nil, errors.New("login is empty")
	}
	encryptLogin, err := encryptor.EncryptData(login)
	if err != nil {
		return nil, err
	}
	err = account.SetLogin(encryptLogin)
	if err != nil {
		return nil, err
	}

	if password == "" {
		return nil, errors.New("password is empty")
	}
	encryptPass, err := encryptor.EncryptData(password)
	if err != nil {
		return nil, err
	}
	err = account.SetPassword(encryptPass)
	if err != nil {
		return nil, err
	}

	err = account.SetURL(url)
	if err != nil {
		return nil, err
	}

	encryptDesc, err := encryptor.EncryptData(description)
	if err != nil {
		return nil, err
	}
	err = account.SetDescription(encryptDesc)
	if err != nil {
		return nil, err
	}

	err = account.SetCreatedAt(time.Now())
	if err != nil {
		return nil, err
	}

	err = account.SetUpdatedAt(time.Now())
	if err != nil {
		return nil, err
	}

	sql := "INSERT INTO accounts (id, uid, login, password, url, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = a.pool.Exec(ctx, sql, account.ID, account.UID, account.Login, account.Password, account.URL, account.Description, account.CreatedAt, account.UpdatedAt)
	if err != nil {
		logger.Error("account insert error", logger.StrArg("error", err.Error()))
		return nil, err
	}

	return account, nil
}

// RemoveAccount удаляет аккаунт по идентификатору
func (a *AccountRepo) RemoveAccount(ctx context.Context, id string, uid string) error {
	if id == "" {
		return errors.New("account id is empty")
	}

	res, err := a.pool.Exec(ctx, "DELETE FROM accounts WHERE id=$1 and uid=$2", id, uid)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("account not found")
	}

	return nil
}

// SearchAccounts возвращает список аккаунтов по критериям поиска
func (a *AccountRepo) SearchAccounts(ctx context.Context, search models.AccountSearchModel) ([]models.AccountModel, error) {
	sql := "SELECT id, uid, login, password, url, description, created_at, updated_at FROM accounts"
	where, values := search.GetSubQuery()
	sql += where

	rows, err := a.pool.Query(ctx, sql, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.AccountModel, 0)

	for rows.Next() {
		var account models.Account
		err = rows.Scan(&account.ID, &account.UID, &account.Login, &account.Password, &account.URL, &account.Description, &account.CreatedAt, &account.UpdatedAt)
		if err != nil {
			logger.Error("account scan error", logger.StrArg("error", err.Error()))
			continue
		}

		decryptLogin, err := encryptor.DecryptData(account.Login)
		if err != nil {
			return nil, err
		}
		account.Login = decryptLogin

		decryptPass, err := encryptor.DecryptData(account.Password)
		if err != nil {
			return nil, err
		}
		account.Password = decryptPass

		decryptDesc, err := encryptor.DecryptData(account.Description)
		if err != nil {
			return nil, err
		}
		account.Description = decryptDesc

		result = append(result, &account)
	}
	rows.Close()

	return result, nil
}
