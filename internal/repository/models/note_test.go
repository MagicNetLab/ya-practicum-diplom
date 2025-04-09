package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestNewNote тест создания нового объекта Note
func TestNewNote(t *testing.T) {
	tests := []struct {
		name        string
		uid         string
		title       string
		content     string
		meta        string
		expectError bool
	}{
		{
			name:        "Проверка успешного создания нового объекта Note",
			uid:         uuid.New().String(),
			title:       "Test Title",
			content:     "Test Content",
			meta:        "Test Meta",
			expectError: false,
		},
		{
			name:        "Проверка попытки создания нового объекта Note с пустым uid",
			uid:         "",
			title:       "Test Title",
			content:     "Test Content",
			meta:        "Test Meta",
			expectError: true,
		},
		{
			name:        "Проверка попытки создания нового объекта Note с пустым title",
			uid:         uuid.New().String(),
			title:       "",
			content:     "Test Content",
			meta:        "Test Meta",
			expectError: true,
		},
		{
			name:        "Проверка попытки создания нового объекта Note с пустым content",
			uid:         uuid.New().String(),
			title:       "Test Title",
			content:     "",
			meta:        "Test Meta",
			expectError: true,
		},
		{
			name:        "Проверка попытки создания нового объекта Note с пустым meta",
			uid:         uuid.New().String(),
			title:       "Test Title",
			content:     "Test Content",
			meta:        "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			note, err := NewNote(tt.uid, tt.title, tt.content, tt.meta)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, note)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, note.ID)
				assert.Equal(t, tt.uid, note.UID)
				assert.Equal(t, tt.title, note.Title)
				assert.Equal(t, tt.content, note.Content)
				assert.Equal(t, tt.meta, note.Meta)
				assert.False(t, note.CreatedAt.IsZero())
				assert.False(t, note.UpdatedAt.IsZero())
			}
		})
	}
}

// TestNote_GetID тест метода GetID
func TestNote_GetID(t *testing.T) {
	id := uuid.New().String()
	note := Note{ID: id}
	assert.Equal(t, id, note.GetID())
}

// TestNote_GetUID тест метода GetUID
func TestNote_GetUID(t *testing.T) {
	uid := uuid.New().String()
	note := Note{UID: uid}
	assert.Equal(t, uid, note.GetUID())
}

// TestNote_GetTitle тест метода GetTitle
func TestNote_GetTitle(t *testing.T) {
	title := "Test Title"
	note := Note{Title: title}
	assert.Equal(t, title, note.GetTitle())
}

// TestNote_GetContent тест метода GetContent
func TestNote_GetContent(t *testing.T) {
	content := "Test Content"
	note := Note{Content: content}
	assert.Equal(t, content, note.GetContent())
}

// TestNote_GetMeta тест метода GetMeta
func TestNote_GetMeta(t *testing.T) {
	meta := "Test Meta"
	note := Note{Meta: meta}
	assert.Equal(t, meta, note.GetMeta())
}

// TestNote_GetCreatedAt тест метода GetCreatedAt
func TestNote_GetCreatedAt(t *testing.T) {
	createdAt := time.Now()
	note := Note{CreatedAt: createdAt}
	assert.Equal(t, createdAt, note.GetCreatedAt())
}

// TestNote_GetUpdatedAt тест метода GetUpdatedAt
func TestNote_GetUpdatedAt(t *testing.T) {
	updatedAt := time.Now()
	note := Note{UpdatedAt: updatedAt}
	assert.Equal(t, updatedAt, note.GetUpdatedAt())
}

// TestNote_SetMethods тест методов установки значений
func TestNote_SetMethods(t *testing.T) {
	t.Run("Проверка метода SetTitle", func(t *testing.T) {
		note := &Note{}

		// Корректный заголовок
		err := note.SetTitle("Test Title")
		assert.NoError(t, err)
		assert.Equal(t, "Test Title", note.GetTitle())

		// Пустой заголовок
		err = note.SetTitle("")
		assert.Error(t, err)

		// Заголовок со спецсимволами
		err = note.SetTitle("Test\nTitle\tWith\rSpecial")
		assert.NoError(t, err)
		assert.Equal(t, "Test\nTitle\tWith\rSpecial", note.GetTitle())
	})

	t.Run("Проверка метода SetContent", func(t *testing.T) {
		note := &Note{}

		// Корректное содержимое
		err := note.SetContent("Test Content")
		assert.NoError(t, err)
		assert.Equal(t, "Test Content", note.GetContent())

		// Пустое содержимое
		err = note.SetContent("")
		assert.Error(t, err)

		// Содержимое со спецсимволами
		err = note.SetContent("Content\nWith\tSpecial\rChars")
		assert.NoError(t, err)
		assert.Equal(t, "Content\nWith\tSpecial\rChars", note.GetContent())
	})

	t.Run("Проверка метода SetMeta", func(t *testing.T) {
		note := &Note{}

		// Корректные метаданные
		err := note.SetMeta("Test Meta")
		assert.NoError(t, err)
		assert.Equal(t, "Test Meta", note.GetMeta())

		// Пустые метаданные
		err = note.SetMeta("")
		assert.Error(t, err)

		// Метаданные со спецсимволами
		err = note.SetMeta("Meta\nWith\tSpecial\rChars")
		assert.NoError(t, err)
		assert.Equal(t, "Meta\nWith\tSpecial\rChars", note.GetMeta())
	})
}

// TestNote_TimeFields тест полей времени
func TestNote_TimeFields(t *testing.T) {
	t.Run("Проверка времени создания и обновления", func(t *testing.T) {
		// Создание новой заметки
		note, err := NewNote(uuid.New().String(), "Test", "Content", "Meta")
		assert.NoError(t, err)

		// Проверка что время создания установлено
		assert.False(t, note.GetCreatedAt().IsZero())
		assert.False(t, note.GetUpdatedAt().IsZero())

		// Проверка что время обновления не раньше времени создания
		assert.True(t, note.GetUpdatedAt().After(note.GetCreatedAt()) ||
			note.GetUpdatedAt().Equal(note.GetCreatedAt()))
	})
}
