package yaml_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLLoader_Load(t *testing.T) {
	dir, _ := ioutil.TempDir("", "translations")
	defer os.RemoveAll(dir)

	fileContent := "hello: Olá\nbye: Adeus"
	ioutil.WriteFile(filepath.Join(dir, "pt-br.yaml"), []byte(fileContent), 0644)

	loader := NewYAMLLoader()
	translations, err := loader.Load(dir)

	assert.Nil(t, err)
	assert.Equal(t, "Olá", translations["pt-br"]["hello"])
	assert.Equal(t, "Adeus", translations["pt-br"]["bye"])
}
