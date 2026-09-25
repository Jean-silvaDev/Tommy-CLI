package system

import "fmt"

// EvaluateAlerts analisa os dados do sistema e gera uma lista de alertas configuráveis
func EvaluateAlerts(s SystemSummary) []AlertInfo {
	var alerts []AlertInfo

	// 1. Alerta de Memória RAM acima de 90%
	if s.Memory.UsedPercent >= 90.0 {
		alerts = append(alerts, AlertInfo{
			Level:   "CRITICAL",
			Message: fmt.Sprintf("⚠️ RAM acima de 90%% (%.1f%% em uso - %s / %s)", s.Memory.UsedPercent, FormatBytes(s.Memory.UsedBytes), FormatBytes(s.Memory.TotalBytes)),
		})
	} else if s.Memory.UsedPercent >= 80.0 {
		alerts = append(alerts, AlertInfo{
			Level:   "WARNING",
			Message: fmt.Sprintf("⚠️ RAM acima de 80%% (%.1f%% em uso)", s.Memory.UsedPercent),
		})
	}

	// 2. Alerta de espaço em discos (Espaço livre < 10%)
	for _, v := range s.Volumes {
		if v.UsedPercent >= 90.0 {
			alerts = append(alerts, AlertInfo{
				Level:   "CRITICAL",
				Message: fmt.Sprintf("⚠️ O disco %s possui apenas %.1f%% de espaço livre (%s livre)", v.Device, 100.0-v.UsedPercent, FormatBytes(v.FreeBytes)),
			})
		}
	}

	return alerts
}
