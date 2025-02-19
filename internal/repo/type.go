package repo

type storage struct {
	data DataStorage
}

// DataStorage интерфейс хранилища данных
type DataStorage interface {
	Close() error
}

// FileStorage интерфейс хранилища файлов
type FileStorage interface {
	Close() error
}
