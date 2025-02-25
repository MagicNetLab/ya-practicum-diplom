package postgres

import "errors"

// UserModel модель пользователя
type UserModel struct {
	UID       string `sql:"uid"`
	Login     string `sql:"login"`
	Password  string `sql:"password"`
	IsRefresh bool   `sql:"is_refresh"`
	Expired   string `sql:"expired"`
}

// GetUID получение UID пользователя
func (u *UserModel) GetUID() (string, error) {
	if u.UID == "" {
		return "", errors.New("uid is not set")
	}
	return u.UID, nil
}

// GetLogin получение логина пользователя
func (u *UserModel) GetLogin() (string, error) {
	if u.Login == "" {
		return "", errors.New("login is not set")
	}
	return u.Login, nil
}

// GetPassword получение пароля пользователя
func (u *UserModel) GetPassword() (string, error) {
	if u.Password == "" {
		return "", errors.New("password is not set")
	}
	return u.Password, nil
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
