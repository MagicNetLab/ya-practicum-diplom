package repository

import (
	"context"
	"errors"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/encryptor"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

// NoteRepository - интерфейс репозитория для работы с заметками
type NoteRepository interface {
	GetNote(ctx context.Context, id string, uid string) (models.NoteModel, error)
	CreateNote(ctx context.Context, note models.NoteModel) error
	RemoveNote(ctx context.Context, id string, uid string) error
	SearchNote(ctx context.Context, search models.NoteSearchModel) ([]models.NoteModel, error)
}

// NewNoteRepository - возвращает новый экземпляр репозитория для работы с заметками
func NewNoteRepository(pool *pgxpool.Pool) NoteRepository {
	return &NoteRepo{pool: pool}
}

// NoteRepo - репозиторий для работы с заметками
type NoteRepo struct {
	pool *pgxpool.Pool
}

// GetNote - получение заметки по id
func (n NoteRepo) GetNote(ctx context.Context, id string, uid string) (models.NoteModel, error) {
	if err := uuid.Validate(id); err != nil {
		return nil, errors.New("invalid note id")
	}

	note := models.Note{}
	sql := "SELECT id, uid, title, content, meta, created_at, updated_at FROM notes WHERE id=$1 and uid=$2"
	row := n.pool.QueryRow(ctx, sql, id, uid)
	err := row.Scan(&note.ID, &note.UID, &note.Title, &note.Content, &note.Meta, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	decryptContent, err := encryptor.DecryptData(note.Content)
	if err != nil {
		return nil, err
	}
	note.Content = decryptContent

	decryptMeta, err := encryptor.DecryptData(note.Meta)
	if err != nil {
		return nil, err
	}
	note.Meta = decryptMeta

	return &note, nil
}

// CreateNote - создание новой заметки
func (n NoteRepo) CreateNote(ctx context.Context, note models.NoteModel) error {
	encryptContent, err := encryptor.EncryptData(note.GetContent())
	if err != nil {
		return err
	}
	err = note.SetContent(encryptContent)
	if err != nil {
		return err
	}

	encryptMeta, err := encryptor.EncryptData(note.GetMeta())
	if err != nil {
		return err
	}
	err = note.SetMeta(encryptMeta)
	if err != nil {
		return err
	}

	sql := "INSERT INTO notes (id, uid, title, content, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err = n.pool.Exec(ctx, sql, note.GetID(), note.GetUID(), note.GetTitle(), note.GetContent(), note.GetMeta(), note.GetCreatedAt(), note.GetUpdatedAt())
	if err != nil {
		return err
	}
	return nil
}

// RemoveNote - удаление заметки
func (n NoteRepo) RemoveNote(ctx context.Context, id string, uid string) error {
	if err := uuid.Validate(id); err != nil {
		return errors.New("invalid note id")
	}

	res, err := n.pool.Exec(ctx, "DELETE FROM notes WHERE id=$1 and uid=$2", id, uid)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("note not found")
	}

	return nil
}

// SearchNote - поиск заметок по критериям
func (n NoteRepo) SearchNote(ctx context.Context, search models.NoteSearchModel) ([]models.NoteModel, error) {
	where, params := search.GetSubQuery()
	sql := "SELECT id, uid, title, content, meta, created_at, updated_at FROM notes" + where

	rows, err := n.pool.Query(ctx, sql, params...)
	if err != nil {
		return nil, err
	}

	result := make([]models.NoteModel, 0)
	for rows.Next() {
		note := &models.Note{}
		err = rows.Scan(&note.ID, &note.UID, &note.Title, &note.Content, &note.Meta, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			logger.Error("error scanning note row", logger.StrArg("error", err.Error()))
			continue
		}
		decryptContent, err := encryptor.DecryptData(note.Content)
		if err != nil {
			logger.Error("error decrypting note content", logger.StrArg("error", err.Error()))
			continue
		}
		note.Content = decryptContent

		decryptMeta, err := encryptor.DecryptData(note.Meta)
		if err != nil {
			logger.Error("error decrypting note meta", logger.StrArg("error", err.Error()))
		}
		note.Meta = decryptMeta

		result = append(result, note)
	}

	return result, nil
}
