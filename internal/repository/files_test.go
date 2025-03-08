package repository

import (
	"context"
	"testing"
	"time"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testFileDSN = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable"

func getFileTestDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), testFileDSN)
	require.NoError(t, err)
	return pool
}

func getFileTestRepo(t *testing.T) (FileRepository, *pgxpool.Pool) {
	pool := getFileTestDB(t)
	return NewFileRepository(pool), pool
}

// TestFileRepo_GetFile проверка получения файла по ID
func TestFileRepo_GetFile(t *testing.T) {
	repo, pgx := getFileTestRepo(t)
	defer pgx.Close()

	ctx := context.Background()

	// Тестовые данные
	id := uuid.New().String()
	uid := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM files WHERE id=$1", id)
	})

	sql := "INSERT INTO files (id, uid, name, path, meta, size, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := pgx.Exec(ctx, sql,
		id,
		uid,
		"test.txt",
		"/path/to/test.txt",
		"text/plain",
		1024,
		time.Now(),
	)
	assert.NoError(t, err)

	t.Run("Проверка поиска существующего файла", func(t *testing.T) {
		res, err := repo.GetFile(ctx, id, uid)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.GetID(), id)
		assert.Equal(t, res.GetUID(), uid)
		assert.Equal(t, res.GetName(), "test.txt")
	})

	t.Run("Проверка поиска с некорректным uid", func(t *testing.T) {
		res, err := repo.GetFile(ctx, id, "invalid-uid")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Проверка поиска с некорректным id", func(t *testing.T) {
		res, err := repo.GetFile(ctx, "invalid-id", uid)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Проверка поиска несуществующего файла", func(t *testing.T) {
		res, err := repo.GetFile(ctx, uuid.New().String(), uid)
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

// TestFileRepo_CreateFile проверка создания файла
func TestFileRepo_CreateFile(t *testing.T) {
	repo, pgx := getFileTestRepo(t)
	ctx := context.Background()

	file := &models.File{
		ID:        uuid.New().String(),
		UID:       uuid.New().String(),
		Name:      "test.txt",
		Path:      "/path/to/test.txt",
		Meta:      "text/plain",
		Size:      1024,
		CreatedAt: time.Now(),
	}

	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM files WHERE id=$1", file.GetID())
	})

	t.Run("Проверка успешной записи о файле", func(t *testing.T) {
		err := repo.CreateFile(ctx, file)
		assert.NoError(t, err)

		var count int
		row := pgx.QueryRow(ctx, "SELECT COUNT(*) FROM files WHERE id=$1", file.GetID())
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Проверка создания файла с некорректными данными", func(t *testing.T) {
		invalidFile := &models.File{}
		err := repo.CreateFile(ctx, invalidFile)
		assert.Error(t, err)
	})

	t.Run("Проверка создания дубликата файла", func(t *testing.T) {
		err := repo.CreateFile(ctx, file)
		assert.Error(t, err)
	})
}

// TestFileRepo_DeleteFile проверка удаления файла
func TestFileRepo_DeleteFile(t *testing.T) {
	repo, pgx := getFileTestRepo(t)
	ctx := context.Background()

	id := uuid.New().String()
	uid := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM files WHERE id=$1", id)
	})

	sql := "INSERT INTO files (id, uid, name, path, meta, size, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := pgx.Exec(ctx, sql,
		id,
		uid,
		"test.txt",
		"/path/to/test.txt",
		"text/plain",
		1024,
		time.Now(),
	)
	assert.NoError(t, err)

	t.Run("Проверка удаления существующего файла", func(t *testing.T) {
		var count int
		row := pgx.QueryRow(ctx, "SELECT COUNT(*) FROM files WHERE id=$1", id)
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		err = repo.DeleteFile(ctx, id, uid)
		assert.NoError(t, err)

		row = pgx.QueryRow(ctx, "SELECT COUNT(*) FROM files WHERE id=$1", id)
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Проверка удаления с некорректным uid", func(t *testing.T) {
		err := repo.DeleteFile(ctx, id, "invalid-uid")
		assert.Error(t, err)
	})

	t.Run("Проверка удаления с некорректным id", func(t *testing.T) {
		err := repo.DeleteFile(ctx, "invalid-id", uid)
		assert.Error(t, err)
	})

	t.Run("Проверка удаления несуществующего файла", func(t *testing.T) {
		err := repo.DeleteFile(ctx, uuid.New().String(), uid)
		assert.Error(t, err)
	})
}

// TestFileRepo_SearchFile проверка поиска файлов
func TestFileRepo_SearchFile(t *testing.T) {
	repo, pgx := getFileTestRepo(t)
	uid1 := uuid.New().String()
	uid2 := uuid.New().String()
	ctx := context.Background()

	testData := []models.File{
		{ID: uuid.New().String(), UID: uid1, Name: "test1.txt", Path: "/path/1.txt", Meta: "text/plain", Size: 1024, CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Name: "test2.txt", Path: "/path/2.txt", Meta: "text/plain", Size: 2048, CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Name: "doc1.pdf", Path: "/path/3.pdf", Meta: "application/pdf", Size: 4096, CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Name: "test3.txt", Path: "/path/4.txt", Meta: "text/plain", Size: 512, CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Name: "doc2.pdf", Path: "/path/5.pdf", Meta: "application/pdf", Size: 8192, CreatedAt: time.Now()},
	}

	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM files WHERE uid IN ($1, $2)", uid1, uid2)
	})

	sql := "INSERT INTO files (id, uid, name, path, meta, size, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	for _, file := range testData {
		_, err := pgx.Exec(ctx, sql, file.ID, file.UID, file.Name, file.Path, file.Meta, file.Size, file.CreatedAt)
		assert.NoError(t, err)
	}

	t.Run("Проверка поиска без условий", func(t *testing.T) {
		search := models.FilesSearch{}
		res, err := repo.SearchFile(ctx, &search)
		assert.NoError(t, err)
		assert.Equal(t, len(testData), len(res))
	})

	t.Run("Проверка поиска по uid", func(t *testing.T) {
		search := models.FilesSearch{UID: uid1}
		res, err := repo.SearchFile(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 3)
	})

	t.Run("Проверка поиска по имени", func(t *testing.T) {
		search := models.FilesSearch{Name: "test"}
		res, err := repo.SearchFile(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 3)

		search = models.FilesSearch{Name: "doc"}
		res, err = repo.SearchFile(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Проверка поиска по нескольким условиям", func(t *testing.T) {
		search := models.FilesSearch{UID: uid1, Name: ".pdf"}
		res, err := repo.SearchFile(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})
}
