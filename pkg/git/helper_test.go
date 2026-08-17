package git

import (
	"testing"
)

func TestGetCurrentBranch(t *testing.T) {
	if !IsGitRepo() {
		t.Skip("Não é repositório Git, pulando teste")
	}

	branch, err := GetCurrentBranch()
	if err != nil {
		t.Fatalf("Erro inesperado ao obter a branch atual: %v", err)
	}
	if branch == "" {
		t.Errorf("Esperado nome de branch não vazio, obteve string vazia")
	}
}
