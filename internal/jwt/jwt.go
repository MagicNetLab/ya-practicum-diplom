package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo/models"
)

// GenerateToken генерирует JWT
func GenerateToken(user models.UserEntity, jwtSecret string) (string, error) {
	uid, err := user.GetUID()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenLifeTime)),
		},
		UID: uid,
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshToken(user models.UserEntity, jwtSecret string) (string, error) {

	return "", nil
}

func ParseToken(token string, jwtSecret string) (*Claims, error) {

	return nil, nil
}

func VerifyToken(token string, jwtSecret string) bool {

	return false
}

// GetRandomSecret возвращает случайный секрет для генерации JWT
func GetRandomSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(b)
}
