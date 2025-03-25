package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestToken_GetID проверяет правильность получения ID токена.
func TestToken_GetID(t *testing.T) {
	token := &Token{id: "test-ID"}
	assert.Equal(t, "test-ID", token.GetID())
}

// TestToken_SetID проверяет правильность установки ID токена.
func TestToken_SetID(t *testing.T) {
	t.Run("Проверка успешной установки ID", func(t *testing.T) {
		token := Token{}
		id := uuid.New().String()
		err := token.SetID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, token.id)
	})

	t.Run("Проверка попытки установки некорректного ID", func(t *testing.T) {
		token := Token{}
		id := "invalid-ID"
		err := token.SetID(id)
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки пустого ID", func(t *testing.T) {
		token := Token{}
		err := token.SetID("")
		assert.Error(t, err)
	})
}

// TestToken_GetToken проверяет правильность получения значения токена.
func TestToken_GetToken(t *testing.T) {
	token := &Token{token: "test-token"}
	assert.Equal(t, "test-token", token.GetToken())
}

// TestToken_SetToken проверяет правильность установки значения токена.
func TestToken_SetToken(t *testing.T) {
	t.Run("Проверка успешной установки значения токена", func(t *testing.T) {
		token := &Token{}
		err := token.SetToken("new-token")
		assert.NoError(t, err)
		assert.Equal(t, "new-token", token.token)
	})

	t.Run("Проверка попытки установки пустого значения токена", func(t *testing.T) {
		token := &Token{}
		err := token.SetToken("")
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки очень длинного значения токена", func(t *testing.T) {
		token := &Token{}
		longToken := string(make([]byte, 1000))
		err := token.SetToken(longToken)
		assert.NoError(t, err)
		assert.Equal(t, longToken, token.token)
	})
}

// TestToken_GetUID проверяет правильность получения UID владельца токена.
func TestToken_GetUID(t *testing.T) {
	token := &Token{uid: "test-UID"}
	assert.Equal(t, "test-UID", token.GetUID())
}

// TestToken_SetUID проверяет правильность установки UID владельца токена.
func TestToken_SetUID(t *testing.T) {
	t.Run("Проверка успешной установки UID", func(t *testing.T) {
		token := &Token{}
		uid := uuid.New().String()
		err := token.SetUID(uid)
		assert.NoError(t, err)
		assert.Equal(t, uid, token.uid)
	})

	t.Run("Проверка попытки установки некорректного UID", func(t *testing.T) {
		token := &Token{}
		uid := "invalid-UID"
		err := token.SetUID(uid)
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки пустого UID", func(t *testing.T) {
		token := &Token{}
		err := token.SetUID("")
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки некорректного формата UUID", func(t *testing.T) {
		token := &Token{}
		uid := "123e4567-e89b-12d3-a456-42661417400" // Неполный UUID
		err := token.SetUID(uid)
		assert.Error(t, err)
	})
}

// TestToken_IsRefresh проверяет правильность получения флага refresh токена.
func TestToken_IsRefresh(t *testing.T) {
	t.Run("Проверка получения true значения", func(t *testing.T) {
		token := &Token{isRefresh: true}
		assert.True(t, token.IsRefresh())
	})

	t.Run("Проверка получения false значения", func(t *testing.T) {
		token := &Token{isRefresh: false}
		assert.False(t, token.IsRefresh())
	})
}

// TestToken_SetIsRefresh проверяет правильность установки флага refresh токена.
func TestToken_SetIsRefresh(t *testing.T) {
	t.Run("Проверка установки true значения", func(t *testing.T) {
		token := &Token{}
		err := token.SetIsRefresh(true)
		assert.NoError(t, err)
		assert.True(t, token.isRefresh)
	})

	t.Run("Проверка установки false значения", func(t *testing.T) {
		token := &Token{}
		err := token.SetIsRefresh(false)
		assert.NoError(t, err)
		assert.False(t, token.isRefresh)
	})
}

// TestToken_GetExpired проверяет правильность получения времени истечения токена.
func TestToken_GetExpired(t *testing.T) {
	expiredTime := time.Now().Add(time.Hour)
	token := &Token{expired: expiredTime}
	assert.Equal(t, expiredTime, token.GetExpired())
}

// TestToken_SetExpired проверяет правильность установки времени истечения токена.
func TestToken_SetExpired(t *testing.T) {
	t.Run("Проверка успешной установки времени истечения", func(t *testing.T) {
		token := &Token{}
		expiredTime := time.Now().Add(time.Hour)
		err := token.SetExpired(expiredTime)
		assert.NoError(t, err)
		assert.Equal(t, expiredTime, token.expired)
	})

	t.Run("Проверка попытки установки нулевого времени истечения", func(t *testing.T) {
		token := &Token{}
		err := token.SetExpired(time.Time{})
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки времени истечения в прошлом", func(t *testing.T) {
		token := &Token{}
		expiredTime := time.Now().Add(-time.Hour)
		err := token.SetExpired(expiredTime)
		assert.Error(t, err)
	})

	t.Run("Проверка установки времени истечения в текущий момент", func(t *testing.T) {
		token := &Token{}
		expiredTime := time.Now()
		err := token.SetExpired(expiredTime)
		assert.Error(t, err)
	})

	t.Run("Проверка установки времени истечения в далеком будущем", func(t *testing.T) {
		token := &Token{}
		expiredTime := time.Now().AddDate(1, 0, 0) // год в будущем
		err := token.SetExpired(expiredTime)
		assert.NoError(t, err)
		assert.Equal(t, expiredTime, token.expired)
	})
}
