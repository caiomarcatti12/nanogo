package i18n_test

import (
	"errors"
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Teste para Factory com erro ao carregar.
func TestFactory_Error(t *testing.T) {
	mockLoader := new(MockLoader)

	mockLoader.On("Load", "./invalid-path").Return(nil, errors.New("load error"))

	service, err := i18n.Factory("./invalid-path")
	assert.NotNil(t, err)
	assert.Nil(t, service)
}

// Teste para Factory com carregamento correto.
func TestFactory_Success(t *testing.T) {
	service, err := i18n.Factory("./translations") // Garanta que esse caminho exista com um YAML válido no teste
	assert.Nil(t, err, "Não deveria retornar erro")
	assert.NotNil(t, service, "Service não deve ser nil")
}
