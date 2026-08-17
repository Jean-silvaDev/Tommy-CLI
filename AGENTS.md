# 🤖 Diretrizes Gerais para Agentes de IA (AGENTS.md)

Este repositório (**Tommy CLI / MeuCLI**) contém regras estritas de arquitetura, desenvolvimento em Go, design de interface TUI e gerenciamento de ambiente que **devem ser seguidas por todos os assistentes de IA** (Antigravity, Claude Code, Gemini CLI, Cursor, Copilot, etc.).

> 🧠 **Memória Secundária do Projeto**: Antes de planejar alterações complexas ou diagnósticos, consulte sempre o arquivo de Memória Secundária:
> [`.agents/MEMORY.md`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.agents/MEMORY.md)

---

## 🎯 1. Princípios Fundamentais

1. **Separação Rígida de Pacotes**:
   - `cmd/`: Comandos Cobra CLI e menus navegáveis TUI. Sem lógica de baixo nível do sistema operacional.
   - `pkg/`: Lógica de negócio, integração de serviços, wrappers e adaptadores de sistema.
   - `pkg/config`: Ponto único de verdade para variáveis de ambiente do arquivo `.env`.
   - `pkg/ui`: Estilos gráficos (Lipgloss) e componentes visuais interativos.

2. **Convenção de Commits (Conventional Commits)**:
   - Todas as mensagens de commit sugeridas ou geradas pela IA devem seguir o padrão:
     - `feat:` Novas funcionalidades.
     - `fix:` Correção de bugs.
     - `docs:` Alterações em documentação.
     - `chore:` Tarefas de manutenção ou dependências.
     - `refactor:` Refatorações sem alteração de comportamento.
     - `test:` Inclusão ou ajuste de testes unitários.

3. **Zero Hardcode de Diretórios Locais**:
   - NUNCA insira caminhos absolutos de máquinas físicas (ex: `C:\Users\...` ou `C:\Program Files\...`).
   - TODA configuração de diretório deve ser lida via `pkg/config/env.go` e ter fallback dinâmico usando `os.Getenv("SystemDrive")`, `os.Getenv("ProgramFiles")`, `os.UserHomeDir()` ou `config.ExpandHome()`.

4. **Idioma Padrão (PT-BR)**:
   - Todas as mensagens da interface do usuário (CLI e TUI), documentação e comentários no código fonte devem estar em **Português do Brasil (PT-BR)**.

---

## 📚 2. Mapeamento de Regras e Skills Sob Demanda

O repositório possui regras adicionais e guias de execução (*runbooks*) no diretório `.agents/`:

### Regras Ativas:
- **[Go Architecture & Errors](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.agents/rules/go-architecture.md)**: Padrões de código Go, tratamento de erros com `%w` e estilos Lipgloss.
- **[Env Guardrails](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.agents/rules/env-guardrails.md)**: Proteção de variáveis de ambiente e arquivo `.env`.

### Skills Sob Demanda (*Runbooks*):
- **[add-cli-command](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.agents/skills/add-cli-command/SKILL.md)**: Passo a passo para adicionar subcomandos Cobra e menus TUI.
- **[verify-and-test](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.agents/skills/verify-and-test/SKILL.md)**: Execução de testes unitários, validação do `.env` e compilação do binário.
- **[promptbuilder-extension](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.agents/skills/promptbuilder-extension/SKILL.md)**: Inclusão de novos tipos de perguntas, entrevistas e regras SecDevOps no PromptBuilder.
- **[docker-cleaner-ops](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.agents/skills/docker-cleaner-ops/SKILL.md)**: Gestão de containers, Stacks Compose, imagens e volumes Docker.
- **[network-diagnostics](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/.agents/skills/network-diagnostics/SKILL.md)**: Diagnósticos de rede, ping, scan de portas, CIDR, flush DNS e kill por porta.

---

## 🧪 3. Protocolo de Encerramento de Tarefas

Antes de declarar qualquer tarefa como concluída, a IA **DEVE**:
1. Executar a suíte de testes: `go test -count=1 ./...`
2. Verificar se o código compila sem erros: `go build -o tommy.exe main.go`
3. Garantir que nenhuma alteração tenha quebrado o suporte ao arquivo `.env` ou `.env.example`.
