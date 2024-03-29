### **Dependências de Projetos Externos**

O **Nanogo** foi construído sobre a solidez e eficiência de várias bibliotecas e ferramentas amplamente reconhecidas na comunidade Golang. Aqui está uma lista das principais dependências externas que potencializam o Nanogo:

#### **Dependências Diretas**

- **[Google UUID](https://github.com/google/uuid) (v1.3.0)**
  - Utilizada para a geração de UUIDs únicos, garantindo identificadores robustos e seguros para os seus recursos.

- **[Gorilla Mux](https://github.com/gorilla/mux) (v1.8.0)**
  - Um potente roteador HTTP e dispatcher para Go, facilitando a criação de APIs RESTful.

- **[GoDotEnv](https://github.com/joho/godotenv) (v1.5.1)**
  - Permite a carga de variáveis de ambiente a partir de um arquivo `.env`, facilitando a gestão de configurações.

- **[Logrus](https://github.com/sirupsen/logrus) (v1.9.3)**
  - Uma biblioteca de logging estruturada para Go, ajudando a manter registros detalhados do comportamento do aplicativo.

- **[Mongo Driver](https://go.mongodb.org/mongo-driver) (v1.12.0)**
  - Facilita a integração com bancos de dados MongoDB, proporcionando uma maneira eficiente de trabalhar com NoSQL.

- **[Go Redis](https://github.com/go-redis/redis) (v6.15.9)**
  - Um cliente Redis para Go, permitindo integrações rápidas e confiáveis com sistemas Redis.

- **[AMQP Streadway](https://github.com/streadway/amqp) (v1.1.0)**
  - Uma implementação do protocolo AMQP para Go, facilitando a comunicação com servidores RabbitMQ.

- **[DGrijalva JWT-Go](https://github.com/dgrijalva/jwt-go) (v3.2.0)**
  - Uma biblioteca poderosa para criar e validar JSON Web Tokens (JWT) no Go.

- **[Swaggo Swag](https://github.com/swaggo/swag) (v1.16.2)**
  - Uma biblioteca para automaticamente gerar documentação de API Swagger 2.0 para Go. Ela analisa os comentários do código e gera a documentação necessária, tornando a manutenção e apresentação de documentações de API mais simples e integrada.

- **[Prometheus Client Golang](https://github.com/prometheus/client_golang) (v1.11.0)**
  - Esta biblioteca fornece funcionalidades para a exportação de métricas para o [Prometheus](https://prometheus.io/), um sistema de monitoramento e alerta de código aberto. Com este cliente, o **Nanogo** pode expor métricas de desempenho internas, tais como contadores, gauges e histograms, que podem ser coletados pelo servidor Prometheus para visualização, alerta e análise de tendências.


