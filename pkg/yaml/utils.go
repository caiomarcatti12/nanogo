package yaml

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

// Função genérica para fazer o unmarshal de um arquivo YAML em um tipo fornecido
func UnmarshalYAMLFile[T any](fileName string, output *T) error {
	// Verifica se o arquivo existe
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return fmt.Errorf("arquivo não encontrado: %s", fileName)
	}

	// Lê o conteúdo do arquivo YAML
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("erro ao ler o arquivo: %v", err)
	}

	// Faz o unmarshal do conteúdo YAML para o tipo genérico
	err = yaml.Unmarshal(data, output)
	if err != nil {
		return fmt.Errorf("erro ao fazer unmarshal do YAML: %v", err)
	}

	return nil
}
