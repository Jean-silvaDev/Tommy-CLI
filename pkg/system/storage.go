package system

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// AnalyzePath analisa recursivamente o tamanho de diretórios e arquivos em um caminho
func AnalyzePath(root string) (*FileNode, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		absRoot = root
	}

	info, err := os.Lstat(absRoot)
	if err != nil {
		return nil, err
	}

	node := &FileNode{
		Path:      absRoot,
		Name:      filepath.Base(absRoot),
		SizeBytes: uint64(info.Size()),
		IsDir:     info.IsDir(),
		ModTime:   info.ModTime(),
	}

	if !info.IsDir() {
		node.Extension = strings.ToLower(filepath.Ext(absRoot))
		return node, nil
	}

	entries, err := os.ReadDir(absRoot)
	if err != nil {
		return node, nil // Retorna nó mesmo se não puder ler filhos (sem permissão)
	}

	var totalDirSize uint64
	for _, entry := range entries {
		// Evitar loops infinitos ignorando Links Simbólicos e Junctions do Windows
		if entry.Type()&os.ModeSymlink != 0 {
			continue
		}

		childPath := filepath.Join(absRoot, entry.Name())
		childNode, errChild := AnalyzePath(childPath)
		if errChild == nil && childNode != nil {
			node.Children = append(node.Children, childNode)
			totalDirSize += childNode.SizeBytes
		}
	}

	node.SizeBytes = totalDirSize
	return node, nil
}

// GetLargestFiles varre um caminho e retorna os N maiores arquivos individuais encontrados
func GetLargestFiles(root string, limit int) ([]FileNode, error) {
	if limit <= 0 {
		limit = 10
	}

	var files []FileNode
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Pula arquivos/pastas com acesso negado
		}
		// Ignorar diretórios e symlinks
		if info.IsDir() || (info.Mode()&os.ModeSymlink != 0) {
			return nil
		}

		files = append(files, FileNode{
			Path:      path,
			Name:      info.Name(),
			SizeBytes: uint64(info.Size()),
			IsDir:     false,
			Extension: strings.ToLower(filepath.Ext(path)),
			ModTime:   info.ModTime(),
		})
		return nil
	})

	if err != nil && len(files) == 0 {
		return nil, err
	}

	// Ordenar os arquivos por tamanho decrescente
	sort.Slice(files, func(i, j int) bool {
		return files[i].SizeBytes > files[j].SizeBytes
	})

	if len(files) > limit {
		files = files[:limit]
	}

	return files, nil
}
