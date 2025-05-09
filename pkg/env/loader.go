package env

import (
	"errors"
	"os"

	"github.com/caiomarcatti12/nanogo/pkg/i18n"
)

// Loader carrega as variáveis de ambiente usando o provedor especificado em ENV_PROVIDER.
// Caso nenhum provedor seja especificado, o padrão é o arquivo .env (FileEnvLoader).
func Loader(i18n i18n.I18N) error {
	provider := os.Getenv("ENV_PROVIDER")

	if provider == "" {
		provider = "ENV_FILE"
	}

	switch provider {
	case "AWS":
		return loadAWSSecrets(i18n)
	case "ENV_FILE":
		return loadFileEnv(i18n)
	default:
		return errors.New(i18n.Get("env.provider_not_found", map[string]interface{}{"provider": provider}))
	}
}

// loadAWSSecrets inicializa e executa o loader para AWS Secrets Manager
func loadAWSSecrets(i18n i18n.I18N) error {
	awsLoader := NewAWSLoader(i18n)
	return awsLoader.LoadSecrets("", "", "", "", "", "")
}

// loadFileEnv inicializa e executa o loader para arquivos .env
func loadFileEnv(i18n i18n.I18N) error {
	fileEnvLoader := NewFileEnvLoader(i18n)
	return fileEnvLoader.Load()
}
