package system

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzePathAndLargestFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Criar arquivos de teste com tamanhos variados
	file1 := filepath.Join(tempDir, "small.txt")
	file2 := filepath.Join(tempDir, "large.bin")

	_ = os.WriteFile(file1, []byte("hello"), 0644)
	_ = os.WriteFile(file2, make([]byte, 1024*1024), 0644) // 1MB

	node, err := AnalyzePath(tempDir)
	if err != nil {
		t.Fatalf("AnalyzePath falhou: %v", err)
	}

	if node.SizeBytes < 1024*1024 {
		t.Errorf("Tamanho do nó diretório deveria ser pelo menos 1MB, obtido %d", node.SizeBytes)
	}

	largest, errL := GetLargestFiles(tempDir, 2)
	if errL != nil {
		t.Fatalf("GetLargestFiles falhou: %v", errL)
	}

	if len(largest) == 0 || largest[0].Name != "large.bin" {
		t.Errorf("O maior arquivo esperado era 'large.bin', obtido: %+v", largest)
	}
}
