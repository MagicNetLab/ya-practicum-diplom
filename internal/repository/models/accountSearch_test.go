package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestAccountSearch_GetUID проверяет правильность получения UID владельца аккаунта
func TestAccountSearch_GetUID(t *testing.T) {
	accountSearch := &AccountSearch{uid: "test-uid"}
	assert.Equal(t, "test-uid", accountSearch.GetUID())
}

// TestAccountSearch_GetLogin проверяет правильность получения логина аккаунта
func TestAccountSearch_GetLogin(t *testing.T) {
	accountSearch := &AccountSearch{login: "test-login"}
	assert.Equal(t, "test-login", accountSearch.GetLogin())
}

// TestAccountSearch_GetURL проверяет правильность получения URL аккаунта
func TestAccountSearch_GetURL(t *testing.T) {
	accountSearch := &AccountSearch{url: "http://test.com"}
	assert.Equal(t, "http://test.com", accountSearch.GetURL())
}

// TestAccountSearch_GetDescription проверяет правильность получения описания аккаунта
func TestAccountSearch_GetDescription(t *testing.T) {
	accountSearch := &AccountSearch{description: "test description"}
	assert.Equal(t, "test description", accountSearch.GetDescription())
}

// TestAccountSearch_GetCreatedFrom проверяет правильность получения минимальной даты создания аккаунта
func TestAccountSearch_GetCreatedFrom(t *testing.T) {
	createdFrom := time.Now()
	accountSearch := &AccountSearch{createdFrom: createdFrom}
	assert.Equal(t, createdFrom, accountSearch.GetCreatedFrom())
}

// TestAccountSearch_GetCreatedTo проверяет правильность получения максимальной даты создания аккаунта
func TestAccountSearch_GetCreatedTo(t *testing.T) {
	createdTo := time.Now()
	accountSearch := &AccountSearch{createdTo: createdTo}
	assert.Equal(t, createdTo, accountSearch.GetCreatedTo())
}

// TestAccountSearch_GetUpdatedFrom проверяет правильность получения минимальной даты обновления аккаунта
func TestAccountSearch_GetUpdatedFrom(t *testing.T) {
	updatedFrom := time.Now()
	accountSearch := &AccountSearch{updatedFrom: updatedFrom}
	assert.Equal(t, updatedFrom, accountSearch.GetUpdatedFrom())
}

// TestAccountSearch_GetUpdatedTo проверяет правильность получения максимальной даты обновления аккаунта
func TestAccountSearch_GetUpdatedTo(t *testing.T) {
	updatedTo := time.Now()
	accountSearch := &AccountSearch{updatedTo: updatedTo}
	assert.Equal(t, updatedTo, accountSearch.GetUpdatedTo())
}

// TestAccountSearch_GetLimit проверяет правильность получения лимита поиска
func TestAccountSearch_GetLimit(t *testing.T) {
	accountSearch := &AccountSearch{limit: 10}
	assert.Equal(t, 10, accountSearch.GetLimit())
}

// TestAccountSearch_GetOffset проверяет правильность получения смещения поиска
func TestAccountSearch_GetOffset(t *testing.T) {
	accountSearch := &AccountSearch{offset: 20}
	assert.Equal(t, 20, accountSearch.GetOffset())
}
