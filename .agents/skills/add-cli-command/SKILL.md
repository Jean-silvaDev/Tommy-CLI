---
name: add-cli-command
description: >-
  Use esta skill quando o usuário solicitar a adição de um novo subcomando CLI Cobra
  ou uma nova opção no menu navegável interativo (TUI) do Tommy CLI.
---

# Skill: Adição de Subcomando CLI & Opção TUI

Este runbook instrui a IA sobre como adicionar novos subcomandos ou opções de menu mantendo o padrão híbrido do **Tommy CLI** (suporte a execução direta via terminal e navegação TUI).

---

## Passo a Passo

### 1. Criar ou Atualizar o Subcomando Cobra em `cmd/`
- Caso seja um novo domínio, crie um novo arquivo em `cmd/<dominio>.go`.
- Caso pertença a um módulo existente (`clean`, `docker`, `network`, `git`, `prompt`), adicione a variável `*cobra.Command`:

```go
var novoSubCmd = &cobra.Command{
    Use:   "acao",
    Short: "Descrição curta em português",
    Run: func(cmd *cobra.Command, args []string) {
        ui.PrintBanner("Título da Ação")
        executarMinhaAcao()
    },
}
```

### 2. Implementar a Lógica de Serviço em `pkg/<modulo>/`
- Não insira lógica de baixo nível no arquivo do pacote `cmd/`.
- Crie ou edite as funções correspondentes em `pkg/<modulo>/`.

### 3. Registrar o Comando no Menu Interativo TUI
- Adicione a chamada do novo comando na função de menu interativo correspondente (ex: `runInteractiveCleanMenu()` em `cmd/clean.go`).
- Certifique-se de incluir a opção de retorno (`Voltar ao Menu Principal`).

### 4. Validar
- Execute os testes e verifique a compilação:
  `go test ./...`
  `go build -o tommy.exe main.go`
