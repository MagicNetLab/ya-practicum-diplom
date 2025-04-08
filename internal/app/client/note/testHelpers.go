package note

import (
	"bytes"
	"io"
	"os"
	"testing"
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

// mockStdin эмуляция ввода данных пользователем
func mockStdin(t *testing.T, input string) func() {
	oldStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	os.Stdin = r

	return func() {
		os.Stdin = oldStdin
	}
}
