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

const testCardDSN = "postgres://gophkeeper:gophkeeper@localhost:5432/gophkeeper?sslmode=disable"

func getCardTestDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), testAccountDSN)
	require.NoError(t, err)
	return pool
}

func getCardTestRepo(t *testing.T) (CardRepository, *pgxpool.Pool) {
	pool := getCardTestDB(t)
	return NewCardRepo(pool), pool
}

// TestCardRepo_GetCardByID проверка получения карты по ID
func TestCardRepo_GetCardByID(t *testing.T) {
	repo, pgx := getCardTestRepo(t)
	defer pgx.Close()

	ctx := context.Background()

	// Тестовые данные
	id := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM cards WHERE id=$1", id)
	})

	sql := "INSERT INTO cards (id, uid, name, number, mask, month, year, cvc, pin, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	_, err := pgx.Exec(ctx, sql,
		id,
		uuid.New().String(),
		"VASYA PUPKIN",
		"4532015112830366",
		"453201******0366",
		1,
		2026,
		"123",
		"1234",
		time.Now(),
	)
	assert.NoError(t, err)

	t.Run("Проверка поиска существующей карты", func(t *testing.T) {
		res, err := repo.GetCardByID(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.GetID(), id)
	})

	t.Run("Проверка поиска несуществующей карты", func(t *testing.T) {
		res, err := repo.GetCardByID(ctx, uuid.New().String())
		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

// TestCardRepo_CreatedCard проверка создания карты
func TestCardRepo_CreatedCard(t *testing.T) {
	repo, pgx := getCardTestRepo(t)
	ctx := context.Background()

	card := models.Card{
		ID:        uuid.New().String(),
		UID:       uuid.New().String(),
		Name:      "VASYA PUPKIN",
		Number:    "4532015112830366",
		Mask:      "453201******0366",
		Month:     1,
		Year:      2026,
		CVC:       "123",
		PIN:       "1234",
		CreatedAt: time.Now(),
	}
	err := card.Validate()
	assert.NoError(t, err)

	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM cards WHERE id=$1", card.GetID())
	})

	t.Run("Проверка создания новой карты", func(t *testing.T) {
		err := repo.CreateCard(ctx, &card)
		assert.NoError(t, err)
		var count int
		row := pgx.QueryRow(ctx, "SELECT COUNT(*) FROM cards WHERE id=$1", card.GetID())
		err = row.Scan(&count)
		assert.NoError(t, err)
	})

	t.Run("Проверка создания карты с не уникальным номером", func(t *testing.T) {
		err := repo.CreateCard(ctx, &card)
		assert.Error(t, err)
	})

	t.Run("Проверка создания карты с некорректными данными", func(t *testing.T) {
		card.UID = ""
		err := repo.CreateCard(ctx, &card)
		assert.Error(t, err)

		card.UID = uuid.New().String()
		card.Name = ""
		err = repo.CreateCard(ctx, &card)
		assert.Error(t, err)

		card.Name = "VASYA PUPKIN"
		card.Number = ""
		err = repo.CreateCard(ctx, &card)
		assert.Error(t, err)
	})

}

// TestCardRepo_DeleteCard проверка удаления карты
func TestCardRepo_DeleteCard(t *testing.T) {
	repo, pgx := getCardTestRepo(t)
	ctx := context.Background()

	id := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM cards WHERE id=$1", id)
	})

	sql := "INSERT INTO cards (id, uid, name, number, mask, month, year, cvc, pin, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	_, err := pgx.Exec(ctx, sql,
		id,
		uuid.New().String(),
		"VASYA PUPKIN",
		"4532015112830366",
		"453201******0366",
		1,
		2026,
		"123",
		"1234",
		time.Now(),
	)
	assert.NoError(t, err)

	t.Run("Проверка удаления существующей карты", func(t *testing.T) {
		var count int
		row := pgx.QueryRow(ctx, "SELECT COUNT(*) FROM cards WHERE id=$1", id)
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, count, 1)

		err = repo.DeleteCard(ctx, id)
		assert.NoError(t, err)

		row = pgx.QueryRow(ctx, "SELECT COUNT(*) FROM cards WHERE id=$1", id)
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, count, 0)
	})

	t.Run("Проверка удаления несуществующей карты", func(t *testing.T) {
		err := repo.DeleteCard(ctx, uuid.New().String())
		assert.Error(t, err)
		assert.Equal(t, "card not found", err.Error())
	})
}

// TestCardRepo_SearchCard проверка поиска карт
func TestCardRepo_SearchCard(t *testing.T) {
	repo, pgx := getCardTestRepo(t)
	uid1 := uuid.New().String()
	uid2 := uuid.New().String()
	ctx := context.Background()

	testData := []models.Card{
		{ID: uuid.New().String(), UID: uid1, Name: "VASYA PUPKIN", Number: "4532015112830366", Mask: "453201******0366", Month: 1, Year: 2026, CVC: "123", PIN: "1", CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Name: "VASYA PUPKIN", Number: "45329999x2830377", Mask: "453299******0377", Month: 2, Year: 2025, CVC: "123", PIN: "1", CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Name: "VASYA PUPKIN", Number: "4532888112830388", Mask: "453288******0388", Month: 3, Year: 2027, CVC: "123", PIN: "1", CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Name: "PETYA PUPKIN", Number: "4532777112830399", Mask: "453277******0399", Month: 4, Year: 2030, CVC: "123", PIN: "1", CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Name: "PETYA PUPKIN", Number: "45320005112830300", Mask: "453200******0300", Month: 5, Year: 2029, CVC: "123", PIN: "1", CreatedAt: time.Now()},
	}

	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM cards WHERE uid IN ($1, $2)", uid1, uid2)
	})

	sql := "INSERT INTO cards (id, uid, name, number, mask, month, year, cvc, pin, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	for _, card := range testData {
		_, err := pgx.Exec(ctx, sql, card.ID, card.UID, card.Name, card.Number, card.Mask, card.Month, card.Year, card.CVC, card.PIN, card.CreatedAt)
		assert.NoError(t, err)
	}

	t.Run("Проверка поиска без условий", func(t *testing.T) {
		search := models.CardSearch{}
		res, err := repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Equal(t, len(res), len(testData))
	})

	t.Run("Проверка поиска c limit", func(t *testing.T) {
		search := models.CardSearch{Limit: 2}
		res, err := repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Equal(t, len(res), 2)
	})

	t.Run("Проверка поиска c offset", func(t *testing.T) {
		search := models.CardSearch{Offset: 2}
		res, err := repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Equal(t, len(res), 3)
	})

	t.Run("Проверка поиска по uid", func(t *testing.T) {
		res, err := repo.SearchCards(ctx, models.CardSearch{UID: uid1})
		assert.NoError(t, err)
		assert.Len(t, res, 3)
	})

	t.Run("Проверка поиска по номеру", func(t *testing.T) {
		search := models.CardSearch{Number: "4532015112830366"}
		res, err := repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, res, 1)

		search = models.CardSearch{Number: "45320"}
		res, err = repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("Проверка поиска по имени", func(t *testing.T) {
		search := models.CardSearch{Name: "VASYA PUPKIN"}
		res, err := repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, res, 3)

		search = models.CardSearch{Name: "PUPKIN"}
		res, err = repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, res, 5)
	})

	t.Run("Проверка поиска по году", func(t *testing.T) {
		search := models.CardSearch{Year: 2025}
		res, err := repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("Проверка поиска c параметрами, limit и offset", func(t *testing.T) {
		search := models.CardSearch{Name: "VASYA PUPKIN", Limit: 2, Offset: 2}
		res, err := repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

}
