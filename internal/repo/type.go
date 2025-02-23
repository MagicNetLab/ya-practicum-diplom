package repo

type storage struct {
	data  DataStorage
	files FileStorage
}

// DataStorage интерфейс хранилища данных
type DataStorage interface {
	Close() error
}

// FileStorage интерфейс хранилища файлов
type FileStorage interface {
	Close() error
}
