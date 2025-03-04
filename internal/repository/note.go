package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

// NoteRepository - интерфейс репозитория для работы с заметками
type NoteRepository interface {
	GetNote(ctx context.Context, id string) (models.NoteModel, error)
	CreateNote(ctx context.Context, note models.NoteModel) error
	UpdateNote(ctx context.Context, note models.NoteModel) error
	RemoveNote(ctx context.Context, id string) error
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
func (n NoteRepo) GetNote(ctx context.Context, id string) (models.NoteModel, error) {
	if err := uuid.Validate(id); err != nil {
		return nil, errors.New("invalid note id")
	}

	note := models.Note{}
	sql := "SELECT id, uid, title, content, meta, created_at, updated_at FROM notes WHERE id=$1"
	row := n.pool.QueryRow(ctx, sql, id)
	err := row.Scan(&note.ID, &note.UID, &note.Title, &note.Content, &note.Meta, &note.CreatedAt, &note.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &note, nil
}

// CreateNote - создание новой заметки
func (n NoteRepo) CreateNote(ctx context.Context, note models.NoteModel) error {
	sql := "INSERT INTO notes (id, uid, title, content, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := n.pool.Exec(ctx, sql, note.GetID(), note.GetUID(), note.GetTitle(), note.GetContent(), note.GetMeta(), note.GetCreatedAt(), note.GetUpdatedAt())
	if err != nil {
		return err
	}
	return nil
}

// UpdateNote - обновление заметки
func (n NoteRepo) UpdateNote(ctx context.Context, note models.NoteModel) error {
	sql := "UPDATE notes SET title=$1, content=$2, meta=$3, updated_at=$4 WHERE id=$5"
	_, err := n.pool.Exec(ctx, sql, note.GetTitle(), note.GetContent(), note.GetMeta(), note.GetUpdatedAt(), note.GetID())
	if err != nil {
		return err
	}
	return nil
}

// RemoveNote - удаление заметки
func (n NoteRepo) RemoveNote(ctx context.Context, id string) error {
	if err := uuid.Validate(id); err != nil {
		return errors.New("invalid note id")
	}

	res, err := n.pool.Exec(ctx, "DELETE FROM notes WHERE id=$1", id)
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

	var result []models.NoteModel
	for rows.Next() {
		note := &models.Note{}
		err = rows.Scan(&note.ID, &note.UID, &note.Title, &note.Content, &note.Meta, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			logger.Error("error scanning note row", logger.StrArg("error", err.Error()))
			continue
		}
		result = append(result, note)
	}

	return result, nil
}
