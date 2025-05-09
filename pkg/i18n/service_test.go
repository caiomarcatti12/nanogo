package i18n_test

import (
	"errors"
	"github.com/caiomarcatti12/nanogo/pkg/i18n"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Loader
type MockLoader struct {
	mock.Mock
}

func (m *MockLoader) Load(path string) (map[string]map[string]string, error) {
	args := m.Called(path)
	return args.Get(0).(map[string]map[string]string), args.Error(1)
}

// Mock Resolver
type MockResolver struct {
	mock.Mock
}

func (m *MockResolver) Resolve(locale, key string) (string, bool) {
	args := m.Called(locale, key)
	return args.String(0), args.Bool(1)
}

// Mock Replacer
type MockReplacer struct {
	mock.Mock
}

func (m *MockReplacer) Replace(text string, vars map[string]interface{}) string {
	args := m.Called(text, vars)
	return args.String(0)
}

// Teste criação do serviço
func TestNewService(t *testing.T) {
	service := i18n.NewService(new(MockLoader), new(MockResolver), new(MockReplacer))
	assert.NotNil(t, service)
}

// Teste SetLanguage com idioma válido carregado
func TestService_SetLanguage_ValidLanguage(t *testing.T) {
	mockLoader := new(MockLoader)
	mockLoader.On("Load", mock.Anything).Return(map[string]map[string]string{"en": {}}, nil)

	service := i18n.NewService(mockLoader, nil, nil)
	service.LoadTranslations("dummy/path")

	service.SetLanguage("en")
	assert.Equal(t, "en", service.GetLanguage())
}

// Teste SetLanguage com idioma inválido (deve definir padrão)
func TestService_SetLanguage_InvalidLanguage(t *testing.T) {
	mockLoader := new(MockLoader)
	mockLoader.On("Load", mock.Anything).Return(map[string]map[string]string{"pt-br": {}}, nil)

	service := i18n.NewService(mockLoader, nil, nil)
	service.LoadTranslations("dummy/path")

	service.SetLanguage("invalid-lang")
	assert.Equal(t, "pt-br", service.GetLanguage())
}

// Teste GetLanguage
func TestService_GetLanguage(t *testing.T) {
	service := i18n.NewService(nil, nil, nil)
	service.SetLanguage("pt-br")
	assert.Equal(t, "pt-br", service.GetLanguage())
}

// Teste GetDefaultLanguage
func TestService_GetDefaultLanguage(t *testing.T) {
	service := i18n.NewService(nil, nil, nil)
	assert.Equal(t, "pt-br", service.GetDefaultLanguage())
}

// Teste LoadTranslations com sucesso
func TestService_LoadTranslations_Success(t *testing.T) {
	mockLoader := new(MockLoader)
	translations := map[string]map[string]string{"pt-br": {"greet": "Olá"}}

	mockLoader.On("Load", "path/to/translations").Return(translations, nil)

	service := i18n.NewService(mockLoader, nil, nil)
	err := service.LoadTranslations("path/to/translations")

	assert.Nil(t, err)
	mockLoader.AssertExpectations(t)
}

// Teste LoadTranslations com falha
func TestService_LoadTranslations_Error(t *testing.T) {
	mockLoader := new(MockLoader)

	// Correção aqui: retorno deve ser um map válido (mesmo que vazio)
	mockLoader.On("Load", "invalid/path").Return(map[string]map[string]string{}, errors.New("load error"))

	service := i18n.NewService(mockLoader, nil, nil)
	err := service.LoadTranslations("invalid/path")

	assert.NotNil(t, err)
	mockLoader.AssertExpectations(t)
}

// Teste Get tradução encontrada
func TestService_Get_TranslationFound(t *testing.T) {
	mockLoader := new(MockLoader)
	mockResolver := new(MockResolver)
	mockReplacer := new(MockReplacer)

	translations := map[string]map[string]string{
		"pt-br": {"hello": "Olá, {{name}}!"},
	}

	mockLoader.On("Load", mock.Anything).Return(translations, nil)

	service := i18n.NewService(mockLoader, mockResolver, mockReplacer)
	service.LoadTranslations("dummy/path")
	service.SetLanguage("pt-br")

	mockResolver.On("Resolve", "pt-br", "hello").Return("Olá, {{name}}!", true)
	mockReplacer.On("Replace", "Olá, {{name}}!", mock.Anything).Return("Olá, João!")

	result := service.Get("hello", map[string]interface{}{"name": "João"})

	assert.Equal(t, "Olá, João!", result)
	mockResolver.AssertExpectations(t)
	mockReplacer.AssertExpectations(t)
}

// Teste Get tradução não encontrada, retorna chave original
func TestService_Get_TranslationNotFound(t *testing.T) {
	mockLoader := new(MockLoader)
	mockResolver := new(MockResolver)

	translations := map[string]map[string]string{"pt-br": {}}

	mockLoader.On("Load", mock.Anything).Return(translations, nil)

	service := i18n.NewService(mockLoader, mockResolver, nil)
	service.LoadTranslations("dummy/path")
	service.SetLanguage("pt-br")

	mockResolver.On("Resolve", "pt-br", "unknown").Return("", false)
	mockResolver.On("Resolve", "pt-br", "unknown").Return("", false)

	result := service.Get("unknown", nil)

	assert.Equal(t, "unknown", result)
	mockResolver.AssertExpectations(t)
}

// Teste Get com fallback para idioma padrão
func TestService_Get_FallbackToDefault(t *testing.T) {
	mockLoader := new(MockLoader)
	mockResolver := new(MockResolver)
	mockReplacer := new(MockReplacer)

	translations := map[string]map[string]string{
		"pt-br": {"hello": "Olá, {{name}}!"},
		"en":    {},
	}

	mockLoader.On("Load", mock.Anything).Return(translations, nil)

	service := i18n.NewService(mockLoader, mockResolver, mockReplacer)
	service.LoadTranslations("dummy/path")
	service.SetLanguage("en")

	mockResolver.On("Resolve", "en", "hello").Return("", false)
	mockResolver.On("Resolve", "pt-br", "hello").Return("Olá, {{name}}!", true)
	mockReplacer.On("Replace", "Olá, {{name}}!", mock.Anything).Return("Olá, Caio!")

	result := service.Get("hello", map[string]interface{}{"name": "Caio"})

	assert.Equal(t, "Olá, Caio!", result)
	mockResolver.AssertExpectations(t)
	mockReplacer.AssertExpectations(t)
}

func TestService_Get_NoVarsProvided(t *testing.T) {
	mockResolver := new(MockResolver)
	mockReplacer := new(MockReplacer)

	mockResolver.On("Resolve", "pt-br", "hello").Return("Olá, Mundo!", true)
	mockReplacer.On("Replace", "Olá, Mundo!", mock.Anything).Return("Olá, Mundo!")

	service := i18n.NewService(nil, mockResolver, mockReplacer)
	service.SetLanguage("pt-br")

	result := service.Get("hello") // Nenhuma variável fornecida
	assert.Equal(t, "Olá, Mundo!", result)

	mockResolver.AssertExpectations(t)
	mockReplacer.AssertExpectations(t)
}
