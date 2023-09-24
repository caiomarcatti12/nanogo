#  JWTManager

O pacote `jwt` fornece uma implementação para gerenciar tokens JWT (JSON Web Tokens) em Go. Ele facilita a criação e validação de tokens JWT, encapsulando as funcionalidades necessárias em uma estrutura `JWTManager`.

Claro! Aqui está a seção de configuração adaptada para o `JWTManager`:

### Configuração

Para configurar o `JWTManager`, você pode optar por definir a seguinte variável de ambiente, embora seja opcional:

```sh
JWT_SECRET=my-secret-key
```

- `JWT_SECRET`: A chave de assinatura utilizada para criar e validar tokens JWT. Se você não configurar essa variável de ambiente, deverá fornecer a chave explicitamente ao criar uma nova instância de `JWTManager`. A ausência da chave em ambos os lugares pode resultar em erro.

#### Métodos

##### NewJWTManager

```go
func NewJWTManager(signingKey ...string) *JWTManager
```
Este método inicializa e retorna uma nova instância de `JWTManager` utilizando a chave de assinatura fornecida.

- **Parâmetros**:
    - `signingKey (string, opcional): A chave de assinatura utilizada para criar e validar tokens JWT. Se não for fornecida, o método tentará obter a variável de ambiente JTW_SECRET.
- **Retorno**:
    - `*JWTManager`: Um ponteiro para uma nova instância de `JWTManager`.

##### GenerateToken

```go
func (manager *JWTManager) GenerateToken(expirationTime time.Duration, data map[string]interface{}) (string, error)

```
Este método cria um novo token JWT com o tempo de expiração e as claims fornecidas, e então o assina com a chave de assinatura do `JWTManager`.
- **Parâmetros**:
    - `expirationTime` (`time.Duration`): A duração até que o token expire.
    - `data` (`map[string]interface{}`): Um mapa de claims que serão incluídas no token.
- **Retorno**:
    - `(string, error)`: O token JWT como uma string ou um erro se a criação do token falhar.

#### ValidateToken

```go
func (manager *JWTManager) ValidateToken(tokenString string) (jwt.Claims, error)
```
Este método valida o token JWT fornecido verificando sua assinatura e tempo de expiração, e retorna suas claims se a validação for bem-sucedida.

- **Parâmetros**:
    - `tokenString` (string): A string do token JWT que deve ser validada.
- **Retorno**:
    - `(jwt.Claims, error)`: As claims do token JWT ou um erro se a validação falhar.

##### DecodeToken

```go
func (manager *JWTManager) DecodeToken(tokenString string) (map[string]interface{}, error)
```
Este método decodifica o token JWT fornecido e retorna suas claims sem validar a assinatura ou o tempo de expiração do token.

- **Parâmetros**:
  - `tokenString` (string): A string do token JWT que deve ser decodificada.
- **Retorno**:
  - `(map[string]interface{}, error)`: Um mapa contendo as claims do token JWT ou um erro se a decodificação falhar.


### Exemplos de Uso

##### Criando um Novo JWTManager

```go
jwtManager := NewJWTManager("my-secret-key")
```

#### Gerando um Token

```go
token, err := jwtManager.GenerateToken(time.Hour, map[string]interface{}{"userID": 123})
if err != nil {
	fmt.Println("Erro ao gerar token:", err)
}
fmt.Println("Token gerado:", token)
```

##### Validando um Token

```go
claims, err := jwtManager.ValidateToken(token)
if err != nil {
	fmt.Println("Token inválido:", err)
} else {
	fmt.Println("Token válido. Claims:", claims)
}
```

##### Decodificando um Token

```go
data, err := jwtManager.DecodeToken(tokenString)
if err != nil {
	fmt.Println("Erro ao decodificar o token:", err)
} else {
	fmt.Println("Dados do token:", data)
}
```

##### Considerações de Segurança

Ao trabalhar com JWTs, é crucial manter a chave de assinatura segura, pois qualquer pessoa com acesso a ela pode criar tokens JWT válidos. É recomendado armazenar a chave de assinatura em um local seguro, como uma variável de ambiente ou um serviço de gerenciamento de segredos.

