package models

import "time"

// UserEntity интерфейс пользователя
type UserEntity interface {
	GetUID() string
	GetLogin() string
	GetPassword() string
	//SetUID(uid string) error
	//SetLogin(login string) error
	//SetPassword(password string) error
}

// TokenEntity интерфейс токена
type TokenEntity interface {
	GetValue() string
	SetValue(value string) error
	GetUID() string
	SetUID(uid string) error
	SetRefresh()
	IsRefresh() bool
	SetExpired(expired time.Time) error
	GetExpired() time.Time
}
