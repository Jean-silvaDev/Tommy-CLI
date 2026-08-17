# 🏗️ Arquitetura e Decisões Técnicas do Tommy CLI

Documentação técnica referente à arquitetura, estrutura de pacotes, gerenciamento de ambiente e padrões de projeto adotados na construção do **Tommy CLI**.

---

## 🧩 Visão Geral da Arquitetura

O **Tommy CLI** foi desenvolvido seguindo os princípios de **Clean Architecture** e **Modularidade em Go**:

```text
Tommy/
├── main.go               # Ponto de entrada (Entrypoint CLI & Env Loader)
├── .env                  # Variáveis de ambiente locais (não versionado)
├── .env.example          # Modelo de variáveis de ambiente para a equipe
├── .gitignore            # Regras de exclusão do repositório Git
├── cmd/                  # Camada de Apresentação & Roteamento (Cobra Framework)
│   ├── root.go           # Comando raiz e Menu Principal Interativo (🚀 TOMMY)
│   ├── prompt.go         # Subcomandos e assistente do Prompt Builder
│   ├── network.go        # Subcomandos e ferramentas de Redes de Computadores
│   ├── docker.go         # Subcomandos e Hub do Docker Especialista
│   ├── git.go            # Subcomandos e fluxo de trabalho do Git
│   └── clean.go          # Subcomandos e Central de Limpeza Geral
├── pkg/                  # Camada de Negócio e Serviços (Domain/Services)
│   ├── config/           # Carregador de variáveis de ambiente (.env)
│   ├── promptbuilder/    # Motor de elicitação de requisitos, catálogo de perguntas e gerador de artefatos com SecDevOps
│   ├── network/          # Motores de rede (IP, Ping, Scan, CIDR, Flush DNS, Process Kill)
│   ├── docker/           # Wrappers Docker (Containers, Stacks, Images, Volumes, Networks, Daemon)
│   ├── git/              # Wrapper Git (Log, Commit, Branch, Tag, Merge Conflicts)
│   ├── cleaner/          # Motores de limpeza resiliente (Temp, NuGet, Docker Prune, Lixeira)
│   └── ui/               # Design System TUI (LipGloss styles, PromptUI templates, Tables)
└── docs/                 # Documentação técnica e manuais de instalação
```

---

## 🛠️ Pacotes Go e Suas Responsabilidades

### 1. `pkg/config/env.go`
- **Leitor & Resolvedor Centralizado de Variáveis de Ambiente:** Carrega automaticamente chave-valor do arquivo `.env` na inicialização do `main.go`.
- **Expansão de Caminho Home (`ExpandHome`):** Converte o caractere `~` no caminho absoluto do diretório de usuário do SO.
- **Getters Tipados com Fallback:** Centraliza a resolução e padrão de `TOMMY_OUTPUT_DIR`, `TOMMY_CONFIG_DIR`, `TOMMY_DOCKER_PATH`, `TOMMY_TEMP_DIR`, `TOMMY_SYS_TEMP_DIR`, `TOMMY_TRASH_DIR` e `TOMMY_AUTHOR`.

### 2. `pkg/promptbuilder/`
- **Elicitação de Requisitos:** Conduz entrevistas interativas para 13 modalidades técnicas.
- **Padrão de Segurança SecDevOps:** Insere por padrão diretrizes OWASP (SQLi com Prepared Statements, XSS, RBAC/IDOR, proteção de segredos e testes automatizados) em todos os prompts.
- **Persistência Flexível:** Salva sessões e presets globais em `~/.tommy/prompt_presets.json`.

### 3. `pkg/cleaner/`
- **Exclusão Resiliente por Item (`tryRemoveItem`):** Desmarca atributos Read-Only e remove arquivos/pastas de forma isolada. Arquivos travados por processos do SO são ignorados graciosamente.

### 4. `pkg/network/`
- **Inspeção de Interfaces & WAN:** Consulta IPs locais/remotos e executa operações de rede (Ping, Scan de Portas, Flush DNS, Kill Process em porta).

### 5. `pkg/docker/`
- **Gerenciamento de Stacks & Pruning:** Identificação de recursos em uso vs órfãos e operações de ciclo de vida do daemon.

### 6. `pkg/ui/`
- **Design System TUI:** Baseado em Lipgloss, com suporte a breadcrumbs, frasário dinâmico por horário (`PrintExitMessage`), botões de navegação uniformes (`Voltar` e `Sair da Aplicação`) e remoção de variação de emojis (`stripVS`).

---

## 📐 Padrões de Design e Decisões Técnicas

1. **Subcomandos Híbridos (TUI + CLI Direto):**
   - **TUI (Interativo):** Executar `tommy` abre o menu navegável por setas.
   - **CLI (Direto):** Executar `tommy clean temp` aceita execução direta scriptável.

2. **Garantia de Não-Versionamento de Dados Sensíveis:**
   - O arquivo `.gitignore` proíbe o envio do `.env`, diretório `./PromptBuilder` e executáveis compilados ao GitHub.
