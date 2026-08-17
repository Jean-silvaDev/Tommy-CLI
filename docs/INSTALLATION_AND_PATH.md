# ⚙️ Guia Completo de Instalação e Configuração do PATH (Tommy CLI)

Este guia orienta como compilar e configurar o **Tommy CLI** para ser executado globalmente de qualquer diretório no seu computador.

---

## 🎯 Objetivo

Após realizar esta configuração simples, você poderá abrir o terminal em qualquer pasta do seu computador (ex: seu diretório de projetos ou pasta de trabalho) e digitar:

```bash
tommy
```

... para abrir o dashboard interativo ou executar subcomandos diretos como `tommy prompt`, `tommy network`, `tommy docker`, `tommy git` ou `tommy clean`.

---

## 🚀 Método 1: No Windows via PowerShell (Recomendado - 1 Clique)

1. Abra o **PowerShell** na pasta raiz onde o projeto **MeuCLI** foi clonado/compilado.
2. Execute o comando abaixo para adicionar dinamicamente o diretório atual do executável `tommy.exe` ao seu `PATH` de Usuário:

```powershell
[Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", "User") + ";" + (Get-Location).Path, "User")
```

3. Feche todas as janelas do terminal e abra uma **nova janela de terminal**.
4. Teste digitando:

```powershell
tommy --help
```

Se o menu de ajuda do **Tommy CLI** for exibido, a configuração foi concluída com sucesso! 🎉

---

## 🛠️ Método 2: No Windows via Interface Gráfica (GUI)

1. Pressione `Win + R`, digite `sysdm.cpl` e aperte **Enter**.
2. Na aba **Avançado**, clique no botão **Variáveis de Ambiente...**.
3. Na seção **Variáveis do Usuário**, selecione a variável **Path** e clique em **Editar...**.
4. Clique em **Novo** e cole o caminho completo da pasta do projeto (exemplo: `C:\Projetos\MeuCLI`).

5. Clique em **OK** em todas as janelas para salvar.
6. Reabra o seu terminal.

---

## 🐧 Método 3: No Linux / macOS

1. Abra o arquivo de configuração do seu shell (`~/.bashrc` ou `~/.zshrc`):

```bash
nano ~/.zshrc
```

2. Adicione a seguinte linha no final do arquivo:

```bash
export PATH="$PATH:/caminho/para/pasta/MeuCLI"
```

3. Salve e recarregue a sessão:

```bash
source ~/.zshrc
```

---

## 💡 Como Recompilar Após Alterar o Código em Go

Sempre que você modificar o código-fonte em Go do **Tommy CLI**, atualize o binário rodando na pasta do projeto:

```bash
go build -o tommy.exe main.go
```

Como a pasta do projeto já está configurada no seu `PATH`, as alterações estarão **imediatamente disponíveis em todo o seu sistema** sem precisar reconfigurar nada! 🎉
