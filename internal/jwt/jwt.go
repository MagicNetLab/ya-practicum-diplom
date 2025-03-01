package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

// GenerateToken генерирует JWT
func GenerateToken(user models.UserModel, jwtSecret string) (string, error) {
	uid := user.GetUID()

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

func ParseToken(tokenString string, jwtSecret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})
	if err != nil {
		return claims, errors.New("failed check token: invalid token")
	}

	if !token.Valid {
		return claims, errors.New("failed check token: invalid token")
	}

	return claims, nil
}

func VerifyToken(token string, jwtSecret string) bool {
	_, err := ParseToken(token, jwtSecret)
	return err == nil
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
