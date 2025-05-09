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
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
)

type AWSLoader struct {
	i18n i18n.I18N
}

func NewAWSLoader(i18n i18n.I18N) *AWSLoader {
	return &AWSLoader{i18n: i18n}
}

// LoadSecrets carrega segredos do AWS Secrets Manager e injeta-os no ambiente.
// Caso os parâmetros sejam vazios, utiliza variáveis de ambiente como fallback.
func (a *AWSLoader) LoadSecrets(region, secretName, roleArn, webIdentityTokenFile, accessKey, secretKey string) error {
	region = a.getOrDefault(region, "AWS_SECRET_MANAGER_REGION")
	secretName = a.getOrDefault(secretName, "AWS_SECRET_MANAGER_NAME")

	a.validateMandatoryEnv("AWS_SECRET_MANAGER_REGION", region)
	a.validateMandatoryEnv("AWS_SECRET_MANAGER_NAME", secretName)

	creds, err := a.getAWSCredentials(region, roleArn, webIdentityTokenFile, accessKey, secretKey)
	if err != nil {
		return err
	}

	return a.retrieveAndInjectSecrets(creds, region, secretName)
}

func (a *AWSLoader) getOrDefault(value, envKey string) string {
	if strings.TrimSpace(value) == "" {
		return os.Getenv(envKey)
	}
	return value
}

func (a *AWSLoader) validateMandatoryEnv(variable, value string) {
	if strings.TrimSpace(value) == "" {
		panic(a.i18n.Get("env.not_found", map[string]interface{}{"variable": variable}))
	}
}

func (a *AWSLoader) getAWSCredentials(region, roleArn, webIdentityTokenFile, accessKey, secretKey string) (*credentials.Credentials, error) {
	roleArn = a.getOrDefault(roleArn, "AWS_ROLE_ARN")
	webIdentityTokenFile = a.getOrDefault(webIdentityTokenFile, "AWS_WEB_IDENTITY_TOKEN_FILE")

	if roleArn != "" && webIdentityTokenFile != "" {
		return a.getWebIdentityCredentials(region, roleArn, webIdentityTokenFile)
	}

	accessKey = a.getOrDefault(accessKey, "AWS_SECRET_MANAGER_ACCESS_KEY")
	secretKey = a.getOrDefault(secretKey, "AWS_SECRET_MANAGER_SECRET_KEY")

	return a.getStaticCredentials(accessKey, secretKey)
}

func (a *AWSLoader) getWebIdentityCredentials(region, roleArn, tokenFile string) (*credentials.Credentials, error) {
	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf(a.i18n.Get("aws.error_reading_token", map[string]interface{}{"error": err}))
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, fmt.Errorf(a.i18n.Get("aws.session_error", map[string]interface{}{"error": err}))
	}

	stsClient := sts.New(sess)
	output, err := stsClient.AssumeRoleWithWebIdentity(&sts.AssumeRoleWithWebIdentityInput{
		RoleArn:          aws.String(roleArn),
		RoleSessionName:  aws.String("nanogo-session"),
		WebIdentityToken: aws.String(string(token)),
	})
	if err != nil {
		return nil, fmt.Errorf(a.i18n.Get("aws.assume_role_error", map[string]interface{}{"error": err}))
	}

	return credentials.NewStaticCredentials(
		*output.Credentials.AccessKeyId,
		*output.Credentials.SecretAccessKey,
		*output.Credentials.SessionToken,
	), nil
}

func (a *AWSLoader) getStaticCredentials(accessKey, secretKey string) (*credentials.Credentials, error) {
	if accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf(a.i18n.Get("aws.missing_static_credentials", nil))
	}
	return credentials.NewStaticCredentials(accessKey, secretKey, ""), nil
}

func (a *AWSLoader) retrieveAndInjectSecrets(creds *credentials.Credentials, region, secretName string) error {
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
	})
	if err != nil {
		return fmt.Errorf(a.i18n.Get("aws.session_error", map[string]interface{}{"error": err}))
	}

	svc := secretsmanager.New(sess)
	result, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		return fmt.Errorf(a.i18n.Get("aws.secrets_retrieval_error", map[string]interface{}{"error": err}))
	}

	if result.SecretString == nil {
		return fmt.Errorf(a.i18n.Get("aws.empty_secret", map[string]interface{}{"secretName": secretName}))
	}

	return a.injectEnvVariables(*result.SecretString)
}

func (a *AWSLoader) injectEnvVariables(jsonData string) error {
	var data map[string]string
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return fmt.Errorf(a.i18n.Get("aws.json_parse_error", map[string]interface{}{"error": err}))
	}

	for key, value := range data {
		if strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return fmt.Errorf(a.i18n.Get("aws.invalid_env_var", map[string]interface{}{"variable": key}))
		}
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf(a.i18n.Get("aws.env_set_error", map[string]interface{}{"variable": key, "error": err}))
		}
	}

	return nil
}
