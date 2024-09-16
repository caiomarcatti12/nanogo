## Feature YamlManager

O módulo `YamlManager` é uma ferramenta desenvolvida para facilitar o trabalho com arquivos YAML em Go. Ele fornece uma API estruturada e simples para carregar, modificar, salvar e gerenciar arquivos YAML. O principal objetivo deste módulo é abstrair a complexidade de manipular estruturas YAML, tornando mais fácil a interação com dados YAML, especialmente para arquivos de configuração ou manifestos do Kubernetes.

### Objetivo

O `YamlManager` foi projetado para ajudar desenvolvedores a:
1. Carregar conteúdo YAML a partir de arquivos.
2. Manipular dados YAML utilizando caminhos dinâmicos.
3. Salvar as alterações de volta em um arquivo ou outro destino.
4. Interagir de forma integrada com o objeto `unstructured.Unstructured` do Kubernetes para gerenciar estruturas de dados YAML.

### Descrição das Funções

- **NewYamlManager()**: Inicializa uma nova instância de `YamlManager` com um mapa de conteúdo vazio.
  
- **LoadYAML(yamlPath string) error**: Carrega um arquivo YAML do caminho fornecido para o mapa `content` interno.

- **LoadFromUnstructured(obj unstructured.Unstructured)**: Carrega dados de um objeto `unstructured.Unstructured` do Kubernetes para o mapa `content`.

- **Save() error**: Salva o conteúdo atual no arquivo YAML de onde ele foi carregado.

- **SaveAs(yamlPath string) error**: Salva o conteúdo atual em um novo arquivo YAML especificado pelo `yamlPath`.

- **SetYAMLValue(path string, newValue interface{}) error**: Define um valor na estrutura YAML por um caminho separado por pontos.

- **RemoveYAMLValue(path string) error**: Remove um valor na estrutura YAML por um caminho separado por pontos.

- **GetYAMLValue(path string, defaultValue interface{}) interface**: Recupera um valor da estrutura YAML por um caminho separado por pontos, retornando `defaultValue` se o caminho não existir.
### Exemplos de Uso

#### 1. Inicializando YamlManager

```go
import "github.com/caiomarcatti12/nanogo/pkg/yaml"

func main() {
    manager := yaml.NewYamlManager()
}
```

#### 2. Carregando um arquivo YAML

```go
import "fmt"

func main() {
    manager := yaml.NewYamlManager()
    
    err := manager.LoadYAML("config.yaml")
    if err != nil {
        fmt.Println("Erro ao carregar YAML:", err)
        return
    }

    fmt.Println("Conteúdo do YAML Carregado:", manager.content)
}
```

#### 3. Carregando de um objeto `unstructured.Unstructured`

```go
import (
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func main() {
    obj := unstructured.Unstructured{
        Object: map[string]interface{}{
            "apiVersion": "v1",
            "kind":       "Pod",
            "metadata": map[string]interface{}{
                "name": "mypod",
            },
        },
    }

    manager := yaml.NewYamlManager()
    manager.LoadFromUnstructured(obj)

    fmt.Println("Conteúdo do YAML Carregado de Unstructured:", manager.content)
}
```

#### 4. Salvando YAML em um arquivo

```go
func main() {
    manager := yaml.NewYamlManager()
    manager.LoadYAML("config.yaml")

    // Modificar o conteúdo, se necessário
    manager.SetYAMLValue("metadata.name", "novoNome")

    err := manager.Save()
    if err != nil {
        fmt.Println("Erro ao salvar YAML:", err)
    }
}
```

#### 5. Salvando YAML em um novo arquivo

```go
func main() {
    manager := yaml.NewYamlManager()
    manager.LoadYAML("config.yaml")

    // Salvar em um novo arquivo
    err := manager.SaveAs("novo_config.yaml")
    if err != nil {
        fmt.Println("Erro ao salvar YAML como novo arquivo:", err)
    }
}
```

#### 6. Modificando um valor YAML usando um caminho

```go
func main() {
    manager := yaml.NewYamlManager()
    manager.LoadYAML("config.yaml")

    // Definir um novo valor
    err := manager.SetYAMLValue("spec.replicas", 3)
    if err != nil {
        fmt.Println("Erro ao definir valor YAML:", err)
    }
}
```

#### 7. Removendo um valor YAML usando um caminho

```go
func main() {
    manager := yaml.NewYamlManager()
    manager.LoadYAML("config.yaml")

    // Remover um valor por caminho
    err := manager.RemoveYAMLValue("metadata.labels.app")
    if err != nil {
        fmt.Println("Erro ao remover valor YAML:", err)
    }
}
```

#### 8. Obtendo um valor YAML com fallback padrão

```go
func main() {
    manager := yaml.NewYamlManager()
    manager.LoadYAML("config.yaml")

    // Recuperar um valor ou usar o valor padrão
    replicas := manager.GetYAMLValue("spec.replicas", 1).(int)
    fmt.Println("Replicas:", replicas)
}
```