package utils

import (
	"gamershub/config"
	"gamershub/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var CFG = config.LoadConfig()

func GenerateJWT(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * time.Duration(CFG.JWTTTL)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(CFG.JWTSecret))
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(CFG.JWTSecret), nil
	})
}
