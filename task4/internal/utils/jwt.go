package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWTToken(secret string, userID uint, username string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"exp":      time.Now().Add(duration).Unix(),
		"iat":      time.Now().Unix(),
		"username": username,
		"user_id":  userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ParseJWTToken(secret, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
