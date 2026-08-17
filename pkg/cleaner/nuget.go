package cleaner

import (
	"fmt"
	"os/exec"
	"strings"

	"tommy/pkg/ui"
)

// CleanNuGetCache executa a limpeza dos caches locais do .NET NuGet
func CleanNuGetCache() error {
	ui.PrintInfo("Executando limpeza de pacotes NuGet (.NET)...")

	cmd := exec.Command("dotnet", "nuget", "locals", "all", "--clear")
	output, err := cmd.CombinedOutput()
	outputStr := strings.TrimSpace(string(output))

	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return fmt.Errorf("O CLI 'dotnet' não foi encontrado no PATH do sistema")
		}
		return fmt.Errorf("falha ao limpar cache NuGet: %s", outputStr)
	}

	ui.PrintSuccess("Cache de pacotes NuGet limpo com sucesso!\n%s", outputStr)
	return nil
}
