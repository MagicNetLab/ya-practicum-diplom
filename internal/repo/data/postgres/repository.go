package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

// Close закрытие соединения с БД
func (r *Repository) Close() error {
	r.pool.Close()
	return nil
}
