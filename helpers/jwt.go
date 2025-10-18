package helpers

import (
	"time"

	"github.com/ahmadalaik/desa-digital/config"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(config.GetEnv("JWT_SECRET", ""))

func GenerateToken(username string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenStr, err
}
