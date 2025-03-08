package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
)

// Repository интерфейс общего репозитория
type Repository interface {
	GetAuthRepo() AuthRepository
	GetAccountRepo() AccountRepository
	GetCardRepo() CardRepository
	GetNoteRepo() NoteRepository
	GetFileRepo() FileRepository
	Close()
}

// NewRepository конструктор репозитория
func NewRepository(cnf config.DataBaseConfigurator) (Repository, error) {
	pool, err := pgxpool.New(context.Background(), cnf.GetDSN())
	if err != nil {
		return nil, errors.Wrap(err, "failed connect to database")
	}

	r := &AppRepository{
		auth:    NewAuthRepository(pool),
		account: NewAccountRepo(pool),
		cards:   NewCardRepo(pool),
		notes:   NewNoteRepository(pool),
		files:   NewFileRepository(pool),
	}

	return r, nil
}

// AppRepository репозиторий приложения
type AppRepository struct {
	auth    AuthRepository
	account AccountRepository
	cards   CardRepository
	notes   NoteRepository
	files   FileRepository
	pool    *pgxpool.Pool
}

// GetAuthRepo возвращает репозиторий авторизации
func (r *AppRepository) GetAuthRepo() AuthRepository {
	return r.auth
}

// GetAccountRepo возвращает репозиторий аккаунта
func (r *AppRepository) GetAccountRepo() AccountRepository {
	return r.account
}

// GetCardRepo возвращает репозиторий карты
func (r *AppRepository) GetCardRepo() CardRepository {
	return r.cards
}

// GetNoteRepo возвращает репозиторий заметки
func (r *AppRepository) GetNoteRepo() NoteRepository {
	return r.notes
}

func (r *AppRepository) GetFileRepo() FileRepository {
	return r.files
}

// Close закрывает соединение с базой данных
func (r *AppRepository) Close() {
	r.pool.Close()
}
