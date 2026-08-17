---
name: promptbuilder-extension
description: >-
  Use esta skill quando precisar adicionar novas modalidades de entrevista, novos blocos de perguntas
  ou alterar o motor de geração de prompts e artefatos SecDevOps no PromptBuilder.
---

# Skill: Extensão do PromptBuilder

Runbook procedimental para incluir novos tipos de prompt ou expandir as perguntas de elicitação de requisitos.

---

## Passos para Adicionar Nova Modalidade de Prompt

1. **Definir a Constante do Tipo**:
   - Abra [`pkg/promptbuilder/types.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/promptbuilder/types.go) e adicione a constante `PromptType` (ex: `TypeMinhaModalidade PromptType = "minha-modalidade"`).

2. **Registrar a Opção no Catálogo**:
   - Adicione o tipo e a descrição amigável em `GetPromptTypes()` em [`pkg/promptbuilder/catalog.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/promptbuilder/catalog.go).

3. **Construir os Blocos de Perguntas**:
   - Na função `GetQuestionBlocks(pType PromptType)` em `catalog.go`, adicione a instrução `switch` correspondente retornando os blocos de perguntas.

4. **Preservar as Diretrizes SecDevOps**:
   - Certifique-se de que a geração de prompts em [`pkg/promptbuilder/generator.go`](file:///c:/Users/jeans/OneDrive/Documentos/Projetos/MeuCLI/pkg/promptbuilder/generator.go) mantenha a seção obrigatória **SecDevOps & OWASP** (Prepared Statements, XSS, RBAC/IDOR e testes).

5. **Testar a Salvação de Artefatos**:
   - Execute `go test ./pkg/promptbuilder` para garantir que `SaveSession` grave os 10 artefatos no diretório configurado via `config.GetOutputDir()`.
