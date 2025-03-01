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
	id          string    `db:"id"`
	uid         string    `db:"uid"`
	login       string    `db:"login"`
	password    string    `db:"password"`
	url         string    `db:"url"`
	description string    `db:"description"`
	createdAt   time.Time `db:"created_at"`
	updatedAt   time.Time `db:"updated_at"`
}

// GetID получение ID аккаунта
func (a *Account) GetID() string {
	return a.id
}

// SetID установка ID аккаунта
func (a *Account) SetID(id string) error {
	if id == "" {
		return errors.New("id is not be empty")
	}

	if err := uuid.Validate(id); err != nil {
		return errors.New("invalid id")
	}

	a.id = id
	return nil
}

// GetUID получение UID владельца аккаунта
func (a *Account) GetUID() string {
	return a.uid
}

// SetUID установка UID владельца аккаунта
func (a *Account) SetUID(uid string) error {
	if uid == "" {
		return errors.New("uid is not be empty")
	}

	if err := uuid.Validate(uid); err != nil {
		return errors.New("invalid uid")
	}

	a.uid = uid
	return nil
}

// GetLogin получение логина аккаунта
func (a *Account) GetLogin() string {
	return a.login
}

// SetLogin установка логина аккаунта
func (a *Account) SetLogin(login string) error {
	if login == "" {
		return errors.New("login is not be empty")
	}

	a.login = login
	return nil
}

// GetPassword получение пароля аккаунта
func (a *Account) GetPassword() string {
	return a.password
}

// SetPassword установка пароля аккаунта
func (a *Account) SetPassword(password string) error {
	if password == "" {
		return errors.New("password is not be empty")
	}
	a.password = password
	return nil
}

// GetURL получение URL аккаунта
func (a *Account) GetURL() string {
	return a.url
}

// SetURL установка URL аккаунта
func (a *Account) SetURL(url string) error {
	if url == "" {
		return errors.New("url is not be empty")
	}

	a.url = url
	return nil
}

// GetDescription получение описания аккаунта
func (a *Account) GetDescription() string {
	return a.description
}

// SetDescription установка описания аккаунта
func (a *Account) SetDescription(description string) error {
	a.description = description
	return nil
}

// GetCreatedAt получение даты создания аккаунта
func (a *Account) GetCreatedAt() time.Time {
	return a.createdAt
}

// SetCreatedAt установка даты создания аккаунта
func (a *Account) SetCreatedAt(createdAt time.Time) error {
	if createdAt.IsZero() {
		return errors.New("created_at is empty")
	}

	if createdAt.After(time.Now()) {
		return errors.New("invalid created_at")
	}

	a.createdAt = createdAt
	return nil
}

// GetUpdatedAt получение даты обновления аккаунта
func (a *Account) GetUpdatedAt() time.Time {
	return a.updatedAt
}

// SetUpdatedAt установка даты обновления аккаунта
func (a *Account) SetUpdatedAt(updatedAt time.Time) error {
	if updatedAt.IsZero() {
		return errors.New("updated_at is empty")
	}

	if updatedAt.After(time.Now()) {
		return errors.New("invalid updated_at")
	}

	a.updatedAt = updatedAt
	return nil
}
