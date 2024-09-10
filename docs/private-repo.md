### Documentação: Habilitando SSH no `go get` para Repositórios Privados

Esta documentação descreve como configurar o `go get` para funcionar com repositórios privados utilizando SSH em vez de HTTPS. Isso é útil para evitar problemas de autenticação que ocorrem quando você tenta acessar repositórios privados no GitHub ou em outras plataformas de controle de versão.

#### 1. **Verifique se o SSH está configurado no GitHub**

Antes de tudo, é necessário garantir que você tenha uma chave SSH configurada no GitHub:

1. **Gerar uma nova chave SSH (se necessário):**

    Se você ainda não possui uma chave SSH, você pode gerar uma com o comando:

    ```bash
    ssh-keygen -t ed25519 -C "seu_email@exemplo.com"
    ```

    Siga as instruções para salvar a chave em uma localização segura.

2. **Adicionar a chave SSH ao agente SSH:**

    Inicie o agente SSH e adicione sua chave:

    ```bash
    eval "$(ssh-agent -s)"
    ssh-add ~/.ssh/id_ed25519
    ```

3. **Adicionar a chave SSH ao GitHub:**

    Copie a chave SSH para o clipboard:

    ```bash
    cat ~/.ssh/id_ed25519.pub | pbcopy
    ```

    Em seguida, acesse [GitHub > Configurações > SSH and GPG keys](https://github.com/settings/keys) e adicione uma nova chave SSH.

#### 2. **Configurar o Git para usar SSH em vez de HTTPS**

Para garantir que o Go utilize SSH ao invés de HTTPS quando fizer o download de módulos privados, você deve configurar o Git para substituir URLs HTTPS por SSH.

1. **Configure o Git para utilizar SSH para o GitHub:**

    No terminal, execute:

    ```bash
    git config --global url."ssh://git@github.com/".insteadOf "https://github.com/"
    ```

    Isso faz com que qualquer operação `go get` que tente acessar `https://github.com/usuario/repo` seja redirecionada para `ssh://git@github.com/usuario/repo`.

#### 3. **Configurar `GOPRIVATE` para repositórios privados**

A variável de ambiente `GOPRIVATE` informa ao Go que ele deve ignorar o proxy de módulos público para certos repositórios privados. Isso impede que informações sensíveis sejam enviadas para proxies públicos.

1. **Definir `GOPRIVATE` para o seu domínio privado:**

    Se você estiver trabalhando com repositórios privados em GitHub, GitLab ou outra plataforma de Git, adicione o seguinte ao seu arquivo de perfil de shell (`~/.bashrc`, `~/.zshrc`, etc.):

    ```bash
    export GOPRIVATE=github.com/usuario/*
    ```

    Substitua `github.com/usuario/*` pelo caminho do repositório privado que você está utilizando.

#### 4. **Verifique a configuração**

Após configurar, é importante verificar se tudo está funcionando corretamente.

1. **Teste com `go get`:**

    Tente buscar um repositório privado utilizando `go get`:

    ```bash
    go get github.com/usuario/repo_privado
    ```

    Se tudo estiver configurado corretamente, o Go deverá baixar o repositório utilizando SSH, sem pedir credenciais ou gerar erros de autenticação.

#### 5. **Dicas Adicionais**

- **SSH Config File:** Você pode configurar alias e comportamentos específicos para seus repositórios privados editando o arquivo `~/.ssh/config`. Um exemplo de configuração para o GitHub pode ser:

    ```bash
    Host github.com
        HostName github.com
        User git
        IdentityFile ~/.ssh/id_ed25519
        IdentitiesOnly yes
    ```

- **Atualizações Futuras:** Se você precisar adicionar mais repositórios privados, basta atualizar a variável `GOPRIVATE` com o novo caminho.
