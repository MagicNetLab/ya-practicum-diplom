package postgres

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPostgresTokenModel_GetValue(t *testing.T) {
	tm := &TokenModel{value: "sample_value"}
	assert.Equal(t, "sample_value", tm.GetValue())
}

func TestPostgresTokenModel_SetValue(t *testing.T) {
	t.Run("Проверка успешной установки значения", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetValue("new_value")
		assert.NoError(t, err)
		assert.Equal(t, "new_value", tm.value)
	})

	t.Run("Проверка попытки установки пустого значения", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetValue("")
		assert.Error(t, err)
	})
}

func TestPostgresTokenModel_GetUID(t *testing.T) {
	tm := &TokenModel{uid: "user_1"}
	assert.Equal(t, "user_1", tm.GetUID())
}

func TestPostgresTokenModel_SetUID(t *testing.T) {
	t.Run("Проверка успешной установки UID", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetUID("user_1")
		assert.NoError(t, err)
		assert.Equal(t, "user_1", tm.uid)
	})

	t.Run("Проверка попытки установки пустого UID", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetUID("")
		assert.Error(t, err)
	})
}

func TestPostgresTokenModel_SetAndIsRefresh(t *testing.T) {
	tm := &TokenModel{}
	// Initially should be false
	assert.False(t, tm.IsRefresh())
	// After setting, flag should be true
	tm.SetRefresh()
	assert.True(t, tm.IsRefresh())
}

func TestPostgresTokenModel_SetExpired(t *testing.T) {
	t.Run("Проверка успешной установки Expired", func(t *testing.T) {
		tm := &TokenModel{}
		expTime := time.Now().Add(time.Hour)
		err := tm.SetExpired(expTime)
		assert.NoError(t, err)
		assert.Equal(t, expTime.Format(time.RFC3339), tm.expired.Format(time.RFC3339))
	})

	t.Run("Проверка попытки установки нулевого Expired", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetExpired(time.Time{})
		assert.Error(t, err)
	})

	t.Run("Проверка попытки установки времени Expired в прошлом", func(t *testing.T) {
		tm := &TokenModel{}
		pastTime := time.Now().Add(-time.Hour)
		err := tm.SetExpired(pastTime)
		assert.Error(t, err)
	})
}

func TestPostgresTokenModel_GetExpired(t *testing.T) {
	t.Run("Проверка получения корректно установленного Expired", func(t *testing.T) {
		expTime := time.Now().Add(time.Hour)
		tm := &TokenModel{expired: expTime}
		got := tm.GetExpired()
		// Compare the formatted strings to avoid minor time differences.
		assert.Equal(t, expTime.Format(time.RFC3339), got.Format(time.RFC3339))
	})

	t.Run("Проверка получения некорректно установленного Expired", func(t *testing.T) {
		tm := &TokenModel{expired: time.Time{}}
		// GetExpired resets to zero time when parsing fails.
		got := tm.GetExpired()
		assert.True(t, got.IsZero())
	})
}
