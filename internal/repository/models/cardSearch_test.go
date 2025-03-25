package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCardSearch_GetSubQuery тест метода GetSubQuery модели CardSearch
func TestCardSearch_GetSubQuery(t *testing.T) {
	t.Run("Проверка пустых условия поиска карт", func(t *testing.T) {
		search := CardSearch{}
		str, values := search.GetSubQuery()
		assert.Equal(t, "", str)
		assert.Empty(t, values)
	})

	t.Run("Проверка поиска с пустой строкой в поле Name", func(t *testing.T) {
		search := CardSearch{Name: ""}
		str, values := search.GetSubQuery()
		assert.Equal(t, "", str)
		assert.Empty(t, values)
	})

	t.Run("Проверка поиска со специальными символами в имени", func(t *testing.T) {
		search := CardSearch{Name: "John%_Doe"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE name ilike $1", str)
		assert.Equal(t, []any{"%John%_Doe%"}, values)
	})

	t.Run("Проверка поиска с отрицательным значением Limit", func(t *testing.T) {
		search := CardSearch{Limit: -10}
		str, values := search.GetSubQuery()
		assert.Equal(t, "", str)
		assert.Empty(t, values)
	})

	t.Run("Проверка поиска с отрицательным значением Offset", func(t *testing.T) {
		search := CardSearch{Offset: -20}
		str, values := search.GetSubQuery()
		assert.Equal(t, "", str)
		assert.Empty(t, values)
	})

	t.Run("Проверка поиска с устновленным UID", func(t *testing.T) {
		search := CardSearch{UID: "test-UID"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1", str)
		assert.Equal(t, []any{"test-UID"}, values)
	})

	t.Run("Проверка условия поиска с установленным именем", func(t *testing.T) {
		search := CardSearch{Name: "John Doe"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE name ilike $1", str)
		assert.Equal(t, []any{"%John Doe%"}, values)
	})

	t.Run("Проверка условия поиска со всеми установленными параметрами", func(t *testing.T) {
		search := CardSearch{
			UID:  "test-UID",
			Name: "John Doe",
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND name ilike $2", str)
		assert.Equal(t, []any{"test-UID", "%John Doe%"}, values)
	})

	t.Run("Проверка условия с установленным LIMIT", func(t *testing.T) {
		search := CardSearch{Limit: 10}
		str, values := search.GetSubQuery()
		assert.Equal(t, " LIMIT $1", str)
		assert.Equal(t, []any{10}, values)
	})

	t.Run("Проверка условий поиска с установленным OFFSET", func(t *testing.T) {
		search := CardSearch{Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " OFFSET $1", str)
		assert.Equal(t, []any{20}, values)
	})

	t.Run("Проверка условия поиска с установленным LIMIT и OFFSET ", func(t *testing.T) {
		search := CardSearch{Limit: 10, Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " LIMIT $1 OFFSET $2", str)
		assert.Equal(t, []any{10, 20}, values)
	})

	t.Run("Проверка условия поиска со всеми установленными параметрами, LIMIT и OFFSET", func(t *testing.T) {
		search := CardSearch{
			UID:    "test-UID",
			Name:   "John Doe",
			Limit:  10,
			Offset: 20,
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND name ilike $2 LIMIT $3 OFFSET $4", str)
		assert.Equal(t, []any{"test-UID", "%John Doe%", 10, 20}, values)
	})
}
