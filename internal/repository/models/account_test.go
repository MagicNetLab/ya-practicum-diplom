package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestAccount_GetID проверяет правильность получения ID аккаунта.
func TestAccount_GetID(t *testing.T) {
	account := &Account{id: "test-id"}
	assert.Equal(t, "test-id", account.GetID())
}

// TestAccount_SetID проверяет правильность установки ID аккаунта.
func TestAccount_SetID(t *testing.T) {
	t.Run("Проверка успешной установки ID", func(t *testing.T) {
		account := Account{}
		id := uuid.New().String()
		err := account.SetID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, account.id)
	})

	t.Run("Проверка попытки установки некорректного ID", func(t *testing.T) {
		account := Account{}
		id := "invalid-id"
		err := account.SetID(id)
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки пустого ID", func(t *testing.T) {
		account := Account{}
		err := account.SetID("")
		assert.Error(t, err)
	})
}

// TestAccount_GetUID проверяет правильность получения UID владельца аккаунта.
func TestAccount_GetUID(t *testing.T) {
	account := &Account{uid: "test-uid"}
	assert.Equal(t, "test-uid", account.GetUID())
}

// TestAccount_SetUID проверяет правильность установки UID владельца аккаунта.
func TestAccount_SetUID(t *testing.T) {
	t.Run("Проверка успешной установки UID", func(t *testing.T) {
		account := Account{}
		uid := uuid.New().String()
		err := account.SetUID(uid)
		assert.NoError(t, err)
		assert.Equal(t, uid, account.uid)
	})

	t.Run("Проверка попытки установки некорректного UID", func(t *testing.T) {
		account := Account{}
		uid := "invalid-uid"
		err := account.SetUID(uid)
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки пустого UID", func(t *testing.T) {
		account := Account{}
		err := account.SetUID("")
		assert.Error(t, err)
	})
}

// TestAccount_GetLogin проверяет правильность получения логина аккаунта.
func TestAccount_GetLogin(t *testing.T) {
	account := &Account{login: "testlogin"}
	assert.Equal(t, "testlogin", account.GetLogin())
}

// TestAccount_SetLogin проверяет правильность установки логина аккаунта.
func TestAccount_SetLogin(t *testing.T) {
	t.Run("Проверка успешной установки логина", func(t *testing.T) {
		account := Account{}
		err := account.SetLogin("newlogin")
		assert.NoError(t, err)
		assert.Equal(t, "newlogin", account.login)
	})

	t.Run("Проверка попытки установки пустого логина", func(t *testing.T) {
		account := Account{}
		err := account.SetLogin("")
		assert.Error(t, err)
	})
}

// TestAccount_GetPassword проверяет правильность получения пароля аккаунта.
func TestAccount_GetPassword(t *testing.T) {
	account := &Account{password: "testpass"}
	assert.Equal(t, "testpass", account.GetPassword())
}

// TestAccount_SetPassword проверяет правильность установки пароля аккаунта.
func TestAccount_SetPassword(t *testing.T) {
	t.Run("Проверка успешной установки пароля", func(t *testing.T) {
		account := Account{}
		err := account.SetPassword("newpass")
		assert.NoError(t, err)
		assert.Equal(t, "newpass", account.password)
	})

	t.Run("Проверка попытки установки пустого пароля", func(t *testing.T) {
		account := Account{}
		err := account.SetPassword("")
		assert.Error(t, err)
	})
}

// TestAccount_GetURL проверяет правильность получения URL аккаунта.
func TestAccount_GetURL(t *testing.T) {
	account := &Account{url: "http://test.com"}
	assert.Equal(t, "http://test.com", account.GetURL())
}

// TestAccount_SetURL проверяет правильность установки URL аккаунта.
func TestAccount_SetURL(t *testing.T) {
	t.Run("Проверка успешной установки URL", func(t *testing.T) {
		account := Account{}
		err := account.SetURL("http://example.com")
		assert.NoError(t, err)
		assert.Equal(t, "http://example.com", account.url)
	})

	t.Run("Проверка попытки установки пустого URL", func(t *testing.T) {
		account := Account{}
		err := account.SetURL("")
		assert.Error(t, err)
	})
}

// TestAccount_GetDescription проверяет правильность получения описания аккаунта.
func TestAccount_GetDescription(t *testing.T) {
	account := &Account{description: "test description"}
	assert.Equal(t, "test description", account.GetDescription())
}

// TestAccount_SetDescription проверяет правильность установки описания аккаунта.
func TestAccount_SetDescription(t *testing.T) {
	t.Run("Проверка успешной установки описания", func(t *testing.T) {
		account := Account{}
		err := account.SetDescription("new description")
		assert.NoError(t, err)
		assert.Equal(t, "new description", account.description)
	})

	t.Run("Проверка установки пустого описания", func(t *testing.T) {
		account := Account{}
		err := account.SetDescription("")
		assert.NoError(t, err)
		assert.Equal(t, "", account.description)
	})
}

// TestAccount_GetCreatedAt проверяет правильность получения даты создания аккаунта.
func TestAccount_GetCreatedAt(t *testing.T) {
	createdAt := time.Now()
	account := &Account{createdAt: createdAt}
	assert.Equal(t, createdAt, account.GetCreatedAt())
}

// TestAccount_SetCreatedAt проверяет правильность установки даты создания аккаунта.
func TestAccount_SetCreatedAt(t *testing.T) {
	t.Run("Проверка успешной установки даты создания", func(t *testing.T) {
		account := Account{}
		createdAt := time.Now()
		err := account.SetCreatedAt(createdAt)
		assert.NoError(t, err)
		assert.Equal(t, createdAt, account.createdAt)
	})

	t.Run("Проверка попытки установки нулевой даты создания", func(t *testing.T) {
		account := Account{}
		err := account.SetCreatedAt(time.Time{})
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки даты создания в будущем", func(t *testing.T) {
		account := Account{}
		createdAt := time.Now().Add(time.Hour)
		err := account.SetCreatedAt(createdAt)
		assert.Error(t, err)
	})
}

// TestAccount_GetUpdatedAt проверяет правильность получения даты обновления аккаунта.
func TestAccount_GetUpdatedAt(t *testing.T) {
	updatedAt := time.Now()
	account := &Account{updatedAt: updatedAt}
	assert.Equal(t, updatedAt, account.GetUpdatedAt())
}

// TestAccount_SetUpdatedAt проверяет правильность установки даты обновления аккаунта.
func TestAccount_SetUpdatedAt(t *testing.T) {
	t.Run("Проверка успешной установки даты обновления", func(t *testing.T) {
		account := Account{}
		updatedAt := time.Now()
		err := account.SetUpdatedAt(updatedAt)
		assert.NoError(t, err)
		assert.Equal(t, updatedAt, account.updatedAt)
	})

	t.Run("Проверка попытки установки нулевой даты обновления", func(t *testing.T) {
		account := Account{}
		err := account.SetUpdatedAt(time.Time{})
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки даты обновления в будущем", func(t *testing.T) {
		account := Account{}
		updatedAt := time.Now().Add(time.Hour)
		err := account.SetUpdatedAt(updatedAt)
		assert.Error(t, err)
	})
}
