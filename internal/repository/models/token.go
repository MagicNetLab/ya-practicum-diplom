package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// TokenModel модель данных токена
type TokenModel interface {
	GetID() string
	SetID(id string) error
	GetToken() string
	SetToken(token string) error
	GetUID() string
	SetUID(uid string) error
	IsRefresh() bool
	SetIsRefresh(refresh bool) error
	GetExpired() time.Time
	SetExpired(expired time.Time) error
}

// Token модель данных токена
type Token struct {
	id        string    `db:"ID"`
	token     string    `db:"token"`
	uid       string    `db:"UID"`
	isRefresh bool      `db:"is_refresh"`
	expired   time.Time `db:"expired"`
}

// GetID возвращает идентификатор токена
func (t *Token) GetID() string {
	return t.id
}

// SetID устанавливает идентификатор токена
func (t *Token) SetID(id string) error {
	if id == "" {
		return errors.New("ID is not be empty")
	}

	if err := uuid.Validate(id); err != nil {
		return errors.New("invalid ID")
	}

	t.id = id
	return nil
}

// GetToken возвращает токен
func (t *Token) GetToken() string {
	return t.token
}

// SetToken устанавливает значение токен
func (t *Token) SetToken(token string) error {
	if token == "" {
		return errors.New("token is not be empty")
	}

	t.token = token
	return nil
}

// GetUID возвращает идентификатор владельца токена
func (t *Token) GetUID() string {
	return t.uid
}

// SetUID устанавливает идентификатор владельца токена
func (t *Token) SetUID(uid string) error {
	if uid == "" {
		return errors.New("UID is not be empty")
	}

	if err := uuid.Validate(uid); err != nil {
		return errors.New("invalid UID")
	}

	t.uid = uid
	return nil
}

// IsRefresh возвращает является ли токен токеном обновления
func (t *Token) IsRefresh() bool {
	return t.isRefresh
}

// SetIsRefresh устанавливает является ли токен токеном обновления
func (t *Token) SetIsRefresh(refresh bool) error {
	t.isRefresh = refresh
	return nil
}

// GetExpired возвращает время истечения токена
func (t *Token) GetExpired() time.Time {
	return t.expired
}

// SetExpired устанавливает время истечения токена
func (t *Token) SetExpired(expired time.Time) error {
	if expired.IsZero() {
		return errors.New("expired is not be empty")
	}

	if expired.Before(time.Now()) {
		return errors.New("expired is not be in the past")
	}

	t.expired = expired
	return nil
}
