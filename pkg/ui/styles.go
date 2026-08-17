package ui

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Paleta de Cores
	PrimaryColor   = lipgloss.Color("#7D56F4") // Roxo moderno
	SecondaryColor = lipgloss.Color("#04B575") // Verde vibrante
	AccentColor    = lipgloss.Color("#FF7675") // Vermelho/Rosa suave
	WarningColor   = lipgloss.Color("#FDCB6E") // Amarelo/Laranja
	InfoColor      = lipgloss.Color("#0984E3") // Azul
	MutedColor     = lipgloss.Color("#B2BEC3") // Cinza

	// Estilos de Texto e Estrutura
	BreadcrumbRootStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(PrimaryColor).
				Padding(0, 1)

	BreadcrumbSubStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(SecondaryColor).
				Padding(0, 1)

	HeaderBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(0, 2).
			Align(lipgloss.Center)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(SecondaryColor)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor)

	WarningStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(WarningColor)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	SubtitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(SecondaryColor)

	MutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	// Status Badges
	BadgeRunning = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#2ECC71")).
			Padding(0, 1)

	BadgeStopped = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#E74C3C")).
			Padding(0, 1)
)

// ClearScreen limpa o terminal para uma transição visual fluida
func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

// GetRandomHeaderQuote retorna uma citação ou dica aleatória para o cabeçalho dos menus
func GetRandomHeaderQuote() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	quotes := []string{
		`💡 "O único modo de fazer um ótimo trabalho é amar o que você faz." - Steve Jobs`,
		`⚡ "Primeiro resolva o problema. Depois, escreva o código." - John Johnson`,
		`🚀 "Simplicidade é pré-requisito para confiabilidade." - Edsger W. Dijkstra`,
		`☕ "Convertendo café e boas ideias em soluções elegantes."`,
		`🎯 "Pequenas melhorias diárias geram grandes resultados no código."`,
		`🛠️ "O código mais limpo é aquele que você não precisa explicar."`,
		`🌌 "Pense duas vezes, code uma vez, teste sempre."`,
		`💻 "Qualidade é fazer certo quando ninguém está olhando." - Henry Ford`,
	}
	return quotes[r.Intn(len(quotes))]
}

// PrintBanner exibe o cabeçalho estilizado com suporte a Breadcrumbs e Frase Inspiradora
func PrintBanner(module string) {
	ClearScreen()

	var breadcrumbs string
	if module == "Menu Principal" || module == "" {
		breadcrumbs = BreadcrumbRootStyle.Render("🚀 TOMMY") + " ❯ " + BreadcrumbSubStyle.Render("MENU PRINCIPAL")
	} else {
		breadcrumbs = BreadcrumbRootStyle.Render("🚀 TOMMY") + " ❯ " + BreadcrumbSubStyle.Render(strings.ToUpper(module))
	}

	quote := GetRandomHeaderQuote()
	boxContent := fmt.Sprintf("%s\n%s\n%s",
		breadcrumbs,
		MutedStyle.Render("Ferramenta de Produtividade & Automação de Tarefas"),
		MutedStyle.Render(quote),
	)
	fmt.Println(HeaderBoxStyle.Render(boxContent))
	fmt.Println()
}

// PrintSuccess imprime uma mensagem de sucesso formatada
func PrintSuccess(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", SuccessStyle.Render("✔"), msg)
}

// PrintError imprime uma mensagem de erro formatada
func PrintError(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", ErrorStyle.Render("✖"), msg)
}

// PrintWarning imprime um aviso formatado
func PrintWarning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", WarningStyle.Render("⚠"), msg)
}

// PrintInfo imprime uma informação formatada
func PrintInfo(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("%s %s\n", InfoStyle.Render("ℹ"), msg)
}

func stripVS(s string) string {
	return strings.ReplaceAll(s, "\ufe0f", "")
}

// PrintExitMessage exibe um card de saída estilizado com frase dinâmica e aleatória baseada no horário do dia
func PrintExitMessage() {
	ClearScreen()
	hour := time.Now().Hour()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var phrases []string
	var greeting string

	switch {
	case hour >= 5 && hour < 12:
		greeting = "🌅 Bom dia!"
		phrases = []string{
			"Que os seus testes passem e o código compile de primeira hoje!",
			"Tenha uma manhã super produtiva e cheia de conquistas no terminal!",
			"Começando o dia com foco total. Nos vemos no próximo commit!",
			"Que a sua lógica de código brilhe nesta manhã! Até logo!",
		}
	case hour >= 12 && hour < 18:
		greeting = "🌤 Boa tarde!"
		phrases = []string{
			"Lembre-se de tomar uma água e fazer uma breve pausa!",
			"Ótimo trabalho nesta tarde. Deploys tranquilos por aqui!",
			"Produtividade em alta nesta tarde! Até a próxima execução!",
			"Mantendo o ritmo de desenvolvimento! Bom trabalho!",
		}
	case hour >= 18 && hour < 24:
		greeting = "🌙 Boa noite!"
		phrases = []string{
			"Hora de salvar as alterações e recarregar as energias!",
			"Mais um dia de código finalizado com sucesso. Até amanhã!",
			"Fechando o terminal por hoje. Aproveite o seu descanso nesta noite!",
			"Excelente dia de entregas! Tenha uma ótima noite.",
		}
	default: // Madrugada (00:00 - 04:59)
		greeting = "🦉 Fala, Dev Coruja!"
		phrases = []string{
			"Codando na madrugada? Lembre-se de ir descansar em breve!",
			"O silêncio da madrugada gera grandes ideias. Bom descanso!",
			"Guerreiro do código noturno! Que o bug tenha sido resolvido.",
			"Deslogando da máquina. Hora de dar um push no sono!",
		}
	}

	selectedPhrase := phrases[r.Intn(len(phrases))]

	headerText := stripVS(MutedStyle.Render("🚀 TOMMY | Produtividade, Automação & Eficiência no Terminal"))
	msgText := stripVS(SuccessStyle.Render(fmt.Sprintf("%s %s", greeting, selectedPhrase)))

	cardContent := fmt.Sprintf("%s\n\n%s", headerText, msgText)

	exitCardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SecondaryColor).
		Padding(0, 2)

	fmt.Println()
	fmt.Println(exitCardStyle.Render(cardContent))
}
