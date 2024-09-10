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

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"github.com/sirupsen/logrus"
// )

// type LoadRemoteEnvParams struct {
// 	Host     string
// 	Token    string
// 	AppName  string
// 	Env      string
// 	Attempts int
// }

// func LoadEnv() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		logrus.Fatal("Erro carregando arquivo .env")
// 	}

// 	logrus.Info("Carregamento do arquivo .env realizado")
// }

// func LoadRemoteEnv(loadRemoteEnvParams ...LoadRemoteEnvParams) {
// 	LoadEnv()

// 	params := getLoadRemoteEnvParams(loadRemoteEnvParams...)

// 	logrus.Trace("Carregando variaveis de ambiente remotamente.")

// 	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/%s", params.Host, params.AppName, params.Env), nil)
// 	if err != nil {
// 		logrus.Fatalf("%s", err)
// 	}

// 	token := GetEnv("CLOUD_PROPERTIES_TOKEN", "")
// 	if token != "" {
// 		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
// 	}

// 	req.Header.Set("Accept", "application/json")

// 	resp, err := http.DefaultClient.Do(req)

// 	if err != nil {
// 		logrus.Fatalf("%s", err)
// 	}

// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		if params.Attempts < 5 {
// 			logrus.Fatalf("Não foi possível obter a configuração: %s, relizando uma nova tentativa", resp.Status)
// 			params.Attempts++
// 			autoRefresh(params)
// 			return
// 		}

// 		logrus.Fatalf("Não foi possível obter a configuração: %s", resp.Status)
// 	}

// 	params.Attempts = 0

// 	var config map[string]interface{}
// 	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
// 		logrus.Fatalf("%s", err)
// 	}

// 	for k, v := range config {
// 		err := os.Setenv(strings.ToUpper(k), fmt.Sprintf("%v", v))
// 		if err != nil {
// 			return
// 		}
// 	}

// 	autoRefresh(params)
// }

// func GetEnv(variable string, default_ ...string) string {
// 	value := os.Getenv(variable)

// 	if value == "" {
// 		if len(default_) > 0 {
// 			return default_[0]
// 		}
// 		logrus.Fatalf("A variavel de ambiente %s não foi definida", variable)
// 	}

// 	return value
// }

// func GetEnvBool(variable string, default_ ...string) bool {
// 	value := GetEnv(variable, default_...)

// 	b, err := strconv.ParseBool(value)
// 	if err != nil {
// 		return false
// 	}
// 	return b
// }

// func getLoadRemoteEnvParams(params ...LoadRemoteEnvParams) LoadRemoteEnvParams {
// 	var p LoadRemoteEnvParams
// 	if len(params) > 0 {
// 		p = params[0]
// 	}

// 	if p.Host == "" {
// 		p.Host = GetEnv("CLOUD_PROPERTIES_HOST")
// 	}
// 	if p.Token == "" {
// 		p.Token = GetEnv("CLOUD_PROPERTIES_TOKEN", "")
// 	}
// 	if p.AppName == "" {
// 		p.AppName = GetEnv("APP_NAME")
// 	}
// 	if p.Env == "" {
// 		p.Env = GetEnv("ENV")
// 	}
// 	if p.Attempts == 0 {
// 		p.Attempts = 1
// 	}

// 	return p
// }

// func autoRefresh(params LoadRemoteEnvParams) {
// 	refreshTime := GetEnv("ENV_REFRESH_TIME", "1")
// 	refreshInterval, err := strconv.Atoi(refreshTime)

// 	if err != nil {
// 		logrus.Fatal("ENV_REFRESH_TIME deve ser um número inteiro")
// 	}

// 	go func() {
// 		time.Sleep(time.Duration(refreshInterval) * time.Minute)
// 		LoadRemoteEnv(params)
// 	}()
// }
