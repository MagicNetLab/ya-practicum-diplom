package inmemory

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepository_AddToken(t *testing.T) {
	repo := &Repository{
		tokens: make(map[string]TokenModel),
	}

	t.Run("Проверка успешного добавления токена", func(t *testing.T) {
		err := repo.AddToken(context.Background(), "test_token", "test_uid", false, time.Now().Add(time.Hour))
		assert.NoError(t, err)
	})

	t.Run("Проверка добавления существующего токена", func(t *testing.T) {
		err := repo.AddToken(context.Background(), "test_token", "test_uid", false, time.Now().Add(time.Hour))
		assert.Error(t, err)
		assert.Equal(t, "token already exists", err.Error())
	})
}

func TestInMemoryRepository_RemoveToken(t *testing.T) {
	repo := &Repository{
		tokens: make(map[string]TokenModel),
	}

	t.Run("Проверка успешного удаления токена", func(t *testing.T) {
		repo.tokens["test_token"] = TokenModel{Value: "test_token"}
		err := repo.RemoveToken(context.Background(), "test_token")
		assert.NoError(t, err)
		assert.NotContains(t, repo.tokens, "test_token")
	})
}

func TestInMemoryRepository_RemoveUserTokens(t *testing.T) {
	repo := &Repository{}

	repo.tokens = map[string]TokenModel{
		"token1": {Value: "token1", UID: "test_uid"},
		"token2": {Value: "token2", UID: "test_uid"},
		"token3": {Value: "token3", UID: "other_uid"},
	}

	t.Run("Проверка успешного удаления токенов пользователя", func(t *testing.T) {
		err := repo.RemoveUserTokens(context.Background(), "test_uid")
		assert.NoError(t, err)
		assert.NotContains(t, repo.tokens, "token1")
		assert.NotContains(t, repo.tokens, "token2")
		assert.Contains(t, repo.tokens, "token3")
	})
}

func TestInMemoryRepository_HasToken(t *testing.T) {
	repo := &Repository{
		tokens: make(map[string]TokenModel),
	}

	t.Run("Успешная проверка существующего токена", func(t *testing.T) {
		repo.tokens["test_token"] = TokenModel{Value: "test_token", UID: "test_uid", Refresh: false}
		found, err := repo.HasToken(context.Background(), "test_token", "test_uid", false)
		assert.NoError(t, err)
		assert.True(t, found)
	})

	t.Run("Не успешная проверка отсутствующего токена", func(t *testing.T) {
		found, err := repo.HasToken(context.Background(), "nonexistent_token", "test_uid", false)
		assert.Error(t, err)
		assert.False(t, found)
		assert.Equal(t, "token not found", err.Error())
	})
}

func TestInMemoryRepository_GetUserByUID(t *testing.T) {
	repo := &Repository{
		users: make(map[string]UserModel),
	}

	t.Run("Успешный поиск пользователя по UID", func(t *testing.T) {
		repo.users["test_uid"] = UserModel{UID: "test_uid", Login: "test_login", Password: "test_password"}
		user, err := repo.GetUserByUID(context.Background(), "test_uid")
		assert.NoError(t, err)
		assert.Equal(t, "test_uid", user.GetUID())
	})

	t.Run("Неуспешный поиск пользователя по несуществующему UID", func(t *testing.T) {
		user, err := repo.GetUserByUID(context.Background(), "nonexistent_uid")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestInMemoryRepository_GetUserByLogin(t *testing.T) {
	repo := &Repository{
		users: make(map[string]UserModel),
	}

	t.Run("Успешный поиск пользователя по логину", func(t *testing.T) {
		repo.users["test_uid"] = UserModel{UID: "test_uid", Login: "test_login", Password: "test_password"}
		user, err := repo.GetUserByLogin(context.Background(), "test_login")
		assert.NoError(t, err)
		assert.Equal(t, "test_login", user.GetLogin())
	})

	t.Run("Неуспешный поиск пользователя по несуществующему логину", func(t *testing.T) {
		user, err := repo.GetUserByLogin(context.Background(), "nonexistent_login")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestInMemoryRepository_GetUserByLoginPassword(t *testing.T) {
	repo := &Repository{
		users: make(map[string]UserModel),
	}

	t.Run("Успешный поиск пользователя по логину/паролю", func(t *testing.T) {
		repo.users["test_uid"] = UserModel{UID: "test_uid", Login: "test_login", Password: "test_password"}
		user, err := repo.GetUserByLoginPassword(context.Background(), "test_login", "test_password")
		assert.NoError(t, err)
		assert.Equal(t, "test_login", user.GetLogin())
	})

	t.Run("Неуспешный поиск пользователя по неверному логину/паролю", func(t *testing.T) {
		user, err := repo.GetUserByLoginPassword(context.Background(), "test_login", "wrong_password")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user not found", err.Error())
	})
}

func TestInMemoryRepository_CreateUser(t *testing.T) {
	repo := &Repository{
		users: make(map[string]UserModel),
	}

	t.Run("Успешное создание пользователя", func(t *testing.T) {
		user, err := repo.CreateUser(context.Background(), "test_login", "test_password")
		assert.NoError(t, err)
		assert.Equal(t, "test_login", user.GetLogin())
	})

	t.Run("Неудачное создание пользователя с существующим логином", func(t *testing.T) {
		user, err := repo.CreateUser(context.Background(), "test_login", "test_password")
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "user already exists", err.Error())
	})
}

func TestInMemoryRepository_RemoveUserByUID(t *testing.T) {
	repo := &Repository{
		users: make(map[string]UserModel),
	}

	t.Run("Успешное удаление пользователя по uid", func(t *testing.T) {
		repo.users["test_uid"] = UserModel{UID: "test_uid"}
		err := repo.RemoveUserByUID(context.Background(), "test_uid")
		assert.NoError(t, err)
		assert.NotContains(t, repo.users, "test_uid")
	})
}

func TestInMemoryRepository_RemoveUserByLogin(t *testing.T) {
	repo := &Repository{
		users: make(map[string]UserModel),
	}

	t.Run("Успешное удаление пользователя по логину", func(t *testing.T) {
		repo.users["test_uid"] = UserModel{UID: "test_uid", Login: "test_login"}
		err := repo.RemoveUserByLogin(context.Background(), "test_login")
		assert.NoError(t, err)
		assert.NotContains(t, repo.users, "test_uid")
	})
}

func TestInMemoryRepository_DumpAndImport(t *testing.T) {
	// Создаем временную директорию для тестов
	tmpDir, err := os.MkdirTemp("", "repo_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	repo := &Repository{
		users:    make(map[string]UserModel),
		tokens:   make(map[string]TokenModel),
		dumpPath: tmpDir,
	}
	ctx := context.Background()

	t.Run("Сохраняем данных в файл и чтение данных из файла", func(t *testing.T) {
		// Создаем тестовые данные
		_, _ = repo.CreateUser(ctx, "testuser", "password")

		// Сохраняем данные
		err = repo.Dump()
		assert.NoError(t, err)

		// Создаем новый репозиторий и импортируем данные
		newRepo := &Repository{
			users:    make(map[string]UserModel),
			tokens:   make(map[string]TokenModel),
			dumpPath: tmpDir,
		}

		err = newRepo.Import()
		os.RemoveAll(tmpDir)
		assert.NoError(t, err)

		// Проверяем, что данные импортировались корректно
		user, err := newRepo.GetUserByLogin(ctx, "testuser")
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.GetLogin())
	})
}
