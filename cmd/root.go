package cmd

import (
	"fmt"
	"os"

	"tommy/pkg/ui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tommy",
	Short: "Tommy - Ferramenta CLI de Produtividade para Desenvolvedores",
	Long: `Tommy é uma aplicação de linha de comando de alta performance para automação de tarefas do dia a dia:
- Elicitação de requisitos e gerador de prompts (Prompt Builder)
- Limpeza de arquivos temporários (%TEMP%), cache do .NET NuGet, Lixeira e Docker
- Gerenciamento de containers Docker (status, start, stop, restart, stacks Compose)
- Fluxo de trabalho Git interativo (Conventional Commits, logs, branches, tags, conflitos e push)
- Ferramentas e diagnósticos de redes de computadores (IP, Ping, Scan de Portas, Flush DNS, Subnet)`,
	Run: func(cmd *cobra.Command, args []string) {
		runMainMenu()
	},
}

// Execute é o ponto de entrada chamado pelo main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runMainMenu() {
	for {
		ui.PrintBanner("Menu Principal")

		options := []string{
			"1. 🧹  Limpeza de Recursos (Temp, NuGet, Docker, Lixeira)",
			"2. 🐳  Gerenciador de Containers Docker",
			"3. 🔀  Assistente de Fluxo de Trabalho Git",
			"4. 🌐  Ferramentas & Diagnóstico de Rede",
			"5. 📝  Prompt Builder (Elicitação de Requisitos & Gerador de Prompts)",
			"6. 🖥️  Gerenciador do Sistema (System Manager)",
			"7. 🚪  Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Selecione um módulo para utilizar", options)
		if err != nil || idx == 6 {
			ui.PrintExitMessage()
			break
		}

		switch idx {
		case 0:
			runInteractiveCleanMenu()
		case 1:
			runDockerDashboard()
		case 2:
			runGitDashboard()
		case 3:
			runNetworkDashboard()
		case 4:
			runPromptBuilderInteractive()
		case 5:
			runSystemDashboard()
		}
	}
}
