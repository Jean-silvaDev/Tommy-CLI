package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"tommy/pkg/promptbuilder"
	"tommy/pkg/ui"

	"github.com/spf13/cobra"
)

var promptCmd = &cobra.Command{
	Use:   "prompt",
	Short: "Assistente de elicitação de requisitos e gerador de prompts (Prompt Builder)",
	Long:  `Conduz uma entrevista estruturada (atuando como Arquiteto de Software / Analista de Sistemas Sênior) para gerar 10 artefatos de documentação e um prompt otimizado para IAs.`,
	Run: func(cmd *cobra.Command, args []string) {
		runPromptBuilderInteractive()
	},
}

var promptListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todas as entrevistas e prompts salvos em ./PromptBuilder",
	Run: func(cmd *cobra.Command, args []string) {
		sessions, err := promptbuilder.ListSessions()
		if err != nil || len(sessions) == 0 {
			ui.PrintInfo("Nenhuma entrevista encontrada em ./PromptBuilder")
			return
		}

		ui.PrintBanner("Entrevistas e Prompts Salvos")
		for i, s := range sessions {
			fmt.Printf(" [%d] 📁 ./PromptBuilder/%s\n", i+1, s)
		}
		fmt.Println()
	},
}

var promptLoadCmd = &cobra.Command{
	Use:   "load [caminho-answers.json]",
	Short: "Carrega uma entrevista salva anteriormente e gera os artefatos novamente",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		session, err := promptbuilder.LoadSession(filePath)
		if err != nil {
			ui.PrintError("Erro ao carregar sessão: %v", err)
			return
		}

		ui.PrintSuccess("Sessão '%s' carregada com sucesso!", session.Title)
		artifacts, err := promptbuilder.GenerateArtifacts(session)
		if err != nil {
			ui.PrintError("Erro ao gerar artefatos: %v", err)
			return
		}

		outDir, err := promptbuilder.SaveSession(session, artifacts)
		if err != nil {
			ui.PrintError("Erro ao salvar artefatos: %v", err)
			return
		}

		ui.PrintSuccess("Artefatos atualizados e salvos em %s", outDir)
	},
}

var promptPresetCmd = &cobra.Command{
	Use:   "preset",
	Short: "Lista os presets de stack técnica salvos em ~/.tommy/prompt_presets.json",
	Run: func(cmd *cobra.Command, args []string) {
		presets, err := promptbuilder.LoadPresets()
		if err != nil || len(presets) == 0 {
			ui.PrintInfo("Nenhum preset encontrado em ~/.tommy/prompt_presets.json")
			return
		}

		ui.PrintBanner("Presets de Stack Técnica Globais")
		for _, p := range presets {
			fmt.Printf(" 📌 %s: %s\n", p.Name, p.Description)
			for k, v := range p.Values {
				fmt.Printf("    - %s: %s\n", k, v)
			}
			fmt.Println()
		}
	},
}

var promptCleanCmd = &cobra.Command{
	Use:     "clean",
	Aliases: []string{"clear", "limpar"},
	Short:   "Esvazia a pasta de entrevistas e prompts salvos em ./PromptBuilder",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Limpeza do PromptBuilder")
		executeCleanPromptBuilderWithConfirmation()
	},
}

func executeCleanPromptBuilderWithConfirmation() {
	absDir, _ := filepath.Abs(filepath.Join(".", "PromptBuilder"))
	ui.PrintWarning("Esta operação apagará TODOS os prompts e entrevistas salvos em:\n   👉 %s", absDir)
	fmt.Println()

	if !ui.ConfirmPrompt("Deseja realmente apagar todo o conteúdo da pasta PromptBuilder?") {
		ui.PrintInfo("Operação de limpeza do PromptBuilder cancelada.")
		return
	}

	if err := promptbuilder.CleanPromptBuilder(); err != nil {
		ui.PrintError("Erro ao limpar PromptBuilder: %v", err)
		return
	}

	ui.PrintSuccess("Pasta PromptBuilder limpa com sucesso!")
}

func init() {
	promptCmd.AddCommand(promptListCmd)
	promptCmd.AddCommand(promptLoadCmd)
	promptCmd.AddCommand(promptPresetCmd)
	promptCmd.AddCommand(promptCleanCmd)
	rootCmd.AddCommand(promptCmd)
}

func runPromptBuilderInteractive() {
	for {
		ui.PrintBanner("Módulo Prompt Builder")

		options := []string{
			"1. 📝 Iniciar Nova Entrevista (Gerar Prompt & Artefatos)",
			"2. 📋 Listar Entrevistas e Prompts Salvos",
			"3. 📌 Ver Presets de Stack Técnica Globais",
			"4. 🧹 Limpar Pasta de Prompts (PromptBuilder)",
			"5. ⬅️  Voltar ao Menu Principal",
			"6. 🚪 Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Selecione uma opção do Prompt Builder", options)
		if err != nil || idx == 4 {
			break
		}

		if idx == 5 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			startNewInterview()
		case 1:
			listSavedSessions()
		case 2:
			listStackPresets()
		case 3:
			executeCleanPromptBuilderWithConfirmation()
			ui.WaitForEnter("Pressione Enter para continuar...")
		}
	}
}

func startNewInterview() {
	_, _, outDir, err := promptbuilder.RunInterviewWizard(nil)
	if err != nil {
		ui.PrintError("Entrevista finalizada: %v", err)
		return
	}

	absLatest, err := filepath.Abs(filepath.Join(".", "PromptBuilder", "latest.md"))
	if err != nil {
		absLatest = filepath.Join(".", "PromptBuilder", "latest.md")
	}

	ui.PrintBanner("Prompt Builder Concluído com Sucesso!")
	ui.PrintSuccess("Seu prompt final foi gerado em:\n   👉 %s\\prompt.md", outDir)
	ui.PrintInfo("Atalho de cópia direta:\n   👉 %s", absLatest)
	fmt.Println()
	ui.WaitForEnter("Pressione Enter para retornar ao menu...")
}

func listSavedSessions() {
	sessions, err := promptbuilder.ListSessions()
	if err != nil || len(sessions) == 0 {
		ui.PrintInfo("Nenhuma entrevista encontrada em ./PromptBuilder")
		ui.WaitForEnter("Pressione Enter para voltar...")
		return
	}

	ui.PrintBanner("Entrevistas e Prompts Salvos")
	for i, s := range sessions {
		fmt.Printf(" [%d] 📁 ./PromptBuilder/%s\n", i+1, s)
	}
	fmt.Println()
	ui.WaitForEnter("Pressione Enter para voltar...")
}

func listStackPresets() {
	presets, err := promptbuilder.LoadPresets()
	if err != nil || len(presets) == 0 {
		ui.PrintInfo("Nenhum preset encontrado em ~/.tommy/prompt_presets.json")
		ui.WaitForEnter("Pressione Enter para voltar...")
		return
	}

	ui.PrintBanner("Presets de Stack Técnica Globais")
	for _, p := range presets {
		fmt.Printf(" 📌 %s: %s\n", p.Name, p.Description)
		for k, v := range p.Values {
			fmt.Printf("    - %s: %s\n", k, v)
		}
		fmt.Println()
	}
	ui.WaitForEnter("Pressione Enter para voltar...")
}
