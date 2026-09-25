package system

import (
	"testing"
)

func TestEvaluateAlerts(t *testing.T) {
	s := SystemSummary{
		Memory: MemoryInfo{UsedPercent: 92.5, UsedBytes: 15000000000, TotalBytes: 16000000000},
		Volumes: []VolumeInfo{
			{Device: "C:", UsedPercent: 95.0, FreeBytes: 1000000000},
		},
	}

	alerts := EvaluateAlerts(s)
	if len(alerts) != 2 {
		t.Errorf("Esperado 2 alertas (RAM e Disco), obtido: %d", len(alerts))
	}
}
