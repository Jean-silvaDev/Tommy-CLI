# Guardrails para Variáveis de Ambiente e Segurança

Regras estritas de gerenciamento de ambiente, caminhos e segredos para o **Tommy CLI**.

---

## 🔒 Regras de Isolamento e Segurança

1. **Centralização em `pkg/config/env.go`**:
   - Todo acesso a variáveis de ambiente deve ser realizado através das funções expostas pelo pacote `pkg/config`.
   - NUNCA espalhe chamadas diretas `os.Getenv("TOMMY_...")` em módulos secundários sem expor um getter limpo em `pkg/config/env.go`.

2. **Expansão de Caminhos de Usuário**:
   - Todo caminho retornado por getters de ambiente que possa conter o caractere `~` deve passar pela função `config.ExpandHome(path)`.

3. **Sincronia entre `.env` e `.env.example`**:
   - Sempre que uma nova variável de ambiente for introduzida na aplicação:
     1. Adicionar o getter tipado em `pkg/config/env.go`.
     2. Documentar a chave com valor de exemplo e comentário no arquivo [`.env.example`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.env.example).
     3. Adicionar a entrada no arquivo local [`.env`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.env).
     4. Atualizar a tabela de variáveis no [`README.md`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/README.md).

4. **Proibição de Versionamento de Dados Locais**:
   - Certifique-se de que o `.env`, o diretório `./PromptBuilder` e o executável `tommy.exe` estejam mantidos no `.gitignore`.
