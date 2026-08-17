package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommitLog representa a estrutura de um commit simplificado
type CommitLog struct {
	Hash    string
	Author  string
	Date    string
	Subject string
}

// GetStatus retorna o status resumido dos arquivos staged e unstaged
func GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--short")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("falha ao obter status do git: %s", string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// StageAll executa 'git add -A'
func StageAll() error {
	cmd := exec.Command("git", "add", "-A")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao adicionar arquivos (git add): %s", string(out))
	}
	return nil
}

// CreateCommit realiza o commit com a mensagem formatada
func CreateCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao realizar commit: %s", string(out))
	}
	return nil
}

// GetRecentLogs obtém os últimos N commits formatados
func GetRecentLogs(count int) ([]CommitLog, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-n%d", count), "--pretty=format:%h|%an|%cr|%s")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("falha ao obter logs do git: %s", string(out))
	}

	var logs []CommitLog
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 4 {
			logs = append(logs, CommitLog{
				Hash:    parts[0],
				Author:  parts[1],
				Date:    parts[2],
				Subject: parts[3],
			})
		}
	}

	return logs, nil
}

// GetBranches lista as branches locais e indica a branch ativa
func GetBranches() ([]string, string, error) {
	cmd := exec.Command("git", "branch", "--list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", fmt.Errorf("falha ao listar branches: %s", string(out))
	}

	var branches []string
	currentBranch := ""

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "*") {
			name := strings.TrimSpace(strings.TrimPrefix(line, "*"))
			currentBranch = name
			branches = append(branches, name)
		} else {
			branches = append(branches, line)
		}
	}

	return branches, currentBranch, nil
}

// CreateBranch cria uma nova branch a partir da atual
func CreateBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", "-b", branchName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao criar branch '%s': %s", branchName, string(out))
	}
	return nil
}

// SwitchBranch altera para uma branch existente
func SwitchBranch(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao alternar para branch '%s': %s", branchName, string(out))
	}
	return nil
}

// RenameBranch renomeia uma branch existente
func RenameBranch(oldName, newName string) error {
	cmd := exec.Command("git", "branch", "-m", oldName, newName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao renomear branch '%s' para '%s': %s", oldName, newName, string(out))
	}
	return nil
}

// DeleteBranch exclui uma branch (use force=true para -D)
func DeleteBranch(branchName string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}
	cmd := exec.Command("git", "branch", flag, branchName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao excluir branch '%s': %s", branchName, string(out))
	}
	return nil
}


// GetTags lista todas as tags do repositório
func GetTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "-l")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("falha ao listar tags: %s", string(out))
	}

	var tags []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			tags = append(tags, line)
		}
	}

	return tags, nil
}

// CreateTag cria uma nova tag anotada no repositório
func CreateTag(tagName string, message string) error {
	cmd := exec.Command("git", "tag", "-a", tagName, "-m", message)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao criar tag '%s': %s", tagName, string(out))
	}
	return nil
}

// GetCurrentBranch retorna o nome da branch ativa
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("falha ao obter branch atual: %s", string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// Push Remote envia commits e tags para o repositório remoto
func Push(includeTags bool) error {
	args := []string{"push"}
	if includeTags {
		args = append(args, "--follow-tags")
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := string(out)
		if strings.Contains(errMsg, "has no upstream branch") || strings.Contains(errMsg, "set-upstream") {
			currBranch, branchErr := GetCurrentBranch()
			if branchErr == nil && currBranch != "" {
				setUpstreamArgs := []string{"push", "--set-upstream", "origin", currBranch}
				if includeTags {
					setUpstreamArgs = append(setUpstreamArgs, "--follow-tags")
				}
				cmdUpstream := exec.Command("git", setUpstreamArgs...)
				outUpstream, errUpstream := cmdUpstream.CombinedOutput()
				if errUpstream == nil {
					return nil
				}
				return fmt.Errorf("falha ao executar git push --set-upstream origin %s: %s", currBranch, string(outUpstream))
			}
		}
		return fmt.Errorf("falha ao executar git push: %s", errMsg)
	}
	return nil
}

// Pull obtém e mescla as alterações do repositório remoto
func Pull() error {
	cmd := exec.Command("git", "pull")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao executar git pull: %s", string(out))
	}
	return nil
}

// GetConflictedFiles retorna a lista de arquivos atualmente em conflito de merge
func GetConflictedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("falha ao verificar conflitos: %s", string(out))
	}

	var files []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// AbortMerge aborta o processo de merge em andamento
func AbortMerge() error {
	cmd := exec.Command("git", "merge", "--abort")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao abortar merge: %s", string(out))
	}
	return nil
}



// IsGitRepo verifica se o diretório atual é um repositório Git válido
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// InitRepo inicializa a pasta atual como repositório Git
func InitRepo() error {
	cmd := exec.Command("git", "init")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao inicializar repositório Git: %s", string(out))
	}
	return nil
}

// GetRemoteURL obtém a URL do repositório remoto 'origin'
func GetRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("nenhum remoto 'origin' configurado")
	}
	return strings.TrimSpace(string(out)), nil
}

// AddOrSetRemote adiciona ou atualiza a URL do remoto 'origin'
func AddOrSetRemote(remoteURL string) error {
	_, err := GetRemoteURL()
	var cmd *exec.Cmd
	if err == nil {
		cmd = exec.Command("git", "remote", "set-url", "origin", remoteURL)
	} else {
		cmd = exec.Command("git", "remote", "add", "origin", remoteURL)
	}
	out, errExec := cmd.CombinedOutput()
	if errExec != nil {
		return fmt.Errorf("falha ao configurar remoto 'origin': %s", string(out))
	}
	return nil
}

