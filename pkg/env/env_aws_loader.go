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
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func awsLoader() error {
	region := os.Getenv("AWS_SECRET_MANAGER_REGION")
	accessKey := os.Getenv("AWS_SECRET_MANAGER_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_MANAGER_SECRET_KEY")
	secretName := os.Getenv("AWS_SECRET_MANAGER_NAME")

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(region),
	})

	if err != nil {
		return fmt.Errorf("erro ao criar sessão: %v", err)
	}

	// Criar um novo cliente do Secrets Manager
	svc := secretsmanager.New(sess)

	// Criar a estrutura de requisição
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	// Chamar o serviço do Secrets Manager para obter o valor do segredo
	result, err := svc.GetSecretValue(input)

	if err != nil {
		return fmt.Errorf("erro ao recuperar o segredo: %v", err)
	}

	if result.SecretString != nil {
		return injectEnvVariables(*result.SecretString)
	} else {
		return fmt.Errorf("segredo não contém string")
	}
}

// Função que injeta valores de configuração no ambiente a partir de uma string JSON
func injectEnvVariables(jsonData string) error {
	var data map[string]string
	err := json.Unmarshal([]byte(jsonData), &data)

	if err != nil {
		return fmt.Errorf("erro ao fazer parse do JSON: %v", err)
	}

	for key, value := range data {
		err := os.Setenv(key, value)

		if err != nil {
			return fmt.Errorf("erro ao definir variável de ambiente: %v", err)
		}
	}

	return nil
}
