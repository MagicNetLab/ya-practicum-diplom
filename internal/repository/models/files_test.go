package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestNewFile Тестирование создания модели файла.
func TestNewFile(t *testing.T) {
	t.Run("Корректное создание модели файла", func(t *testing.T) {
		uid := uuid.New().String()
		name := "test_file.txt"
		meta := "test meta"
		size := 1024

		model, err := NewFile(uid, name, meta, size)
		assert.Nil(t, err)
		assert.Equal(t, uid, model.GetUID())
		assert.Equal(t, name, model.GetName())
		assert.Equal(t, meta, model.GetMeta())
		assert.Equal(t, size, model.GetSize())
	})

	t.Run("Попытка создать модель файла с некорректными данными", func(t *testing.T) {
		// Некорректный UID
		model, err := NewFile("invalid_uid", "text.txt", "test meta", 1024)
		assert.Error(t, err)
		assert.Nil(t, model)

		// Отрицательный размер
		model, err = NewFile(uuid.New().String(), "text.txt", "test meta", -1)
		assert.Error(t, err)
		assert.Nil(t, model)

		// Нулевой размер
		model, err = NewFile(uuid.New().String(), "text.txt", "test meta", 0)
		assert.Error(t, err)
		assert.Nil(t, model)

		// Пустое имя файла
		model, err = NewFile(uuid.New().String(), "", "test meta", 1024)
		assert.Error(t, err)
		assert.Nil(t, model)

		// Пустые мета-данные (допустимо)
		model, err = NewFile(uuid.New().String(), "test.txt", "", 1024)
		assert.NoError(t, err)
		assert.NotNil(t, model)
	})

	t.Run("Проверка генерации пути к файлу", func(t *testing.T) {
		uid := uuid.New().String()
		name := "test.txt"
		model, err := NewFile(uid, name, "test meta", 1024)
		assert.NoError(t, err)
		expectedPath := "/files/" + uid + "/" + name
		assert.Equal(t, expectedPath, model.GetPath())
	})
}

// TestFile_SetMethods Тестирование методов установки значений модели файла.
func TestFile_SetMethods(t *testing.T) {
	model := &File{}

	t.Run("Установка ID", func(t *testing.T) {
		// Корректный ID
		id := uuid.New().String()
		err := model.SetID(id)
		assert.NoError(t, err)
		assert.Equal(t, id, model.GetID())

		// Некорректный ID
		err = model.SetID("")
		assert.NoError(t, err)
	})

	t.Run("Установка UID", func(t *testing.T) {
		// Корректный UID
		uid := uuid.New().String()
		err := model.SetUID(uid)
		assert.NoError(t, err)
		assert.Equal(t, uid, model.GetUID())

		// Некорректный UID
		err = model.SetUID("")
		assert.NoError(t, err)
	})

	t.Run("Установка имени файла", func(t *testing.T) {
		// Корректное имя
		name := "test.txt"
		err := model.SetName(name)
		assert.NoError(t, err)
		assert.Equal(t, name, model.GetName())

		// Пустое имя
		err = model.SetName("")
		assert.NoError(t, err)
	})

	t.Run("Установка мета-информации", func(t *testing.T) {
		// Корректные мета-данные
		meta := "test meta"
		err := model.SetMeta(meta)
		assert.NoError(t, err)
		assert.Equal(t, meta, model.GetMeta())

		// Пустые мета-данные
		err = model.SetMeta("")
		assert.NoError(t, err)
	})

	t.Run("Установка пути к файлу", func(t *testing.T) {
		// Корректный путь
		path := "/path/to/file"
		err := model.SetPath(path)
		assert.NoError(t, err)
		assert.Equal(t, path, model.GetPath())

		// Пустой путь
		err = model.SetPath("")
		assert.NoError(t, err)
	})

	t.Run("Установка размера файла", func(t *testing.T) {
		// Корректный размер
		size := 1024
		err := model.SetSize(size)
		assert.NoError(t, err)
		assert.Equal(t, size, model.GetSize())

		// Отрицательный размер
		err = model.SetSize(-1)
		assert.NoError(t, err)

		// Нулевой размер
		err = model.SetSize(0)
		assert.NoError(t, err)
	})

	t.Run("Установка времени создания", func(t *testing.T) {
		// Текущее время
		createdAt := time.Now()
		err := model.SetCreatedAt(createdAt)
		assert.NoError(t, err)
		assert.Equal(t, createdAt, model.GetCreatedAt())

		// Время в прошлом
		pastTime := time.Now().Add(-24 * time.Hour)
		err = model.SetCreatedAt(pastTime)
		assert.NoError(t, err)

		// Нулевое время
		err = model.SetCreatedAt(time.Time{})
		assert.NoError(t, err)
	})
}

// TestFile_AllMethods Тестирование всех методов модели файла.
func TestFile_AllMethods(t *testing.T) {
	id := uuid.New().String()
	uid := uuid.New().String()
	name := "test_file.txt"
	meta := "test meta"
	path := "path/to/file"
	size := 1024
	createdAt := time.Now().Add(-time.Hour * 24)

	model := File{
		ID:        id,
		UID:       uid,
		Name:      name,
		Meta:      meta,
		Path:      path,
		Size:      size,
		CreatedAt: createdAt,
	}

	t.Run("Проверка получения id файла", func(t *testing.T) {
		assert.Equal(t, id, model.GetID())
	})

	t.Run("Проверка получения uid файла", func(t *testing.T) {
		assert.Equal(t, uid, model.GetUID())
	})

	t.Run("Проверка получения имени файла", func(t *testing.T) {
		assert.Equal(t, name, model.GetName())
	})

	t.Run("Проверка получения метаинформации файла", func(t *testing.T) {
		assert.Equal(t, meta, model.GetMeta())
	})

	t.Run("Проверка получения пути к файлу", func(t *testing.T) {
		assert.Equal(t, path, model.GetPath())
	})

	t.Run("Проверка получения размера файла", func(t *testing.T) {
		assert.Equal(t, size, model.GetSize())
	})

	t.Run("Проверка получения времени создания файла", func(t *testing.T) {
		assert.Equal(t, createdAt, model.GetCreatedAt())
	})

}

// TestFile_IsValid Тестирование метода проверки валидности модели файла.
func TestFile_IsValid(t *testing.T) {
	tests := []struct {
		name string
		file File
		want bool
	}{
		{
			name: "Корректная модель файла",
			file: File{
				ID:        uuid.New().String(),
				UID:       uuid.New().String(),
				Name:      "test.txt",
				Path:      "/path/to/file",
				Size:      100,
				CreatedAt: time.Now(),
			},
			want: true,
		},
		{
			name: "Некорректный ID",
			file: File{
				ID:        "invalid-uuid",
				UID:       uuid.New().String(),
				Name:      "test.txt",
				Path:      "/path/to/file",
				Size:      100,
				CreatedAt: time.Now(),
			},
			want: false,
		},
		{
			name: "Некорректный UID",
			file: File{
				ID:        uuid.New().String(),
				UID:       "invalid-uuid",
				Name:      "test.txt",
				Path:      "/path/to/file",
				Size:      100,
				CreatedAt: time.Now(),
			},
			want: false,
		},
		{
			name: "Пустое имя файла",
			file: File{
				ID:        uuid.New().String(),
				UID:       uuid.New().String(),
				Name:      "",
				Path:      "/path/to/file",
				Size:      100,
				CreatedAt: time.Now(),
			},
			want: false,
		},
		{
			name: "Пустой путь к файлу",
			file: File{
				ID:        uuid.New().String(),
				UID:       uuid.New().String(),
				Name:      "test.txt",
				Path:      "",
				Size:      100,
				CreatedAt: time.Now(),
			},
			want: false,
		},
		{
			name: "Некорректный размер файла",
			file: File{
				ID:        uuid.New().String(),
				UID:       uuid.New().String(),
				Name:      "test.txt",
				Path:      "/path/to/file",
				Size:      0,
				CreatedAt: time.Now(),
			},
			want: false,
		},
		{
			name: "Нулевое время создания",
			file: File{
				ID:        uuid.New().String(),
				UID:       uuid.New().String(),
				Name:      "test.txt",
				Path:      "/path/to/file",
				Size:      100,
				CreatedAt: time.Time{},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.file.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}
