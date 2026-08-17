# 🧠 Memória Secundária do Projeto (Tommy CLI / MeuCLI)

Este arquivo serve como a **Memória Secundária Persistente** para assistentes de IA (Antigravity, Claude Code, Gemini CLI, Cursor, etc.). Ele registra decisões arquiteturais, mapeamento do ecossistema do código, lições aprendidas e regras de negócio essenciais.

---

## 🏛️ 1. Decisões Arquiteturais Fundamentais

1. **Separação Rígida de Camadas**:
   - `cmd/`: Apenas roteamento Cobra CLI e orquestração dos menus interativos TUI. Nenhuma lógica de sistema operacional ou negócios deve residir aqui.
   - `pkg/`: Lógica de negócio pura, adaptadores de sistema operacional, gerenciadores de serviços e design system (`pkg/ui/`).
   - `main.go`: Ponto de entrada leve que inicializa o leitor de ambiente (`config.LoadEnv()`) e dispara `cmd.Execute()`.

2. **Gerenciamento Desacoplado de Ambiente (`pkg/config/env.go`)**:
   - Toda leitura de variáveis de ambiente deve utilizar os getters centralizados de `pkg/config`:
     - `GetOutputDir()` -> `TOMMY_OUTPUT_DIR` (padrão: `./PromptBuilder`)
     - `GetConfigDir()` -> `TOMMY_CONFIG_DIR` (padrão: `~/.tommy`)
     - `GetDockerPath()` -> `TOMMY_DOCKER_PATH` (auto-detectado se vazio)
     - `GetTempDirs()` -> `TOMMY_TEMP_DIR`, `TOMMY_SYS_TEMP_DIR` + temporários padrão do SO
     - `GetTrashDir()` -> `TOMMY_TRASH_DIR` (Lixeira do SO ou pasta Trash)
     - `GetAuthor()` -> `TOMMY_AUTHOR` (fallback para `USERNAME` -> `USER` -> `"Desenvolvedor"`)
   - Suporte a expansão de til (`~`) ativado via `ExpandHome(path)`.

3. **Design System & UX/UI TUI (`pkg/ui/`)**:
   - Todas as mensagens de interface devem estar em **Português (PT-BR)**.
   - Utilização de estilos do `Lipgloss` em `pkg/ui/styles.go`.
   - Navegação interativa uniforme com `Voltar ao Menu Principal` e `Sair da Aplicação`.

---

## 🗺️ 2. Mapeamento do Código e Módulos

| Pacote | Função Principal | Principais Arquivos |
|---|---|---|
| `pkg/config` | Leitura e getters do `.env` com expansão `~` | [`env.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/config/env.go) |
| `pkg/promptbuilder` | Entrevistas IA com persona de Arquiteto Sênior e geração de artefatos com SecDevOps/OWASP | [`engine.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/promptbuilder/engine.go), [`storage.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/promptbuilder/storage.go), [`generator.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/promptbuilder/generator.go) |
| `pkg/cleaner` | Limpeza resiliente de temporários (%TEMP%, Windows Temp, var/tmp), NuGet cache, Docker e Lixeira | [`temp.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/cleaner/temp.go), [`recyclebin.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/cleaner/recyclebin.go), [`nuget.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/cleaner/nuget.go) |
| `pkg/docker` | Gestão de containers, Stacks Compose, imagens, volumes, redes e limpeza | [`manager.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/docker/manager.go) |
| `pkg/network` | Inspeção de IP local/WAN, ping, port scan, CIDR, flush DNS e encerramento de processo por porta | [`network.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/network/network.go) |
| `pkg/git` | Assistente Git (commits interativos no padrão Conventional Commits, histórico, branches, tags, remote) | [`helper.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/git/helper.go) |
| `pkg/ui` | Componentes visuais TUI, tabelas, caixas de diálogo, banner, quotes e exit card | [`styles.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/ui/styles.go), [`selector.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/ui/selector.go) |

---

## 🚨 3. Lições Aprendidas & Gotchas Evitados (Histórico)

- **Nunca Hardcodar Caminhos de Disco**:
  - Evite usar `C:\Program Files`, `C:\Windows` ou `C:\Users\username`. Em vez disso, use `os.Getenv("ProgramFiles")`, `os.Getenv("SystemDrive")` ou a função `config.ExpandHome`.
- **Inclusão do PATH no Windows**:
  - A documentação de instalação no Windows deve instruir o comando PowerShell dinâmico:
    `[Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", "User") + ";" + (Get-Location).Path, "User")`
- **Remoção de Temporários com Atributo Read-Only no Windows**:
  - No Windows, arquivos temporários podem falhar ao serem excluídos se estiverem com atributo somente leitura. A função `tryRemoveItem` em `pkg/cleaner/temp.go` aplica `os.Chmod(path, 0666)` antes de `os.RemoveAll`.
- **Padrão de Commit**:
  - O projeto utiliza a convenção **Conventional Commits** (`feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`, `style:`).

---

## 🧪 4. Protocolo Obrigatório de Validação para IAs

Sempre que a IA realizar modificações no código Go:
1. Executar os testes automatizados sem cache:
   ```bash
   go test -count=1 ./...
   ```
2. Garantir que `go.mod` e `go.sum` permaneçam íntegros.
3. Verificar se `.env.example` reflete qualquer nova variável de ambiente adicionada ao código.
