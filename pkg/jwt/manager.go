/*
 * Copyright 2023 Caio Matheus Marcatti Calimério
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package jwt

// func NewJWTManager(signingKey ...string) *JWTManager {
// 	var key []byte
// 	if len(signingKey) > 0 && signingKey[0] != "" {
// 		key = []byte(signingKey[0])
// 	} else {
// 		key = []byte(env.GetEnv("JTW_SECRET"))
// 	}
// 	return &JWTManager{signingKey: []byte(key)}
// }

// func (manager *JWTManager) GenerateToken(expirationTime time.Duration, data map[string]interface{}) (string, error) {
// 	claims := jwt.MapClaims{
// 		"exp": time.Now().Add(expirationTime).Unix(),
// 	}

// 	for k, v := range data {
// 		claims[k] = v
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	return token.SignedString(manager.signingKey)
// }

// func (manager *JWTManager) ValidateToken(tokenString string) (jwt.Claims, error) {
// 	claims := jwt.MapClaims{}

// 	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 		return manager.signingKey, nil
// 	})

// 	if err != nil || !token.Valid {
// 		return nil, fmt.Errorf("token inválido")
// 	}

// 	return claims, nil
// }

// func (manager *JWTManager) DecodeToken(tokenString string) (map[string]interface{}, error) {
// 	claims := jwt.MapClaims{}

// 	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 		return manager.signingKey, nil
// 	})

// 	if err != nil {
// 		return nil, fmt.Errorf("erro ao decodificar o token: %w", err)
// 	}

// 	decodedData := make(map[string]interface{})
// 	for k, v := range claims {
// 		decodedData[k] = v
// 	}

// 	return decodedData, nil
// }
