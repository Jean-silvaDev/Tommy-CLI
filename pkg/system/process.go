package system

import (
	"sort"
	"strings"
)

var protectedProcessNames = map[string]bool{
	"system":              true,
	"system idle process": true,
	"smss.exe":            true,
	"csrss.exe":           true,
	"wininit.exe":         true,
	"services.exe":        true,
	"lsass.exe":           true,
	"svchost.exe":         true,
	"winlogon.exe":        true,
	"dwm.exe":             true,
	"fontdrvhost.exe":     true,
	"explorer.exe":        true,
	"conhost.exe":         true,
}

// IsProtectedProcess verifica se um processo é vital para o sistema operacional Windows
func IsProtectedProcess(pid int, name string) bool {
	if pid <= 4 {
		return true
	}
	cleanName := strings.ToLower(strings.TrimSpace(name))
	return protectedProcessNames[cleanName]
}

// SortProcessesByMemory ordena uma lista de processos por consumo de memória (decrescente)
func SortProcessesByMemory(procs []ProcessInfo) []ProcessInfo {
	sorted := make([]ProcessInfo, len(procs))
	copy(sorted, procs)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].MemoryBytes > sorted[j].MemoryBytes
	})
	return sorted
}

// FilterProcesses filtra uma lista de processos por nome ou PID
func FilterProcesses(procs []ProcessInfo, query string) []ProcessInfo {
	if strings.TrimSpace(query) == "" {
		return procs
	}
	q := strings.ToLower(strings.TrimSpace(query))

	var filtered []ProcessInfo
	for _, p := range procs {
		if strings.Contains(strings.ToLower(p.Name), q) ||
			strings.Contains(p.Path, q) ||
			strings.Contains(string(rune(p.PID)), q) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}
