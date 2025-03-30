package client

import "os"

type FileManager interface {
	Read(path string) ([]byte, error)
}

type FileReader struct{}

func (r FileReader) Read(path string) ([]byte, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
