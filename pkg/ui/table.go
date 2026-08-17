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

	// Sanitizar quebras de linha nos dados
	sanitRows := make([][]string, len(rows))
	for rIdx, row := range rows {
		sanitRow := make([]string, len(row))
		for cIdx, val := range row {
			cleanVal := strings.ReplaceAll(val, "\r\n", " ")
			cleanVal = strings.ReplaceAll(cleanVal, "\n", " ")
			sanitRow[cIdx] = cleanVal
		}
		sanitRows[rIdx] = sanitRow
	}

	// Calcular a largura visível real de cada coluna baseada em cabeçalhos e valores
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

	// Aplicar limites máximos por coluna para evitar quebras de linha desastrosas no terminal
	for i, h := range headers {
		hUpper := strings.ToUpper(h)
		maxCap := 35 // limite padrão

		switch {
		case hUpper == "ID":
			maxCap = 14
		case hUpper == "ESTADO" || hUpper == "STATUS":
			if len(headers) >= 6 {
				maxCap = 18
			} else {
				maxCap = 25
			}
		case hUpper == "STACK" || hUpper == "PROJETO":
			maxCap = 18
		case hUpper == "NOME":
			if len(headers) >= 6 {
				maxCap = 22
			} else {
				maxCap = 30
			}
		case hUpper == "IMAGEM":
			if len(headers) >= 6 {
				maxCap = 28
			} else {
				maxCap = 35
			}
		case hUpper == "PORTA" || hUpper == "TIPO" || hUpper == "PROTO" || hUpper == "PID":
			maxCap = 12
		}

		if colWidths[i] > maxCap {
			colWidths[i] = maxCap
		}
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(PrimaryColor)

	cellStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ECF0F1"))

	// Renderizar cabeçalho unificado como uma barra contínua roxa
	var headerLineBuilder strings.Builder
	headerLineBuilder.WriteString(" ")
	for i, h := range headers {
		truncatedH := truncateString(h, colWidths[i])
		visibleLen := lipgloss.Width(truncatedH)
		padLen := colWidths[i] - visibleLen
		if padLen < 0 {
			padLen = 0
		}
		headerLineBuilder.WriteString(truncatedH)
		headerLineBuilder.WriteString(strings.Repeat(" ", padLen))
		if i < len(headers)-1 {
			headerLineBuilder.WriteString("  ") // 2 espaços de separação entre colunas
		}
	}
	headerLineBuilder.WriteString(" ")

	fmt.Println()
	fmt.Println(headerStyle.Render(headerLineBuilder.String()))

	// Linha divisória
	totalW := lipgloss.Width(headerLineBuilder.String())
	fmt.Println(MutedStyle.Render(strings.Repeat("─", totalW)))

	// Renderizar linhas
	for _, row := range sanitRows {
		var rowLineBuilder strings.Builder
		rowLineBuilder.WriteString(" ")
		for i := 0; i < len(headers); i++ {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			truncatedVal := truncateString(val, colWidths[i])
			visibleLen := lipgloss.Width(truncatedVal)
			padLen := colWidths[i] - visibleLen
			if padLen < 0 {
				padLen = 0
			}

			rowLineBuilder.WriteString(cellStyle.Render(truncatedVal))
			rowLineBuilder.WriteString(strings.Repeat(" ", padLen))
			if i < len(headers)-1 {
				rowLineBuilder.WriteString("  ") // 2 espaços de separação entre colunas
			}
		}
		rowLineBuilder.WriteString(" ")
		fmt.Println(rowLineBuilder.String())
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

// RenderGitLog renderiza o histórico de commits em uma tabela limpa e alinhada com cabeçalhos roxos estilizados
func RenderGitLog(logs []GitLogItem) {
	if len(logs) == 0 {
		return
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(PrimaryColor)

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

	hashW := 10
	authorW := 14
	dateW := 14
	subjectW := 55

	for _, l := range logs {
		if lipgloss.Width(l.Author) > authorW {
			authorW = lipgloss.Width(l.Author)
		}
		if lipgloss.Width(l.Date) > dateW {
			dateW = lipgloss.Width(l.Date)
		}
	}

	if authorW > 20 {
		authorW = 20
	}
	if dateW > 18 {
		dateW = 18
	}

	// Renderizar cabeçalho unificado roxo sem lacunas entre as colunas
	headerStr := fmt.Sprintf(" %-*s  %-*s  %-*s  %-*s ",
		hashW, "HASH",
		authorW, "AUTOR",
		dateW, "DATA",
		subjectW, "MENSAGEM / ASSUNTO",
	)

	fmt.Println()
	fmt.Println(headerStyle.Render(headerStr))

	totalW := lipgloss.Width(headerStr)
	fmt.Println(MutedStyle.Render(strings.Repeat("─", totalW)))

	// Renderizar cada commit com alinhamento rigoroso
	for _, l := range logs {
		cleanSubject := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(l.Subject, "\r\n", " "), "\n", " "))

		hashVal := truncateString(l.Hash, hashW)
		authorVal := truncateString(l.Author, authorW)
		dateVal := truncateString(l.Date, dateW)
		subjectVal := truncateString(cleanSubject, subjectW)

		hashPadded := fmt.Sprintf("%-*s", hashW, hashVal)
		authorPadded := fmt.Sprintf("%-*s", authorW, authorVal)
		datePadded := fmt.Sprintf("%-*s", dateW, dateVal)
		subjectPadded := fmt.Sprintf("%-*s", subjectW, subjectVal)

		line := fmt.Sprintf(" %s  %s  %s  %s ",
			hashStyle.Render(hashPadded),
			authorStyle.Render(authorPadded),
			dateStyle.Render(datePadded),
			subjectStyle.Render(subjectPadded),
		)
		fmt.Println(line)
	}
	fmt.Println()
}

func truncateString(s string, maxLen int) string {
	w := lipgloss.Width(s)
	if w <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s
	}

	if strings.Contains(s, "\x1b") {
		return lipgloss.NewStyle().MaxWidth(maxLen).Render(s)
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}
