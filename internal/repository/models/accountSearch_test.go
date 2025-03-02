package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAccountSearch_GetUID проверяет правильность получения UID владельца аккаунта
func TestAccountSearch_GetSubQuery(t *testing.T) {
	t.Run("Проверка пустого запроса", func(t *testing.T) {
		search := AccountSearch{}

		str, values := search.GetSubQuery()
		assert.Equal(t, "", str)
		assert.Empty(t, values)
	})

	t.Run("Проверка запроса по UID", func(t *testing.T) {
		search := AccountSearch{UID: "test-UID"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE UID = $1", str)
		assert.Equal(t, []any{"test-UID"}, values)
	})

	t.Run("Проверка запроса с несколькими параметрами", func(t *testing.T) {
		search := AccountSearch{UID: "test-UID", Login: "test-Login", URL: "test-URL"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE UID = $1 AND Login ilike '%$2%' AND URL ilike '%$3%'", str)
		assert.Equal(t, []any{"test-UID", "test-Login", "test-URL"}, values)
	})

	t.Run("Проверка запроса с параметром Limit", func(t *testing.T) {
		search := AccountSearch{Limit: 10}
		str, values := search.GetSubQuery()
		assert.Equal(t, " LIMIT $1", str)
		assert.Equal(t, []any{10}, values)
	})

	t.Run("Проверка запроса с параметром Offset", func(t *testing.T) {
		search := AccountSearch{Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " OFFSET $1", str)
		assert.Equal(t, []any{20}, values)
	})

	t.Run("Проверка запроса с Limit и Offset", func(t *testing.T) {
		search := AccountSearch{Limit: 10, Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " LIMIT $1 OFFSET $2", str)
		assert.Equal(t, []any{10, 20}, values)
	})

	t.Run("Проверка запроса с несколькими параметрам, Limit и Offset", func(t *testing.T) {
		search := AccountSearch{UID: "test-UID", Login: "test-Login", URL: "test-URL", Description: "test-Description", Limit: 10, Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE UID = $1 AND Login ilike '%$2%' AND URL ilike '%$3%' AND Description ilike '%$4%' LIMIT $5 OFFSET $6", str)
		assert.Equal(t, []any{"test-UID", "test-Login", "test-URL", "test-Description", 10, 20}, values)
	})
}
