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
package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// fileLoader verifica se o arquivo .env existe e carrega suas variáveis
func fileLoader() error {
	const envFilePath = "../../configs/.env"

	// Verifica se o arquivo .env existe
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		return fmt.Errorf("Arquivo %s não encontrado", envFilePath)
	}

	// Carrega o arquivo .env
	if err := godotenv.Load(envFilePath); err != nil {
		return fmt.Errorf("Erro ao carregar o arquivo %s: %v", envFilePath, err)
	}

	return nil
}
