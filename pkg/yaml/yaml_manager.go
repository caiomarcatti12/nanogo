package yaml

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// YamlManager estrutura que implementa a interface IYamlManager
type YamlManager struct {
	content  map[string]interface{}
	yamlPath string
}

func NewYamlManager() YamlManager {
	return YamlManager{
		content:  make(map[string]interface{}),
		yamlPath: "",
	}
}

// loadYAML loads a YAML file and returns an unstructured object
func (y *YamlManager) LoadYAML(yamlPath string) error {
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return err
	}
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)

	y.content = result
	y.yamlPath = yamlPath

	return nil
}

// LoadFromUnstructured loads data from an unstructured.Unstructured object into YamlManager
func (y *YamlManager) LoadFromUnstructured(obj unstructured.Unstructured) {
	y.content = obj.Object
}

// saveYAML saves an unstructured object to a YAML file
func (y *YamlManager) Save() error {
	yamlData, err := yaml.Marshal(y.content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(y.yamlPath, yamlData, 0644)
}

// saveYAML saves an unstructured object to a YAML file
func (y *YamlManager) SaveAs(yamlPath string) error {
	// Abre o arquivo para escrita
	file, err := os.Create(yamlPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Cria um novo encoder YAML e define a indentação
	yamlEncoder := yaml.NewEncoder(file)
	yamlEncoder.SetIndent(2)

	// Codifica os dados YAML e escreve no arquivo
	if err := yamlEncoder.Encode(y.content); err != nil {
		return err
	}

	return nil
}

func (y *YamlManager) SetYAMLValue(path string, newValue interface{}) error {
	fields := strings.Split(path, ".")
	var current interface{}

	current = y.content

	for i, key := range fields {
		if index, err := strconv.Atoi(key); err == nil {
			if list, ok := current.([]interface{}); ok {
				if index >= 0 && index < len(list) {
					current = list[index]
				} else if i < (len(fields) - 1) {
					zero := make([]interface{}, 0)
					setIndex(&list, index, zero)
					current = &list
				} else if (len(fields) - 1) == i {
					setIndex(&list, index, newValue)
					current = &list
				}
			} else {
				return nil
			}
		} else {
			if valCurrent, ok := current.(map[string]interface{}); ok {
				if val, ok := valCurrent[key]; ok {
					if i == len(fields)-1 {
						valCurrent[key] = newValue
					}
					if nextMap, ok := val.(map[string]interface{}); ok {
						current = nextMap
					} else if nextList, ok := val.([]interface{}); ok {
						current = nextList
					} else {
						return nil
					}
				} else if i < (len(fields) - 1) {
					valCurrent[key] = make(map[string]interface{})
					current = valCurrent[key]
				} else if (len(fields) - 1) == i {
					valCurrent[key] = newValue
				}
			} else {
				return nil
			}
		}
	}

	return nil
}

// removeYAMLValue removes a value in the YAML at the given path
func (y *YamlManager) RemoveYAMLValue(path string) error {
	fields := strings.Split(path, ".")
	var current interface{}

	current = y.content

	for i, key := range fields {
		if index, err := strconv.Atoi(key); err == nil {
			if list, ok := current.([]interface{}); ok {
				if index >= 0 && index < len(list) {
					current = list[index]
				} else {
					return nil
				}
			} else {
				return nil
			}
		} else {
			if valCurrent, ok := current.(map[string]interface{}); ok {
				if val, ok := valCurrent[key]; ok {
					if i == len(fields)-1 {
						delete(valCurrent, fields[len(fields)-1])
					}
					if nextMap, ok := val.(map[string]interface{}); ok {
						current = nextMap
					} else if nextList, ok := val.([]interface{}); ok {
						current = nextList
					} else {
						return nil
					}
				}
			} else {
				return nil
			}
		}
	}

	return nil
}

// getYAMLValue gets a value from the YAML at the given path
func (y *YamlManager) GetYAMLValue(path string, defaultValue interface{}) interface{} {
	fields := strings.Split(path, ".")
	var current interface{}

	current = y.content

	for i, key := range fields {
		if index, err := strconv.Atoi(key); err == nil {
			if list, ok := current.([]interface{}); ok {
				if i == len(fields)-1 {
					return list[index]
				} else if index >= 0 && index < len(list) {
					current = list[index]
				} else {
					return defaultValue
				}
			} else {
				return defaultValue
			}
		} else {
			if valCurrent, ok := current.(map[string]interface{}); ok {
				if val, ok := valCurrent[key]; ok {
					if i == len(fields)-1 {
						return val
					}
					if nextMap, ok := val.(map[string]interface{}); ok {
						current = nextMap
					} else if nextList, ok := val.([]interface{}); ok {
						current = nextList
					} else {
						return defaultValue
					}
				}
			} else {
				return defaultValue
			}
		}
	}
	return defaultValue
}

func setIndex(slice *[]interface{}, index int, value interface{}) {
	// Verifica se o índice é maior que o comprimento do slice
	if index >= len(*slice) {
		// Cria um novo slice com capacidade suficiente
		newSlice := make([]interface{}, index+1)
		// Copia os elementos do slice original para o novo slice
		copy(newSlice, *slice)

		newSlice[index] = value
		// Atualiza o slice original para o novo slice
		*slice = newSlice
	}
}
