package cleaner

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"tommy/pkg/config"
	"tommy/pkg/ui"
)

// GetRecycleBinItems consulta os nomes dos arquivos/pastas atualmente na Lixeira
func GetRecycleBinItems() ([]string, error) {
	if customTrash := config.GetTrashDir(); customTrash != "" {
		if entries, err := os.ReadDir(customTrash); err == nil {
			var items []string
			for _, entry := range entries {
				items = append(items, entry.Name())
			}
			return items, nil
		}
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-NoProfile", "-Command",
			"(New-Object -ComObject Shell.Application).NameSpace(0x0a).Items() | Select-Object -First 30 | ForEach-Object { $_.Name }")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}

		var items []string
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				items = append(items, line)
			}
		}
		return items, nil
	}

	// Linux / macOS
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	trashDir := filepath.Join(home, ".local", "share", "Trash", "files")

	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return nil, nil
	}

	var items []string
	for _, entry := range entries {
		items = append(items, entry.Name())
	}
	return items, nil
}

// CleanRecycleBin esvazia a Lixeira do sistema (Windows Recycle Bin ou Trash no Linux/macOS)
func CleanRecycleBin() error {
	ui.PrintInfo("Esvaziando a Lixeira do sistema...")

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "Clear-RecycleBin -Force -Confirm:$false -ErrorAction SilentlyContinue")
		_ = cmd.Run()

		// Varredura complementar graciosa em diretórios $Recycle.Bin nos discos locais
		var drives []string
		for _, d := range []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"} {
			bin := fmt.Sprintf(`%s:\$Recycle.Bin`, d)
			if _, err := os.Stat(bin); err == nil {
				drives = append(drives, bin)
			}
		}

		for _, binPath := range drives {
			entries, err := os.ReadDir(binPath)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				itemPath := filepath.Join(binPath, entry.Name())
				_ = os.Chmod(itemPath, 0666)
				_ = os.RemoveAll(itemPath)
			}
		}

		ui.PrintSuccess("Lixeira do Windows esvaziada com sucesso (arquivos bloqueados pelo SO foram ignorados).")
		return nil
	}

	// Linux / macOS
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("não foi possível obter diretório home: %w", err)
	}

	trashPaths := []string{
		filepath.Join(home, ".local", "share", "Trash", "files"),
		filepath.Join(home, ".Trash"),
	}

	for _, trashPath := range trashPaths {
		if _, err := os.Stat(trashPath); err == nil {
			entries, err := os.ReadDir(trashPath)
			if err != nil {
				continue
			}
			for _, entry := range entries {
				itemPath := filepath.Join(trashPath, entry.Name())
				_ = os.Chmod(itemPath, 0666)
				_ = os.RemoveAll(itemPath)
			}
		}
	}

	ui.PrintSuccess("Lixeira do sistema esvaziada com sucesso!")
	return nil
}
