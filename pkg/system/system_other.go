//go:build !windows

package system

import (
	"errors"
)

type OtherSystemProvider struct{}

func NewSystemProvider() SystemProvider {
	return &OtherSystemProvider{}
}

func (p *OtherSystemProvider) GetMemoryInfo() (MemoryInfo, error) {
	return MemoryInfo{}, errors.New("plataforma não suportada nativamente para leitura de memória")
}

func (p *OtherSystemProvider) ListVolumes() ([]VolumeInfo, error) {
	return nil, errors.New("plataforma não suportada nativamente para leitura de volumes")
}

func (p *OtherSystemProvider) ListProcesses() ([]ProcessInfo, error) {
	return nil, errors.New("plataforma não suportada nativamente para leitura de processos")
}

func (p *OtherSystemProvider) KillProcess(pid int) error {
	return errors.New("plataforma não suportada nativamente para encerramento de processos")
}

func (p *OtherSystemProvider) ListApplications() ([]AppInfo, error) {
	return nil, errors.New("plataforma não suportada nativamente para leitura de aplicativos")
}

func (p *OtherSystemProvider) UninstallApplication(app AppInfo) error {
	return errors.New("plataforma não suportada nativamente para desinstalação de aplicativos")
}

func (p *OtherSystemProvider) AnalyzePath(path string) (*FileNode, error) {
	return AnalyzePath(path)
}

func (p *OtherSystemProvider) GetLargestFiles(path string, limit int) ([]FileNode, error) {
	return GetLargestFiles(path, limit)
}

func (p *OtherSystemProvider) GetCPUUsage() (float64, error) {
	return 0.0, errors.New("plataforma não suportada nativamente para leitura de CPU")
}

func (p *OtherSystemProvider) GetSystemSummary() (SystemSummary, error) {
	return SystemSummary{}, errors.New("plataforma não suportada")
}
