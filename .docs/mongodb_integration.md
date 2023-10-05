## Camada de Repository com MongoDB no Framework nanogo

Neste documento, discutiremos como configurar e usar a camada de repository para interagir com o MongoDB usando o framework nanogo. A estrutura do nanogo facilita essa conexão e manipulação de dados com estruturas Go nativas.

### Configuração de Ambiente

Para iniciar, é necessário definir as seguintes variáveis de ambiente:

```plaintext
MONGO_URI=mongodb+srv://user:password@host/database
MONGO_DBNAME=database
```

### Conexão com MongoDB

Para conectar sua aplicação ao MongoDB, você pode usar o seguinte código:

```go
package main

import (
	"github.com/caiomarcatti12/nanogo/v2/config/env"
	"github.com/caiomarcatti12/nanogo/v2/config/mongodb"
)

func main() {
	env.LoadEnv()
	mongodb.ConnectMongoDB()
}
```

### Definindo o Modelo

Defina o modelo da sua entidade implementando a interface `repository.Identifier`. Aqui está um modelo simplificado baseado no exemplo fornecido:

```go
package model

import (
	"github.com/caiomarcatti12/nanogo/v2/config/repository"
	"github.com/google/uuid"
)

var _ repository.Identifier = &MyEntity{}

type MyEntity struct {
	ID       uuid.UUID `bson:"_id,omitempty"`
	Name     string     `bson:"name"`
}

func (e *MyEntity) GetID() uuid.UUID {
	return e.ID
}

func (e *MyEntity) SetID(id uuid.UUID) {
	e.ID = id
}
```

### Criando a Interface do Repositório

Agora, crie uma interface para o repositório que estende `repository.Repository` com o tipo do seu modelo:

```go
package repository

import (
	"github.com/caiomarcatti12/nanogo/v2/config/repository"
	"yourproject/model"
)

type MyEntityRepositoryInterface interface {
	repository.Repository[*model.MyEntity]
}
```

### Implementando o Repositório

Por fim, implemente a camada de repositório encapsulando a comunicação com o MongoDB:

```go
package repository

import (
	"yourproject/model"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/caiomarcatti12/nanogo/v2/config/mongodb"
	"github.com/google/uuid"
)

var _ repository.Repository[*model.MyEntity] = &MyEntityMongoRepository{}

type MyEntityMongoRepository struct {
	collection mongodb.MongoRepository[*model.MyEntity]
}

func NewMyEntityMongoRepository() MyEntityMongoRepository {
	return MyEntityMongoRepository{
		collection: mongodb.NewMongoRepository("MyEntity", &model.MyEntity{}),
	}
}

// Implementando o método Insert da interface Repository
func (r MyEntityMongoRepository) Insert(document *model.MyEntity) (*model.MyEntity, error) {
	return r.collection.Insert(document)
}

// Implementando o método Update da interface Repository
func (r MyEntityMongoRepository) Update(document *model.MyEntity) (bool, error) {
	return r.collection.Update(document)
}

// Implementando o método Delete da interface Repository
func (r MyEntityMongoRepository) Delete(document *model.MyEntity) (bool, error) {
	return r.collection.Delete(document)
}

// Implementando o método DeleteById da interface Repository
func (r MyEntityMongoRepository) DeleteById(id uuid.UUID) (bool, error) {
	return r.collection.DeleteById(id)
}

// Implementando o método FindById da interface Repository
func (r MyEntityMongoRepository) FindById(id uuid.UUID) (*model.MyEntity, error) {
	return r.collection.FindById(id)
}

// Implementando o método FindAll da interface Repository
func (r MyEntityMongoRepository) FindAll() ([]*model.MyEntity, error) {
	return r.collection.FindAll()
}

// Implementando o método RawQueryParseRsql da interface Repository
func (r MyEntityMongoRepository) RawQueryParseRsql(filter rsql.QueryFilter) ([]*model.MyEntity, int64, error) {
	return r.collection.RawQueryParseRsql(filter)
}

```

Neste código:

- A estrutura `MyEntityMongoRepository` está encapsulando a camada de MongoDB.
- Criamos um novo repositório MongoDB para nossa entidade com `mongodb.NewMongoRepository`.
- Todos os métodos necessários da interface `repository.Repository` são implementados, chamando os métodos correspondentes na camada de coleção encapsulada.

Com estas etapas, você terá configurado com sucesso a camada de repositório para se comunicar com o MongoDB através do framework nanogo, e estará pronto para realizar operações CRUD e consultas no banco de dados de maneira organizada e eficiente.

