package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"tommy/pkg/system"
	"tommy/pkg/ui"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var dryRunFlag bool

var systemCmd = &cobra.Command{
	Use:     "system",
	Aliases: []string{"sys"},
	Short:   "Gerenciamento e monitoramento de recursos do sistema (RAM, Discos, Processos, Apps)",
	Long:    `Subcomando modular para monitorar e gerenciar a saúde, consumo de memória RAM, uso de discos, processos ativos, aplicativos e rotinas de limpeza do computador.`,
	Run: func(cmd *cobra.Command, args []string) {
		runSystemDashboard()
	},
}

var systemMemoryCmd = &cobra.Command{
	Use:   "memory",
	Short: "Exibe o consumo detalhado da memória RAM e principais processos",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Gerenciador de Memória RAM")
		displayMemoryMenu(system.NewSystemProvider())
	},
}

var systemProcessesCmd = &cobra.Command{
	Use:     "processes",
	Aliases: []string{"process", "ps"},
	Short:   "Gerenciador interativo de processos ativos",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Gerenciador de Processos")
		displayProcessMenu(system.NewSystemProvider())
	},
}

var systemDiskCmd = &cobra.Command{
	Use:     "disk",
	Aliases: []string{"disks", "drive"},
	Short:   "Exibe volumes de armazenamento e saúde dos discos",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Gerenciador de Discos")
		displayDisksMenu(system.NewSystemProvider())
	},
}

var systemAppsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Gerenciador e instalador/desinstalador de aplicativos",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Aplicativos Instalados")
		displayAppsMenu(system.NewSystemProvider())
	},
}

var systemAppsUninstallCmd = &cobra.Command{
	Use:   "uninstall [nome_do_app]",
	Short: "Desinstala um aplicativo pelo nome ou ID registrado",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Desinstalação de Aplicativo")
		provider := system.NewSystemProvider()
		apps, err := provider.ListApplications()
		if err != nil {
			ui.PrintError("%v", err)
			return
		}

		query := ""
		if len(args) > 0 {
			query = args[0]
		} else {
			query, _ = ui.PromptText("Informe o nome do aplicativo para desinstalar", "")
		}

		if strings.TrimSpace(query) == "" {
			return
		}

		filtered := system.FilterApps(apps, query)
		if len(filtered) == 0 {
			ui.PrintWarning("Nenhum aplicativo encontrado correspondente a '%s'.", query)
			return
		}

		selectedApp := filtered[0]
		executeUninstallAppWithConfirmation(provider, selectedApp)
	},
}

var systemAnalyzeCmd = &cobra.Command{
	Use:   "analyze [caminho]",
	Short: "Analisa a ocupação de espaço em diretórios",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Analisador de Espaço em Disco")
		targetPath := "C:\\"
		if len(args) > 0 {
			targetPath = args[0]
		}
		executeStorageAnalysis(system.NewSystemProvider(), targetPath)
	},
}

var systemLargestCmd = &cobra.Command{
	Use:   "largest [caminho]",
	Short: "Exibe os maiores arquivos em um diretório",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Maiores Arquivos no Disco")
		targetPath := "C:\\"
		if len(args) > 0 {
			targetPath = args[0]
		}
		executeLargestFiles(system.NewSystemProvider(), targetPath)
	},
}

var systemMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitora recursos de sistema (CPU, RAM, Discos) em tempo real",
	Run: func(cmd *cobra.Command, args []string) {
		displayRealtimeMonitorMenu(system.NewSystemProvider())
	},
}

var systemCleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Central de limpeza integrada (%TEMP%, NuGet, Docker, Lixeira)",
	Run: func(cmd *cobra.Command, args []string) {
		executeCleanAllWithConfirmation()
	},
}

func runSystemDashboard() {
	provider := system.NewSystemProvider()

	for {
		ui.PrintBanner("Gerenciador do Sistema")

		ui.PrintInfo("Coletando telemetria e estatísticas do sistema...")
		summary, err := provider.GetSystemSummary()

		if err != nil {
			ui.PrintError("Falha ao obter dados do sistema: %v", err)
		} else {
			renderSystemSummaryBox(summary)
		}

		options := []string{
			"1. 🧠  Memória RAM (Uso & Processos Consumidores)",
			"2. ⚙️  Gerenciador de Processos (Finalização Segura)",
			"3. 💾  Discos & Volumes (Estatísticas & Saúde)",
			"4. 📦  Aplicativos Instalados (Busca & Desinstalação)",
			"5. 📁  Analisador de Espaço em Disco",
			"6. 📦  Maiores Arquivos em Disco",
			"7. 🧹  Central de Limpeza Integrada do Sistema",
			"8. 📊  Monitoramento em Tempo Real",
			"9. ⬅️  Voltar ao Menu Principal",
			"10. 🚪  Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Selecione um módulo do System Manager", options)
		if err != nil || idx == 8 {
			break
		}

		if idx == 9 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			displayMemoryMenu(provider)
		case 1:
			displayProcessMenu(provider)
		case 2:
			displayDisksMenu(provider)
		case 3:
			displayAppsMenu(provider)
		case 4:
			target, _ := ui.PromptText("Informe o caminho para análise (ex: C:\\, D:\\Projetos)", "C:\\")
			if strings.TrimSpace(target) != "" {
				executeStorageAnalysis(provider, target)
			}
		case 5:
			target, _ := ui.PromptText("Informe o caminho para varrer maiores arquivos", "C:\\")
			if strings.TrimSpace(target) != "" {
				executeLargestFiles(provider, target)
			}
		case 6:
			executeCleanAllWithConfirmation()
		case 7:
			displayRealtimeMonitorMenu(provider)
		}
	}
}

func renderSystemSummaryBox(s system.SystemSummary) {
	ramBar := system.RenderProgressBar(s.Memory.UsedPercent, 20)
	ramDetail := fmt.Sprintf("%s / %s", system.FormatBytes(s.Memory.UsedBytes), system.FormatBytes(s.Memory.TotalBytes))

	var volumeLines []string
	for _, v := range s.Volumes {
		vBar := system.RenderProgressBar(v.UsedPercent, 18)
		vDetail := fmt.Sprintf("%s (%s) %s | %s / %s livre",
			v.Device, v.DiskType, vBar, system.FormatBytes(v.FreeBytes), system.FormatBytes(v.TotalBytes))
		volumeLines = append(volumeLines, " 💾 "+vDetail)
	}

	cpuBar := system.RenderProgressBar(s.CPU.UsagePercent, 20)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ui.PrimaryColor).
		Padding(1, 2)

	content := fmt.Sprintf("🧠 RAM: %s  (%s)\n⚙️ CPU: %s\n", ramBar, ramDetail, cpuBar)
	if len(volumeLines) > 0 {
		content += "\n"
		for i := 0; i < len(volumeLines); i++ {
			content += volumeLines[i] + "\n"
		}
	}

	if len(s.Alerts) > 0 {
		content += "\n"
		for _, alert := range s.Alerts {
			content += alert.Message + "\n"
		}
	}

	fmt.Println(boxStyle.Render(content))
	fmt.Println()
}

func displayMemoryMenu(p system.SystemProvider) {
	info, err := p.GetMemoryInfo()
	if err != nil {
		ui.PrintError("%v", err)
		ui.WaitForEnter("")
		return
	}

	fmt.Println(ui.SubtitleStyle.Render("Resumo da Memória RAM:"))
	fmt.Printf("  • Memória Total:       %s\n", system.FormatBytes(info.TotalBytes))
	fmt.Printf("  • Memória Em Uso:      %s\n", ui.WarningStyle.Render(system.FormatBytes(info.UsedBytes)))
	fmt.Printf("  • Memória Disponível:  %s\n", ui.SuccessStyle.Render(system.FormatBytes(info.AvailableBytes)))
	fmt.Printf("  • Percentual Utilizado: %s\n\n", system.RenderProgressBar(info.UsedPercent, 25))

	procs, errProc := p.ListProcesses()
	if errProc == nil {
		sortedProcs := system.SortProcessesByMemory(procs)
		limit := 10
		if len(sortedProcs) < limit {
			limit = len(sortedProcs)
		}

		headers := []string{"PID", "PROCESSO", "MEMÓRIA", "% RAM", "STATUS"}
		var rows [][]string

		for i := 0; i < limit; i++ {
			pr := sortedProcs[i]
			statusStr := pr.Status
			if pr.IsProtected {
				statusStr = "Protegido (SO)"
			}
			rows = append(rows, []string{
				fmt.Sprintf("%d", pr.PID),
				pr.Name,
				system.FormatBytes(pr.MemoryBytes),
				fmt.Sprintf("%.1f%%", pr.MemoryPercent),
				statusStr,
			})
		}

		fmt.Println(ui.SubtitleStyle.Render("Principais Processos Consumidores de Memória:"))
		ui.RenderTable(headers, rows)
	}

	ui.WaitForEnter("Pressione Enter para voltar...")
}

func displayProcessMenu(p system.SystemProvider) {
	for {
		ui.PrintBanner("Gerenciador de Processos")
		procs, err := p.ListProcesses()
		if err != nil {
			ui.PrintError("%v", err)
			ui.WaitForEnter("")
			return
		}

		sorted := system.SortProcessesByMemory(procs)

		options := []string{
			"1. 📋  Listar Processos Ativos (Ordenado por Memória)",
			"2. 🔍  Pesquisar Processo por Nome / PID",
			"3. ⛔  Encerrar Processo (Kill Seguro)",
			"4. ⬅️  Voltar ao System Manager",
		}

		idx, _, errOpt := ui.SelectOption("Opções de gerenciamento de processos", options)
		if errOpt != nil || idx == 3 {
			break
		}

		switch idx {
		case 0:
			renderProcessTable(sorted, 15)
			ui.WaitForEnter("")
		case 1:
			query, _ := ui.PromptText("Informe o nome do processo ou PID para buscar", "")
			if query != "" {
				filtered := system.FilterProcesses(sorted, query)
				if len(filtered) == 0 {
					ui.PrintWarning("Nenhum processo correspondente a '%s'.", query)
				} else {
					renderProcessTable(filtered, 20)
				}
				ui.WaitForEnter("")
			}
		case 2:
			pidStr, _ := ui.PromptText("Informe o PID do processo que deseja encerrar", "")
			pid, errConv := strconv.Atoi(strings.TrimSpace(pidStr))
			if errConv != nil || pid <= 0 {
				ui.PrintError("PID inválido!")
				ui.WaitForEnter("")
				continue
			}

			executeKillProcessWithConfirmation(p, pid)
		}
	}
}

func renderProcessTable(procs []system.ProcessInfo, limit int) {
	if len(procs) < limit {
		limit = len(procs)
	}

	headers := []string{"PID", "PROCESSO", "THREADS", "MEMÓRIA", "% RAM", "STATUS"}
	var rows [][]string

	for i := 0; i < limit; i++ {
		pr := procs[i]
		statusStr := pr.Status
		if pr.IsProtected {
			statusStr = "Protegido (SO)"
		}

		rows = append(rows, []string{
			fmt.Sprintf("%d", pr.PID),
			pr.Name,
			fmt.Sprintf("%d", pr.Threads),
			system.FormatBytes(pr.MemoryBytes),
			fmt.Sprintf("%.1f%%", pr.MemoryPercent),
			statusStr,
		})
	}

	ui.RenderTable(headers, rows)
}

func executeKillProcessWithConfirmation(p system.SystemProvider, pid int) {
	procs, err := p.ListProcesses()
	if err != nil {
		ui.PrintError("%v", err)
		ui.WaitForEnter("")
		return
	}

	var targetProc *system.ProcessInfo
	for _, pr := range procs {
		if pr.PID == pid {
			targetProc = &pr
			break
		}
	}

	if targetProc == nil {
		ui.PrintError("Processo com PID %d não foi localizado no sistema.", pid)
		ui.WaitForEnter("")
		return
	}

	if targetProc.IsProtected {
		ui.PrintError("⚠️ O processo '%s' (PID %d) é vital para o sistema operacional Windows e está protegido contra encerramento.", targetProc.Name, pid)
		ui.WaitForEnter("")
		return
	}

	ui.PrintWarning("ATENÇÃO: Você está prestes a encerrar o processo '%s' (PID %d) ocupando %s de RAM.",
		targetProc.Name, targetProc.PID, system.FormatBytes(targetProc.MemoryBytes))
	fmt.Println()

	if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente ENCERRAR o processo '%s' (PID %d)?", targetProc.Name, targetProc.PID)) {
		ui.PrintInfo("Operação de encerramento de processo cancelada pelo usuário.")
		ui.WaitForEnter("")
		return
	}

	errKill := p.KillProcess(pid)
	if errKill != nil {
		ui.PrintError("Erro ao encerrar o processo: %v", errKill)
	} else {
		ui.PrintSuccess("Processo '%s' (PID %d) encerrado com sucesso!", targetProc.Name, pid)
	}

	ui.WaitForEnter("")
}

func displayDisksMenu(p system.SystemProvider) {
	vols, err := p.ListVolumes()
	if err != nil {
		ui.PrintError("%v", err)
		ui.WaitForEnter("")
		return
	}

	headers := []string{"DISCO", "RÓTULO", "TIPO", "FS", "TOTAL", "USADO", "LIVRE", "USO", "SAÚDE"}
	var rows [][]string

	for _, v := range vols {
		usageBadge := system.RenderProgressBar(v.UsedPercent, 12)
		rows = append(rows, []string{
			v.Device,
			v.Label,
			v.DiskType,
			v.FileSystem,
			system.FormatBytes(v.TotalBytes),
			system.FormatBytes(v.UsedBytes),
			system.FormatBytes(v.FreeBytes),
			usageBadge,
			v.HealthStatus,
		})
	}

	ui.RenderTable(headers, rows)
	ui.WaitForEnter("")
}

func displayAppsMenu(p system.SystemProvider) {
	for {
		ui.PrintBanner("Aplicativos Instalados")
		ui.PrintInfo("Lendo registros de softwares instalados no Windows...")
		apps, err := p.ListApplications()
		if err != nil {
			ui.PrintError("%v", err)
			ui.WaitForEnter("")
			return
		}

		options := []string{
			"1. 📋  Listar Aplicativos Instalados",
			"2. 🔍  Pesquisar Aplicativo por Nome / Desenvolvedor",
			"3. 📊  Ordenar Aplicativos por Tamanho em Disco",
			"4. 💾  Filtrar Programas Ocupando Espaço no SSD/Disco",
			"5. 🗑️  Desinstalar Aplicativo",
			"6. ⬅️  Voltar ao System Manager",
		}

		idx, _, errOpt := ui.SelectOption("Gerenciamento de aplicativos", options)
		if errOpt != nil || idx == 5 {
			break
		}

		switch idx {
		case 0:
			renderAppsTable(apps, 15)
			ui.WaitForEnter("")
		case 1:
			query, _ := ui.PromptText("Informe o termo para busca", "")
			if query != "" {
				filtered := system.FilterApps(apps, query)
				renderAppsTable(filtered, 20)
				ui.WaitForEnter("")
			}
		case 2:
			sorted := system.SortAppsBySize(apps)
			renderAppsTable(sorted, 15)
			ui.WaitForEnter("")
		case 3:
			drive, _ := ui.PromptText("Informe a letra do disco (ex: C:, D:)", "C:")
			filtered := system.FilterAppsByDrive(apps, drive)
			renderAppsTable(filtered, 15)
			ui.WaitForEnter("")
		case 4:
			appName, _ := ui.PromptText("Informe o nome do aplicativo que deseja desinstalar", "")
			if appName != "" {
				filtered := system.FilterApps(apps, appName)
				if len(filtered) == 0 {
					ui.PrintWarning("Nenhum aplicativo encontrado correspondente a '%s'.", appName)
					ui.WaitForEnter("")
				} else {
					executeUninstallAppWithConfirmation(p, filtered[0])
				}
			}
		}
	}
}

func renderAppsTable(apps []system.AppInfo, limit int) {
	if len(apps) == 0 {
		ui.PrintWarning("Nenhum aplicativo encontrado.")
		return
	}
	if len(apps) < limit {
		limit = len(apps)
	}

	headers := []string{"NOME", "VERSÃO", "FABRICANTE", "TAMANHO", "ORIGEM"}
	var rows [][]string

	for i := 0; i < limit; i++ {
		app := apps[i]
		sizeStr := "Não informado"
		if app.SizeBytes > 0 {
			sizeStr = system.FormatBytes(app.SizeBytes)
		}

		rows = append(rows, []string{
			app.Name,
			app.Version,
			app.Publisher,
			sizeStr,
			app.Source,
		})
	}

	ui.RenderTable(headers, rows)
}

func executeUninstallAppWithConfirmation(p system.SystemProvider, app system.AppInfo) {
	ui.PrintWarning("⚠️ ATENÇÃO: Você está prestes a iniciar a desinstalação do aplicativo:")
	fmt.Printf("  • Nome:       %s\n", ui.SubtitleStyle.Render(app.Name))
	fmt.Printf("  • Versão:     %s\n", app.Version)
	fmt.Printf("  • Fabricante: %s\n", app.Publisher)
	if app.SizeBytes > 0 {
		fmt.Printf("  • Tamanho:    %s\n", system.FormatBytes(app.SizeBytes))
	}
	fmt.Println()

	if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente DESINSTALAR o aplicativo '%s'?", app.Name)) {
		ui.PrintInfo("Operação de desinstalação cancelada.")
		ui.WaitForEnter("")
		return
	}

	ui.PrintInfo("Iniciando rotina nativa de desinstalação...")
	err := p.UninstallApplication(app)
	if err != nil {
		ui.PrintError("Falha na desinstalação: %v", err)
	} else {
		ui.PrintSuccess("Desinstalação concluída com sucesso!")
	}

	ui.WaitForEnter("")
}

func executeStorageAnalysis(p system.SystemProvider, path string) {
	ui.PrintInfo("Analisando estrutura de armazenamento em '%s' (varredura profunda sem symlinks)...", path)
	node, err := p.AnalyzePath(path)
	if err != nil {
		ui.PrintError("Erro na análise de armazenamento: %v", err)
		ui.WaitForEnter("")
		return
	}

	fmt.Println()
	fmt.Println(ui.SubtitleStyle.Render(fmt.Sprintf("Resultado do Analisador de Espaço para '%s':", node.Path)))
	fmt.Printf("  • Espaço Total Mapeado: %s\n\n", ui.SuccessStyle.Render(system.FormatBytes(node.SizeBytes)))

	if len(node.Children) > 0 {
		headers := []string{"NOME DA PASTA / ARQUIVO", "TIPO", "TAMANHO EM DISCO"}
		var rows [][]string

		limit := 20
		if len(node.Children) < limit {
			limit = len(node.Children)
		}

		for i := 0; i < limit; i++ {
			child := node.Children[i]
			typeStr := "Arquivo"
			if child.IsDir {
				typeStr = "Diretório"
			}
			rows = append(rows, []string{
				child.Name,
				typeStr,
				system.FormatBytes(child.SizeBytes),
			})
		}
		ui.RenderTable(headers, rows)
	}

	ui.WaitForEnter("")
}

func executeLargestFiles(p system.SystemProvider, path string) {
	ui.PrintInfo("Identificando os maiores arquivos contidos em '%s'...", path)
	files, err := p.GetLargestFiles(path, 15)
	if err != nil {
		ui.PrintError("%v", err)
		ui.WaitForEnter("")
		return
	}

	if len(files) == 0 {
		ui.PrintWarning("Nenhum arquivo encontrado no caminho especificado.")
		ui.WaitForEnter("")
		return
	}

	headers := []string{"TAMANHO", "NOME DO ARQUIVO", "CAMINHO COMPLETO"}
	var rows [][]string

	for _, f := range files {
		rows = append(rows, []string{
			system.FormatBytes(f.SizeBytes),
			f.Name,
			f.Path,
		})
	}

	ui.RenderTable(headers, rows)

	if dryRunFlag {
		ui.PrintInfo("[MODO DRY-RUN] Nenhuma exclusão será realizada.")
	}

	ui.WaitForEnter("")
}

func displayRealtimeMonitorMenu(p system.SystemProvider) {
	ui.PrintInfo("Iniciando monitoramento em tempo real... (Pressione Ctrl+C para sair)")
	time.Sleep(1 * time.Second)

	for i := 0; i < 5; i++ {
		ui.ClearScreen()
		ui.PrintBanner("Monitor de Recursos em Tempo Real")

		summary, err := p.GetSystemSummary()
		if err == nil {
			renderSystemSummaryBox(summary)
		}
		time.Sleep(2 * time.Second)
	}

	ui.PrintSuccess("Sessão de monitoramento finalizada.")
	ui.WaitForEnter("")
}

func init() {
	systemLargestCmd.Flags().BoolVarP(&dryRunFlag, "dry-run", "d", false, "Simula operações sem alterar arquivos")

	systemAppsCmd.AddCommand(systemAppsUninstallCmd)

	systemCmd.AddCommand(systemMemoryCmd)
	systemCmd.AddCommand(systemProcessesCmd)
	systemCmd.AddCommand(systemDiskCmd)
	systemCmd.AddCommand(systemAppsCmd)
	systemCmd.AddCommand(systemAnalyzeCmd)
	systemCmd.AddCommand(systemLargestCmd)
	systemCmd.AddCommand(systemMonitorCmd)
	systemCmd.AddCommand(systemCleanupCmd)

	rootCmd.AddCommand(systemCmd)
}
