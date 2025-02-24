package inmemory

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

type Repository struct {
	users    map[string]User
	tokens   map[string]Token
	dumpPath string
}

// Close закрытие репозитория
func (r *Repository) Close() error {
	err := r.Dump()
	if err != nil {
		logger.Error("failed save repository to dump", logger.StrArg("error", err.Error()))
	}
	return nil
}

// Import импорт данных из дампа
func (r *Repository) Import() error {
	if r.dumpPath != "" {
		usersData := make(map[string]User, 0)
		userFile := r.dumpPath + "/users.json"
		if _, err := os.Stat(userFile); err == nil {
			f, err := os.OpenFile(r.dumpPath+"/users.json", os.O_RDONLY|os.O_WRONLY|os.O_CREATE, 0666)
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
				row := User{}
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
			rowData := User{
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
