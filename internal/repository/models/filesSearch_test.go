package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFilesSearch_GetSubQuery тест метода GetSubQuery модели FilesSearch
func TestFilesSearch_GetSubQuery(t *testing.T) {
	t.Run("Проверка пустых условий поиска файлов", func(t *testing.T) {
		search := FilesSearch{}
		str, values := search.GetSubQuery()
		assert.Equal(t, "", str)
		assert.Empty(t, values)
	})

	t.Run("Проверка поиска с установленным UID", func(t *testing.T) {
		search := FilesSearch{UID: "test-UID"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1", str)
		assert.Equal(t, []any{"test-UID"}, values)
	})

	t.Run("Проверка поиска с установленным именем файла", func(t *testing.T) {
		search := FilesSearch{Name: "document.pdf"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE name ilike $1", str)
		assert.Equal(t, []any{"%document.pdf%"}, values)
	})

	t.Run("Проверка поиска с установленными метаданными", func(t *testing.T) {
		search := FilesSearch{Meta: "image/jpeg"}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE meta ilike $1", str)
		assert.Equal(t, []any{"%image/jpeg%"}, values)
	})

	t.Run("Проверка поиска с установленными UID и именем файла", func(t *testing.T) {
		search := FilesSearch{
			UID:  "test-UID",
			Name: "document.pdf",
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND name ilike $2", str)
		assert.Equal(t, []any{"test-UID", "%document.pdf%"}, values)
	})

	t.Run("Проверка поиска с установленными UID и метаданными", func(t *testing.T) {
		search := FilesSearch{
			UID:  "test-UID",
			Meta: "image/jpeg",
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND meta ilike $2", str)
		assert.Equal(t, []any{"test-UID", "%image/jpeg%"}, values)
	})

	t.Run("Проверка поиска с установленными именем файла и метаданными", func(t *testing.T) {
		search := FilesSearch{
			Name: "document.pdf",
			Meta: "image/jpeg",
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE name ilike $1 AND meta ilike $2", str)
		assert.Equal(t, []any{"%document.pdf%", "%image/jpeg%"}, values)
	})

	t.Run("Проверка поиска со всеми установленными параметрами", func(t *testing.T) {
		search := FilesSearch{
			UID:  "test-UID",
			Name: "document.pdf",
			Meta: "image/jpeg",
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND name ilike $2 AND meta ilike $3", str)
		assert.Equal(t, []any{"test-UID", "%document.pdf%", "%image/jpeg%"}, values)
	})

	t.Run("Проверка условия с установленным LIMIT", func(t *testing.T) {
		search := FilesSearch{Limit: 10}
		str, values := search.GetSubQuery()
		assert.Equal(t, " LIMIT $1", str)
		assert.Equal(t, []any{int32(10)}, values)
	})

	t.Run("Проверка условий поиска с установленным OFFSET", func(t *testing.T) {
		search := FilesSearch{Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " OFFSET $1", str)
		assert.Equal(t, []any{int32(20)}, values)
	})

	t.Run("Проверка условия поиска с установленным LIMIT и OFFSET", func(t *testing.T) {
		search := FilesSearch{Limit: 10, Offset: 20}
		str, values := search.GetSubQuery()
		assert.Equal(t, " LIMIT $1 OFFSET $2", str)
		assert.Equal(t, []any{int32(10), int32(20)}, values)
	})

	t.Run("Проверка условия поиска со всеми установленными параметрами, LIMIT и OFFSET", func(t *testing.T) {
		search := FilesSearch{
			UID:    "test-UID",
			Name:   "document.pdf",
			Meta:   "image/jpeg",
			Limit:  10,
			Offset: 20,
		}
		str, values := search.GetSubQuery()
		assert.Equal(t, " WHERE uid = $1 AND name ilike $2 AND meta ilike $3 LIMIT $4 OFFSET $5", str)
		assert.Equal(t, []any{"test-UID", "%document.pdf%", "%image/jpeg%", int32(10), int32(20)}, values)
	})
}
