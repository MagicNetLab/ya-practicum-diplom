package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

// FilesModel интерфейс модели файла
type FilesModel interface {
	GetID() string
	GetUID() string
	GetName() string
	GetSize() int
	GetPath() string
	GetMeta() string
	GetCreatedAt() time.Time
	IsValid() bool
}

// NewFile создание модели файла
func NewFile(uid string, name string, meta string, size int) (FilesModel, error) {
	if uid == "" || name == "" || size <= 0 {
		return nil, errors.New("invalid params")
	}

	if err := uuid.Validate(uid); err != nil {
		return nil, errors.New("invalid uid")
	}

	path := "/files/" + uid + "/" + name
	return &File{
		ID:        uuid.New().String(),
		UID:       uid,
		Name:      name,
		Meta:      meta,
		Path:      path,
		Size:      size,
		CreatedAt: time.Now(),
	}, nil
}

type File struct {
	ID        string    `db:"id"`
	UID       string    `db:"uid"`
	Name      string    `db:"name"`
	Meta      string    `db:"meta"`
	Path      string    `db:"path"`
	Size      int       `db:"size"`
	CreatedAt time.Time `db:"created_at"`
}

// GetID возвращает ID файла
func (f File) GetID() string {
	return f.ID
}

// GetUID возвращает UID файла
func (f File) GetUID() string {
	return f.UID
}

// GetName возвращает имя файла
func (f File) GetName() string {
	return f.Name
}

// GetMeta возвращает мета данные файла
func (f File) GetMeta() string {
	return f.Meta
}

// GetSize возвращает размер файла
func (f File) GetSize() int {
	return f.Size
}

// GetPath возвращает путь к файлу
func (f File) GetPath() string {
	return f.Path
}

// GetCreatedAt возвращает дату создания файла
func (f File) GetCreatedAt() time.Time {
	return f.CreatedAt
}

// IsValid валидация модели файла
func (f File) IsValid() bool {
	if uuid.Validate(f.ID) != nil {
		return false
	}

	if uuid.Validate(f.UID) != nil {
		return false
	}

	if f.Name == "" || f.Path == "" || f.Size <= 0 {
		return false
	}

	if f.CreatedAt.IsZero() {
		return false
	}

	return true

}
