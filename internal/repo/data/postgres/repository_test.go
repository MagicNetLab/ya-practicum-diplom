package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	_ "time"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/stretchr/testify/require"
)

func getRepository(t *testing.T) *Repository {

	cnf, err := conf.GetCnf()
	require.NoError(t, err)

	r, err := NewRepository(cnf)
	require.NoError(t, err)
	return r
}

// TestPostgresRepository_GetUserByUID проверка функции получения пользователя из базы по UID
func TestPostgresRepository_GetUserByUID(t *testing.T) {
	ctx := context.Background()
	r := getRepository(t)

	testLogin := "test_login"
	testPassword := "test_password"

	t.Cleanup(func() {
		_ = r.RemoveUserByLogin(ctx, testLogin)
	})

	// Create a user using the repository instance that wraps the transaction.
	user, err := r.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("Проверка успешного получения пользователя из базы по UID", func(t *testing.T) {
		got, err := r.GetUserByUID(ctx, user.GetUID())
		assert.NoError(t, err)
		assert.Equal(t, user.GetUID(), got.GetUID())
	})

	t.Run("Проверка попытки получения пользователя по несуществующему UID", func(t *testing.T) {
		_, err := r.GetUserByUID(ctx, uuid.New().String())
		assert.Error(t, err)
	})
}

// TestPostgresRepository_GetUserByLogin проверка функции получения пользователя из базы по логину
func TestPostgresRepository_GetUserByLogin(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testLogin := "test_login"
	testPassword := "test_password"

	t.Cleanup(func() {
		_ = r.RemoveUserByLogin(ctx, testLogin)
	})

	// подготовка тестового пользователя
	user, err := r.CreateUser(ctx, testLogin, testPassword)
	if err != nil {
		_ = r.RemoveUserByLogin(ctx, testLogin)
	}
	require.NoError(t, err)

	t.Run("Проверка успешного получения пользователя из базы по логину", func(t *testing.T) {
		got, err := r.GetUserByLogin(ctx, user.GetLogin())
		if err != nil {
			_ = r.RemoveUserByLogin(ctx, testLogin)
		}
		assert.NoError(t, err)
		assert.Equal(t, user.GetLogin(), got.GetLogin())
	})

	t.Run("Проверка попытки получения пользователя по несуществующему логину", func(t *testing.T) {
		_, err := r.GetUserByLogin(ctx, "non_existent_login")
		assert.Error(t, err)
	})
}

// TestPostgresRepository_GetUserByLoginPassword проверка функции получения пользователя из базы по логину и паролю
func TestPostgresRepository_GetUserByLoginPassword(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testLogin := "test_login"
	testPassword := "test_password"
	t.Cleanup(func() {
		_ = r.RemoveUserByLogin(ctx, testLogin)
	})

	// подготовка тестового пользователя
	user, err := r.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("Проверка успешного получения пользователя из базы по логину и паролю", func(t *testing.T) {
		got, err := r.GetUserByLoginPassword(ctx, testLogin, testPassword)
		assert.NoError(t, err)
		assert.Equal(t, user.GetLogin(), got.GetLogin())
	})

	t.Run("Проверка попытки получения пользователя из базы с некорректной парой логин/пароль", func(t *testing.T) {
		_, err := r.GetUserByLoginPassword(ctx, testLogin, "wrong_password")
		assert.Error(t, err)
	})
}

// TestPostgresRepository_CreateUser проверка функции создания пользователя в базеs
func TestPostgresRepository_CreateUser(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testLogin := "test_login"
	testPassword := "test_password"
	t.Cleanup(func() {
		_ = r.RemoveUserByLogin(ctx, testLogin)
	})

	t.Run("Проверка успешного создания пользователя", func(t *testing.T) {
		user, err := r.CreateUser(ctx, testLogin, testPassword)
		assert.NoError(t, err)
		assert.Equal(t, testLogin, user.GetLogin())
	})

	t.Run("Проверка попытки создания пользователя с уже существующим логином", func(t *testing.T) {
		_, err := r.CreateUser(ctx, testLogin, testPassword)
		assert.Error(t, err)
	})
}

// TestPostgresRepository_RemoveUserByUID проверка функции удаления пользователя из базы по UID
func TestPostgresRepository_RemoveUserByUID(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testLogin := "test_login"
	testPassword := "test_password"
	t.Cleanup(func() {
		_, _ = r.pool.Exec(ctx, "DELETE FROM users WHERE login=$1", testLogin)
	})

	// Подготовка тестового пользователя
	user, err := r.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("Проверка успешного удаления пользователя по UID", func(t *testing.T) {
		err := r.RemoveUserByUID(ctx, user.GetUID())
		assert.NoError(t, err)

		_, err = r.GetUserByUID(ctx, user.GetUID())
		assert.Error(t, err)
	})
}

// TestPostgresRepository_RemoveUserByLogin проверка функции удаления пользователя из базы по логину
func TestPostgresRepository_RemoveUserByLogin(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testLogin := "test_login"
	testPassword := "test_password"
	t.Cleanup(func() {
		_, _ = r.pool.Exec(ctx, "DELETE FROM users WHERE login=$1", testLogin)
	})

	// Создание тестового пользователя
	user, err := r.CreateUser(ctx, testLogin, testPassword)
	require.NoError(t, err)

	t.Run("Проверка успешного удаления пользователя по логину", func(t *testing.T) {
		err := r.RemoveUserByLogin(ctx, user.GetLogin())
		assert.NoError(t, err)
		_, err = r.GetUserByLogin(ctx, user.GetLogin())
		assert.Error(t, err)
	})
}

// TestPostgresRepository_AddToken проверка функции добавления токена в базу
func TestPostgresRepository_AddToken(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testUid := uuid.New().String()
	testToken := "test_token"
	t.Cleanup(func() {
		_, _ = r.pool.Exec(ctx, "DELETE FROM tokens WHERE token=$1", testToken)
	})

	t.Run("Проверка успешного добавления токена в базу", func(t *testing.T) {
		err := r.AddToken(ctx, testToken, testUid, false, time.Now().Add(time.Hour))
		assert.NoError(t, err)
	})

	t.Run("Проверка попытки добавления существующего токена в базу", func(t *testing.T) {
		err := r.AddToken(ctx, testToken, testUid, false, time.Now().Add(time.Hour))
		assert.Error(t, err)
	})
}

// TestPostgresRepository_RemoveToken проверка функции удаления токена из базы
func TestPostgresRepository_RemoveToken(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testUid := uuid.New().String()
	testToken := "test_token"
	t.Cleanup(func() {
		_, _ = r.pool.Exec(ctx, "DELETE FROM tokens WHERE token=$1", testToken)
	})

	err := r.AddToken(ctx, testToken, testUid, false, time.Now().Add(time.Hour))
	require.NoError(t, err)

	t.Run("Проверка успешного удаления токена", func(t *testing.T) {
		err := r.RemoveToken(ctx, testToken)
		assert.NoError(t, err)

		row := r.pool.QueryRow(ctx, "SELECT count(*) FROM tokens WHERE token=$1", testToken)
		var count int
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)

	})

}

// TestPostgresRepository_RemoveUserTokens проверка функции удаления всех токенов пользователя из базы
func TestPostgresRepository_RemoveUserTokens(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testUid1 := uuid.New().String()
	testUid2 := uuid.New().String()
	testToken1 := "test_token1"
	testToken2 := "test_token2"
	testToken3 := "test_token3"

	t.Cleanup(func() {
		_, _ = r.pool.Exec(ctx, "DELETE FROM tokens WHERE token IN ($1,$2,$3)", testToken1, testToken2, testToken3)
	})

	err := r.AddToken(ctx, testToken1, testUid1, false, time.Now().Add(time.Hour))
	require.NoError(t, err)
	err = r.AddToken(ctx, testToken2, testUid1, false, time.Now().Add(time.Hour))
	require.NoError(t, err)
	err = r.AddToken(ctx, testToken3, testUid2, false, time.Now().Add(time.Hour))
	require.NoError(t, err)

	t.Run("Проверка успешного удаления токенов пользователя", func(t *testing.T) {
		err := r.RemoveUserTokens(ctx, testUid1)
		assert.NoError(t, err)

		row := r.pool.QueryRow(ctx, "SELECT count(*) FROM tokens WHERE uid=$1", testUid1)
		var count int
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)

		row = r.pool.QueryRow(ctx, "SELECT count(*) FROM tokens WHERE uid=$1", testUid2)
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

// TestPostgresRepository_HasToken проверка функции проверки наличия токена в базе
func TestPostgresRepository_HasToken(t *testing.T) {
	r := getRepository(t)
	ctx := context.Background()

	testToken := "test_token"
	testUid := uuid.New().String()

	t.Cleanup(func() {
		_, _ = r.pool.Exec(ctx, "DELETE FROM tokens WHERE token=$1", testToken)
	})

	err := r.AddToken(ctx, testToken, testUid, false, time.Now().Add(time.Hour))
	require.NoError(t, err)

	t.Run("Проверка успешной проверки наличия токена в базе", func(t *testing.T) {
		exists, err := r.HasToken(ctx, testToken, testUid, false)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Проверка попытки проверки наличия несущего токена в базе", func(t *testing.T) {
		exists, err := r.HasToken(ctx, "non_existent_token", testUid, false)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
