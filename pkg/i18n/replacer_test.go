package i18n_test

import (
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewReplacer verifica a criação correta da instância de DefaultReplacer.
func TestNewReplacer(t *testing.T) {
	replacer := i18n.NewReplacer()
	assert.NotNil(t, replacer, "o replacer não deve ser nulo após criação")
}

// TestReplace verifica se as variáveis são substituídas corretamente.
func TestDefaultReplacer_Replace_Success(t *testing.T) {
	replacer := i18n.NewReplacer()

	template := "Olá, {{name}}! Bem-vindo ao {{place}}."
	vars := map[string]interface{}{
		"name":  "Caio",
		"place": "Brasil",
	}

	expected := "Olá, Caio! Bem-vindo ao Brasil."
	result := replacer.Replace(template, vars)

	assert.Equal(t, expected, result, "a substituição das variáveis deve estar correta")
}

// TestReplace_MissingVars verifica o comportamento quando variáveis faltam no mapa.
func TestDefaultReplacer_Replace_MissingVars(t *testing.T) {
	replacer := i18n.NewReplacer()

	template := "Olá, {{name}}! Seu código é {{code}}."
	vars := map[string]interface{}{
		"name": "João",
	}

	expected := "Olá, João! Seu código é {{code}}."
	result := replacer.Replace(template, vars)

	assert.Equal(t, expected, result, "variáveis não fornecidas devem permanecer intactas")
}

// TestReplace_NoVars verifica comportamento quando nenhum placeholder existe.
func TestDefaultReplacer_Replace_NoVars(t *testing.T) {
	replacer := i18n.NewReplacer()

	template := "Bem-vindo ao sistema."
	vars := map[string]interface{}{
		"name": "João",
	}

	expected := "Bem-vindo ao sistema."
	result := replacer.Replace(template, vars)

	assert.Equal(t, expected, result, "string sem placeholders deve permanecer intacta")
}

// TestReplace_EmptyVars verifica comportamento com o mapa de variáveis vazio.
func TestDefaultReplacer_Replace_EmptyVars(t *testing.T) {
	replacer := i18n.NewReplacer()

	template := "Olá, {{name}}!"
	vars := map[string]interface{}{}

	expected := "Olá, {{name}}!"
	result := replacer.Replace(template, vars)

	assert.Equal(t, expected, result, "quando o mapa de variáveis estiver vazio, o template deve permanecer igual")
}

func TestDefaultResolver_Resolve_KeyNotFound(t *testing.T) {
	translations := map[string]map[string]string{
		"pt-br": {"hello": "Olá"},
	}

	resolver := i18n.NewResolver(translations)

	// Testando chave não existente
	val, ok := resolver.Resolve("pt-br", "invalid_key")

	assert.False(t, ok, "Deve retornar false para chave não encontrada")
	assert.Equal(t, "", val, "Valor deve ser vazio quando chave não encontrada")
}
