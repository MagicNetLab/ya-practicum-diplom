package postgres

import "time"

// UserModel модель пользователя
type UserModel struct {
	UID       string    `db:"uid"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// GetUID получение UID пользователя
func (u UserModel) GetUID() string {
	return u.UID
}

// GetLogin получение логина пользователя
func (u UserModel) GetLogin() string {
	return u.Login
}

// GetPassword получение пароля пользователя
func (u UserModel) GetPassword() string {
	return u.Password
}
