package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestAccountModel test account model
func TestAccountModel(t *testing.T) {
	t.Run("Test ID methods", func(t *testing.T) {
		account := &Account{}
		id := uuid.New().String()

		// Test SetID
		err := account.SetID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, account.GetID())

		// Test SetID with empty string
		err = account.SetID("")
		assert.Error(t, err)

		// Test SetID with invalid UUID
		err = account.SetID("invalid-uuid")
		assert.Error(t, err)
	})

	t.Run("Test UID methods", func(t *testing.T) {
		account := &Account{}
		uid := uuid.New().String()

		// Test SetUID
		err := account.SetUID(uid)
		assert.NoError(t, err)
		assert.Equal(t, uid, account.GetUID())

		// Test SetUID with empty string
		err = account.SetUID("")
		assert.Error(t, err)

		// Test SetUID with invalid UUID
		err = account.SetUID("invalid-uuid")
		assert.Error(t, err)
	})

	t.Run("Test Login methods", func(t *testing.T) {
		account := &Account{}
		login := "test-login"

		// Test SetLogin
		err := account.SetLogin(login)
		assert.NoError(t, err)
		assert.Equal(t, login, account.GetLogin())

		// Test SetLogin with empty string
		err = account.SetLogin("")
		assert.Error(t, err)
	})

	t.Run("Test Password methods", func(t *testing.T) {
		account := &Account{}
		password := "test-password"

		// Test SetPassword
		err := account.SetPassword(password)
		assert.NoError(t, err)
		assert.Equal(t, password, account.GetPassword())

		// Test SetPassword with empty string
		err = account.SetPassword("")
		assert.Error(t, err)
	})

	t.Run("Test URL methods", func(t *testing.T) {
		account := &Account{}
		url := "https://example.com"

		// Test SetURL
		err := account.SetURL(url)
		assert.NoError(t, err)
		assert.Equal(t, url, account.GetURL())

		// Test SetURL with empty string
		err = account.SetURL("")
		assert.Error(t, err)
	})

	t.Run("Test Description methods", func(t *testing.T) {
		account := &Account{}
		description := "test description"

		// Test SetDescription
		err := account.SetDescription(description)
		assert.NoError(t, err)
		assert.Equal(t, description, account.GetDescription())

		// Test SetDescription with empty string (should not return error)
		err = account.SetDescription("")
		assert.NoError(t, err)
	})

	t.Run("Test CreatedAt methods", func(t *testing.T) {
		account := &Account{}
		createdAt := time.Now().Add(-time.Hour) // 1 hour ago

		// Test SetCreatedAt
		err := account.SetCreatedAt(createdAt)
		assert.NoError(t, err)
		assert.Equal(t, createdAt, account.GetCreatedAt())

		// Test SetCreatedAt with zero time
		err = account.SetCreatedAt(time.Time{})
		assert.Error(t, err)

		// Test SetCreatedAt with future time
		futureTime := time.Now().Add(time.Hour)
		err = account.SetCreatedAt(futureTime)
		assert.Error(t, err)
	})

	t.Run("Test UpdatedAt methods", func(t *testing.T) {
		account := &Account{}
		updatedAt := time.Now().Add(-time.Hour) // 1 hour ago

		// Test SetUpdatedAt
		err := account.SetUpdatedAt(updatedAt)
		assert.NoError(t, err)
		assert.Equal(t, updatedAt, account.GetUpdatedAt())

		// Test SetUpdatedAt with zero time
		err = account.SetUpdatedAt(time.Time{})
		assert.Error(t, err)

		// Test SetUpdatedAt with future time
		futureTime := time.Now().Add(time.Hour)
		err = account.SetUpdatedAt(futureTime)
		assert.Error(t, err)
	})
}
