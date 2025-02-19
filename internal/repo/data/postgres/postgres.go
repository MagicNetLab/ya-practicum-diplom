package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRepository(cnf conf.Configurator) (*Repository, error) {
	r := &Repository{}

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cnf.DBUser(),
		cnf.DBPassword(),
		cnf.DBHost(),
		cnf.DBPort(),
		cnf.DBName())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		logger.Error("failed to create pgsql connection pool", logger.StrArg("error", err.Error()))
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		logger.Error("failed to ping pgsql connection pool", logger.StrArg("error", err.Error()))
		return nil, err
	}

	r.pool = pool

	return r, nil
}
