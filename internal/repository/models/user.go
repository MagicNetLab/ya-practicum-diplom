package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// UserModel модель данных пользователя
type UserModel interface {
	GetUID() string
	SetUID(uid string) error
	GetLogin() string
	SetLogin(login string) error
	GetPassword() string
	SetPassword(password string) error
	GetCreatedAt() time.Time
	SetCreatedAt(createdAt time.Time) error
	GetUpdatedAt() time.Time
	SetUpdatedAt(updatedAt time.Time) error
}

type User struct {
	UID       string    `json:"ID"`
	Login     string    `json:"Login"`
	Password  string    `json:"Password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetUID получение UID пользователя
func (u *User) GetUID() string {
	return u.UID
}

// SetUID установка UID пользователя
func (u *User) SetUID(uid string) error {
	if uid == "" {
		return errors.New("UID is empty")
	}

	if err := uuid.Validate(uid); err != nil {
		return errors.New("invalid UID")
	}

	u.UID = uid
	return nil
}

// GetLogin получение логина пользователя
func (u *User) GetLogin() string {
	return u.Login
}

// SetLogin установка логина пользователя
func (u *User) SetLogin(login string) error {
	if login == "" {
		return errors.New("login is empty")
	}
	u.Login = login
	return nil
}

// GetPassword получение пароля пользователя
func (u *User) GetPassword() string {
	return u.Password
}

// SetPassword установка пароля пользователя
func (u *User) SetPassword(password string) error {
	if password == "" {
		return errors.New("password is empty")
	}
	u.Password = password
	return nil
}

// GetCreatedAt получение даты создания пользователя
func (u *User) GetCreatedAt() time.Time {
	if u.CreatedAt.IsZero() {
		return time.Now()
	}
	return u.CreatedAt
}

// SetCreatedAt установка даты создания пользователя
func (u *User) SetCreatedAt(createdAt time.Time) error {
	if createdAt.IsZero() {
		return errors.New("created_at is empty")
	}

	if createdAt.After(time.Now()) {
		return errors.New("invalid created_at")
	}

	u.CreatedAt = createdAt
	return nil
}

// GetUpdatedAt получение даты обновления пользователя
func (u *User) GetUpdatedAt() time.Time {
	if u.UpdatedAt.IsZero() {
		return time.Now()
	}
	return u.UpdatedAt
}

// SetUpdatedAt установка даты обновления пользователя
func (u *User) SetUpdatedAt(updatedAt time.Time) error {
	if updatedAt.IsZero() {
		return errors.New("updated_at is empty")
	}

	if updatedAt.After(time.Now()) {
		return errors.New("invalid updated_at")
	}

	u.UpdatedAt = updatedAt
	return nil
}
