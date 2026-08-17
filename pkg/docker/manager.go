package docker

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
)

// ContainerInfo representa informações essenciais de um container Docker
type ContainerInfo struct {
	ID      string
	Name    string
	Image   string
	Status  string
	State   string // "running", "exited", etc.
	Project string // Nome da Stack / Docker Compose Project
}

// DockerStack agrupa containers pertencentes à mesma Stack Compose
type DockerStack struct {
	Name         string
	Containers   []ContainerInfo
	RunningCount int
	TotalCount   int
}

// IsDockerRunning verifica se o daemon do Docker está ativo e respondendo
func IsDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	err := cmd.Run()
	return err == nil
}

// StartDockerEngine tenta inicializar o serviço / app do Docker Desktop
func StartDockerEngine() error {
	if customPath := config.GetDockerPath(); customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			cmd := exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf("Start-Process '%s'", customPath))
			if errExec := cmd.Run(); errExec == nil {
				return nil
			}
		}
	}

	if runtime.GOOS == "windows" {
		var possiblePaths []string

		pf := os.Getenv("ProgramFiles")
		if pf != "" {
			possiblePaths = append(possiblePaths, filepath.Join(pf, "Docker", "Docker", "Docker Desktop.exe"))
		}
		if sysDrive := os.Getenv("SystemDrive"); sysDrive != "" {
			possiblePaths = append(possiblePaths, filepath.Join(sysDrive+`\`, "Program Files", "Docker", "Docker", "Docker Desktop.exe"))
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			possiblePaths = append(possiblePaths, filepath.Join(localAppData, "Programs", "Docker", "Docker", "Docker Desktop.exe"))
		}



		for _, exePath := range possiblePaths {
			if exePath != "" {
				if _, err := os.Stat(exePath); err == nil {
					cmd := exec.Command("powershell", "-NoProfile", "-Command", fmt.Sprintf("Start-Process '%s'", exePath))
					if errExec := cmd.Run(); errExec == nil {
						return nil
					}
				}
			}
		}

		// Tentar via registro de atalho do PowerShell
		cmdPs := exec.Command("powershell", "-NoProfile", "-Command", "Start-Process 'Docker Desktop'")
		if errPs := cmdPs.Run(); errPs == nil {
			return nil
		}

		// Tentar via serviço do Windows net start
		cmdNet := exec.Command("net", "start", "com.docker.service")
		if errNet := cmdNet.Run(); errNet == nil {
			return nil
		}

		return fmt.Errorf("não foi possível localizar ou iniciar o 'Docker Desktop.exe' no sistema")
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("open", "-a", "Docker")
		return cmd.Run()
	} else {
		cmd := exec.Command("sudo", "systemctl", "start", "docker")
		return cmd.Run()
	}
}

// StopDockerEngine encerra o serviço / processo do Docker Desktop
func StopDockerEngine() error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("taskkill", "/F", "/IM", "Docker Desktop.exe", "/IM", "com.docker.backend.exe")
		out, err := cmd.CombinedOutput()
		if err != nil && !strings.Contains(string(out), "não encontrado") {
			return fmt.Errorf("erro ao parar Docker Desktop: %s", string(out))
		}
		return nil
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("osascript", "-e", `quit app "Docker"`)
		return cmd.Run()
	} else {
		cmd := exec.Command("sudo", "systemctl", "stop", "docker")
		return cmd.Run()
	}
}

// ListContainers lista todos os containers do Docker (rodando e parados) com informação de Stack
func ListContainers() ([]ContainerInfo, error) {
	cmd := exec.Command("docker", "ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Image}}|{{.Status}}|{{.State}}|{{.Label \"com.docker.compose.project\"}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			return nil, fmt.Errorf("O CLI 'docker' não está instalado ou no PATH")
		}
		return nil, fmt.Errorf("falha ao listar containers: %s", string(out))
	}

	var containers []ContainerInfo
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 5 {
			proj := ""
			if len(parts) >= 6 {
				proj = strings.TrimSpace(parts[5])
			}

			cName := parts[1]
			if proj == "" && strings.Contains(cName, "-") {
				subParts := strings.Split(cName, "-")
				if len(subParts) >= 3 {
					proj = strings.Join(subParts[:len(subParts)-2], "-")
				} else if len(subParts) == 2 {
					proj = subParts[0]
				}
			}

			if proj == "" {
				proj = "Avulsos"
			}

			containers = append(containers, ContainerInfo{
				ID:      parts[0],
				Name:    cName,
				Image:   parts[2],
				Status:  parts[3],
				State:   parts[4],
				Project: proj,
			})
		}
	}

	return containers, nil
}

// GetDockerStacks agrupa os containers por Stack (projeto Compose)
func GetDockerStacks(containers []ContainerInfo) []DockerStack {
	stackMap := make(map[string][]ContainerInfo)
	var stackNames []string

	for _, c := range containers {
		proj := c.Project
		if proj == "" {
			proj = "Avulsos"
		}
		if _, exists := stackMap[proj]; !exists {
			stackNames = append(stackNames, proj)
		}
		stackMap[proj] = append(stackMap[proj], c)
	}

	var result []DockerStack
	for _, name := range stackNames {
		cList := stackMap[name]
		running := 0
		for _, c := range cList {
			if strings.ToLower(c.State) == "running" {
				running++
			}
		}
		result = append(result, DockerStack{
			Name:         name,
			Containers:   cList,
			RunningCount: running,
			TotalCount:   len(cList),
		})
	}

	return result
}

// StartStack inicia todos os containers de uma Stack
func StartStack(containerIDs []string) error {
	if len(containerIDs) == 0 {
		return nil
	}
	args := append([]string{"start"}, containerIDs...)
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao iniciar stack: %s", string(out))
	}
	return nil
}

// StopStack para todos os containers de uma Stack
func StopStack(containerIDs []string) error {
	if len(containerIDs) == 0 {
		return nil
	}
	args := append([]string{"stop"}, containerIDs...)
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao parar stack: %s", string(out))
	}
	return nil
}

// RestartStack reinicia todos os containers de uma Stack
func RestartStack(containerIDs []string) error {
	if len(containerIDs) == 0 {
		return nil
	}
	args := append([]string{"restart"}, containerIDs...)
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao reiniciar stack: %s", string(out))
	}
	return nil
}

// StartContainer inicia um container parado por ID ou Nome
func StartContainer(idOrName string) error {
	cmd := exec.Command("docker", "start", idOrName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao iniciar container %s: %s", idOrName, string(out))
	}
	return nil
}

// StopContainer para um container rodando por ID ou Nome
func StopContainer(idOrName string) error {
	cmd := exec.Command("docker", "stop", idOrName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao parar container %s: %s", idOrName, string(out))
	}
	return nil
}

// RestartContainer reinicia um container por ID ou Nome
func RestartContainer(idOrName string) error {
	cmd := exec.Command("docker", "restart", idOrName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao reiniciar container %s: %s", idOrName, string(out))
	}
	return nil
}

// ImageInfo representa dados de uma imagem Docker
type ImageInfo struct {
	ID         string
	Repository string
	Tag        string
	Size       string
	Created    string
	InUse      bool
}

// VolumeInfo representa dados de um volume Docker
type VolumeInfo struct {
	Name   string
	Driver string
	Scope  string
	InUse  bool
}

// NetworkInfo representa dados de uma rede Docker
type NetworkInfo struct {
	ID          string
	Name        string
	Driver      string
	Scope       string
	IsProtected bool
	InUse       bool
}

// ListImages lista todas as imagens Docker e identifica se estão em uso por algum container
func ListImages() ([]ImageInfo, error) {
	containers, _ := ListContainers()
	usedImages := make(map[string]bool)
	for _, c := range containers {
		usedImages[c.Image] = true
	}

	cmd := exec.Command("docker", "images", "--format", "{{.Repository}}|{{.Tag}}|{{.ID}}|{{.Size}}|{{.CreatedAt}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("falha ao listar imagens: %s", string(out))
	}

	var images []ImageInfo
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 4 {
			repo := parts[0]
			tag := parts[1]
			id := parts[2]
			size := parts[3]
			created := ""
			if len(parts) >= 5 {
				created = parts[4]
			}

			fullRef := repo + ":" + tag
			inUse := usedImages[fullRef] || usedImages[repo] || usedImages[id]

			images = append(images, ImageInfo{
				ID:         id,
				Repository: repo,
				Tag:        tag,
				Size:       size,
				Created:    created,
				InUse:      inUse,
			})
		}
	}

	return images, nil
}

// RemoveImage exclui uma imagem Docker por ID ou Repositório:Tag
func RemoveImage(idOrRef string, force bool) error {
	args := []string{"rmi"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, idOrRef)
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao remover imagem %s: %s", idOrRef, string(out))
	}
	return nil
}

// PruneImages executa docker image prune [-a] -f
func PruneImages(all bool) error {
	args := []string{"image", "prune", "-f"}
	if all {
		args = append(args, "-a")
	}
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao limpar imagens: %s", string(out))
	}
	return nil
}

// ListVolumes lista todos os volumes do Docker e verifica se estão anexados a containers
func ListVolumes() ([]VolumeInfo, error) {
	cmd := exec.Command("docker", "volume", "ls", "--format", "{{.Name}}|{{.Driver}}|{{.Scope}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("falha ao listar volumes: %s", string(out))
	}

	var volumes []VolumeInfo
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 2 {
			name := parts[0]
			driver := parts[1]
			scope := "local"
			if len(parts) >= 3 && parts[2] != "" {
				scope = parts[2]
			}

			cmdInspect := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("volume=%s", name), "--format", "{{.ID}}")
			outInsp, _ := cmdInspect.CombinedOutput()
			inUse := strings.TrimSpace(string(outInsp)) != ""

			volumes = append(volumes, VolumeInfo{
				Name:   name,
				Driver: driver,
				Scope:  scope,
				InUse:  inUse,
			})
		}
	}

	return volumes, nil
}

// RemoveVolume remove um volume por Nome
func RemoveVolume(name string, force bool) error {
	args := []string{"volume", "rm"}
	if force {
		args = append(args, "-f")
	}
	args = append(args, name)
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao remover volume %s: %s", name, string(out))
	}
	return nil
}

// PruneVolumes executa docker volume prune -a -f e garante a remoção de volumes órfãos identificados
func PruneVolumes() error {
	cmd := exec.Command("docker", "volume", "prune", "-a", "-f")
	_, _ = cmd.CombinedOutput()

	// Remover explicitamente qualquer volume órfão remanescente
	volumes, err := ListVolumes()
	if err == nil {
		for _, v := range volumes {
			if !v.InUse {
				_ = RemoveVolume(v.Name, true)
			}
		}
	}
	return nil
}

// ListNetworks lista todas as redes Docker
func ListNetworks() ([]NetworkInfo, error) {
	cmd := exec.Command("docker", "network", "ls", "--format", "{{.ID}}|{{.Name}}|{{.Driver}}|{{.Scope}}")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("falha ao listar redes: %s", string(out))
	}

	var networks []NetworkInfo
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			name := parts[1]
			isProtected := name == "bridge" || name == "host" || name == "none"
			networks = append(networks, NetworkInfo{
				ID:          parts[0],
				Name:        name,
				Driver:      parts[2],
				Scope:       parts[3],
				IsProtected: isProtected,
			})
		}
	}

	return networks, nil
}

// RemoveNetwork exclui uma rede customizada
func RemoveNetwork(name string) error {
	cmd := exec.Command("docker", "network", "rm", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao remover rede %s: %s", name, string(out))
	}
	return nil
}

// PruneNetworks remove redes não utilizadas
func PruneNetworks() error {
	cmd := exec.Command("docker", "network", "prune", "-f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao limpar redes: %s", string(out))
	}
	return nil
}

// PruneContainers remove containers parados
func PruneContainers() error {
	cmd := exec.Command("docker", "container", "prune", "-f")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao limpar containers parados: %s", string(out))
	}
	return nil
}

// PruneSystem executa docker system prune -f [--all] [--volumes]
func PruneSystem(all, volumes bool) error {
	args := []string{"system", "prune", "-f"}
	if all {
		args = append(args, "-a")
	}
	if volumes {
		args = append(args, "--volumes")
	}
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro na limpeza do sistema docker: %s", string(out))
	}
	return nil
}

// GetDockerUsageSummary obtém a saída formatada de docker system df
func GetDockerUsageSummary() (string, error) {
	cmd := exec.Command("docker", "system", "df")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("falha ao consultar espaço do docker: %s", string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// CreateNetwork cria uma nova rede Docker
func CreateNetwork(name, driver, subnet string) error {
	args := []string{"network", "create"}
	if driver != "" {
		args = append(args, "--driver", driver)
	}
	if subnet != "" {
		args = append(args, "--subnet", subnet)
	}
	args = append(args, name)

	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao criar rede %s: %s", name, string(out))
	}
	return nil
}

// InspectNetwork retorna os detalhes formatados de uma rede Docker (inspect)
func InspectNetwork(name string) (string, error) {
	cmd := exec.Command("docker", "network", "inspect", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro ao inspecionar rede %s: %s", name, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// ConnectNetwork conecta um container a uma rede Docker
func ConnectNetwork(networkName, containerID string) error {
	cmd := exec.Command("docker", "network", "connect", networkName, containerID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao conectar container à rede: %s", string(out))
	}
	return nil
}

// DisconnectNetwork desconecta um container de uma rede Docker
func DisconnectNetwork(networkName, containerID string) error {
	cmd := exec.Command("docker", "network", "disconnect", networkName, containerID)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao desconectar container da rede: %s", string(out))
	}
	return nil
}

// GetContainerLogs obtém as últimas N linhas de logs de um container
func GetContainerLogs(idOrName string, tailLines int) (string, error) {
	cmd := exec.Command("docker", "logs", "--tail", fmt.Sprintf("%d", tailLines), idOrName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro ao obter logs do container %s: %s", idOrName, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// InspectContainer obtém a saída formatada do docker inspect em um container
func InspectContainer(idOrName string) (string, error) {
	cmd := exec.Command("docker", "inspect", idOrName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("erro ao inspecionar container %s: %s", idOrName, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// PullImage faz o download de uma imagem do Docker Hub
func PullImage(imageRef string) error {
	cmd := exec.Command("docker", "pull", imageRef)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao baixar imagem %s: %s", imageRef, string(out))
	}
	return nil
}

// CreateVolume cria um novo volume nomeado
func CreateVolume(name, driver string) error {
	args := []string{"volume", "create"}
	if driver != "" {
		args = append(args, "--driver", driver)
	}
	args = append(args, name)

	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("erro ao criar volume %s: %s", name, string(out))
	}
	return nil
}
