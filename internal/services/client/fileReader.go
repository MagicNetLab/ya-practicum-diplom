package client

import "os"

// FileReader интерфейс сервис чтения файлов из файловой системы.
type FileReader interface {
	Read(path string) ([]byte, error)
}

// Reader сервис чтения файлов из файловой системы.
type Reader struct{}

// Read чтение содержимого файла
func (r Reader) Read(path string) ([]byte, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
