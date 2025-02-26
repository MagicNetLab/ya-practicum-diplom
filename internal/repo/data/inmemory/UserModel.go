package inmemory

import "errors"

// UserModel модель пользователя
type UserModel struct {
	UID      string `json:"uid"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

// GetUID получение UID пользователя
func (u *UserModel) GetUID() string {
	return u.UID
}

// GetLogin получение логина пользователя
func (u *UserModel) GetLogin() string {
	return u.Login
}

// GetPassword получение пароля пользователя
func (u *UserModel) GetPassword() string {
	return u.Password
}

// SetUID установка UID пользователя
func (u *UserModel) SetUID(uid string) error {
	if uid == "" {
		return errors.New("uid is not be empty")
	}
	u.UID = uid
	return nil
}

// SetLogin установка логина пользователя
func (u *UserModel) SetLogin(login string) error {
	if login == "" {
		return errors.New("login is not be empty")
	}
	u.Login = login
	return nil
}

// SetPassword установка пароля пользователя
func (u *UserModel) SetPassword(password string) error {
	if password == "" {
		return errors.New("password is not be empty")
	}
	u.Password = password
	return nil
}
