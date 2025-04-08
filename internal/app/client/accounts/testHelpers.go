package accounts

import (
	"bytes"
	"io"
	"os"
)

// Временно переопределяем функцию ReadTokenFromFile для тестирования
var originalReadTokenFromFile = readTokenFromFile

// Восстанавливаем оригинальную функцию после тестов
func restoreReadTokenFromFile() {
	readTokenFromFile = originalReadTokenFromFile
}

// captureStdoutOutput перехватывает вывод в stdout
func captureStdoutOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		if err != nil {
			panic(err)
		}
		outC <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = oldStdout
	return <-outC
}
