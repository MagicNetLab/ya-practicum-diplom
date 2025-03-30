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
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/encryptor"
)

const testAccountDSN = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable"

func getAccountTestDB(t *testing.T) *pgxpool.Pool {
	t.Setenv("ENCRYPT_KEY", "test-encryption-key-32-bytes-length!")

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

	// Encrypt sensitive data
	encryptedLogin, err := encryptor.EncryptData(login)
	assert.NoError(t, err)
	encryptedPassword, err := encryptor.EncryptData(password)
	assert.NoError(t, err)
	encryptedDescription, err := encryptor.EncryptData(description)
	assert.NoError(t, err)

	// Тестовые данные в базе данных
	sql := "INSERT INTO accounts (id, uid, login, password, url, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err = pgx.Exec(
		context.Background(),
		sql,
		id, uid, encryptedLogin, encryptedPassword, url, encryptedDescription, createdAt, updatedAt)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pgx.Exec(context.Background(), "DELETE FROM accounts WHERE id=$1", id)
	})

	t.Run("Проверка успешного получения и дешифрования данных аккаунта", func(t *testing.T) {
		account, err := repo.GetAccount(context.Background(), id, uid)
		assert.NoError(t, err)
		assert.Equal(t, id, account.GetID())
		assert.Equal(t, uid, account.GetUID())
		assert.Equal(t, login, account.GetLogin())
		assert.Equal(t, password, account.GetPassword())
		assert.Equal(t, url, account.GetURL())
		assert.Equal(t, description, account.GetDescription())
		assert.Equal(t, createdAt.Format(time.DateTime), account.GetCreatedAt().Format(time.DateTime))
		assert.Equal(t, updatedAt.Format(time.DateTime), account.GetUpdatedAt().Format(time.DateTime))

		// Проверка что данные действительно зашифрованы в БД
		var dbLogin, dbPassword, dbDescription string
		err = pgx.QueryRow(context.Background(), "SELECT login, password, description FROM accounts WHERE id=$1", id).Scan(&dbLogin, &dbPassword, &dbDescription)
		assert.NoError(t, err)
		assert.NotEqual(t, login, dbLogin)
		assert.NotEqual(t, password, dbPassword)
		assert.NotEqual(t, description, dbDescription)
	})

	t.Run("Проверка попытки получения несуществующего аккаунта", func(t *testing.T) {
		account, err := repo.GetAccount(context.Background(), "nonexistent-id", uid)
		assert.Error(t, err)
		assert.Nil(t, account)

		account, err = repo.GetAccount(context.Background(), id, uuid.New().String())
		assert.Error(t, err)
		assert.Nil(t, account)
	})

	t.Run("Проверка ошибки дешифрования данных", func(t *testing.T) {
		// Вставляем некорректно зашифрованные данные
		_, err = pgx.Exec(
			context.Background(),
			"UPDATE accounts SET login='invalid-encrypted-data' WHERE id=$1",
			id)
		assert.NoError(t, err)

		account, err := repo.GetAccount(context.Background(), id, uid)
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
		login := "test-login"
		password := "test-password"
		url := "test-url"
		description := "test-description"

		account, err := repo.CreateAccount(context.Background(), uid, login, password, url, description)
		assert.NoError(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, uid, account.GetUID())

		assert.Equal(t, url, account.GetURL())

		decryptedDescription, err := encryptor.DecryptData(account.GetDescription())
		assert.NoError(t, err)
		assert.Equal(t, description, decryptedDescription)
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
	uid := uuid.New().String()
	login := "test-login"
	password := "test-password"
	url := "test-url"
	description := "test-description"

	// Создаем тестовый аккаунт
	account, err := repo.CreateAccount(context.Background(), uid, login, password, url, description)
	assert.NoError(t, err)
	t.Cleanup(func() {
		_, _ = pgx.Exec(context.Background(), "DELETE FROM accounts WHERE uid=$1", uid)
	})

	t.Run("Проверка успешного удаления аккаунта", func(t *testing.T) {
		err := repo.RemoveAccount(context.Background(), account.GetID(), uid)
		assert.NoError(t, err)

		// Проверяем что аккаунт действительно удален
		_, err = repo.GetAccount(context.Background(), account.GetID(), uid)
		assert.Error(t, err)
	})

	t.Run("Проверка удаления с пустым ID", func(t *testing.T) {
		err := repo.RemoveAccount(context.Background(), "", uid)
		assert.Error(t, err)
	})

	t.Run("Проверка удаления с некорректным UID", func(t *testing.T) {
		// Создаем новый аккаунт для теста
		account, err := repo.CreateAccount(context.Background(), uid, login, password, url, description)
		assert.NoError(t, err)

		// Пытаемся удалить с неправильным UID
		err = repo.RemoveAccount(context.Background(), account.GetID(), uuid.New().String())
		assert.Error(t, err)

		// Проверяем что аккаунт не был удален
		_, err = repo.GetAccount(context.Background(), account.GetID(), uid)
		assert.NoError(t, err)
	})

	t.Run("Проверка удаления несуществующего аккаунта", func(t *testing.T) {
		err := repo.RemoveAccount(context.Background(), uuid.New().String(), uid)
		assert.Error(t, err)
	})
}

// TestAccountRepo_SearchAccounts проверяет поиск аккаунтов
func TestAccountRepo_SearchAccounts(t *testing.T) {
	repo, pgx := getAccountTestRepo(t)
	ctx := context.Background()

	// Тестовые данные
	uid1 := uuid.New().String()
	uid2 := uuid.New().String()
	testData := make([]models.Account, 0)
	testData = append(testData, models.Account{ID: uuid.New().String(), UID: uid1, Login: "test-login-1", Password: "test-password-1", URL: "test-url-1", Description: "test-description-1", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	testData = append(testData, models.Account{ID: uuid.New().String(), UID: uid1, Login: "test-login-2", Password: "test-password-2", URL: "test-url-2", Description: "test-description-2", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	testData = append(testData, models.Account{ID: uuid.New().String(), UID: uid1, Login: "test-login-3", Password: "test-password-3", URL: "test-url-3", Description: "test-description-3", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	testData = append(testData, models.Account{ID: uuid.New().String(), UID: uid2, Login: "test-login-4", Password: "test-password-4", URL: "test-url-4", Description: "test-description-4", CreatedAt: time.Now(), UpdatedAt: time.Now()})
	testData = append(testData, models.Account{ID: uuid.New().String(), UID: uid2, Login: "test-login-5", Password: "test-password-5", URL: "test-url-5", Description: "test-description-5", CreatedAt: time.Now(), UpdatedAt: time.Now()})

	// Encrypt test data
	for _, account := range testData {
		_, err := repo.CreateAccount(ctx, account.UID, account.Login, account.Password, account.URL, account.Description)
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

	// поиск по url
	t.Run("Проверка поиска по тексту", func(t *testing.T) {
		search := models.AccountSearch{Search: "test-url"}
		results, err := repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 5)

		search = models.AccountSearch{Search: "test-url-1"}
		results, err = repo.SearchAccounts(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, testData[0].URL, results[0].GetURL())
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

	t.Run("Проверка поиска с limit и offset", func(t *testing.T) {
		search := models.AccountSearch{
			UID:    uid1,
			Limit:  2,
			Offset: 1,
		}
		result, err := repo.SearchAccounts(context.Background(), search)
		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})

	t.Run("Проверка поиска по несуществующему UID", func(t *testing.T) {
		search := models.AccountSearch{
			UID: uuid.New().String(),
		}
		result, err := repo.SearchAccounts(context.Background(), search)
		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})

	t.Run("Проверка дешифрования данных при поиске", func(t *testing.T) {
		search := models.AccountSearch{
			UID: uid1,
		}
		result, err := repo.SearchAccounts(context.Background(), search)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)

		// Проверяем что данные расшифрованы
		for _, acc := range result {
			assert.Contains(t, acc.GetLogin(), "test-login")
			assert.Contains(t, acc.GetDescription(), "test-description")
		}
	})

}
