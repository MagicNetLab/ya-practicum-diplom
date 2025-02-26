package inmemory

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/models"
)

type Repository struct {
	users    map[string]UserModel
	tokens   map[string]TokenModel
	dumpPath string
}

// AddToken добавление токена в бд
func (r *Repository) AddToken(ctx context.Context, value, uid string, isRefresh bool, expired time.Time) error {
	if _, ok := r.tokens[value]; ok {
		return errors.New("token already exists")
	}
	token := TokenModel{
		Value:   value,
		UID:     uid,
		Refresh: isRefresh,
		Expired: expired.Format(time.RFC3339),
	}

	r.tokens[value] = token
	return nil
}

// RemoveToken удаление токена из бд
func (r *Repository) RemoveToken(ctx context.Context, value string) error {
	delete(r.tokens, value)
	return nil
}

// RemoveUserTokens удаление всех токенов пользователя
func (r *Repository) RemoveUserTokens(ctx context.Context, uid string) error {
	for k, v := range r.tokens {
		if v.UID == uid {
			delete(r.tokens, k)
		}
	}

	return nil
}

// HasToken проверка наличия токена
func (r *Repository) HasToken(ctx context.Context, value, uid string, isRefresh bool) (bool, error) {
	token, ok := r.tokens[value]
	if !ok {
		return false, errors.New("token not found")
	}

	if token.Refresh == isRefresh && token.UID == uid {
		return true, nil
	}

	return false, nil
}

// GetUserByUID получение пользователя по UID
func (r *Repository) GetUserByUID(ctx context.Context, uid string) (models.UserEntity, error) {
	user, ok := r.users[uid]
	if !ok {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetUserByLogin получение пользователя по логину
func (r *Repository) GetUserByLogin(ctx context.Context, login string) (models.UserEntity, error) {
	for _, v := range r.users {
		if v.Login == login {
			return &v, nil
		}
	}
	return nil, errors.New("user not found")
}

// GetUserByLoginPassword получение пользователя по логину и паролю
func (r *Repository) GetUserByLoginPassword(ctx context.Context, login, password string) (models.UserEntity, error) {
	for _, v := range r.users {
		if v.Login == login && v.Password == password {
			return &v, nil
		}
	}
	return nil, errors.New("user not found")
}

// CreateUser создание пользователя
func (r *Repository) CreateUser(ctx context.Context, login, password string) (models.UserEntity, error) {
	for _, v := range r.users {
		if v.Login == login {
			return nil, errors.New("user already exists")
		}
	}

	uid := uuid.New().String()
	user := UserModel{
		UID:      uid,
		Login:    login,
		Password: password,
	}
	r.users[uid] = user

	return &user, nil
}

// RemoveUserByUID удаление пользователя по UID
func (r *Repository) RemoveUserByUID(ctx context.Context, uid string) error {
	delete(r.users, uid)
	return nil
}

// RemoveUserByLogin удаление пользователя по логину
func (r *Repository) RemoveUserByLogin(ctx context.Context, login string) error {
	for k, v := range r.users {
		if v.Login == login {
			delete(r.users, k)
			return nil
		}
	}

	return nil
}

func (r *Repository) Close(ctx context.Context) error {
	err := r.Dump()
	if err != nil {
		logger.Error("failed save repository to dump", logger.StrArg("error", err.Error()))
	}
	return nil
}

// Import импорт данных из дампа
func (r *Repository) Import() error {
	if r.dumpPath != "" {
		usersData := make(map[string]UserModel, 0)
		userFile := r.dumpPath + "/users.json"
		if _, err := os.Stat(userFile); err == nil {
			f, err := os.OpenFile(userFile, os.O_RDONLY, 0666)
			if err != nil {
				logger.Error("failed import data from users.json", logger.StrArg("error", err.Error()))
			}
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					logger.Error("failed import close users.json", logger.StrArg("error", err.Error()))
				}
			}(f)

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				row := UserModel{}
				if err := json.Unmarshal(scanner.Bytes(), &row); err != nil {
					logger.Error("failed parse  users dump file", logger.StrArg("error", err.Error()))
					return err
				}
				usersData[row.UID] = row
			}
			r.users = usersData
		}
	}

	return nil
}

// Dump экспорт данных в дамп
func (r *Repository) Dump() error {
	if r.dumpPath != "" {
		file, err := os.OpenFile(r.dumpPath+"/users.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger.Error("failed open local storage file", logger.StrArg("error", err.Error()))
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				logger.Error("failed close local storage file", logger.StrArg("error", err.Error()))
			}
		}(file)

		writer := bufio.NewWriter(file)

		for _, user := range r.users {
			rowData := UserModel{
				UID:      user.UID,
				Login:    user.Login,
				Password: user.Password,
			}

			rowString, err := json.Marshal(rowData)
			if err != nil {
				logger.Error("failed serialize cache data", logger.StrArg("error", err.Error()))
				return err
			}

			_, err = writer.WriteString(string(rowString) + "\n")
			if err != nil {
				logger.Error("failed write local storage file", logger.StrArg("error", err.Error()))
				return err
			}

			if err := writer.Flush(); err != nil {
				logger.Error("failed flush local storage file", logger.StrArg("error", err.Error()))
				return err
			}
		}
	}

	return nil
}
