package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type jwtTokenGenarator struct {
	secretKey string
}

func NewJWTToKenGenarator(secret string) *jwtTokenGenarator {
	return &jwtTokenGenarator{secretKey: secret}
}

func (j *jwtTokenGenarator) GenerateToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
