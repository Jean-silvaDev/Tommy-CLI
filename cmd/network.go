package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"tommy/pkg/network"
	"tommy/pkg/ui"

	"github.com/spf13/cobra"
)

var netCmd = &cobra.Command{
	Use:     "network",
	Aliases: []string{"net", "ip"},
	Short:   "Ferramentas e diagnósticos de rede (IP, Ping, Portas, DNS, Subnet)",
	Long:    `Subcomando para análise de conectividade de rede, verificação de IP local e público, teste de portas de desenvolvimento, consulta DNS e calculadoras de sub-rede.`,
	Run: func(cmd *cobra.Command, args []string) {
		runNetworkDashboard()
	},
}

func runNetworkDashboard() {
	for {
		ui.PrintBanner("Diagnóstico de Rede")

		ui.PrintInfo("Analisando conexões de rede locais e públicas...")
		summary := network.GetNetworkSummary()

		onlineBadge := ui.BadgeStopped.Render("OFFLINE")
		if summary.IsOnline {
			onlineBadge = ui.BadgeRunning.Render("ONLINE")
		}

		fmt.Println()
		fmt.Printf(" 🌐 IP Local: %s (%s) | %s\n",
			ui.SuccessStyle.Render(summary.LocalIP),
			summary.InterfaceName,
			onlineBadge)
		fmt.Printf(" 🔗 IP Público: %s | Provedor: %s\n\n",
			ui.InfoStyle.Render(summary.PublicIP),
			ui.MutedStyle.Render(summary.ISP))

		options := []string{
			" 1. 📊 Ver Resumo Completo de Conectividade (IPs, MAC, ISP)",
			" 2. 🔍 Listar Adaptadores & Placas de Rede",
			" 3. ⚡ Testar Latência & Conectividade (Ping)",
			" 4. 🔌 Testador de Porta TCP (Check Host:Porta)",
			" 5. 💥 Liberar Porta em Uso (Encerrar Processo na Porta)",
			" 6. 🚀 Scanner de Portas de Dev Locais (3000, 8080, 5432...)",
			" 7. 🌐 Consulta & Resolução DNS (NSLookup)",
			" 8. 🧹 Limpar Cache DNS do Sistema (Flush DNS)",
			" 9. 🧮 Calculadora de Sub-rede IP (CIDR)",
			"10. ⬅️  Voltar ao Menu Principal",
			"11. 🚪 Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Selecione uma ferramenta de rede", options)
		if err != nil || idx == 9 {
			break // Voltar ao Menu Principal
		}

		if idx == 10 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			displayFullSummary(summary)
		case 1:
			displayInterfacesTable()
		case 2:
			executePingTests(summary)
		case 3:
			executeCheckCustomPort()
		case 4:
			executeFreePort()
		case 5:
			executeDevPortScan()
		case 6:
			executeDNSLookup()
		case 7:
			executeFlushDNS()
		case 8:
			executeSubnetCalculator()
		}
	}
}

func displayFullSummary(s network.NetworkSummary) {
	fmt.Println(ui.SubtitleStyle.Render("Resumo Geral de Conectividade:"))
	fmt.Printf("  • IP Local (IPv4):     %s\n", ui.SuccessStyle.Render(s.LocalIP))
	fmt.Printf("  • IP Local (IPv6):     %s\n", s.LocalIPv6)
	fmt.Printf("  • Interface Ativa:     %s\n", s.InterfaceName)
	fmt.Printf("  • Endereço MAC:        %s\n", s.MACAddress)
	fmt.Printf("  • IP Público (WAN):    %s\n", ui.InfoStyle.Render(s.PublicIP))
	fmt.Printf("  • Provedor / Org:      %s\n", s.ISP)
	fmt.Printf("  • Localização Aprox:   %s\n", s.Location)
	fmt.Println()
	ui.WaitForEnter("")
}

func displayInterfacesTable() {
	ifaces, err := network.GetAllInterfaces()
	if err != nil {
		ui.PrintError("%v", err)
		ui.WaitForEnter("")
		return
	}

	headers := []string{"INTERFACE", "STATUS", "ENDEREÇO MAC", "IPv4", "IPv6"}
	var rows [][]string

	for _, iface := range ifaces {
		statusBadge := ui.BadgeStopped.Render("DOWN")
		if iface.Status == "UP" {
			statusBadge = ui.BadgeRunning.Render("UP")
		}

		ipv4 := iface.IPv4
		if ipv4 == "" {
			ipv4 = "-"
		}
		ipv6 := iface.IPv6
		if ipv6 == "" {
			ipv6 = "-"
		}

		rows = append(rows, []string{
			iface.Name,
			statusBadge,
			iface.MAC,
			ipv4,
			ipv6,
		})
	}

	ui.RenderTable(headers, rows)
	ui.WaitForEnter("")
}

func executePingTests(s network.NetworkSummary) {
	ui.PrintInfo("Testando latência de resposta TCP para destinos chave...")
	fmt.Println()

	targets := []struct {
		Name string
		Host string
		Port int
	}{
		{"Google DNS", "8.8.8.8", 53},
		{"Cloudflare DNS", "1.1.1.1", 53},
		{"Servidor Web (Google.com)", "google.com", 443},
	}

	for _, t := range targets {
		duration, err := network.TestPing(t.Host, t.Port)
		if err != nil {
			fmt.Printf("  ❌ %-25s (%s): %s\n", t.Name, t.Host, ui.ErrorStyle.Render("TIMEOUT / FALHA"))
		} else {
			ms := duration.Milliseconds()
			statusStr := ui.SuccessStyle.Render(fmt.Sprintf("%d ms (Ótima)", ms))
			if ms > 100 {
				statusStr = ui.WarningStyle.Render(fmt.Sprintf("%d ms (Lenta)", ms))
			}
			fmt.Printf("  ✔ %-25s (%s): %s\n", t.Name, t.Host, statusStr)
		}
	}

	fmt.Println()
	ui.WaitForEnter("")
}

func executeCheckCustomPort() {
	host, err := ui.PromptText("Informe o Host/IP para testar (ex: localhost, 192.168.1.1, google.com)", "localhost")
	if err != nil || host == "" {
		return
	}

	portStr, errP := ui.PromptText("Informe a Porta TCP (ex: 80, 443, 8080, 5432)", "8080")
	if errP != nil || portStr == "" {
		return
	}

	port, errConv := strconv.Atoi(portStr)
	if errConv != nil {
		ui.PrintError("Porta inválida!")
		ui.WaitForEnter("")
		return
	}

	ui.PrintInfo("Testando conexão TCP em %s:%d...", host, port)
	isOpen := network.CheckPort(host, port)
	if isOpen {
		ui.PrintSuccess("A porta %d em '%s' está ABERTA e respondendo!", port, host)
	} else {
		ui.PrintError("A porta %d em '%s' está FECHADA ou inacessível (Timeout).", port, host)
	}

	ui.WaitForEnter("")
}

func executeFreePort() {
	portStr, err := ui.PromptText("Informe a porta que deseja liberar (ex: 3000, 8080, 5432)", "3000")
	if err != nil || strings.TrimSpace(portStr) == "" {
		return
	}

	port, errConv := strconv.Atoi(strings.TrimSpace(portStr))
	if errConv != nil || port <= 0 {
		ui.PrintError("Número de porta inválido!")
		ui.WaitForEnter("")
		return
	}

	ui.PrintInfo("Buscando processo ocupando a porta %d...", port)
	procInfo, errProc := network.GetProcessOnPort(port)
	if errProc != nil {
		ui.PrintWarning("%v", errProc)
		ui.WaitForEnter("")
		return
	}

	fmt.Println()
	ui.PrintWarning("Processo Identificado na Porta %d:\n  • Executável: %s\n  • PID:        %d",
		port, procInfo.ProcessName, procInfo.PID)
	fmt.Println()

	if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente ENCERRAR '%s' (PID %d) para liberar a porta %d?", procInfo.ProcessName, procInfo.PID, port)) {
		ui.PrintInfo("Operação cancelada.")
		ui.WaitForEnter("")
		return
	}

	ui.PrintInfo("Encerrando processo %s (PID %d)...", procInfo.ProcessName, procInfo.PID)
	killedInfo, errKill := network.KillProcessOnPort(port)
	if errKill != nil {
		ui.PrintError("%v", errKill)
	} else {
		ui.PrintSuccess("Porta %d liberada com sucesso! O processo '%s' (PID %d) foi encerrado.", port, killedInfo.ProcessName, killedInfo.PID)
	}

	ui.WaitForEnter("")
}

func executeDevPortScan() {
	ui.PrintInfo("Verificando portas de serviços de desenvolvimento local (127.0.0.1)...")
	fmt.Println()

	results := network.ScanDevPorts()

	headers := []string{"PORTA", "SERVIÇO DEV", "STATUS DA PORTA"}
	var rows [][]string
	var listeningCount int

	for _, r := range results {
		statusBadge := ui.MutedStyle.Render("LIVRE")
		if r.IsListening {
			statusBadge = ui.SuccessStyle.Render("✔ EM USO (OUVINDO)")
			listeningCount++
		}

		rows = append(rows, []string{
			fmt.Sprintf("%d", r.Port),
			r.Service,
			statusBadge,
		})
	}

	ui.RenderTable(headers, rows)

	if listeningCount > 0 {
		if ui.ConfirmPrompt(fmt.Sprintf("Foram encontradas %d porta(s) em uso. Deseja liberar alguma porta agora?", listeningCount)) {
			executeFreePort()
			return
		}
	}

	ui.WaitForEnter("")
}

func executeDNSLookup() {
	domain, err := ui.PromptText("Informe o domínio para consulta DNS (ex: google.com, github.com)", "google.com")
	if err != nil || domain == "" {
		return
	}

	ui.PrintInfo("Consultando registros DNS para '%s'...", domain)
	res, errL := network.LookupDNS(domain)
	if errL != nil {
		ui.PrintError("%v", errL)
		ui.WaitForEnter("")
		return
	}

	fmt.Println()
	fmt.Println(ui.SubtitleStyle.Render(fmt.Sprintf("Resultado DNS para '%s':", res.Domain)))
	if res.CNAME != "" {
		fmt.Printf("  • CNAME: %s\n", res.CNAME)
	}
	for _, ip := range res.IPs {
		fmt.Printf("  • Registro IP: %s\n", ui.SuccessStyle.Render(ip))
	}

	fmt.Println()
	ui.WaitForEnter("")
}

func executeFlushDNS() {
	if !ui.ConfirmPrompt("Deseja realmente limpar o cache de DNS do sistema operacional?") {
		ui.PrintInfo("Operação cancelada.")
		return
	}

	ui.PrintInfo("Executando limpeza de cache DNS...")
	err := network.FlushDNS()
	if err != nil {
		ui.PrintError("%v", err)
	} else {
		ui.PrintSuccess("Cache de DNS do sistema limpo com sucesso!")
	}

	ui.WaitForEnter("")
}

func executeSubnetCalculator() {
	cidr, err := ui.PromptText("Informe o IP com máscara no formato CIDR (ex: 192.168.1.0/24)", "192.168.1.0/24")
	if err != nil || cidr == "" {
		return
	}

	info, errCalc := network.CalculateSubnet(cidr)
	if errCalc != nil {
		ui.PrintError("%v", errCalc)
		ui.WaitForEnter("")
		return
	}

	fmt.Println()
	fmt.Println(ui.SubtitleStyle.Render(fmt.Sprintf("Cálculo de Sub-rede para %s:", info.CIDR)))
	fmt.Printf("  • Máscara de Sub-rede:   %s\n", info.Netmask)
	fmt.Printf("  • IP da Rede (Network):  %s\n", info.NetworkIP)
	fmt.Printf("  • IP de Broadcast:       %s\n", info.BroadcastIP)
	fmt.Printf("  • Primeiro IP Utilizável: %s\n", ui.SuccessStyle.Render(info.FirstUsableIP))
	fmt.Printf("  • Último IP Utilizável:   %s\n", ui.SuccessStyle.Render(info.LastUsableIP))
	fmt.Printf("  • Hosts Válidos (Úteis):  %s\n", ui.InfoStyle.Render(fmt.Sprintf("%d", info.TotalHosts)))

	fmt.Println()
	ui.WaitForEnter("")
}

func init() {
	rootCmd.AddCommand(netCmd)
}
