package utils

import (
	"errors"
	"fmt"
	"gamershub/config"
	"gamershub/internal/models"
	"gamershub/internal/types"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var c = config.LoadConfig()

type JWTClaims struct {
	UserId uint
	Role   types.Role
	jwt.RegisteredClaims
}

func GenerateJWT(user *models.User) (string, error) {
	claims := JWTClaims{
		UserId: user.Id,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.JWTSecret))
}

func ParseJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
