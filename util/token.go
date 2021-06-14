package util

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(userID string, jwtKey string, expireDuration time.Duration) (string, error) {
	expire := time.Now().Add(expireDuration)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.UnixNano(),
		},
	})

	return token.SignedString([]byte(jwtKey))
}

func ParseToken(tokenStr, jwtKey string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(tk *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("Unauthorized")
	}

	return claims.ID, nil
}

type Claims struct {
	ID string
	jwt.StandardClaims
}
