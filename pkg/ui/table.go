package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderTable renderiza uma tabela formatada com alinhamento preciso (ANSI-aware) no terminal
func RenderTable(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Sanitizar quebras de linha e limitar tamanhos para evitar estouro de tabela
	sanitRows := make([][]string, len(rows))
	for rIdx, row := range rows {
		sanitRow := make([]string, len(row))
		for cIdx, val := range row {
			cleanVal := strings.ReplaceAll(val, "\r\n", " ")
			cleanVal = strings.ReplaceAll(cleanVal, "\n", " ")
			if cIdx == len(row)-1 && lipgloss.Width(cleanVal) > 75 {
				cleanVal = truncateString(cleanVal, 75)
			}
			sanitRow[cIdx] = cleanVal
		}
		sanitRows[rIdx] = sanitRow
	}

	// Calcular largura visível real de cada coluna (ignorando códigos ANSI)
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = lipgloss.Width(h)
	}

	for _, row := range sanitRows {
		for i, val := range row {
			if i < len(colWidths) {
				w := lipgloss.Width(val)
				if w > colWidths[i] {
					colWidths[i] = w
				}
			}
		}
	}

	// Adicionar margem entre colunas para ótima legibilidade
	for i := range colWidths {
		colWidths[i] += 2
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(PrimaryColor).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECF0F1"))

	// Renderizar cabeçalho
	var headerLine strings.Builder
	for i, h := range headers {
		visibleLen := lipgloss.Width(h)
		padLen := colWidths[i] - visibleLen
		if padLen < 0 {
			padLen = 0
		}
		cellStr := fmt.Sprintf("%s%s", h, strings.Repeat(" ", padLen))
		headerLine.WriteString(headerStyle.Render(cellStr))
	}
	fmt.Println(headerLine.String())

	// Linha divisória
	totalW := 0
	for _, w := range colWidths {
		totalW += w + 2 // +2 pelo padding (0,1) do headerStyle
	}
	fmt.Println(MutedStyle.Render(strings.Repeat("─", totalW)))

	// Renderizar linhas
	for _, row := range sanitRows {
		var rowLine strings.Builder
		for i, val := range row {
			if i < len(colWidths) {
				visibleLen := lipgloss.Width(val)
				padLen := colWidths[i] - visibleLen
				if padLen < 0 {
					padLen = 0
				}
				cellStr := fmt.Sprintf(" %s%s ", val, strings.Repeat(" ", padLen))
				rowLine.WriteString(cellStyle.Render(cellStr))
			}
		}
		fmt.Println(rowLine.String())
	}
	fmt.Println()
}

// GitLogItem representa uma entrada de commit formatada
type GitLogItem struct {
	Hash    string
	Author  string
	Date    string
	Subject string
}

// RenderGitLog renderiza o histórico de commits em uma tabela limpa com cabeçalhos roxos estilizados
func RenderGitLog(logs []GitLogItem) {
	if len(logs) == 0 {
		return
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(PrimaryColor).
		Padding(0, 1)

	hashStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F1C40F")) // Amarelo Dourado

	authorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#DFE6E9")) // Branco Suave

	dateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00CEC9")) // Ciano / Azul

	subjectStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")) // Branco Brilhante

	maxAuthorW := 14
	maxDateW := 14
	for _, l := range logs {
		if lipgloss.Width(l.Author) > maxAuthorW {
			maxAuthorW = lipgloss.Width(l.Author)
		}
		if lipgloss.Width(l.Date) > maxDateW {
			maxDateW = lipgloss.Width(l.Date)
		}
	}
	if maxAuthorW > 22 {
		maxAuthorW = 22
	}
	if maxDateW > 20 {
		maxDateW = 20
	}

	// Renderizar cabeçalho roxo alinhado das colunas
	hashHeader := fmt.Sprintf("%-10s", "HASH")
	authorHeader := fmt.Sprintf("%-*s", maxAuthorW+2, "AUTOR")
	dateHeader := fmt.Sprintf("%-*s", maxDateW+2, "DATA")
	msgHeader := fmt.Sprintf("%-80s", "MENSAGEM / ASSUNTO")

	headerStr := fmt.Sprintf(" %s %s %s %s",
		headerStyle.Render(hashHeader),
		headerStyle.Render(authorHeader),
		headerStyle.Render(dateHeader),
		headerStyle.Render(msgHeader),
	)

	fmt.Println()
	fmt.Println(headerStr)

	totalW := 10 + (maxAuthorW + 2) + (maxDateW + 2) + 80 + 10
	fmt.Println(MutedStyle.Render(strings.Repeat("─", totalW)))

	// Renderizar cada commit
	for _, l := range logs {
		cleanSubject := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(l.Subject, "\r\n", " "), "\n", " "))
		if lipgloss.Width(cleanSubject) > 80 {
			cleanSubject = truncateString(cleanSubject, 80)
		}

		hashPadded := fmt.Sprintf("%-10s", l.Hash)
		authorPadded := fmt.Sprintf("%-*s", maxAuthorW+2, truncateString(l.Author, maxAuthorW))
		datePadded := fmt.Sprintf("%-*s", maxDateW+2, truncateString(l.Date, maxDateW))

		line := fmt.Sprintf("  %s  %s  %s  %s",
			hashStyle.Render(hashPadded),
			authorStyle.Render(authorPadded),
			dateStyle.Render(datePadded),
			subjectStyle.Render(cleanSubject),
		)
		fmt.Println(line)
	}
	fmt.Println()
}

func truncateString(s string, maxLen int) string {
	if lipgloss.Width(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

