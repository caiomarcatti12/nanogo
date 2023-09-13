##  JWTManager

O pacote `jwt` fornece uma implementação para gerenciar tokens JWT (JSON Web Tokens) em Go. Ele facilita a criação e validação de tokens JWT, encapsulando as funcionalidades necessárias em uma estrutura `JWTManager`.


### Métodos

#### NewJWTManager

```go
func NewJWTManager(signingKey string) *JWTManager
```

- **Parâmetros**:
    - `signingKey` (string): A chave de assinatura utilizada para criar e validar tokens JWT.
- **Retorno**:
    - `*JWTManager`: Um ponteiro para uma nova instância de `JWTManager`.
- **Descrição**: Este método inicializa e retorna uma nova instância de `JWTManager` utilizando a chave de assinatura fornecida.

#### GenerateToken

```go
func (manager *JWTManager) GenerateToken(expirationTime time.Duration, data map[string]interface{}) (string, error)
```

- **Parâmetros**:
    - `expirationTime` (`time.Duration`): A duração até que o token expire.
    - `data` (`map[string]interface{}`): Um mapa de claims que serão incluídas no token.
- **Retorno**:
    - `(string, error)`: O token JWT como uma string ou um erro se a criação do token falhar.
- **Descrição**: Este método cria um novo token JWT com o tempo de expiração e as claims fornecidas, e então o assina com a chave de assinatura do `JWTManager`.

#### ValidateToken

```go
func (manager *JWTManager) ValidateToken(tokenString string) (jwt.Claims, error)
```

- **Parâmetros**:
    - `tokenString` (string): A string do token JWT que deve ser validada.
- **Retorno**:
    - `(jwt.Claims, error)`: As claims do token JWT ou um erro se a validação falhar.
- **Descrição**: Este método valida o token JWT fornecido verificando sua assinatura e tempo de expiração, e retorna suas claims se a validação for bem-sucedida.

#### DecodeToken

```go
func (manager *JWTManager) DecodeToken(tokenString string) (map[string]interface{}, error)
```

- **Parâmetros**:
  - `tokenString` (string): A string do token JWT que deve ser decodificada.
- **Retorno**:
  - `(map[string]interface{}, error)`: Um mapa contendo as claims do token JWT ou um erro se a decodificação falhar.
- **Descrição**: Este método decodifica o token JWT fornecido e retorna suas claims sem validar a assinatura ou o tempo de expiração do token.


### Exemplos de Uso

#### Criando um Novo JWTManager

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

#### Validando um Token

```go
claims, err := jwtManager.ValidateToken(token)
if err != nil {
	fmt.Println("Token inválido:", err)
} else {
	fmt.Println("Token válido. Claims:", claims)
}
```

#### Decodificando um Token

```go
data, err := jwtManager.DecodeToken(tokenString)
if err != nil {
	fmt.Println("Erro ao decodificar o token:", err)
} else {
	fmt.Println("Dados do token:", data)
}
```

### Considerações de Segurança

Ao trabalhar com JWTs, é crucial manter a chave de assinatura segura, pois qualquer pessoa com acesso a ela pode criar tokens JWT válidos. É recomendado armazenar a chave de assinatura em um local seguro, como uma variável de ambiente ou um serviço de gerenciamento de segredos.

