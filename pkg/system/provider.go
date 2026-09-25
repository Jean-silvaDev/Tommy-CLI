package system

// MemoryProvider abstrai a obtenção de dados de memória RAM do sistema operacional
type MemoryProvider interface {
	GetMemoryInfo() (MemoryInfo, error)
}

// ProcessProvider abstrai a leitura e encerramento de processos
type ProcessProvider interface {
	ListProcesses() ([]ProcessInfo, error)
	KillProcess(pid int) error
}

// DiskProvider abstrai a leitura de volumes e partições de disco
type DiskProvider interface {
	ListVolumes() ([]VolumeInfo, error)
}

// AppProvider abstrai a descoberta e desinstalação de aplicativos instalados
type AppProvider interface {
	ListApplications() ([]AppInfo, error)
	UninstallApplication(app AppInfo) error
}

// StorageAnalyzerProvider abstrai a análise de tamanho de diretórios e arquivos
type StorageAnalyzerProvider interface {
	AnalyzePath(path string) (*FileNode, error)
	GetLargestFiles(path string, limit int) ([]FileNode, error)
}

// CPUProvider abstrai a leitura de utilização da CPU
type CPUProvider interface {
	GetCPUUsage() (float64, error)
}

// SystemProvider agrupa as interfaces de leitura e gerenciamento do sistema operacional
type SystemProvider interface {
	MemoryProvider
	ProcessProvider
	DiskProvider
	AppProvider
	StorageAnalyzerProvider
	CPUProvider
	GetSystemSummary() (SystemSummary, error)
}
