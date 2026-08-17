---
name: verify-and-test
description: >-
  Use esta skill sempre que concluir alterações no código fonte Go, refatorações ou modificações de ambiente
  para validar a integridade dos testes e a compilação do binário do Tommy CLI.
---

# Skill: Protocolo de Verificação e Testes

Runbook para execução de testes unitários automatizados, checagem do binário e validação de configurações.

---

## Passos de Verificação

### 1. Executar Testes Unitários Sem Cache
No terminal da raiz do projeto, execute:

```powershell
go test -count=1 ./...
```

Verifique se a saída indica `ok` para todos os pacotes (`cleaner`, `config`, `promptbuilder`, etc.).

### 2. Validar Compilação do Binário Executável
Gere o executável compilado para garantir que não haja erros de build:

```powershell
go build -o tommy.exe main.go
```

### 3. Verificar Sincronismo do `.env`
- Confirme se todas as variáveis utilizadas em `pkg/config/env.go` possuem correspondência documentada em `.env.example`.
- Garanta que nenhuma senha ou caminho físico estático tenha sido adicionado ao código.
