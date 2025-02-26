package inmemory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryUserModel_GetUID(t *testing.T) {
	u := &UserModel{UID: "test_uid"}
	assert.Equal(t, "test_uid", u.GetUID())
}

func TestInMemoryUserModel_GetLogin(t *testing.T) {
	u := &UserModel{Login: "test_login"}
	assert.Equal(t, "test_login", u.GetLogin())
}

func TestInMemoryUserModel_GetPassword(t *testing.T) {
	u := &UserModel{Password: "test_password"}
	assert.Equal(t, "test_password", u.GetPassword())
}

func TestInMemoryUserModel_SetUID(t *testing.T) {
	t.Run("Проверка успешной установки значения", func(t *testing.T) {
		u := &UserModel{}
		err := u.SetUID("test_uid")
		assert.NoError(t, err)
		assert.Equal(t, "test_uid", u.UID)
	})

	t.Run("Проверка попытки установки пустого значения", func(t *testing.T) {
		u := &UserModel{}
		err := u.SetUID("")
		assert.Error(t, err)
	})
}

func TestInMemoryUserModel_SetLogin(t *testing.T) {
	t.Run("Проверка успешной установки логина", func(t *testing.T) {
		u := &UserModel{}
		err := u.SetLogin("test_login")
		assert.NoError(t, err)
		assert.Equal(t, "test_login", u.Login)
	})

	t.Run("Проверка попытки установки пустого логина", func(t *testing.T) {
		u := &UserModel{}
		err := u.SetLogin("")
		assert.Error(t, err)
	})
}

func TestUserModel_SetPassword(t *testing.T) {
	t.Run("Проверка успешной установки пароля", func(t *testing.T) {
		u := &UserModel{}
		err := u.SetPassword("test_password")
		assert.NoError(t, err)
		assert.Equal(t, "test_password", u.Password)
	})

	t.Run("Проверка попытки установки пустого пароля", func(t *testing.T) {
		u := &UserModel{}
		err := u.SetPassword("")
		assert.Error(t, err)
	})
}
