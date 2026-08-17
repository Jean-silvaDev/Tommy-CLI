package network

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// NetworkSummary contém o resumo de conectividade local e pública
type NetworkSummary struct {
	LocalIP       string
	LocalIPv6     string
	InterfaceName string
	MACAddress    string
	Gateway       string
	PublicIP      string
	ISP           string
	Location      string
	IsOnline      bool
}

// InterfaceInfo guarda os dados de um adaptador de rede
type InterfaceInfo struct {
	Name       string
	MAC        string
	IPv4       string
	IPv6       string
	Status     string
	IsLoopback bool
}

// DevPortStatus representa o estado de uma porta de desenvolvimento
type DevPortStatus struct {
	Port        int
	Service     string
	IsListening bool
}

// DNSResult representa o resultado de uma consulta de nomes
type DNSResult struct {
	Domain string
	IPs    []string
	CNAME  string
}

// SubnetInfo contém as estatísticas de um bloco CIDR
type SubnetInfo struct {
	CIDR          string
	Netmask       string
	NetworkIP     string
	BroadcastIP   string
	FirstUsableIP string
	LastUsableIP  string
	TotalHosts    uint64
}

// fetchPublicIPInfo tenta consultar o IP público e dados ISP em múltiplos provedores com fallback automático
func fetchPublicIPInfo() (string, string, string) {
	client := &http.Client{
		Timeout: 4 * time.Second,
	}

	// Provedor 1: ipinfo.io
	if ip, isp, loc := tryFetchIPInfo(client, "https://ipinfo.io/json", func(b []byte) (string, string, string) {
		var d struct {
			IP      string `json:"ip"`
			Org     string `json:"org"`
			City    string `json:"city"`
			Country string `json:"country"`
		}
		if json.Unmarshal(b, &d) == nil && d.IP != "" {
			l := ""
			if d.City != "" {
				l = fmt.Sprintf("%s, %s", d.City, d.Country)
			}
			return d.IP, d.Org, l
		}
		return "", "", ""
	}); ip != "" {
		return ip, isp, loc
	}

	// Provedor 2: ipapi.co
	if ip, isp, loc := tryFetchIPInfo(client, "https://ipapi.co/json/", func(b []byte) (string, string, string) {
		var d struct {
			IP      string `json:"ip"`
			Org     string `json:"org"`
			City    string `json:"city"`
			Country string `json:"country_name"`
		}
		if json.Unmarshal(b, &d) == nil && d.IP != "" {
			l := ""
			if d.City != "" {
				l = fmt.Sprintf("%s, %s", d.City, d.Country)
			}
			return d.IP, d.Org, l
		}
		return "", "", ""
	}); ip != "" {
		return ip, isp, loc
	}

	// Provedor 3: api.ipify.org
	if ip, isp, loc := tryFetchIPInfo(client, "https://api.ipify.org?format=json", func(b []byte) (string, string, string) {
		var d struct {
			IP string `json:"ip"`
		}
		if json.Unmarshal(b, &d) == nil && d.IP != "" {
			return d.IP, "", ""
		}
		return "", "", ""
	}); ip != "" {
		return ip, isp, loc
	}

	// Provedor 4: checkip.amazonaws.com (Plain text fallback)
	req, err := http.NewRequest("GET", "https://checkip.amazonaws.com", nil)
	if err == nil {
		req.Header.Set("User-Agent", "TommyCLI/1.0")
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			defer resp.Body.Close()
			buf := make([]byte, 64)
			n, _ := resp.Body.Read(buf)
			ipStr := strings.TrimSpace(string(buf[:n]))
			if net.ParseIP(ipStr) != nil {
				return ipStr, "", ""
			}
		}
	}

	return "", "", ""
}

func tryFetchIPInfo(client *http.Client, urlStr string, parseFn func([]byte) (string, string, string)) (string, string, string) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", "", ""
	}
	req.Header.Set("User-Agent", "TommyCLI/1.0")

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return "", "", ""
	}
	defer resp.Body.Close()

	var buf [4096]byte
	n, _ := resp.Body.Read(buf[:])
	if n == 0 {
		return "", "", ""
	}

	return parseFn(buf[:n])
}

// GetNetworkSummary descobre as informações de rede locais e públicas
func GetNetworkSummary() NetworkSummary {
	summary := NetworkSummary{
		LocalIP:   "Indisponível",
		LocalIPv6: "Indisponível",
		Gateway:   "Desconhecido",
		PublicIP:  "Sem Conexão Externa",
		ISP:       "Desconhecido",
		Location:  "Desconhecido",
		IsOnline:  false,
	}

	// 1. Obter IP Local conectando via UDP a um endereço de referência (não envia dados reais)
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err == nil {
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		summary.LocalIP = localAddr.IP.String()
		summary.IsOnline = true
	}

	// 2. Procurar interface de rede correspondente ao IP local
	interfaces, errIf := net.Interfaces()
	if errIf == nil {
		for _, iface := range interfaces {
			addrs, errAddr := iface.Addrs()
			if errAddr != nil {
				continue
			}
			for _, addr := range addrs {
				ipnet, ok := addr.(*net.IPNet)
				if ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.String() == summary.LocalIP {
						summary.InterfaceName = iface.Name
						summary.MACAddress = iface.HardwareAddr.String()
					} else if ipnet.IP.To4() == nil && summary.LocalIPv6 == "Indisponível" {
						summary.LocalIPv6 = ipnet.IP.String()
					}
				}
			}
		}
	}

	// 3. Obter IP público com resiliência em múltiplos provedores
	if summary.IsOnline {
		pubIP, isp, loc := fetchPublicIPInfo()
		if pubIP != "" {
			summary.PublicIP = pubIP
			if isp != "" {
				summary.ISP = isp
			}
			if loc != "" {
				summary.Location = loc
			}
		}
	}

	return summary
}

// GetAllInterfaces retorna a lista formatada de todos os adaptadores de rede
func GetAllInterfaces() ([]InterfaceInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var result []InterfaceInfo
	for _, iface := range interfaces {
		status := "DOWN"
		if iface.Flags&net.FlagUp != 0 {
			status = "UP"
		}

		info := InterfaceInfo{
			Name:       iface.Name,
			MAC:        iface.HardwareAddr.String(),
			Status:     status,
			IsLoopback: iface.Flags&net.FlagLoopback != 0,
		}

		addrs, errAddr := iface.Addrs()
		if errAddr == nil {
			for _, a := range addrs {
				ipnet, ok := a.(*net.IPNet)
				if ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						info.IPv4 = ipnet.IP.String()
					} else if info.IPv6 == "" {
						info.IPv6 = ipnet.IP.String()
					}
				}
			}
		}

		result = append(result, info)
	}

	return result, nil
}

// TestPing realiza um teste de conexão TCP rápido para um host
func TestPing(host string, port int) (time.Duration, error) {
	target := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", target, 2*time.Second)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return time.Since(start), nil
}

// CheckPort verifica se um Host:Porta específico está aberto
func CheckPort(host string, port int) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", target, 1500*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// ScanDevPorts faz uma varredura nas portas populares de desenvolvimento local
func ScanDevPorts() []DevPortStatus {
	devPorts := []struct {
		Port    int
		Service string
	}{
		{3000, "Node / React / Next.js / Vue"},
		{5173, "Vite / Svelte"},
		{8000, "FastAPI / Django / Python"},
		{8080, "Spring Boot / Express / Nginx"},
		{5432, "PostgreSQL Database"},
		{3306, "MySQL / MariaDB"},
		{27017, "MongoDB Database"},
		{6379, "Redis Cache"},
		{1433, "Microsoft SQL Server"},
		{9090, "Prometheus / Metrics"},
	}

	var results []DevPortStatus
	for _, p := range devPorts {
		isListening := CheckPort("127.0.0.1", p.Port)
		results = append(results, DevPortStatus{
			Port:        p.Port,
			Service:     p.Service,
			IsListening: isListening,
		})
	}

	return results
}

// LookupDNS realiza resolução de nomes DNS para um domínio
func LookupDNS(domain string) (DNSResult, error) {
	domain = strings.TrimSpace(domain)
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.Split(domain, "/")[0]

	result := DNSResult{
		Domain: domain,
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		return result, fmt.Errorf("falha ao resolver DNS para '%s': %v", domain, err)
	}

	for _, ip := range ips {
		result.IPs = append(result.IPs, ip.String())
	}

	cname, errCname := net.LookupCNAME(domain)
	if errCname == nil && cname != domain+"." {
		result.CNAME = strings.TrimSuffix(cname, ".")
	}

	return result, nil
}

// FlushDNS esvazia o cache DNS do sistema operacional
func FlushDNS() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ipconfig", "/flushdns")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("sudo", "killall", "-HUP", "mDNSResponder")
	} else {
		cmd = exec.Command("sudo", "systemd-resolve", "--flush-caches")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao limpar cache DNS: %s", string(out))
	}
	return nil
}

// CalculateSubnet calcula os dados de um bloco CIDR
func CalculateSubnet(cidrStr string) (SubnetInfo, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return SubnetInfo{}, fmt.Errorf("formato CIDR inválido (exemplo correto: 192.168.1.0/24)")
	}

	mask := ipnet.Mask
	ip := ipnet.IP.To4()
	if ip == nil {
		return SubnetInfo{}, fmt.Errorf("suporte apenas para cálculo IPv4")
	}

	netIP := ip.Mask(mask)
	broadcast := make(net.IP, len(netIP))
	for i := range netIP {
		broadcast[i] = netIP[i] | ^mask[i]
	}

	firstUsable := make(net.IP, len(netIP))
	copy(firstUsable, netIP)
	firstUsable[3]++

	lastUsable := make(net.IP, len(broadcast))
	copy(lastUsable, broadcast)
	lastUsable[3]--

	ones, bits := mask.Size()
	totalHosts := uint64(1 << (bits - ones))
	usableHosts := totalHosts - 2
	if totalHosts <= 2 {
		usableHosts = 0
	}

	netmaskStr := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])

	return SubnetInfo{
		CIDR:          cidrStr,
		Netmask:       netmaskStr,
		NetworkIP:     netIP.String(),
		BroadcastIP:   broadcast.String(),
		FirstUsableIP: firstUsable.String(),
		LastUsableIP:  lastUsable.String(),
		TotalHosts:    usableHosts,
	}, nil
}

// ProcessPortInfo guarda os dados de um processo em uma porta
type ProcessPortInfo struct {
	PID         int
	ProcessName string
	Port        int
}

// GetProcessOnPort descobre qual processo está utilizando a porta especificada
func GetProcessOnPort(port int) (*ProcessPortInfo, error) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", fmt.Sprintf("netstat -ano | findstr LISTENING | findstr :%d", port))
		out, err := cmd.CombinedOutput()
		if err != nil || len(out) == 0 {
			return nil, fmt.Errorf("nenhum processo encontrado ouvindo na porta %d", port)
		}

		lines := strings.Split(string(out), "\n")
		var targetPID int
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) >= 5 {
				if strings.HasSuffix(fields[1], fmt.Sprintf(":%d", port)) {
					pid, errConv := strconv.Atoi(fields[len(fields)-1])
					if errConv == nil && pid > 0 {
						targetPID = pid
						break
					}
				}
			}
		}

		if targetPID == 0 {
			return nil, fmt.Errorf("nenhum processo ativo encontrado na porta %d", port)
		}

		procName := getProcessNameByPID(targetPID)

		return &ProcessPortInfo{
			PID:         targetPID,
			ProcessName: procName,
			Port:        port,
		}, nil
	} else {
		cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", port), "-t")
		out, err := cmd.CombinedOutput()
		if err != nil || len(out) == 0 {
			return nil, fmt.Errorf("nenhum processo encontrado na porta %d", port)
		}

		pidStr := strings.TrimSpace(strings.Split(string(out), "\n")[0])
		pid, _ := strconv.Atoi(pidStr)

		return &ProcessPortInfo{
			PID:         pid,
			ProcessName: fmt.Sprintf("Processo (PID %d)", pid),
			Port:        port,
		}, nil
	}
}

func getProcessNameByPID(pid int) string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", fmt.Sprintf("tasklist /FI \"PID eq %d\" /FO CSV /NH", pid))
		out, err := cmd.CombinedOutput()
		if err == nil && len(out) > 0 {
			parts := strings.Split(string(out), ",")
			if len(parts) > 0 {
				name := strings.Trim(parts[0], "\" \r\n")
				if name != "" && !strings.Contains(name, "Nenhum") {
					return name
				}
			}
		}
	}
	return fmt.Sprintf("PID %d", pid)
}

// KillProcessOnPort encerra o processo que está ocupando a porta especificada
func KillProcessOnPort(port int) (*ProcessPortInfo, error) {
	info, err := GetProcessOnPort(port)
	if err != nil {
		return nil, err
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/F", "/PID", fmt.Sprintf("%d", info.PID))
	} else {
		cmd = exec.Command("kill", "-9", fmt.Sprintf("%d", info.PID))
	}

	out, errExec := cmd.CombinedOutput()
	if errExec != nil {
		return nil, fmt.Errorf("falha ao encerrar processo '%s' (PID %d): %s", info.ProcessName, info.PID, string(out))
	}

	return info, nil
}
