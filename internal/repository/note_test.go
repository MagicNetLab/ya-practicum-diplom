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

const testNoteDSN = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable"

func getNoteTestDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), testNoteDSN)
	require.NoError(t, err)
	return pool
}

func getNoteTestRepo(t *testing.T) (NoteRepository, *pgxpool.Pool) {
	pool := getNoteTestDB(t)
	return NewNoteRepository(pool), pool
}

// TestNoteRepo_GetNoteByID проверка получения заметки по ID
func TestNoteRepo_GetNoteByID(t *testing.T) {
	repo, pgx := getNoteTestRepo(t)
	ctx := context.Background()
	defer pgx.Close()

	// Тестовые данные
	id := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM notes WHERE id=$1", id)
	})

	sql := "INSERT INTO notes (id, uid, title, content, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := pgx.Exec(ctx, sql,
		id,
		uuid.New().String(),
		"Test Note",
		"Test Content",
		"Test Meta",
		time.Now(),
		time.Now(),
	)
	assert.NoError(t, err)

	t.Run("Проверка поиска существующей заметки", func(t *testing.T) {
		res, err := repo.GetNote(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.GetID(), id)
	})

	t.Run("Проверка поиска несуществующей заметки", func(t *testing.T) {
		res, err := repo.GetNote(ctx, uuid.New().String())
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

// TestNoteRepo_CreateNote проверка создания заметки
func TestNoteRepo_CreateNote(t *testing.T) {
	repo, pgx := getNoteTestRepo(t)
	ctx := context.Background()

	note := models.Note{
		ID:        uuid.New().String(),
		UID:       uuid.New().String(),
		Title:     "Test Note",
		Content:   "Test Content",
		Meta:      "Test Meta",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM notes WHERE id=$1", note.GetID())
	})

	t.Run("Проверка создания новой заметки", func(t *testing.T) {
		err := repo.CreateNote(ctx, &note)
		assert.NoError(t, err)
		var count int
		row := pgx.QueryRow(ctx, "SELECT COUNT(*) FROM notes WHERE id=$1", note.GetID())
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("Проверка создания заметки с некорректными данными", func(t *testing.T) {
		note.UID = ""
		err := repo.CreateNote(ctx, &note)
		assert.Error(t, err)

		note.UID = uuid.New().String()
		note.Title = ""
		err = repo.CreateNote(ctx, &note)
		assert.Error(t, err)

		note.Title = "Test Note"
		note.Content = ""
		err = repo.CreateNote(ctx, &note)
		assert.Error(t, err)
	})
}

// TestNoteRepo_DeleteNote проверка удаления заметки
func TestNoteRepo_DeleteNote(t *testing.T) {
	repo, pgx := getNoteTestRepo(t)
	ctx := context.Background()

	id := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM notes WHERE id=$1", id)
	})

	sql := "INSERT INTO notes (id, uid, title, content, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := pgx.Exec(ctx, sql,
		id,
		uuid.New().String(),
		"Test Note",
		"Test Content",
		"Test Meta",
		time.Now(),
		time.Now(),
	)
	assert.NoError(t, err)

	t.Run("Проверка удаления существующей заметки", func(t *testing.T) {
		var count int
		row := pgx.QueryRow(ctx, "SELECT COUNT(*) FROM notes WHERE id=$1", id)
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 1, count)

		err = repo.RemoveNote(ctx, id)
		assert.NoError(t, err)

		row = pgx.QueryRow(ctx, "SELECT COUNT(*) FROM notes WHERE id=$1", id)
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Проверка удаления несуществующей заметки", func(t *testing.T) {
		err := repo.RemoveNote(ctx, uuid.New().String())
		assert.Error(t, err)
		assert.Equal(t, "note not found", err.Error())
	})
}

// TestNoteRepo_SearchNotes проверка поиска заметок
func TestNoteRepo_SearchNotes(t *testing.T) {
	repo, pgx := getNoteTestRepo(t)
	uid1 := uuid.New().String()
	uid2 := uuid.New().String()
	ctx := context.Background()

	testData := []models.Note{
		{ID: uuid.New().String(), UID: uid1, Title: "Note 1", Content: "Content 1", Meta: "Meta 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Title: "Note 2", Content: "Content 2", Meta: "Meta 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Title: "Note 3", Content: "Content 3", Meta: "Meta 3", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Title: "Other Note 1", Content: "Other Content 1", Meta: "Other Meta 1", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Title: "Other Note 2", Content: "Other Content 2", Meta: "Other Meta 2", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM notes WHERE uid IN ($1, $2)", uid1, uid2)
	})

	sql := "INSERT INTO notes (id, uid, title, content, meta, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	for _, note := range testData {
		_, err := pgx.Exec(ctx, sql, note.ID, note.UID, note.Title, note.Content, note.Meta, note.CreatedAt, note.UpdatedAt)
		assert.NoError(t, err)
	}

	t.Run("Проверка поиска без условий", func(t *testing.T) {
		search := models.NoteSearch{}
		res, err := repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Equal(t, len(testData), len(res))
	})

	t.Run("Проверка поиска c limit", func(t *testing.T) {
		search := models.NoteSearch{Limit: 2}
		res, err := repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(res))
	})

	t.Run("Проверка поиска c offset", func(t *testing.T) {
		search := models.NoteSearch{Offset: 2}
		res, err := repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Equal(t, 3, len(res))
	})

	t.Run("Проверка поиска по uid", func(t *testing.T) {
		res, err := repo.SearchNote(ctx, &models.NoteSearch{UID: uid1})
		assert.NoError(t, err)
		assert.Len(t, res, 3)
	})

	t.Run("Проверка поиска по заголовку", func(t *testing.T) {
		search := models.NoteSearch{Title: "Note 1"}
		res, err := repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)

		search = models.NoteSearch{Title: "Other"}
		res, err = repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Проверка поиска по содержимому", func(t *testing.T) {
		search := models.NoteSearch{Content: "Content 1"}
		res, err := repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)

		search = models.NoteSearch{Content: "Other Content"}
		res, err = repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Проверка поиска по мета-данным", func(t *testing.T) {
		search := models.NoteSearch{Meta: "Meta 1"}
		res, err := repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)

		search = models.NoteSearch{Meta: "Other Meta"}
		res, err = repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Проверка поиска c параметрами, limit и offset", func(t *testing.T) {
		search := models.NoteSearch{Title: "Note", Limit: 3, Offset: 3}
		res, err := repo.SearchNote(ctx, &search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})
}
