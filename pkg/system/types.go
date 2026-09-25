package system

import "time"

// MemoryInfo armazena estatísticas de memória RAM
type MemoryInfo struct {
	TotalBytes     uint64
	UsedBytes      uint64
	AvailableBytes uint64
	UsedPercent    float64
}

// CPUInfo armazena métricas básicas de processador
type CPUInfo struct {
	UsagePercent float64
}

// VolumeInfo armazena estatísticas de um volume ou disco
type VolumeInfo struct {
	Device       string // Ex: "C:"
	Label        string // Rótulo do volume
	DriveType    string // Ex: "Fixo", "Removível", "Rede", "CDROM"
	DiskType     string // Ex: "SSD/HDD", "Pendrive/Removível", "Unidade de Rede"
	TotalBytes   uint64
	UsedBytes    uint64
	FreeBytes    uint64
	UsedPercent  float64
	FileSystem   string // NTFS, FAT32, exFAT, etc.
	Status       string // "OK", "Aviso", "Erro"
	HealthStatus string // Ex: "OK", "98% (Ótimo)", "Informação não disponível"
	TempCelsius  int    // Temperatura se disponível (0 se indisponível)
}

// ProcessInfo armazena dados de um processo ativo no SO
type ProcessInfo struct {
	PID           int
	Name          string
	Path          string
	User          string
	MemoryBytes   uint64
	MemoryPercent float64
	CPUPercent    float64
	Threads       int
	Status        string
	IsProtected   bool
}

// AppInfo armazena dados de um aplicativo instalado no sistema
type AppInfo struct {
	ID             string    // Identificador único ou nome do pacote
	Name           string    // Nome de exibição do programa
	Version        string    // Versão instalada
	Publisher      string    // Fabricante / Desenvolvedor
	SizeBytes      uint64    // Tamanho estimado/calculado no disco
	InstallDate    time.Time // Data de instalação
	InstallLocation string   // Caminho de instalação
	UninstallString string   // Comando nativo de desinstalação
	Source         string    // "Winget", "Registry (HKLM)", "Registry (HKCU)"
	UpdateAvailable bool     // Se possui atualização disponível no Winget
	NewVersion     string    // Nova versão se disponível
}

// FileNode armazena estatísticas de um arquivo ou diretório analisado
type FileNode struct {
	Path        string
	Name        string
	SizeBytes   uint64
	IsDir       bool
	Extension   string
	ModTime     time.Time
	Children    []*FileNode
}

// AlertInfo representa um aviso ou alerta emitido pelo System Manager
type AlertInfo struct {
	Level   string // "CRITICAL", "WARNING", "INFO"
	Message string
}

// SystemSummary agrega os dados principais para o Dashboard do System Manager
type SystemSummary struct {
	Memory  MemoryInfo
	CPU     CPUInfo
	Volumes []VolumeInfo
	Alerts  []AlertInfo
}
