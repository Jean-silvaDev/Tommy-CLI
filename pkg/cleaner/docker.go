package cleaner

import (
	"fmt"
	"os/exec"
	"strings"

	"tommy/pkg/ui"
)

// GetDockerUsageSummary obtém um resumo do uso de disco do Docker (df)
func GetDockerUsageSummary() (string, error) {
	cmd := exec.Command("docker", "system", "df")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return "", fmt.Errorf("O CLI 'docker' não foi encontrado no PATH do sistema")
		}
		return "", fmt.Errorf("falha ao consultar espaço do docker: %s", string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// CleanDockerResources executa a limpeza de imagens, containers parados, caches e volumes do Docker
func CleanDockerResources() error {
	ui.PrintInfo("Executando limpeza profunda do Docker (containers parados, imagens sem uso e caches)...")

	// 1. System Prune
	cmdSystem := exec.Command("docker", "system", "prune", "-af")
	outSystem, err := cmdSystem.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return fmt.Errorf("O CLI 'docker' não foi encontrado no PATH do sistema")
		}
		return fmt.Errorf("erro no docker system prune: %s", string(outSystem))
	}

	ui.PrintSuccess("Limpeza de sistema e imagens do Docker concluída.")

	// 2. Volume Prune
	ui.PrintInfo("Executando limpeza de volumes não utilizados...")
	cmdVolume := exec.Command("docker", "volume", "prune", "-a", "-f")
	outVolume, err := cmdVolume.CombinedOutput()
	if err != nil {
		ui.PrintWarning("Erro ao limpar volumes do Docker: %s", string(outVolume))
	} else {
		ui.PrintSuccess("Limpeza de volumes concluída.")
	}

	return nil
}
