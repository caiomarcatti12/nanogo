package yaml

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type YAMLLoader struct{}

func NewYAMLLoader() *YAMLLoader {
	return &YAMLLoader{}
}

func (l *YAMLLoader) Load(path string) (map[string]map[string]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	translations := make(map[string]map[string]string)

	for _, file := range files {
		if ext := filepath.Ext(file.Name()); ext == ".yaml" || ext == ".yml" {
			content, err := ioutil.ReadFile(filepath.Join(path, file.Name()))
			if err != nil {
				return nil, err
			}

			var rawData map[string]interface{}
			if err := yaml.Unmarshal(content, &rawData); err != nil {
				return nil, err
			}

			flatMap := flatten(rawData, "", make(map[string]string))
			locale := strings.TrimSuffix(file.Name(), ext)
			translations[locale] = flatMap
		}
	}

	return translations, nil
}

// flatten converte mapas aninhados em um mapa plano com notação ponto
func flatten(data map[string]interface{}, prefix string, result map[string]string) map[string]string {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch typedValue := value.(type) {
		case map[interface{}]interface{}:
			convertedMap := make(map[string]interface{})
			for k, v := range typedValue {
				convertedMap[k.(string)] = v
			}
			flatten(convertedMap, fullKey, result)
		case map[string]interface{}:
			flatten(typedValue, fullKey, result)
		case string:
			result[fullKey] = typedValue
		default:
			result[fullKey] = ""
		}
	}
	return result
}
