package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// LoadEnv carrega variáveis de ambiente a partir do arquivo .env sem sobrescrever variáveis do SO
func LoadEnv() {
	envPath := ".env"
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		if exe, err := os.Executable(); err == nil {
			exeDir := filepath.Dir(exe)
			altPath := filepath.Join(exeDir, ".env")
			if _, err := os.Stat(altPath); err == nil {
				envPath = altPath
			} else {
				return
			}
		} else {
			return
		}
	}

	file, err := os.Open(envPath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)

			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}

// GetEnv obtém o valor da variável de ambiente com valor fallback de segurança
func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// ExpandHome substitui o prefixo '~' pelo diretório home do usuário
func ExpandHome(path string) string {
	if path == "" {
		return ""
	}
	if path == "~" || strings.HasPrefix(path, "~/") || strings.HasPrefix(path, `~\`) {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		if path == "~" {
			return home
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// GetOutputDir retorna o diretório configurado para saída dos prompts e artefatos
func GetOutputDir() string {
	if envDir := os.Getenv("TOMMY_OUTPUT_DIR"); envDir != "" {
		return ExpandHome(envDir)
	}
	return filepath.Join(".", "PromptBuilder")
}

// GetConfigDir retorna o diretório para presets e configurações globais
func GetConfigDir() string {
	if envDir := os.Getenv("TOMMY_CONFIG_DIR"); envDir != "" {
		return ExpandHome(envDir)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".tommy")
	}
	return filepath.Join(home, ".tommy")
}

// GetDockerPath retorna o caminho customizado do executável do Docker Desktop, se definido
func GetDockerPath() string {
	return ExpandHome(os.Getenv("TOMMY_DOCKER_PATH"))
}

// GetTempDirs retorna diretórios adicionais/customizados de arquivos temporários configurados via env
func GetTempDirs() []string {
	var dirs []string
	if val := os.Getenv("TOMMY_TEMP_DIR"); val != "" {
		for _, d := range strings.Split(val, ",") {
			if trimmed := strings.TrimSpace(d); trimmed != "" {
				dirs = append(dirs, ExpandHome(trimmed))
			}
		}
	}
	if val := os.Getenv("TOMMY_SYS_TEMP_DIR"); val != "" {
		for _, d := range strings.Split(val, ",") {
			if trimmed := strings.TrimSpace(d); trimmed != "" {
				dirs = append(dirs, ExpandHome(trimmed))
			}
		}
	}
	return dirs
}

// GetTrashDir retorna o caminho customizado do diretório de lixeira, se definido
func GetTrashDir() string {
	return ExpandHome(os.Getenv("TOMMY_TRASH_DIR"))
}

// GetAuthor retorna o nome do autor configurado ou o usuário do sistema operacional
func GetAuthor() string {
	if val := os.Getenv("TOMMY_AUTHOR"); val != "" {
		return val
	}
	if val := os.Getenv("USERNAME"); val != "" {
		return val
	}
	if val := os.Getenv("USER"); val != "" {
		return val
	}
	return "Desenvolvedor"
}

