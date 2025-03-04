package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNoteSearch_GetSubQuery тестирование функции GetSubQuery для NoteSearch
func TestNoteSearch_GetSubQuery(t *testing.T) {
	t.Run("Проверка пустого поиска", func(t *testing.T) {
		search := NoteSearch{}
		str, values := search.GetSubQuery()
		assert.Equal(t, "", str)
		assert.Empty(t, values)
	})

	t.Run("Проверка поиска по UID", func(t *testing.T) {
		search := NoteSearch{UID: "test-uid"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1", str)
		assert.Equal(t, []any{"test-uid"}, values)
	})

	t.Run("Проверка поиска по Title", func(t *testing.T) {
		search := NoteSearch{Title: "test-title"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE title ilike $1", str)
		assert.Equal(t, []any{"%test-title%"}, values)
	})

	t.Run("Проверка поиска по Content", func(t *testing.T) {
		search := NoteSearch{Content: "test-content"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE content ilike $1", str)
		assert.Equal(t, []any{"%test-content%"}, values)
	})

	t.Run("Проверка поиска по Meta", func(t *testing.T) {
		search := NoteSearch{Meta: "test-meta"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE meta ilike $1", str)
		assert.Equal(t, []any{"%test-meta%"}, values)
	})

	t.Run("Проверка поиска по всем полям", func(t *testing.T) {
		search := NoteSearch{
			UID:     "test-uid",
			Title:   "test-title",
			Content: "test-content",
			Meta:    "test-meta",
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND title ilike $2 AND content ilike $3 AND meta ilike $4", str)
		assert.Equal(t, []any{"test-uid", "%test-title%", "%test-content%", "%test-meta%"}, values)
	})

	t.Run("Проверка поиска с  limit", func(t *testing.T) {
		search := NoteSearch{Limit: 10}
		str, values := search.GetSubQuery()
		assert.Equal(t, " LIMIT $1", str)
		assert.Equal(t, []any{int32(10)}, values)
	})

	t.Run("Проверка поиска с offset", func(t *testing.T) {
		search := NoteSearch{Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " OFFSET $1", str)
		assert.Equal(t, []any{int32(20)}, values)
	})

	t.Run("Проверка поиска с limit и offset", func(t *testing.T) {
		search := NoteSearch{Limit: 10, Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " LIMIT $1 OFFSET $2", str)
		assert.Equal(t, []any{int32(10), int32(20)}, values)
	})

	t.Run("Проверка поиска с limit, offset и другими параметрами", func(t *testing.T) {
		search := NoteSearch{
			UID:     "test-uid",
			Title:   "test-title",
			Content: "test-content",
			Meta:    "test-meta",
			Limit:   10,
			Offset:  20,
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND title ilike $2 AND content ilike $3 AND meta ilike $4 LIMIT $5 OFFSET $6", str)
		assert.Equal(t, []any{"test-uid", "%test-title%", "%test-content%", "%test-meta%", int32(10), int32(20)}, values)
	})
}
