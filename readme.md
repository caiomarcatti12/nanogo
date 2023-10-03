### Nanogo: Framework de Desenvolvimento Golang

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](license) ![Static Badge](https://img.shields.io/badge/N%C3%A3o%20pronto%20para%20produ%C3%A7%C3%A3o-red)

Bem-vindo ao universo do **Nanogo**, o framework Golang concebido para simplificar e agilizar o ciclo de desenvolvimento do seu software sem sacrificar a robustez e a flexibilidade. Inspirado pelos princípios da Arquitetura Limpa (Clean Architecture) e Arquitetura Hexagonal, o Nanogo encapsula todas as funcionalidades essenciais de infraestrutura, permitindo aos desenvolvedores focar no que realmente importa: o domínio do software.

#### **Características Principais**

1. **Integração Simplificada com Sistemas Essenciais**
    - **Redis**: Conexão facilitada com sistemas Redis, favorecendo um armazenamento de dados ágil e eficiente.
    - **Mongo**: Implemente bancos de dados MongoDB com facilidade, tirando proveito de uma das soluções de banco de dados NoSQL mais populares e eficazes disponíveis atualmente.
    - **RabbitMQ**: Integre-se com RabbitMQ para facilitar a comunicação entre diferentes serviços ou componentes do seu aplicativo.

2. **Gestão de Configurações Facilitada**
    - **Carregamento de Arquivos `.env`**: Gerencie configurações e secretas com eficiência, carregando-as facilmente através de arquivos `.env`.

3. **API WebServer**
    - **Desenvolvimento Ágil de APIs**: Acelere o desenvolvimento de APIs robustas e seguras com a API WebServer integrada do Nanogo.

4. **Compliance com Arquiteturas Modernas**
    - **Clean Architecture**: O Nanogo segue o princípio da Clean Architecture, promovendo a separação de preocupações e facilitando a manutenção do código.
    - **Arquitetura Hexagonal**: Adote uma arquitetura hexagonal para garantir a flexibilidade e a intercambialidade dos componentes do seu software, facilitando testes e integrações.

#### **Benefícios**

- **Redução Significativa de Código**: O Nanogo foi projetado para minimizar a quantidade de código necessária para configurar e manter sua aplicação, tornando o processo de desenvolvimento mais rápido e menos propenso a erros.
- **Implementação Simplificada**: Com foco na simplicidade e na eficiência, o Nanogo facilita a configuração e a implementação de sistemas complexos, permitindo que os desenvolvedores criem soluções robustas sem o aborrecimento associado à gestão de infraestruturas intricadas.
- **Adaptação Rápida**: Graças à sua estrutura intuitiva e bem organizada, os desenvolvedores podem aprender e começar a usar o Nanogo em pouco tempo, sem uma curva de aprendizado íngreme.

**Nanogo** surge como uma ferramenta indispensável para desenvolvedores Golang modernos, promovendo uma experiência de desenvolvimento mais fluida, rápida e menos complexa, sem comprometer a qualidade e a funcionalidade final do software.

### **Dependências Externas 🛠**

O Nanogo é potencializado por várias bibliotecas e ferramentas robustas da comunidade Golang. Para garantir um desenvolvimento fluido e funcionalidades avançadas, incorporamos uma série de dependências externas no projeto.

Para uma descrição detalhada de cada dependência e como elas são usadas no framework Nanogo, consulte o nosso guia detalhado de dependências:

📄 [Consulte o Guia de Dependências](./.docs/dependencies.md)

### **Arquitetura 🔥**

- **[Arquitetura Limpa (Clean Architecture)](./.docs/clean_architecture.md)**
    - Mantenha seu código organizado e fácil de manter aderindo aos princípios da Arquitetura Limpa.

- **[Arquitetura Hexagonal](./.docs/hexagonal_architecture.md)**
    - Facilite a manutenção e o teste de seus aplicativos através da implementação da Arquitetura Hexagonal.
      Claro, podemos adicionar um novo tópico à sua documentação README.md principal que destaca a "Arquitetura de Repository". Aqui está uma sugestão:

- **[Arquitetura de Repository](./.docs/repository_architecture.md)**
    - Descubra como a estrutura de repository no nanogo facilita a abstração do acesso ao banco de dados, promovendo um código mais limpo, flexível e testável. Explore como o princípio SOLID de inversão de dependência é central para esta arquitetura, permitindo um design de software robusto e de fácil manutenção.


### **Instalando a Biblioteca Nanogo**

**Adicione a Biblioteca como Dependência**: Abra o terminal no diretório do seu projeto e execute o comando a seguir para adicionar a biblioteca Nanogo como uma dependência:

```
go get github.com/caiomarcatti12/nanogo/v2
```

### **Funcionalidades 🔥**

O framework **Nanogo** oferece um conjunto robusto de funcionalidades projetadas para facilitar e acelerar o desenvolvimento de software. Abaixo, você encontrará uma lista de funcionalidades-chave juntamente com links para suas respectivas documentações:

- **[API Webserver](./.docs/api_webserver.md)**
   - Construa APIs poderosas e escaláveis com nossa funcionalidade de servidor web integrado.

- **[Conexão com Redis](./.docs/redis_cache.md)**
   - Facilite a integração e a manipulação de bancos de dados Redis em seus projetos.

- **[Integração com MongoDB](./.docs/mongodb_integration.md)**
   - Implemente soluções de banco de dados NoSQL rapidamente com nossa integração nativa com MongoDB.

- **[Suporte para RabbitMQ](./.docs/rabbitmq_support.md)**
   - Orquestre microserviços eficientemente com nosso suporte integrado para RabbitMQ.

- **[Carregamento de Arquivos .env](./.docs/local_env_loading.md)**
   - Gerencie configurações de aplicativos com facilidade através do suporte para carregamento de arquivos .env.

- **[Carregamento dinamico de Arquivos .env](./.docs/remote_env_loading.md)**
    - Gerencie configurações de aplicativos dinamicamente através do suporte para carregamento por api

- **[JWT Manager](./.docs/jwt.md)**
    - Crie e valide tokens JWT de maneira fácil e segura com o nosso JWT Manager.

- **[Logger Integrado](./.docs/logger.md)**
    - Monitore e rastreie as operações e eventos de seu aplicativo com eficiência e precisão usando nosso logger integrado. Com suporte a múltiplos níveis de log e Correlation ID.

Cada funcionalidade foi meticulosamente desenvolvida para fornecer a melhor experiência possível aos desenvolvedores, poupando tempo e esforço ao criar softwares incríveis com **Nanogo**.

### **Como Contribuir**

Estamos sempre de portas abertas para novas contribuições! Se você deseja auxiliar no crescimento e aprimoramento do projeto — seja através da correção de bugs, propostas de melhorias ou incorporação de novas funcionalidades — sua ajuda será sempre bem-vinda.

Acesse nosso [Guia de Contribuição](contributing.md) para entender melhor como você pode fazer parte desse processo e garantir que sua contribuição seja integrada da maneira mais eficaz possível.

### **Código de Conduta**

Estamos profundamente comprometidos em construir e manter uma comunidade inclusiva e acolhedora. Para isso, esperamos que todos os colaboradores sigam nosso [Código de Conduta](./.docs/code_of_conduct.md), que estabelece diretrizes claras para garantir um ambiente respeitoso e produtivo para todos.

Faça um favor a si e à comunidade: dedique um momento para ler e internalizar o código de conduta.

### **Licença**

O projeto Nanogo está disponível sob a licença Apache 2.0, uma licença permissiva e de código aberto que permite a liberdade de usar o software para qualquer finalidade, respeitando as limitações estabelecidas na [LICENÇA](license).

Dessa forma, você pode contribuir, modificar e distribuir o projeto, estando protegido juridicamente e respeitando os direitos e esforços dos outros contribuintes.

