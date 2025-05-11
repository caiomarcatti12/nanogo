package env

import (
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

// LoadVaultHashcorp inicializa o cliente Vault, obtém o segredo e define variáveis de ambiente.
func loadVaultHashcorp(i18n i18n.I18N) error {
	// Validação das variáveis de ambiente necessárias
	vaultAddr := os.Getenv("VAULT_HASHICORP_HOST")
	if vaultAddr == "" {
		log.Fatal(i18n.Get("vault.addr_not_defined", nil))
	}

	vaultToken := os.Getenv("VAULT_HASHICORP_TOKEN")
	if vaultToken == "" {
		log.Fatal(i18n.Get("vault.token_not_defined", nil))
	}

	secretPath := os.Getenv("VAULT_HASHICORP_SECRET_PATH")
	if secretPath == "" {
		log.Fatal(i18n.Get("vault.secret_path_not_defined", nil))
	}

	// Configura o cliente Vault com o endereço obtido das variáveis de ambiente
	config := api.DefaultConfig()
	config.Address = vaultAddr

	// Inicializa o cliente Vault
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatal(i18n.Get("vault.client_creation_error", map[string]interface{}{"error": err}))
	}

	// Define o token de autenticação do Vault
	client.SetToken(vaultToken)

	// Obtém o segredo do Vault no caminho especificado
	secret, err := client.Logical().Read(secretPath)
	if err != nil {
		log.Fatal(i18n.Get("vault.secret_read_error", map[string]interface{}{"error": err}))
	}

	if secret == nil || secret.Data == nil {
		log.Fatal(i18n.Get("vault.secret_not_found", map[string]interface{}{"path": secretPath}))
	}

	// Trata os dados especificamente para KV versão 2 (com campo "data")
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatal(i18n.Get("vault.secret_invalid_format", map[string]interface{}{"path": secretPath}))
	}

	// Percorre os valores do segredo e os define como variáveis de ambiente
	for key, value := range data {
		strValue, ok := value.(string)
		if !ok {
			log.Fatal(i18n.Get("vault.secret_value_not_string", map[string]interface{}{"key": key}))
		}
		if err := os.Setenv(key, strValue); err != nil {
			log.Fatal(i18n.Get("vault.env_set_error", map[string]interface{}{"key": key, "error": err}))
		}
	}

	log.Println(i18n.Get("vault.load_success", map[string]interface{}{"path": secretPath}))
	return nil
}
