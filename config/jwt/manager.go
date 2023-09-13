package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func NewJWTManager(signingKey string) *JWTManager {
	return &JWTManager{signingKey: []byte(signingKey)}
}

func (manager *JWTManager) GenerateToken(expirationTime time.Duration, data map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(expirationTime).Unix(),
	}

	for k, v := range data {
		claims[k] = v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(manager.signingKey)
}

func (manager *JWTManager) ValidateToken(tokenString string) (jwt.Claims, error) {
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return manager.signingKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token inválido")
	}

	return claims, nil
}
