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

	t.Run("Проверка поиска с устновленным UID", func(t *testing.T) {
		search := CardSearch{UID: "test-UID"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1", str)
		assert.Equal(t, []any{"test-UID"}, values)
	})

	t.Run("Проверка поиска с установленным номером карты", func(t *testing.T) {
		search := CardSearch{Number: "1234"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE number ilike $1", str)
		assert.Equal(t, []any{"%1234%"}, values)
	})

	t.Run("Проверка условия поиска с установленным имененем", func(t *testing.T) {
		search := CardSearch{Name: "John Doe"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE name ilike $1", str)
		assert.Equal(t, []any{"%John Doe%"}, values)
	})

	t.Run("Проверка условия поиска с установленным годом", func(t *testing.T) {
		search := CardSearch{Year: 2025}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE year = $1", str)
		assert.Equal(t, []any{"2025"}, values)
	})

	t.Run("Проверка условия поиска со всеми установленными параметрами", func(t *testing.T) {
		search := CardSearch{
			UID:    "test-UID",
			Number: "1234",
			Name:   "John Doe",
			Year:   2025,
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND number ilike $2 AND name ilike $3 AND year = $4", str)
		assert.Equal(t, []any{"test-UID", "%1234%", "%John Doe%", "2025"}, values)
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
			Number: "1234",
			Name:   "John Doe",
			Year:   2025,
			Limit:  10,
			Offset: 20,
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND number ilike $2 AND name ilike $3 AND year = $4 LIMIT $5 OFFSET $6", str)
		assert.Equal(t, []any{"test-UID", "%1234%", "%John Doe%", "2025", 10, 20}, values)
	})
}
