package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testDSN = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable"

func getTestDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), testDSN)
	require.NoError(t, err)
	return pool
}

func getTestRepo(t *testing.T) AuthRepository {
	pool := getTestDB(t)
	return NewAuthRepository(pool)
}

// TestAuthRepo_GetUserByUID проверка получения пользователя по UID
func TestAuthRepo_GetUserByUID(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testLogin := "test_login"
	testPassword := "test_password"

	// Создание тестового пользователя
	user, err := repo.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = repo.RemoveUser(ctx, user.GetUID())
	})

	t.Run("Проверка успешного получения пользователя по UID", func(t *testing.T) {
		got, err := repo.GetUserByUID(ctx, user.GetUID())
		assert.NoError(t, err)
		assert.Equal(t, user.GetUID(), got.GetUID())
		assert.Equal(t, testLogin, got.GetLogin())
	})

	t.Run("Проверка ошибки получения пользователя с несуществующим UID", func(t *testing.T) {
		_, err := repo.GetUserByUID(ctx, uuid.New().String())
		assert.Error(t, err)
	})
}

// TestAuthRepo_GetUserByLogin проверка получения пользователя по логину
func TestAuthRepo_GetUserByLogin(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testLogin := "test_login"
	testPassword := "test_password"

	// создание тестового пользователя
	user, err := repo.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = repo.RemoveUser(ctx, user.GetUID())
	})

	t.Run("Проверка успешного получения пользователя по логину", func(t *testing.T) {
		got, err := repo.GetUserByLogin(ctx, testLogin)
		assert.NoError(t, err)
		assert.Equal(t, testLogin, got.GetLogin())
	})

	t.Run("Проверка ошибки получения пользователя с несуществующим логином", func(t *testing.T) {
		_, err := repo.GetUserByLogin(ctx, "non_existent_login")
		assert.Error(t, err)
	})
}

// TestAuthRepo_GetUserByLoginAndPassword проверка получения пользователя по логину и паролю
func TestAuthRepo_GetUserByLoginAndPassword(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testLogin := "test_login"
	testPassword := "test_password"

	// создание тестового пользователя
	user, err := repo.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = repo.RemoveUser(ctx, user.GetUID())
	})

	t.Run("Проверка успешного получения пользователя по логину и паролю", func(t *testing.T) {
		got, err := repo.GetUserByLoginAndPassword(ctx, testLogin, testPassword)
		assert.NoError(t, err)
		assert.Equal(t, testLogin, got.GetLogin())
	})

	t.Run("Проверка ошибки получения пользователя с некорректным паролем", func(t *testing.T) {
		_, err := repo.GetUserByLoginAndPassword(ctx, testLogin, "wrong_password")
		assert.Error(t, err)
	})

	t.Run("Проверка ошибки получения пользователя с несуществующим логином", func(t *testing.T) {
		_, err := repo.GetUserByLoginAndPassword(ctx, "non_existent_login", testPassword)
		assert.Error(t, err)
	})
}

// TestAuthRepo_CreateUser проверка создания пользователя
func TestAuthRepo_CreateUser(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testLogin := "test_login"
	testPassword := "test_password"

	t.Run("Проверка успешного создания пользователя", func(t *testing.T) {
		user, err := repo.CreateUser(ctx, testLogin, testPassword)
		assert.NoError(t, err)
		assert.Equal(t, testLogin, user.GetLogin())

		// Cleanup
		t.Cleanup(func() {
			_ = repo.RemoveUser(ctx, user.GetUID())
		})
	})

	t.Run("Проверка ошибки создания пользователя с уже существующим логином", func(t *testing.T) {
		// Создание первого пользователя
		user1, err := repo.CreateUser(ctx, testLogin, testPassword)
		require.NoError(t, err)

		// Попытка создать второго пользователя с тем же логином
		_, err = repo.CreateUser(ctx, testLogin, "different_password")
		assert.Error(t, err)

		t.Cleanup(func() {
			_ = repo.RemoveUser(ctx, user1.GetUID())
		})
	})
}

// TestAuthRepo_RemoveUser проверка удаления пользователя
func TestAuthRepo_RemoveUser(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testLogin := "test_login"
	testPassword := "test_password"

	// создание тестового пользователя
	user, err := repo.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("Удаление существующего пользователя", func(t *testing.T) {
		err := repo.RemoveUser(ctx, user.GetUID())
		assert.NoError(t, err)

		// Verify user is removed
		_, err = repo.GetUserByUID(ctx, user.GetUID())
		assert.Error(t, err)
	})

	t.Run("Попытка удаления несуществующего пользователя", func(t *testing.T) {
		err := repo.RemoveUser(ctx, uuid.New().String())
		assert.NoError(t, err)
	})
}

// TestAuthRepo_CreateToken проверка создания токена
func TestAuthRepo_CreateToken(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testUID := uuid.New().String()
	testToken := "test_token"
	expiredTime := time.Now().Add(time.Hour)

	t.Run("Успешное создание токена", func(t *testing.T) {
		err := repo.CreateToken(ctx, testUID, testToken, false, expiredTime)
		assert.NoError(t, err)

		exists, err := repo.HasToken(ctx, testToken)
		assert.NoError(t, err)
		assert.True(t, exists)

		t.Cleanup(func() {
			_ = repo.RemoveToken(ctx, testToken)
		})
	})

	t.Run("Проверка создания дубликата токена", func(t *testing.T) {
		err := repo.CreateToken(ctx, testUID, testToken, false, expiredTime)
		require.NoError(t, err)

		err = repo.CreateToken(ctx, testUID, testToken, false, expiredTime)
		assert.Error(t, err)

		t.Cleanup(func() {
			_ = repo.RemoveToken(ctx, testToken)
		})
	})
}

// TestAuthRepo_RemoveToken проверка удаления токена
func TestAuthRepo_RemoveToken(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testUID := uuid.New().String()
	testToken := "test_token"
	expiredTime := time.Now().Add(time.Hour)

	err := repo.CreateToken(ctx, testUID, testToken, false, expiredTime)
	require.NoError(t, err)

	t.Run("Проверка успешного удаления существующего токена", func(t *testing.T) {
		err := repo.RemoveToken(ctx, testToken)
		assert.NoError(t, err)

		exists, err := repo.HasToken(ctx, testToken)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Проверка попытки удаления несуществующего токена", func(t *testing.T) {
		err := repo.RemoveToken(ctx, "non_existent_token")
		assert.NoError(t, err)
	})
}

// TestAuthRepo_RemoveUserTokens проверка удаления всех токенов пользователя
func TestAuthRepo_RemoveUserTokens(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testUID1 := uuid.New().String()
	testUID2 := uuid.New().String()
	testToken1 := "test_token1"
	testToken2 := "test_token2"
	testToken3 := "test_token3"
	expiredTime := time.Now().Add(time.Hour)

	require.NoError(t, repo.CreateToken(ctx, testUID1, testToken1, false, expiredTime))
	require.NoError(t, repo.CreateToken(ctx, testUID1, testToken2, false, expiredTime))
	require.NoError(t, repo.CreateToken(ctx, testUID2, testToken3, false, expiredTime))

	t.Cleanup(func() {
		_ = repo.RemoveToken(ctx, testToken1)
		_ = repo.RemoveToken(ctx, testToken2)
		_ = repo.RemoveToken(ctx, testToken3)
	})

	t.Run("Успешное удаление всех токенов пользователя", func(t *testing.T) {
		err := repo.RemoveUserTokens(ctx, testUID1)
		assert.NoError(t, err)

		exists, err := repo.HasToken(ctx, testToken1)
		assert.NoError(t, err)
		assert.False(t, exists)

		exists, err = repo.HasToken(ctx, testToken2)
		assert.NoError(t, err)
		assert.False(t, exists)

		exists, err = repo.HasToken(ctx, testToken3)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}

// TestAuthRepo_HasToken Проверка наличия токена в базе данных
func TestAuthRepo_HasToken(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testUID := uuid.New().String()
	testToken := "test_token"
	expiredTime := time.Now().Add(time.Hour)

	err := repo.CreateToken(ctx, testUID, testToken, false, expiredTime)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = repo.RemoveToken(ctx, testToken)
	})

	t.Run("Проверка наличия существующего токена", func(t *testing.T) {
		exists, err := repo.HasToken(ctx, testToken)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Проверка отсутствия несуществующего токена", func(t *testing.T) {
		exists, err := repo.HasToken(ctx, "non_existent_token")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Проверка наличия истекшего токена", func(t *testing.T) {
		// Create expired token
		expiredToken := "expired_token"
		err := repo.CreateToken(ctx, testUID, expiredToken, false, time.Now().Add(-time.Hour))
		require.NoError(t, err)

		exists, err := repo.HasToken(ctx, expiredToken)
		assert.NoError(t, err)
		assert.False(t, exists)

		t.Cleanup(func() {
			_ = repo.RemoveToken(ctx, expiredToken)
		})
	})
}

// TestAuthRepo_HasLogin проверка наличия логина в базе данных
func TestAuthRepo_HasLogin(t *testing.T) {
	ctx := context.Background()
	repo := getTestRepo(t)

	testLogin := "test_login"
	testPassword := "test_password"

	user, err := repo.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = repo.RemoveUser(ctx, user.GetUID())
	})

	t.Run("Успешная проверка наличия существующего логина", func(t *testing.T) {
		exists, err := repo.HasLogin(ctx, testLogin)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Проверка попытки проверки наличия несуществующего логина", func(t *testing.T) {
		exists, err := repo.HasLogin(ctx, "non_existent_login")
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
