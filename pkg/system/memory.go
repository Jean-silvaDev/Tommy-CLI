package system

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// FormatBytes formata bytes em unidades legíveis (B, KB, MB, GB, TB)
func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
	if exp < len(units) {
		return fmt.Sprintf("%.2f %s", float64(b)/float64(div), units[exp])
	}
	return fmt.Sprintf("%d B", b)
}

// CalculatePercent calcula o percentual seguro entre dois valores
func CalculatePercent(used, total uint64) float64 {
	if total == 0 {
		return 0.0
	}
	pct := (float64(used) / float64(total)) * 100.0
	if pct > 100.0 {
		return 100.0
	}
	return pct
}

// RenderProgressBar gera uma barra de progresso visual estilizada para o terminal
func RenderProgressBar(percent float64, width int) string {
	if width <= 0 {
		width = 20
	}
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filledLen := int((percent / 100.0) * float64(width))
	if filledLen > width {
		filledLen = width
	}
	emptyLen := width - filledLen

	filledStr := strings.Repeat("█", filledLen)
	emptyStr := strings.Repeat("░", emptyLen)

	var colorStyle lipgloss.Style
	switch {
	case percent >= 90.0:
		colorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7675")) // Vermelho/Alerta
	case percent >= 75.0:
		colorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FDCB6E")) // Amarelo/Aviso
	default:
		colorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#04B575")) // Verde
	}

	bar := colorStyle.Render(filledStr) + lipgloss.NewStyle().Foreground(lipgloss.Color("#636E72")).Render(emptyStr)
	return fmt.Sprintf("%s %.1f%%", bar, percent)
}
