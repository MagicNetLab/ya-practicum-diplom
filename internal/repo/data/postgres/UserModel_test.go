package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresUserModel_GetUID(t *testing.T) {
	u := &UserModel{UID: "test_uid"}
	assert.Equal(t, "test_uid", u.GetUID())
}

func TestPostgresUserModel_GetLogin(t *testing.T) {
	u := &UserModel{Login: "test_login"}
	assert.Equal(t, "test_login", u.GetLogin())
}

func TestPostgresUserModel_GetPassword(t *testing.T) {
	u := &UserModel{Password: "test_password"}
	assert.Equal(t, "test_password", u.GetPassword())
}
