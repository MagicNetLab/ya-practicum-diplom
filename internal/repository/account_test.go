package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

const testAccountDSN = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable"

func getAccountTestDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), testAccountDSN)
	require.NoError(t, err)
	return pool
}

func getAccountTestRepo(t *testing.T) (AccountRepository, *pgxpool.Pool) {
	pool := getAccountTestDB(t)
	return NewAccountRepo(pool), pool
}

// TestAccountRepo_GetAccount проверяет получение аккаунта
func TestAccountRepo_GetAccount(t *testing.T) {
	repo, pgx := getAccountTestRepo(t)
	id := uuid.New().String()
	uid := uuid.New().String()
	login := "test-login"
	password := "test-password"
	url := "test-url"
	description := "test-description"
	createdAt := time.Now()
	updatedAt := time.Now()

	// Тестовые данные в базе данных
	sql := "INSERT INTO accounts (id, uid, login, password, url, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := pgx.Exec(
		context.Background(),
		sql,
		id, uid, login, password, url, description, createdAt, updatedAt)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pgx.Exec(context.Background(), "DELETE FROM accounts WHERE id=$1", id)
	})

	t.Run("Проверка успешного получения аккаунта по идентификатору", func(t *testing.T) {
		account, err := repo.GetAccount(context.Background(), id, uid)
		assert.NoError(t, err)
		assert.Equal(t, id, account.GetID())
		assert.Equal(t, uid, account.GetUID())
		assert.Equal(t, login, account.GetLogin())
		assert.Equal(t, password, account.GetPassword())
		assert.Equal(t, description, account.GetDescription())
		assert.Equal(t, createdAt.Format(time.DateTime), account.GetCreatedAt().Format(time.DateTime))
		assert.Equal(t, updatedAt.Format(time.DateTime), account.GetUpdatedAt().Format(time.DateTime))
	})

	t.Run("Проверка попытки получения несуществующего аккаунта", func(t *testing.T) {
		account, err := repo.GetAccount(context.Background(), "nonexistent-id", uid)
		assert.Error(t, err)
		assert.Nil(t, account)

		account, err = repo.GetAccount(context.Background(), id, uuid.New().String())
		assert.Error(t, err)
		assert.Nil(t, account)
	})

}

// TestAccountRepo_CreateAccount проверяет создание аккаунта
func TestAccountRepo_CreateAccount(t *testing.T) {
	repo, pgx := getAccountTestRepo(t)
	uid := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(context.Background(), "DELETE FROM accounts WHERE uid=$1", uid)
	})

	t.Run("Проверка успешного создания аккаунта", func(t *testing.T) {
		account, err := repo.CreateAccount(context.Background(), uid, "test-login", "test-password", "test-url", "test-description")
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, uid, account.GetUID())
		assert.Equal(t, "test-login", account.GetLogin())
		assert.Equal(t, "test-password", account.GetPassword())
		assert.Equal(t, "test-url", account.GetURL())
		assert.Equal(t, "test-description", account.GetDescription())

		var count int
		err = pgx.QueryRow(context.Background(), "SELECT COUNT(*) FROM accounts WHERE uid=$1", uid).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Проверка попытки создания аккаунта с некорректными данными", func(t *testing.T) {
		account, err := repo.CreateAccount(context.Background(), uid, "", "test-password", "test-url", "test-description")
		assert.Error(t, err)
		assert.Nil(t, account)

		account, err = repo.CreateAccount(context.Background(), uid, "test-login", "", "test-url", "test-description")
		assert.Error(t, err)
		assert.Nil(t, account)

		account, err = repo.CreateAccount(context.Background(), uid, "test-login", "test-password", "", "test-description")
		assert.Error(t, err)
		assert.Nil(t, account)
	})

}

// TestAccountRepo_RemoveAccount проверяет удаление аккаунта
func TestAccountRepo_RemoveAccount(t *testing.T) {
	repo, pgx := getAccountTestRepo(t)
	id := uuid.New().String()
	uid := uuid.New().String()
	login := "test-login"
	password := "test-password"
	url := "test-url"
	description := "test-description"
	createdAt := time.Now()
	updatedAt := time.Now()

	// Создание тестового аккаунта
	sql := "INSERT INTO accounts (id, uid, login, password, url, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := pgx.Exec(
		context.Background(),
		sql,
		id, uid, login, password, url, description, createdAt, updatedAt)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pgx.Exec(context.Background(), "DELETE FROM accounts WHERE id=$1", id)
	})

	t.Run("Проверка удаления несуществующего аккаунта", func(t *testing.T) {
		var count int
		err := pgx.QueryRow(context.Background(), "SELECT COUNT(*) FROM accounts").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		err = repo.RemoveAccount(context.Background(), "nonexistent-id", uid)
		assert.Error(t, err)

		err = pgx.QueryRow(context.Background(), "SELECT COUNT(*) FROM accounts").Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Проверка успешного удаления существующего аккаунта", func(t *testing.T) {
		var count int
		err := pgx.QueryRow(context.Background(), "SELECT COUNT(*) FROM accounts WHERE id=$1", id).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		err = repo.RemoveAccount(context.Background(), id, uid)
		assert.NoError(t, err)

		err = pgx.QueryRow(context.Background(), "SELECT COUNT(*) FROM accounts WHERE id=$1", id).Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

}

// TestAccountRepo_SearchAccounts проверяет поиск аккаунтов
func TestAccountRepo_SearchAccounts(t *testing.T) {
	repo, pgx := getAccountTestRepo(t)
	ctx := context.Background()

	// Тестовые данные
	uid1 := uuid.New().String()
	uid2 := uuid.New().String()
	testData := []models.Account{
		{ID: uuid.New().String(), UID: uid1, Login: "test-login-1", Password: "test-password-1", URL: "test-url-1", Description: "test-description-1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Login: "test-login-2", Password: "test-password-2", URL: "test-url-2", Description: "test-description-2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Login: "test-login-3", Password: "test-password-3", URL: "test-url-3", Description: "test-description-3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Login: "test-login-4", Password: "test-password-4", URL: "test-url-4", Description: "test-description-4", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Login: "test-login-5", Password: "test-password-5", URL: "test-url-5", Description: "test-description-5", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	sql := "INSERT INTO accounts (id, uid, login, password, url, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	for _, account := range testData {
		_, err := pgx.Exec(context.Background(), sql, account.ID, account.UID, account.Login, account.Password, account.URL, account.Description, account.CreatedAt, account.UpdatedAt)
		assert.NoError(t, err)
	}

	t.Cleanup(func() {
		_, _ = pgx.Exec(context.Background(), "DELETE FROM accounts WHERE uid=$1", uid1)
		_, _ = pgx.Exec(context.Background(), "DELETE FROM accounts WHERE uid=$1", uid2)
	})

	t.Run("Проверка поиска по uid", func(t *testing.T) {
		search := models.AccountSearch{UID: uid2}
		results, err := repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		search = models.AccountSearch{UID: uid1}
		results, err = repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("Проверка поиска по login", func(t *testing.T) {
		search := models.AccountSearch{Login: "test-login"}
		results, err := repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 5)

		search = models.AccountSearch{Login: "test-login-1"}
		results, err = repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	// поиск по url
	t.Run("Проверка поиска по url", func(t *testing.T) {
		search := models.AccountSearch{URL: "test-url"}
		results, err := repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 5)

		search = models.AccountSearch{URL: "test-url-1"}
		results, err = repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	t.Run("Проверка поиска по description", func(t *testing.T) {
		search := models.AccountSearch{Description: "test-description"}
		results, err := repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 5)

		search = models.AccountSearch{Description: "test-description-1"}
		results, err = repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

	// поиск по uid с limit и offset
	t.Run("Проверка поиска по uid с limit и offset", func(t *testing.T) {
		search := models.AccountSearch{UID: uid1, Limit: 5}
		results, err := repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 3)

		search = models.AccountSearch{UID: uid1, Limit: 2}
		results, err = repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 2)

		search = models.AccountSearch{UID: uid1, Limit: 2, Offset: 2}
		results, err = repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})

}
