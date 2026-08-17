package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
)

// SelectOption exibe um menu de seleção interativo e retorna o índice e valor selecionados
func SelectOption(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  12,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . | bold }}:",
			Active:   "▶ {{ . | cyan | bold }}",
			Inactive: "  {{ . }}",
			Selected: "✔ {{ . | green | bold }}",
			Help:     "Use as setas do teclado para navegar: ↓ ↑ → ←",
		},
	}

	index, result, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}
	return index, result, err
}

// PromptText solicita ao usuário que digite um texto (suporta mensagens longas sem bugs de re-render)
func PromptText(label string, defaultVal string) (string, error) {
	if defaultVal != "" {
		fmt.Printf("%s %s [%s]: ", InfoStyle.Render("?"), lipgloss.NewStyle().Bold(true).Render(label), MutedStyle.Render(defaultVal))
	} else {
		fmt.Printf("%s %s: ", InfoStyle.Render("?"), lipgloss.NewStyle().Bold(true).Render(label))
	}

	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	result = strings.TrimSpace(result)
	if result == "" && defaultVal != "" {
		result = defaultVal
	}

	return result, nil
}

// PromptInput solicita um texto simples do usuário sem valor padrão
func PromptInput(label string) string {
	res, _ := PromptText(label, "")
	return res
}

// ConfirmPrompt faz uma pergunta Sim/Não (S/n) ao usuário de forma limpa e confiável
func ConfirmPrompt(question string) bool {
	fmt.Printf("%s %s %s: ",
		WarningStyle.Render("❓"),
		lipgloss.NewStyle().Bold(true).Render(question),
		MutedStyle.Render("(S/n)"),
	)

	reader := bufio.NewReader(os.Stdin)
	res, err := reader.ReadString('\n')
	if err != nil {
		return true
	}

	res = strings.TrimSpace(strings.ToLower(res))
	if res == "" || res == "y" || res == "s" || res == "sim" || res == "yes" {
		return true
	}
	return false
}

// WaitForEnter exibe uma mensagem e bloqueia até o usuário pressionar Enter.
func WaitForEnter(message string) {
	if strings.TrimSpace(message) == "" {
		message = "Pressione Enter para voltar..."
	}
	fmt.Printf("%s ", MutedStyle.Render(message))
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}
