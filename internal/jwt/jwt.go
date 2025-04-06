package jwt

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

// TokenLifeTime время жизни токена
const TokenLifeTime = time.Hour * 3

// Claims структура токена
type Claims struct {
	jwt.RegisteredClaims
	UID string
}

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

// ParseToken парсинг JWT токена
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

// VerifyToken проверяет JWT
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

// ExtractToken извлекает токен из контекста
func ExtractToken(ctx context.Context) (string, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("failed to parse metadata")
	}

	val := meta.Get("token")
	if len(val) == 0 {
		return "", errors.New("no token provided")
	}

	token := val[0]
	return token, nil
}

// GetUIDFromContext извлекает UID из контекста
func GetUIDFromContext(ctx context.Context, jwtSecret string) (string, error) {
	tokenString, err := ExtractToken(ctx)
	if err != nil {
		return "", err
	}

	tokenData, err := ParseToken(tokenString, jwtSecret)
	if err != nil {
		return "", err
	}

	if tokenData.UID == "" {
		return "", errors.New("invalid token")
	}

	if err := uuid.Validate(tokenData.UID); err != nil {
		return "", err
	}

	return tokenData.UID, nil
}

// ReadTokenFromFile парсинг JWT токена из файла
func ReadTokenFromFile() (string, error) {
	fileName := "token.txt"
	token := ""
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		token = strings.TrimSpace(scanner.Text())
		break
	}

	if token == "" {
		return "", fmt.Errorf("token not found")
	}

	return token, nil
}

// SaveTokenToFile сохраняет токен в файл
func SaveTokenToFile(token string) error {
	fileName := "token.txt"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(token)
	if err != nil {
		return err
	}

	return nil
}
