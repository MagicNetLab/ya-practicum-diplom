package repository

import (
	"context"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// FileRepository - интерфейс репозитория для работы с файлами.
type FileRepository interface {
	GetFile(ctx context.Context, fileID string, uid string) (models.FilesModel, error)
	CreateFile(ctx context.Context, file models.FilesModel) error
	DeleteFile(ctx context.Context, fileID string, uid string) error
	SearchFile(ctx context.Context, search models.FilesSearchModel) ([]models.FilesModel, error)
}

func NewFileRepository(pool *pgxpool.Pool) FileRepository {
	return &FileRepo{pool: pool}
}

// FileRepo - репозиторий для работы с файлами.
type FileRepo struct {
	pool *pgxpool.Pool
}

// GetFile - получить данные файла из БД
func (f FileRepo) GetFile(ctx context.Context, fileID string, uid string) (models.FilesModel, error) {
	if err := uuid.Validate(uid); err != nil {
		return nil, fmt.Errorf("uid validation failed: %w", err)
	}

	if err := uuid.Validate(fileID); err != nil {
		return nil, fmt.Errorf("id validation failed: %w", err)
	}

	model := &models.File{}
	sql := `SELECT id, uid, name, path, meta, size, created_at FROM files WHERE id=$1 AND uid=$2`
	row := f.pool.QueryRow(ctx, sql, fileID, uid)
	err := row.Scan(&model.ID, &model.UID, &model.Name, &model.Path, &model.Meta, &model.Size, &model.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	return model, nil
}

// CreateFile - запись данных о файле в БД
func (f FileRepo) CreateFile(ctx context.Context, file models.FilesModel) error {
	if !file.IsValid() {
		return fmt.Errorf("invalid model params")
	}

	sql := `INSERT INTO files (id, uid, name, path, meta, size, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	res, err := f.pool.Exec(ctx, sql, file.GetID(), file.GetUID(), file.GetName(), file.GetPath(), file.GetMeta(), file.GetSize(), file.GetCreatedAt())
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	if res.RowsAffected() != 1 {
		return fmt.Errorf("failed to create file")
	}

	return nil
}

// DeleteFile - удаление данных о файле из БД
func (f FileRepo) DeleteFile(ctx context.Context, fileID string, uid string) error {
	if err := uuid.Validate(uid); err != nil {
		return fmt.Errorf("uid validation failed: %w", err)
	}

	if err := uuid.Validate(fileID); err != nil {
		return fmt.Errorf("id validation failed: %w", err)
	}

	sql := `DELETE FROM files WHERE id=$1 AND uid=$2`
	res, err := f.pool.Exec(ctx, sql, fileID, uid)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	if res.RowsAffected() != 1 {
		return fmt.Errorf("failed to delete file")
	}

	return nil
}

// SearchFile - поиск данных о файлах в БД
func (f FileRepo) SearchFile(ctx context.Context, search models.FilesSearchModel) ([]models.FilesModel, error) {
	where, args := search.GetSubQuery()
	sql := `SELECT id, uid, name, path, meta, size, created_at FROM files` + where
	rows, err := f.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search file: %w", err)
	}

	var result []models.FilesModel
	for rows.Next() {
		model := &models.File{}
		err = rows.Scan(&model.ID, &model.UID, &model.Name, &model.Path, &model.Meta, &model.Size, &model.CreatedAt)
		if err != nil {
			logger.Error("failed to search file", logger.StrArg("error", err.Error()), logger.StrArg("uid", search.GetUID()))
			continue
		}
		result = append(result, model)
	}

	return result, nil
}
