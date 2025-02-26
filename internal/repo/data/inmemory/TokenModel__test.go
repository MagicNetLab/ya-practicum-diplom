package inmemory

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryTokenModel_GetValue(t *testing.T) {
	tm := &TokenModel{Value: "sample_value"}
	assert.Equal(t, "sample_value", tm.GetValue())
}

func TestInMemoryTokenModel_SetValue(t *testing.T) {
	t.Run("Проверка успешной установки значения", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetValue("new_value")
		assert.NoError(t, err)
		assert.Equal(t, "new_value", tm.Value)
	})

	t.Run("Проверка попытки установки пустого значения", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetValue("")
		assert.Error(t, err)
	})
}

func TestInMemoryTokenModel_GetUID(t *testing.T) {
	tm := &TokenModel{UID: "user_1"}
	assert.Equal(t, "user_1", tm.GetUID())
}

func TestInMemoryTokenModel_SetUID(t *testing.T) {
	t.Run("Проверка успешной установки UID", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetUID("user_1")
		assert.NoError(t, err)
		assert.Equal(t, "user_1", tm.UID)
	})

	t.Run("Проверка попытки установки пустого UID", func(t *testing.T) {
		tm := &TokenModel{}
		err := tm.SetUID("")
		assert.Error(t, err)
	})
}

func TestInMemoryTokenModel_SetAndIsRefresh(t *testing.T) {
	tm := &TokenModel{}
	// Initially should be false
	assert.False(t, tm.IsRefresh())
	// After setting, flag should be true
	tm.SetRefresh()
	assert.True(t, tm.IsRefresh())
}

func TestInMemoryTokenModel_SetExpired(t *testing.T) {
	t.Run("Проверка успешной установки Expired", func(t *testing.T) {
		tm := &TokenModel{}
		expTime := time.Now().Add(time.Hour)
		err := tm.SetExpired(expTime)
		assert.NoError(t, err)
		assert.Equal(t, expTime.Format(time.RFC3339), tm.Expired)
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

func TestInMemoryTokenModel_GetExpired(t *testing.T) {
	t.Run("Проверка получения Expired при корректном установленном значении", func(t *testing.T) {
		expTime := time.Now().Add(time.Hour)
		tm := &TokenModel{Expired: expTime.Format(time.RFC3339)}
		got := tm.GetExpired()
		// Compare the formatted strings to avoid minor time differences.
		assert.Equal(t, expTime.Format(time.RFC3339), got.Format(time.RFC3339))
	})

	t.Run("Проверка получения Expired при некорректном установленном значении", func(t *testing.T) {
		tm := &TokenModel{Expired: "invalid_time"}
		// GetExpired resets to now+3 hours when parsing fails.
		got := tm.GetExpired()
		expected := time.Now().Add(time.Hour * 3)
		assert.WithinDuration(t, expected, got, time.Minute)
	})
}
