package cmd

import (
	"fmt"
	"os"

	"tommy/pkg/cleaner"
	"tommy/pkg/ui"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Limpa arquivos desnecessários (Temp, NuGet, Docker, Lixeira)",
	Long:  `Subcomando para varredura e remoção de arquivos temporários do SO, pacotes em cache do .NET NuGet, Lixeira e imagens/volumes sem uso no Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveCleanMenu()
	},
}

var cleanTempCmd = &cobra.Command{
	Use:   "temp",
	Short: "Limpa arquivos temporários do SO (Usuário + Sistema)",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Limpeza de Arquivos Temporários")
		executeCleanTempWithConfirmation()
	},
}

var cleanNuGetCmd = &cobra.Command{
	Use:   "nuget",
	Short: "Limpa o cache local de pacotes do .NET NuGet",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Limpeza de Cache NuGet")
		executeCleanNuGetWithConfirmation()
	},
}

var cleanDockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Limpa containers parados, imagens não utilizadas e volumes do Docker",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Limpeza de Recursos Docker")
		executeCleanDockerWithConfirmation()
	},
}

var cleanTrashCmd = &cobra.Command{
	Use:     "trash",
	Short:   "Esvazia a Lixeira do sistema (Recycle Bin)",
	Aliases: []string{"lixeira", "recyclebin"},
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Limpeza da Lixeira")
		executeCleanTrashWithConfirmation()
	},
}

var cleanAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Executa a limpeza completa (Temp + NuGet + Docker + Lixeira)",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Limpeza Completa do Sistema")
		executeCleanAllWithConfirmation()
	},
}

func executeCleanTempWithConfirmation() {
	ui.PrintInfo("Analisando arquivos temporários do sistema...")
	scanRes, err := cleaner.ScanTempFiles()
	if err != nil {
		ui.PrintError("%v", err)
		return
	}

	if scanRes.FilesCount == 0 {
		ui.PrintWarning("Nenhum arquivo temporário pendente para exclusão nos diretórios do sistema.")
		ui.WaitForEnter("")
		return
	}

	fmt.Println()
	if len(scanRes.FileList) > 0 {
		fmt.Println(ui.SubtitleStyle.Render("Arquivos / Pastas Temporários Encontrados:"))
		limit := 15
		if len(scanRes.FileList) < limit {
			limit = len(scanRes.FileList)
		}
		for i := 0; i < limit; i++ {
			fmt.Printf("  • %s\n", scanRes.FileList[i])
		}
		if len(scanRes.FileList) > limit {
			fmt.Printf("  ... e mais %d itens.\n", len(scanRes.FileList)-limit)
		}
		fmt.Println()
	}

	ui.PrintWarning("Foram encontrados %d arquivos/pastas (%s ocupados) elegíveis para exclusão nos diretórios temporários.",
		scanRes.FilesCount, cleaner.FormatBytes(scanRes.BytesFreed))
	fmt.Println()

	if !ui.ConfirmPrompt("Deseja realmente prosseguir com a exclusão destes arquivos temporários?") {
		ui.PrintInfo("Operação de limpeza temporária cancelada pelo usuário.")
		return
	}

	res, err := cleaner.CleanTempFiles()
	if err != nil {
		ui.PrintError("Erro na limpeza temp: %v", err)
		ui.WaitForEnter("Pressione Enter para continuar...")
		return
	}
	ui.PrintSuccess("Limpeza concluída! Removidos: %d arquivos (%s liberados). Ignorados: %d em uso pelo SO.",
		res.FilesRemoved, cleaner.FormatBytes(res.BytesFreed), res.ErrorsCount)
	ui.WaitForEnter("Pressione Enter para continuar...")
}

func executeCleanNuGetWithConfirmation() {
	ui.PrintWarning("Esta operação limpará TODO o cache global de pacotes do .NET NuGet no seu computador.")
	fmt.Println()

	if !ui.ConfirmPrompt("Deseja realmente limpar o cache de pacotes NuGet?") {
		ui.PrintInfo("Operação de limpeza NuGet cancelada pelo usuário.")
		return
	}

	err := cleaner.CleanNuGetCache()
	if err != nil {
		ui.PrintError("%v", err)
	}
	ui.WaitForEnter("Pressione Enter para continuar...")
}

func executeCleanDockerWithConfirmation() {
	runDockerCleanCenter()
}

func executeCleanTrashWithConfirmation() {
	ui.PrintInfo("Analisando arquivos contidos na Lixeira do sistema...")
	items, err := cleaner.GetRecycleBinItems()

	fmt.Println()
	if err != nil || len(items) == 0 {
		ui.PrintInfo("Nenhum item detectado ou a Lixeira do sistema já está vazia.")
		fmt.Println()
		ui.WaitForEnter("Pressione Enter para voltar...")
		return
	}

	fmt.Println(ui.SubtitleStyle.Render("Itens / Arquivos Encontrados na Lixeira:"))
	limit := 15
	if len(items) < limit {
		limit = len(items)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("  • %s\n", items[i])
	}
	if len(items) > limit {
		fmt.Printf("  ... e mais %d itens.\n", len(items)-limit)
	}
	fmt.Println()

	ui.PrintWarning("Esta operação apagará PERMANENTEMENTE todos os arquivos contidos na Lixeira do sistema.")
	fmt.Println()

	if !ui.ConfirmPrompt("Deseja realmente esvaziar a Lixeira do sistema?") {
		ui.PrintInfo("Esvaziamento de Lixeira cancelado pelo usuário.")
		return
	}

	err = cleaner.CleanRecycleBin()
	if err != nil {
		ui.PrintError("%v", err)
	}
	ui.WaitForEnter("Pressione Enter para continuar...")
}

func executeCleanAllWithConfirmation() {
	ui.PrintWarning("A limpeza completa afetará: Arquivos Temporários (%%TEMP%%), Caches NuGet, Recursos Docker e Lixeira do Sistema.")
	fmt.Println()

	if !ui.ConfirmPrompt("Deseja realmente executar a ROTINA COMPLETA de limpeza?") {
		ui.PrintInfo("Limpeza completa cancelada pelo usuário.")
		return
	}

	fmt.Println("1/4 [SO Temp]")
	executeCleanTempWithConfirmation()
	fmt.Println()

	fmt.Println("2/4 [.NET NuGet]")
	executeCleanNuGetWithConfirmation()
	fmt.Println()

	fmt.Println("3/4 [Docker Prune]")
	executeCleanDockerWithConfirmation()
	fmt.Println()

	fmt.Println("4/4 [Lixeira do Sistema]")
	executeCleanTrashWithConfirmation()
	fmt.Println()

	ui.PrintSuccess("Rotina completa de limpeza finalizada!")
	ui.WaitForEnter("Pressione Enter para continuar...")
}

func runInteractiveCleanMenu() {
	for {
		ui.PrintBanner("Módulo de Limpeza")

		options := []string{
			"1. 🧹 Limpar Arquivos Temporários (%TEMP% + C:\\Windows\\Temp)",
			"2. 📦 Limpar Cache de Pacotes .NET NuGet",
			"3. 🐳 Limpar Recursos Não Utilizados do Docker",
			"4. 🗑️  Esvaziar Lixeira do Sistema (Recycle Bin)",
			"5. ⚡ Executar Limpeza Completa (Tudo)",
			"6. ⬅️  Voltar ao Menu Principal",
			"7. 🚪 Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Selecione uma opção de limpeza", options)
		if err != nil || idx == 5 {
			break // Voltar ao Menu Principal
		}

		if idx == 6 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			executeCleanTempWithConfirmation()
		case 1:
			executeCleanNuGetWithConfirmation()
		case 2:
			executeCleanDockerWithConfirmation()
		case 3:
			executeCleanTrashWithConfirmation()
		case 4:
			executeCleanAllWithConfirmation()
		}
	}
}

func init() {
	cleanCmd.AddCommand(cleanTempCmd)
	cleanCmd.AddCommand(cleanNuGetCmd)
	cleanCmd.AddCommand(cleanDockerCmd)
	cleanCmd.AddCommand(cleanTrashCmd)
	cleanCmd.AddCommand(cleanAllCmd)
	rootCmd.AddCommand(cleanCmd)
}
