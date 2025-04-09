package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestUser_GetUID проверяет правильность получения UID пользователя.
func TestUser_GetUID(t *testing.T) {
	user := &User{UID: "test-UID"}
	assert.Equal(t, "test-UID", user.GetUID())
}

// TestUser_SetUID проверяет правильность установки UID пользователя.
func TestUser_SetUID(t *testing.T) {
	t.Run("Проверка успешной установки UID", func(t *testing.T) {
		user := &User{}
		uid := uuid.New().String()
		err := user.SetUID(uid)
		assert.NoError(t, err)
		assert.Equal(t, uid, user.UID)
	})

	t.Run("Проверка попытки установки некорректного UID", func(t *testing.T) {
		user := &User{}
		uid := "invalid-UID"
		err := user.SetUID(uid)
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки пустого UID", func(t *testing.T) {
		user := &User{}
		err := user.SetUID("")
		assert.Error(t, err)
	})
}

// TestUser_GetLogin проверяет правильность получения логина пользователя.
func TestUser_GetLogin(t *testing.T) {
	user := &User{Login: "testuser"}
	assert.Equal(t, "testuser", user.GetLogin())
}

// TestUser_SetLogin проверяет правильность установки логина пользователя.
func TestUser_SetLogin(t *testing.T) {
	t.Run("Проверка успешной установки логина", func(t *testing.T) {
		user := &User{}
		err := user.SetLogin("newuser")
		assert.NoError(t, err)
		assert.Equal(t, "newuser", user.Login)
	})

	t.Run("Проверка попытки установки пустого логина", func(t *testing.T) {
		user := &User{}
		err := user.SetLogin("")
		assert.Error(t, err)
	})

	t.Run("Проверка установки длинного логина", func(t *testing.T) {
		user := &User{}
		longLogin := string(make([]byte, 1000))
		err := user.SetLogin(longLogin)
		assert.NoError(t, err)
		assert.Equal(t, longLogin, user.Login)
	})
}

// TestUser_GetPassword проверяет правильность получения пароля пользователя.
func TestUser_GetPassword(t *testing.T) {
	user := &User{Password: "testpass"}
	assert.Equal(t, "testpass", user.GetPassword())
}

// TestUser_SetPassword проверяет правильность установки пароля пользователя.
func TestUser_SetPassword(t *testing.T) {
	t.Run("Проверка успешной установки пароля", func(t *testing.T) {
		user := &User{}
		err := user.SetPassword("newpass")
		assert.NoError(t, err)
		assert.Equal(t, "newpass", user.Password)
	})

	t.Run("Проверка попытки установки пустого пароля", func(t *testing.T) {
		user := &User{}
		err := user.SetPassword("")
		assert.Error(t, err)
	})

	t.Run("Проверка установки длинного пароля", func(t *testing.T) {
		user := &User{}
		longPassword := string(make([]byte, 1000))
		err := user.SetPassword(longPassword)
		assert.NoError(t, err)
		assert.Equal(t, longPassword, user.Password)
	})
}

// TestUser_GetCreatedAt проверяет правильность получения времени создания пользователя.
func TestUser_GetCreatedAt(t *testing.T) {
	t.Run("Проверка правильности получения времени создания пользователя", func(t *testing.T) {
		createdAt := time.Now().Add(-time.Hour)
		user := &User{CreatedAt: createdAt}
		assert.Equal(t, createdAt, user.GetCreatedAt())
	})

	t.Run("Проверка правильности получения времени создания пользователя с нулевым значением", func(t *testing.T) {
		user := &User{CreatedAt: time.Time{}}
		createdAt := user.GetCreatedAt()
		assert.True(t, time.Since(createdAt) < time.Second)
	})
}

// TestUser_SetCreatedAt проверяет правильность установки времени создания пользователя.
func TestUser_SetCreatedAt(t *testing.T) {
	t.Run("Проверка успешной установки времени создания пользователя", func(t *testing.T) {
		user := &User{}
		createdAt := time.Now()
		err := user.SetCreatedAt(createdAt)
		assert.NoError(t, err)
		assert.Equal(t, createdAt, user.CreatedAt)
	})

	t.Run("Проверка попытки установки нулевого времени создания пользователя", func(t *testing.T) {
		user := &User{}
		createdAt := time.Time{}
		err := user.SetCreatedAt(createdAt)
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки времени создания пользователя в будущем", func(t *testing.T) {
		user := &User{}
		createdAt := time.Now().Add(time.Minute)
		err := user.SetCreatedAt(createdAt)
		assert.Error(t, err)
	})
}

// TestUser_GetUpdatedAt проверяет правильность получения времени обновления пользователя.
func TestUser_GetUpdatedAt(t *testing.T) {
	t.Run("Проверка правильности получения времени обновления пользователя", func(t *testing.T) {
		updatedAt := time.Now().Add(-time.Hour)
		user := &User{UpdatedAt: updatedAt}
		assert.Equal(t, updatedAt, user.GetUpdatedAt())
	})

	t.Run("Проверка правильности получения времени обновления пользователя с нулевым значением", func(t *testing.T) {
		user := &User{UpdatedAt: time.Time{}}
		updatedAt := user.GetUpdatedAt()
		assert.True(t, time.Since(updatedAt) < time.Second)
	})
}

// TestUser_SetUpdatedAt проверяет правильность установки времени обновления пользователя.
func TestUser_SetUpdatedAt(t *testing.T) {
	t.Run("Проверка успешной установки времени обновления пользователя", func(t *testing.T) {
		user := &User{}
		updatedAt := time.Now()
		err := user.SetUpdatedAt(updatedAt)
		assert.NoError(t, err)
		assert.Equal(t, updatedAt, user.UpdatedAt)
	})

	t.Run("Проверка попытки установки нулевого времени обновления пользователя", func(t *testing.T) {
		user := &User{}
		updatedAt := time.Time{}
		err := user.SetUpdatedAt(updatedAt)
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки времени обновления пользователя в будущем", func(t *testing.T) {
		user := &User{}
		updatedAt := time.Now().Add(time.Minute)
		err := user.SetUpdatedAt(updatedAt)
		assert.Error(t, err)
	})

	t.Run("Проверка установки времени обновления раньше времени создания", func(t *testing.T) {
		user := &User{}
		createdAt := time.Now()
		err := user.SetCreatedAt(createdAt)
		assert.NoError(t, err)

		updatedAt := createdAt.Add(-time.Hour)
		err = user.SetUpdatedAt(updatedAt)
		assert.Error(t, err)
	})

	t.Run("Проверка установки времени обновления раньше времени создания", func(t *testing.T) {
		user := &User{}
		createdAt := time.Now()
		err := user.SetCreatedAt(createdAt)
		assert.NoError(t, err)

		updatedAt := createdAt.Add(-time.Hour)
		err = user.SetUpdatedAt(updatedAt)
		assert.Error(t, err)
	})
}
