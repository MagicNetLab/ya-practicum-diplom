package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenLifeTime время жизни токена
const TokenLifeTime = time.Hour * 3

// Claims структура токена
type Claims struct {
	jwt.RegisteredClaims
	UID string
}
