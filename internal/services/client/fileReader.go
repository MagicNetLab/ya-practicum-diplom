package client

import "os"

// FileReader интерфейс сервис чтения файлов из файловой системы.
type FileManager interface {
	Read(path string) ([]byte, error)
}

// FileReader сервис чтения файлов из файловой системы.
type FileReader struct{}

// Read чтение содержимого файла
func (r FileReader) Read(path string) ([]byte, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
