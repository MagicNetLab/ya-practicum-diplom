package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// AccountModel интерфейс модели данных аккаунта
type AccountModel interface {
	GetID() string
	SetID(id string) error
	GetUID() string
	SetUID(uid string) error
	GetLogin() string
	SetLogin(login string) error
	GetPassword() string
	SetPassword(password string) error
	GetURL() string
	SetURL(url string) error
	GetDescription() string
	SetDescription(description string) error
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time) error
	GetUpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time) error
}

// Account модель данных аккаунта
type Account struct {
	ID          string    `db:"ID"`
	UID         string    `db:"UID"`
	Login       string    `db:"Login"`
	Password    string    `db:"Password"`
	URL         string    `db:"URL"`
	Description string    `db:"Description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

// GetID получение ID аккаунта
func (a *Account) GetID() string {
	return a.ID
}

// SetID установка ID аккаунта
func (a *Account) SetID(id string) error {
	if id == "" {
		return errors.New("ID is not be empty")
	}

	if err := uuid.Validate(id); err != nil {
		return errors.New("invalid ID")
	}

	a.ID = id
	return nil
}

// GetUID получение UID владельца аккаунта
func (a *Account) GetUID() string {
	return a.UID
}

// SetUID установка UID владельца аккаунта
func (a *Account) SetUID(uid string) error {
	if uid == "" {
		return errors.New("UID is not be empty")
	}

	if err := uuid.Validate(uid); err != nil {
		return errors.New("invalid UID")
	}

	a.UID = uid
	return nil
}

// GetLogin получение логина аккаунта
func (a *Account) GetLogin() string {
	return a.Login
}

// SetLogin установка логина аккаунта
func (a *Account) SetLogin(login string) error {
	if login == "" {
		return errors.New("Login is not be empty")
	}

	a.Login = login
	return nil
}

// GetPassword получение пароля аккаунта
func (a *Account) GetPassword() string {
	return a.Password
}

// SetPassword установка пароля аккаунта
func (a *Account) SetPassword(password string) error {
	if password == "" {
		return errors.New("Password is not be empty")
	}
	a.Password = password
	return nil
}

// GetURL получение URL аккаунта
func (a *Account) GetURL() string {
	return a.URL
}

// SetURL установка URL аккаунта
func (a *Account) SetURL(url string) error {
	if url == "" {
		return errors.New("URL is not be empty")
	}

	a.URL = url
	return nil
}

// GetDescription получение описания аккаунта
func (a *Account) GetDescription() string {
	return a.Description
}

// SetDescription установка описания аккаунта
func (a *Account) SetDescription(description string) error {
	a.Description = description
	return nil
}

// GetCreatedAt получение даты создания аккаунта
func (a *Account) GetCreatedAt() time.Time {
	return a.CreatedAt
}

// SetCreatedAt установка даты создания аккаунта
func (a *Account) SetCreatedAt(createdAt time.Time) error {
	if createdAt.IsZero() {
		return errors.New("created_at is empty")
	}

	if createdAt.After(time.Now()) {
		return errors.New("invalid created_at")
	}

	a.CreatedAt = createdAt
	return nil
}

// GetUpdatedAt получение даты обновления аккаунта
func (a *Account) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}

// SetUpdatedAt установка даты обновления аккаунта
func (a *Account) SetUpdatedAt(updatedAt time.Time) error {
	if updatedAt.IsZero() {
		return errors.New("updated_at is empty")
	}

	if updatedAt.After(time.Now()) {
		return errors.New("invalid updated_at")
	}

	a.UpdatedAt = updatedAt
	return nil
}
