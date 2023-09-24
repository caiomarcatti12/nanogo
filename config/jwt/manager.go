package jwt

import (
	"fmt"
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func NewJWTManager(signingKey ...string) *JWTManager {
	var key []byte
	if len(signingKey) > 0 && signingKey[0] != "" {
		key = []byte(signingKey[0])
	} else {
		key = []byte(env.GetEnv("JTW_SECRET"))
	}
	return &JWTManager{signingKey: []byte(key)}
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

func (manager *JWTManager) DecodeToken(tokenString string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return manager.signingKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("erro ao decodificar o token: %w", err)
	}

	decodedData := make(map[string]interface{})
	for k, v := range claims {
		decodedData[k] = v
	}

	return decodedData, nil
}
