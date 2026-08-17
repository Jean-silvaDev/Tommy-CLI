---
name: docker-cleaner-ops
description: >-
  Use esta skill quando o usuário solicitar operações de inspeção, agrupamento de Stacks Compose,
  pruning de recursos Docker ou resolução de problemas com a inicialização do Docker Desktop.
---

# Skill: Operações Docker & Pruning

Runbook procedimental para interação com o engine e CLI do Docker através de `pkg/docker/manager.go`.

---

## Passos Procedimentais

1. **Checar Status do Daemon Docker**:
   - Chame `docker.IsDockerRunning()`. Se estiver inativo, chame `docker.StartDockerEngine()`.
   - O `StartDockerEngine()` prioriza a busca do executável configurado em `config.GetDockerPath()`.

2. **Agrupar Containers por Stacks Compose**:
   - Utilize `docker.ListContainers()` e `docker.GetDockerStacks(containers)`. O agrupamento utiliza a label `com.docker.compose.project`.

3. **Pruning de Recursos**:
   - Para remoção segura de containers, imagens, volumes ou redes sem uso, utilize as funções expostas em `pkg/docker/manager.go`:
     - `PruneContainers()`
     - `PruneImages(all bool)`
     - `PruneVolumes()`
     - `PruneNetworks()`
     - `PruneSystem(all, volumes bool)`

4. **Tratamento de Exceções**:
   - Caso o executável `docker` não esteja no PATH, exiba instruções amigáveis ao usuário via `ui.PrintError(...)`.
