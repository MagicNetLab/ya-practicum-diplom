package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// NoteModel интерфейс модели заметки
type NoteModel interface {
	GetID() string
	GetUID() string
	GetTitle() string
	GetContent() string
	GetMeta() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

// NewNote возвращает новую заметку
func NewNote(uid, title, content, meta string) (Note, error) {
	if uid == "" || title == "" || content == "" {
		return Note{}, errors.New("invalid data")
	}

	return Note{
		ID:        uuid.New().String(),
		UID:       uid,
		Title:     title,
		Content:   content,
		Meta:      meta,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Note модель заметки
type Note struct {
	ID        string    `db:"id"`
	UID       string    `db:"uid"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Meta      string    `db:"meta"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// GetID возвращает ID заметки
func (n *Note) GetID() string {
	return n.ID
}

// GetUID возвращает UID заметки
func (n *Note) GetUID() string {
	return n.UID
}

// GetTitle возвращает заголовок заметки
func (n *Note) GetTitle() string {
	return n.Title
}

// GetContent возвращает содержимое заметки
func (n *Note) GetContent() string {
	return n.Content
}

// GetMeta возвращает метаданные заметки
func (n *Note) GetMeta() string {
	return n.Meta
}

// GetCreatedAt возвращает дату создания заметки
func (n *Note) GetCreatedAt() time.Time {
	return n.CreatedAt
}

// GetUpdatedAt возвращает дату обновления заметки
func (n *Note) GetUpdatedAt() time.Time {
	return n.UpdatedAt
}
