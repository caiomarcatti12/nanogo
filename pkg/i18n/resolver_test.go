package i18n_test

import (
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultResolver_Resolve(t *testing.T) {
	translations := map[string]map[string]string{
		"pt-br": {"greet": "Olá"},
	}

	resolver := i18n.NewResolver(translations)

	val, ok := resolver.Resolve("pt-br", "greet")
	assert.True(t, ok)
	assert.Equal(t, "Olá", val)

	val, ok = resolver.Resolve("pt-br", "invalid")
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

// Testa método resolver quando chave não existe.
func TestDefaultResolver_Resolve_NotFound(t *testing.T) {
	translations := map[string]map[string]string{
		"pt-br": {"greet": "Olá"},
	}

	resolver := i18n.NewResolver(translations)

	val, ok := resolver.Resolve("pt-br", "farewell")
	assert.False(t, ok)
	assert.Equal(t, "", val)
}

func TestDefaultResolver_Resolve_MultiLevelNotFound(t *testing.T) {
	translations := map[string]map[string]string{
		"pt-br": {"greeting.hello": "Olá"},
	}

	resolver := i18n.NewResolver(translations)

	val, ok := resolver.Resolve("pt-br", "greeting.invalid")
	assert.False(t, ok, "Deve retornar falso ao não encontrar a chave multinível")
	assert.Equal(t, "", val, "Deve retornar string vazia ao não encontrar a chave")
}
