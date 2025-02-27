package models

import "time"

// UserEntity интерфейс пользователя
type UserEntity interface {
	GetUID() string
	GetLogin() string
	GetPassword() string
}

// TokenEntity интерфейс токена
type TokenEntity interface {
	GetValue() string
	GetUID() string
	IsRefresh() bool
	GetExpired() time.Time
}
