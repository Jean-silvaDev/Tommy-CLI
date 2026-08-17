package cleaner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"tommy/pkg/config"
	"tommy/pkg/ui"
)

// CleanResult armazena o resultado de uma operação de limpeza
type CleanResult struct {
	FilesCount   int
	FilesRemoved int
	BytesFreed   int64
	ErrorsCount  int
	FileList     []string // Lista de caminhos/nomes de arquivos encontrados
}

// GetTempDirectories retorna os caminhos de pastas temporárias do usuário e do sistema
func GetTempDirectories() []string {
	var dirs []string
	seen := make(map[string]bool)

	addDir := func(d string) {
		if d == "" {
			return
		}
		abs, err := filepath.Abs(d)
		if err != nil {
			abs = d
		}
		if !seen[abs] {
			seen[abs] = true
			dirs = append(dirs, d)
		}
	}

	// 1. Diretórios customizados configurados no .env (TOMMY_TEMP_DIR / TOMMY_SYS_TEMP_DIR)
	for _, customDir := range config.GetTempDirs() {
		addDir(customDir)
	}

	// 2. Pasta Temp do Usuário (%TEMP% ou $TMPDIR)
	userTemp := os.TempDir()
	addDir(userTemp)

	// 3. Pasta Temp Geral do Sistema (Windows: SystemRoot\Temp)
	if runtime.GOOS == "windows" {
		sysRoot := os.Getenv("SystemRoot")
		if sysRoot == "" {
			sysDrive := os.Getenv("SystemDrive")
			if sysDrive == "" {
				sysDrive = "C:"
			}
			sysRoot = filepath.Join(sysDrive+`\`, "Windows")
		}
		winTemp := filepath.Join(sysRoot, "Temp")
		addDir(winTemp)
	} else {
		// Linux/macOS
		if _, err := os.Stat("/var/tmp"); err == nil {
			addDir("/var/tmp")
		}
	}

	return dirs
}

// FormatBytes formata bytes em uma string legível (KB, MB, GB)
func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// ScanTempFiles analisa os diretórios de temporários (Usuário e Sistema) e calcula a contagem, tamanho e lista de arquivos
func ScanTempFiles() (CleanResult, error) {
	var result CleanResult
	tempDirs := GetTempDirectories()

	for _, dir := range tempDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue // Pula caso não tenha permissão de leitura
		}

		for _, entry := range entries {
			result.FilesCount++
			fullPath := filepath.Join(dir, entry.Name())
			result.FileList = append(result.FileList, fullPath)

			info, err := entry.Info()
			if err == nil {
				result.BytesFreed += info.Size()
			}
		}
	}

	return result, nil
}

// CleanTempFiles limpa os diretórios de arquivos temporários do usuário e do sistema
func CleanTempFiles() (CleanResult, error) {
	var result CleanResult
	tempDirs := GetTempDirectories()

	for _, dir := range tempDirs {
		ui.PrintInfo("Executando exclusão no diretório: %s", dir)

		entries, err := os.ReadDir(dir)
		if err != nil {
			ui.PrintWarning("Não foi possível ler o diretório %s: %v", dir, err)
			continue
		}

		for _, entry := range entries {
			result.FilesCount++
			fullPath := filepath.Join(dir, entry.Name())

			rem, bytesFreed, errs := tryRemoveItem(fullPath)
			result.FilesRemoved += rem
			result.BytesFreed += bytesFreed
			result.ErrorsCount += errs
		}
	}

	return result, nil
}

// tryRemoveItem tenta excluir um arquivo ou diretório uma vez, removendo a proteção de leitura se necessário.
// Se for um diretório e a exclusão completa falhar (por conter arquivos travados pelo SO), tenta excluir cada subitem individualmente.
func tryRemoveItem(path string) (int, int64, int) {
	info, err := os.Lstat(path)
	if err != nil {
		return 0, 0, 1
	}

	fileSize := info.Size()

	// Remover atributo Read-Only no Windows para evitar bloqueio por permissão de arquivo
	if runtime.GOOS == "windows" {
		_ = os.Chmod(path, 0666)
	}

	// Tentar remover diretamente
	err = os.RemoveAll(path)
	if err == nil {
		return 1, fileSize, 0
	}

	// Se falhou e for um diretório, varre e apaga os subitens individualmente
	if info.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return 0, 0, 1
		}

		var removed int
		var bytesFreed int64
		var errs int

		for _, entry := range entries {
			subPath := filepath.Join(path, entry.Name())
			r, b, e := tryRemoveItem(subPath)
			removed += r
			bytesFreed += b
			errs += e
		}

		// Tenta remover o diretório se tiver ficado vazio
		_ = os.Remove(path)

		return removed, bytesFreed, errs
	}

	return 0, 0, 1
}
