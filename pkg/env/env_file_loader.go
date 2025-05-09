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
	"log"
	"os"

	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"github.com/caiomarcatti12/nanogo/pkg/util"
	"github.com/joho/godotenv"
)

type FileEnvLoader struct {
	i18n i18n.I18N
}

func NewFileEnvLoader(i18n i18n.I18N) *FileEnvLoader {
	return &FileEnvLoader{i18n: i18n}
}

// Load carrega variáveis de ambiente a partir de um arquivo .env
// Verifica múltiplos locais possíveis até encontrar o arquivo.
func (f *FileEnvLoader) Load() error {
	possiblePaths := []string{
		util.GetExecutableAbsolutePath(".env"),
		util.GetExecutableAbsolutePath("configs/.env"),
		os.Getenv("ENV_FILE_PATH"),
	}

	for _, path := range possiblePaths {
		log.Println(fmt.Sprintf("Loading env file: %s", path))
		if _, err := os.Stat(path); err == nil {
			if err := godotenv.Load(path); err != nil {
				log.Fatalf(f.i18n.Get("env.load_error", map[string]interface{}{"path": path, "error": err}))
			}
			log.Printf(f.i18n.Get("env.load_success", map[string]interface{}{"path": path}))
			return nil
		}
	}

	log.Fatal(f.i18n.Get("env.not_found", nil))
	return nil
}
