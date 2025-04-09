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

// getCardTestDB возвращает подключение к тестовой базе данных
func getCardTestDB(t *testing.T) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), testCardDSN)
	require.NoError(t, err)
	return pool
}

// getCardTestRepo возвращает тестовый репозиторий карт
func getCardTestRepo(t *testing.T) (CardRepository, *pgxpool.Pool) {
	t.Setenv("ENCRYPT_KEY", "test-encryption-key-32-bytes-length!")
	pool := getCardTestDB(t)
	return NewCardRepo(pool), pool
}

// TestCardRepo_GetCardByID проверка получения карты по ID
func TestCardRepo_GetCardByID(t *testing.T) {
	repo, pgx := getCardTestRepo(t)

	ctx := context.Background()

	// Тестовые данные
	id := uuid.New().String()
	uid := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM cards WHERE id=$1", id)
	})

	account := models.Card{
		ID:        id,
		UID:       uid,
		Name:      "VASYA PUPKIN",
		Number:    "4532015112830366",
		Mask:      "453201******0366",
		Month:     1,
		Year:      2026,
		CVC:       "123",
		PIN:       "1234",
		Meta:      "skjkasjdlkjaskldjalk lsak dlask dklasd",
		CreatedAt: time.Now(),
	}
	err := repo.CreateCard(ctx, &account)
	assert.NoError(t, err)

	t.Run("Проверка поиска существующей карты", func(t *testing.T) {
		res, err := repo.GetCardByID(ctx, id, uid)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, res.GetID(), id)
	})

	t.Run("Проверка поиска несуществующей карты", func(t *testing.T) {
		res, err := repo.GetCardByID(ctx, id, uuid.New().String())
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Проверка ошибки дешифрования данных", func(t *testing.T) {
		// Вставляем некорректно зашифрованные данные
		_, err = pgx.Exec(ctx, "UPDATE cards SET number='invalid-encrypted-data' WHERE id=$1", id)
		assert.NoError(t, err)

		res, err := repo.GetCardByID(ctx, id, uid)
		assert.Error(t, err)
		assert.Nil(t, res)

		// Восстанавливаем корректные данные
		_, err = pgx.Exec(ctx, "UPDATE cards SET number=$1 WHERE id=$2", account.Number, id)
		assert.NoError(t, err)
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
		Meta:      "skjkasjdlkjaskldjalk lsak dlask dklasd",
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

		card.Number = "4532015112830366"
		card.Month = 13
		err = repo.CreateCard(ctx, &card)
		assert.Error(t, err)

		card.Month = 1
		card.Year = 2020
		err = repo.CreateCard(ctx, &card)
		assert.Error(t, err)

		card.Year = 2026
		card.CVC = ""
		err = repo.CreateCard(ctx, &card)
		assert.Error(t, err)
	})

}

// TestCardRepo_DeleteCard проверка удаления карты
func TestCardRepo_DeleteCard(t *testing.T) {
	repo, pgx := getCardTestRepo(t)
	ctx := context.Background()

	id := uuid.New().String()
	uid := uuid.New().String()
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM cards WHERE id=$1", id)
	})

	sql := "INSERT INTO cards (id, uid, name, number, mask, month, year, cvc, pin, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	_, err := pgx.Exec(ctx, sql,
		id,
		uid,
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

		err = repo.DeleteCard(ctx, id, uid)
		assert.NoError(t, err)

		row = pgx.QueryRow(ctx, "SELECT COUNT(*) FROM cards WHERE id=$1", id)
		err = row.Scan(&count)
		assert.NoError(t, err)
		assert.Equal(t, count, 0)
	})

	t.Run("Проверка удаления несуществующей карты", func(t *testing.T) {
		err := repo.DeleteCard(ctx, uuid.New().String(), uid)
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
		{ID: uuid.New().String(), UID: uid1, Name: "VASYA PUPKIN", Number: "4532015112830366", Mask: "453201******0366", Month: 1, Year: 2026, CVC: "123", PIN: "1234", CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Name: "VASYA PUPKIN", Number: "5555555555554444", Mask: "555555******4444", Month: 2, Year: 2025, CVC: "123", PIN: "1234", CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid1, Name: "VASYA PUPKIN", Number: "378282246310005", Mask: "378282******0005", Month: 3, Year: 2027, CVC: "123", PIN: "1234", CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Name: "PETYA PUPKIN", Number: "4111111111111111", Mask: "411111******1111", Month: 4, Year: 2030, CVC: "123", PIN: "1234", CreatedAt: time.Now()},
		{ID: uuid.New().String(), UID: uid2, Name: "PETYA PUPKIN", Number: "5105105105105100", Mask: "510510******5100", Month: 5, Year: 2029, CVC: "123", PIN: "1234", CreatedAt: time.Now()},
	}

	for _, card := range testData {
		err := repo.CreateCard(ctx, &card)
		assert.NoError(t, err)
	}
	t.Cleanup(func() {
		_, _ = pgx.Exec(ctx, "DELETE FROM cards WHERE uid IN ($1, $2)", uid1, uid2)
	})

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

	t.Run("Проверка поиска c параметрами, limit и offset", func(t *testing.T) {
		search := models.CardSearch{Name: "VASYA PUPKIN", Limit: 2, Offset: 2}
		res, err := repo.SearchCards(ctx, search)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
	})

	t.Run("Проверка поиска с некорректно зашифрованными данными", func(t *testing.T) {
		// Вставляем некорректно зашифрованные данные
		_, err := pgx.Exec(ctx, "UPDATE cards SET number='invalid-encrypted-data' WHERE mask = '453201******0366'")
		assert.NoError(t, err)

		res, err := repo.SearchCards(ctx, models.CardSearch{UID: uid1})
		assert.NoError(t, err)
		assert.Len(t, res, 2) // Должно вернуть только карты с корректными данными
	})
}
