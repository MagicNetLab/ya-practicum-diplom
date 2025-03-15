package repository

import (
	"context"
	"errors"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewCardRepo создает новый репозиторий карт
func NewCardRepo(pool *pgxpool.Pool) CardRepository {
	return &CardRepo{pool: pool}
}

// CardRepository интерфейс репозитория карт
type CardRepository interface {
	// GetCardByID возвращает карту по ID
	GetCardByID(ctx context.Context, id string, uid string) (models.CardModel, error)
	// GetCardByNumber возвращает карту по номеру
	CreateCard(ctx context.Context, card models.CardModel) error
	// DeleteCard удаляет карту по ID
	DeleteCard(ctx context.Context, id string, uid string) error
	// SearchCards возвращает карты по запросу
	SearchCards(ctx context.Context, search models.CardSearch) ([]models.CardModel, error)
}

// CardRepo репозиторий карт
type CardRepo struct {
	pool *pgxpool.Pool
}

// GetCardByID возвращает карту по ID
func (r *CardRepo) GetCardByID(ctx context.Context, id string, uid string) (models.CardModel, error) {
	model := models.Card{}

	sql := `SELECT id, uid, name, number, meta, month, year, cvc, pin, created_at FROM cards WHERE id = $1 AND uid = $2`
	row := r.pool.QueryRow(ctx, sql, id, uid)
	err := row.Scan(&model.ID, &model.UID, &model.Name, &model.Number, &model.Meta, &model.Month, &model.Year, &model.CVC, &model.PIN, &model.CreatedAt)
	if err != nil {
		return nil, errors.New("card not found")
	}

	return &model, nil
}

// CreateCard создает новую карту
func (r *CardRepo) CreateCard(ctx context.Context, card models.CardModel) error {
	sql := "INSERT INTO cards (id, uid, name, meta, number, mask, month, year, cvc, pin, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"
	_, err := r.pool.Exec(ctx, sql, card.GetID(), card.GetUID(), card.GetName(), card.GetMeta(), card.GetNumber(), card.GetMask(), card.GetMonth(), card.GetYear(), card.GetCVC(), card.GetPIN(), card.GetCreatedAt())
	return err
}

// DeleteCard удаляет карту по ID
func (r *CardRepo) DeleteCard(ctx context.Context, id string, uid string) error {
	res, err := r.pool.Exec(ctx, "DELETE FROM cards WHERE id = $1 AND uid = $2", id, uid)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return errors.New("card not found")
	}

	return err
}

// SearchCards возвращает карты по запросу
func (r *CardRepo) SearchCards(ctx context.Context, search models.CardSearch) ([]models.CardModel, error) {
	sql := "SELECT id, uid, name, number, mask, meta, month, year, cvc, pin, created_at FROM cards"
	where, values := search.GetSubQuery()
	sql += where

	rows, err := r.pool.Query(ctx, sql, values...)
	if err != nil {
		return nil, errors.New("cards search error")
	}
	defer rows.Close()

	result := make([]models.CardModel, 0)
	for rows.Next() {
		var card = models.Card{}
		err := rows.Scan(&card.ID, &card.UID, &card.Name, &card.Number, &card.Mask, &card.Meta, &card.Month, &card.Year, &card.CVC, &card.PIN, &card.CreatedAt)
		if err != nil {
			logger.Error("account scan error", logger.StrArg("error", err.Error()))
			continue
		}
		result = append(result, &card)
	}

	return result, nil
}
