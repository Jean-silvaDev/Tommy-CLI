# 🚀 Tommy - CLI de Produtividade, Redes & Docker para Desenvolvedores

**Tommy** é uma ferramenta de linha de comando (CLI) moderna, modular e de alta performance escrita em **Go**, projetada para automatizar tarefas cotidianas do fluxo de desenvolvimento, elicitação de requisitos com IA (**PromptBuilder**), gerenciamento completo do ambiente **Docker**, diagnósticos de **Redes de Computadores** e produtividade no **Git**.

---

## 📋 Recursos Principais

### 📝 1. Prompt Builder & Elicitação de Requisitos (`tommy prompt`)
- **Entrevista Interativa (Persona Arquiteto Sênior):** Conduz perguntas inteligentes orientadas ao tipo de projeto (Novo Projeto, Feature, Bugfix, Refatoração, Arquitetura, Segurança, DevOps, etc.).
- **Geração de 10 Artefatos + Prompt IA Consolidado:** Gera prompts prontos para produção com a seção padrão **SecDevOps & OWASP** (prevenção contra SQL Injection com Prepared Statements, XSS, RBAC/IDOR, proteção de segredos e bateria de testes automatizados).
- **Gerenciamento & Limpeza de Prompts:** Listagem, recarregamento e limpeza de prompts salvos em `./PromptBuilder` (criados e recriados automaticamente se a pasta for removida).
- **Presets de Stack Técnica Globais (`~/.tommy/prompt_presets.json`):** Salva e reutiliza stacks de tecnologia preferidas.

### 🌐 2. Módulo de Redes de Computadores & Diagnóstico (`tommy network` / `net`)
- **Inspetor de Interfaces de Rede:** Exibe IPv4 local, IPv6, Endereço MAC, Gateway padrão e IP Público WAN (com provedor/ISP).
- **Liberar Porta em Uso (`KillProcessOnPort`):** Mapeia portas escutadas para o PID e nome do executável e encerra o processo de forma interativa.
- **Checador de Latência Ping:** Medição de tempo de resposta em milissegundos para hosts locais ou remotos.
- **Scanner de Portas de Dev:** Varredura em portas comuns de desenvolvimento local (`3000`, `5173`, `8080`, `5432`, `3306`, `27017`, `6379`).
- **Calculadora de Sub-rede CIDR:** Cálculo de máscara de rede, IP de rede, IP de broadcast e faixa de IPs válidos a partir de notações CIDR (ex: `192.168.1.0/24`).
- **Flush DNS (`ipconfig /flushdns`):** Limpeza rápida de cache DNS do sistema operacional.

### 🐳 3. Hub Completo Docker & Recursos de Especialista (`tommy docker`)
- **Containers & Stacks Compose:**
  - Agrupamento automático de containers por projeto Docker Compose (`com.docker.compose.project`).
  - Operações de inicialização, parada e reinicialização em lote por Stack.
  - Leitura de Logs em Tempo Real (`docker logs --tail 100`) e Inspeção Detalhada (`docker inspect`).
- **Gerenciador de Imagens Docker:**
  - Listagem com badges em tempo real: **`✔ EM USO`** vs **`⚠️ ÓRFÃ (SEM USO)`**.
  - Remoção seletiva ou forçada (`-f`), download de novas imagens do Docker Hub (`docker pull`) e descarte em massa (`docker image prune -a`).
- **Gerenciador de Volumes Docker:**
  - Identificação de uso (**`✔ EM USO`** vs **`⚠️ ÓRFÃO`**).
  - Exclusão individual e descarte de volumes nomeados de Compose (`docker volume prune -a -f`).
- **Gerenciador de Redes Docker:**
  - Classificação de redes em **`CUSTOMIZADA`** vs **`SISTEMA (PROTEGIDA)`** (`bridge`, `host`, `none`).
  - Inspeção detalhada de redes (`docker network inspect`) com containers e IPs internos vinculados.
- **Central de Limpeza & Pruning Docker:**
  - Painel de disco em tempo real (`docker system df`) e limpeza seletiva ou profunda do sistema (`docker system prune -a --volumes`).

### 🔀 4. Assistente de Fluxo de Trabalho Git (`tommy git`)
- **Visualizador de Histórico de Commits (`tommy git log`):** Tabela elegante com colunas `HASH`, `AUTOR`, `DATA` e `ASSUNTO`.
- **Assistente de Resolução de Conflitos de Merge:** Identifica arquivos em conflito, abre resolução visual ou permite abortar o merge.
- **Commit Interativo (Conventional Commits):** Assistente passo a passo para estagiar arquivos (`git add -A`) e estruturar mensagens no padrão *Conventional Commits* (`feat:`, `fix:`, `chore:`, `docs:`, `refactor:`, `test:`, `style:`).
- **Gerenciador de Branches e Tags:** Listagem, alternância, criação e exclusão de branches e tags anotadas.
- **Operações Remotas:** `git pull`, `git push` e vínculo de remoto `origin`.

### 🧹 5. Módulo de Limpeza Geral (`tommy clean`)
- **Exclusão Resiliente de Temporários:** Varredura em `%TEMP%` (Usuário) e `C:\Windows\Temp` (Sistema) tentando remoção item a item (desmarcando atributos de leitura e ignorando graciosamente arquivos em uso pelo SO).
- **Cache do .NET NuGet:** Limpeza de caches globais de pacotes NuGet (`dotnet nuget locals all --clear`).
- **Lixeira do Sistema (`tommy clean trash`):** Esvaziamento permanente da Lixeira do Windows (Recycle Bin) sem mensagens de travamento.
- **Limpeza Completa (`tommy clean all`):** Execução sequencial de todas as rotinas de limpeza do computador com pausa e relatório de encerramento.

### 🖥️ 6. System Manager & Monitoramento (`tommy system` / `sys`)
- **Dashboard Interativo:** Cartão visual estilizado com acompanhamento em tempo real da Memória RAM (total/usada/% com barra colorida), uso estimado da CPU e capacidade de armazenamento em todos os volumes de disco (`C:`, `D:`, etc.).
- **Gerenciador de Memória RAM (`tommy system memory`):** Métricas de RAM total, usada e disponível e listagem dos maiores processos consumidores.
- **Gerenciador de Processos Ativos (`tommy system processes`):** Tabela de processos com PID, uso de memória, threads e status. Permite busca por nome/PID e encerramento seguro (`KillProcess`) com proteção automática contra fechamento acidental de processos críticos do sistema operacional (`csrss.exe`, `lsass.exe`, `services.exe`, `explorer.exe`, etc.) e diálogo de confirmação.
- **Gerenciador de Discos & Volumes (`tommy system disk`):** Classificação dos discos (SSD, HDD, NVMe, Removível, Rede), sistema de arquivos, total/usado/livre, uso percentual e status de saúde.
- **Aplicativos Instalados & Desinstalação (`tommy system apps`):** Descoberta nativa via Registro do Windows (`HKLM`/`HKCU`) e `winget`. Suporta busca por fabricante/nome, ordenação por tamanho em disco, filtro por unidade e desinstalação segura com caixa de detalhes de confirmação.
- **Analisador de Espaço (`tommy system analyze <caminho>`):** Mapeamento profundo do tamanho ocupado por diretórios e arquivos (protegido contra loops de symlinks/junctions).
- **Maiores Arquivos (`tommy system largest <caminho>`):** Localiza os maiores arquivos individuais no disco com suporte à flag `--dry-run`.
- **Monitoramento em Tempo Real (`tommy system monitor`):** Atualização contínua de recursos sem sobrecarregar a CPU.

---

## ⚙️ Variáveis de Ambiente (`.env`)

O **Tommy CLI** suporta configuração via arquivo `.env`. Copie o modelo `.env.example` para `.env` na raiz da aplicação:

```bash
# Copiar o exemplo para criar suas configurações locais
cp .env.example .env
```

### Variáveis Suportadas:

| Variável | Valor Padrão | Descrição |
|---|---|---|
| `TOMMY_ENV` | `development` | Ambiente de execução (`development` / `production`). |
| `TOMMY_OUTPUT_DIR` | `./PromptBuilder` | Pasta onde os prompts e artefatos são salvos. |
| `TOMMY_CONFIG_DIR` | `~/.tommy` | Pasta de configurações globais e presets de stack técnica. |
| `TOMMY_AUTHOR` | `os.Getenv("USERNAME")` | Nome do autor atribuído às sessões do PromptBuilder. |
| `TOMMY_DOCKER_PATH` | *Auto-detectado* | Caminho customizado para o executável `Docker Desktop.exe`. |
| `TOMMY_TEMP_DIR` | `%TEMP%` / `$TMPDIR` | Pasta(s) temporária(s) customizada(s) para limpeza. |
| `TOMMY_SYS_TEMP_DIR` | `SystemRoot\Temp` | Pasta temporária de sistema adicional para limpeza. |
| `TOMMY_TRASH_DIR` | `~/.local/share/Trash/files` / Lixeira | Pasta de Lixeira customizada para esvaziamento. |
| `TOMMY_LOG_LEVEL` | `info` | Nível de detalhamento de logs. |
| `TOMMY_LANGUAGE` | `pt-BR` | Idioma das mensagens da interface. |

> 🔒 **Segurança:** O arquivo `.env` e a pasta `./PromptBuilder` estão listados no `.gitignore` e **nunca** serão enviados ao GitHub.

---

## 🎨 Sistema de Design & Experiência de Uso (UX/UI)

- **Citações Inspiradoras nos Cabeçalhos:** Frases motivacionais e dicas de tecnologia sorteadas aleatoriamente (`GetRandomHeaderQuote`).
- **Exit Card Dinâmico por Horário (`PrintExitMessage`):** Card estilizado com frases motivacionais adaptadas ao horário do dia (Manhã, Tarde, Noite, Madrugada) e bordas simétricas alinhadas.
- **Navegação Uniforme:** Opção `⬅️  Voltar ao Menu Principal` e `🚪 Sair da Aplicação` em todos os submenus interativos.
- **Interface 100% em Português (PT-BR).**

---

## 🌐 Como Adicionar ao PATH e Usar em Qualquer Pasta

Para utilizar o `tommy` em **qualquer diretório** do seu computador:

### No Windows (PowerShell)
Abra o PowerShell na pasta do projeto e execute:
```powershell
[Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", "User") + ";" + (Get-Location).Path, "User")
```

Feche e reabra o terminal. Agora você pode executar `tommy` de qualquer pasta!

> 📘 **Guia Detalhado de Instalação e PATH:** Veja o passo a passo completo no documento [`docs/INSTALLATION_AND_PATH.md`](docs/INSTALLATION_AND_PATH.md).

---

## 💻 Guia Rápido de Uso

### Modo Interativo (Dashboard Principal)
Basta digitar em qualquer terminal:

```bash
tommy
```

### Compilação do Executável

```bash
# Compilar / Atualizar o executável tommy.exe
go build -o tommy.exe main.go
```

### Modo Direto por Comandos

```bash
# Prompt Builder
tommy prompt list       # Lista entrevistas e prompts salvos
tommy prompt clean      # Limpa a pasta de prompts salvos

# Redes
tommy network info      # Exibe informações de IP, MAC e Gateway
tommy network ping      # Teste de latência
tommy network scan      # Varredura de portas dev locais
tommy network release   # Liberar porta em uso (kill process)
tommy network flushdns  # Limpar cache de DNS do SO

# Docker
tommy docker list       # Exibe a tabela de containers
tommy docker start web  # Inicia o container 'web'
tommy docker stop web   # Para o container 'web'

# Limpeza
tommy clean temp        # Limpa arquivos temporários do SO
tommy clean nuget       # Limpa o cache de pacotes .NET NuGet
tommy clean trash       # Esvazia a Lixeira do SO
tommy clean all         # Executa a limpeza completa

# Git
tommy git commit        # Abre o assistente de commit
tommy git log           # Exibe o histórico de commits
tommy git branch        # Gerencia branches
tommy git push          # Executa git push no remoto
```

---

## 🧪 Testes Unitários

Para rodar os testes automatizados da aplicação:

```bash
go test ./... -v
```

---

## 📄 Licença
Distribuído sob a licença MIT. Sinta-se livre para adaptar e reutilizar em seus projetos!
