package models

import "time"

// UserEntity интерфейс пользователя
type UserEntity interface {
	GetUID() (string, error)
	GetLogin() (string, error)
	GetPassword() (string, error)
	SetUID(uid string) error
	SetLogin(login string) error
	SetPassword(password string) error
}

// TokenEntity интерфейс токена
type TokenEntity interface {
	GetValue() (string, error)
	SetValue(value string) error
	GetUID() (string, error)
	SetUID(uid string) error
	SetRefresh() error
	IsRefresh() (bool, error)
	SetExpired(expired time.Time) error
	GetExpired() (time.Time, error)
}
