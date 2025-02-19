package jwt

import (
	"crypto/rand"
	"encoding/hex"
)

// GetRandomSecret возвращает случайный секрет для генерации JWT
func GetRandomSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
