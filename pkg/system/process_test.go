package system

import (
	"testing"
)

func TestIsProtectedProcess(t *testing.T) {
	if !IsProtectedProcess(0, "System Idle Process") {
		t.Errorf("PID 0 deveria ser protegido")
	}
	if !IsProtectedProcess(4, "System") {
		t.Errorf("PID 4 (System) deveria ser protegido")
	}
	if !IsProtectedProcess(1234, "csrss.exe") {
		t.Errorf("csrss.exe deveria ser protegido")
	}
	if IsProtectedProcess(9999, "my_custom_app.exe") {
		t.Errorf("my_custom_app.exe NAO deveria ser protegido")
	}
}

func TestSortProcessesByMemory(t *testing.T) {
	procs := []ProcessInfo{
		{PID: 1, Name: "app1", MemoryBytes: 100},
		{PID: 2, Name: "app2", MemoryBytes: 500},
		{PID: 3, Name: "app3", MemoryBytes: 250},
	}

	sorted := SortProcessesByMemory(procs)
	if sorted[0].MemoryBytes != 500 || sorted[1].MemoryBytes != 250 || sorted[2].MemoryBytes != 100 {
		t.Errorf("Ordenação incorreta por memória: %+v", sorted)
	}
}

func TestFilterProcesses(t *testing.T) {
	procs := []ProcessInfo{
		{PID: 101, Name: "chrome.exe"},
		{PID: 202, Name: "docker.exe"},
		{PID: 303, Name: "code.exe"},
	}

	res := FilterProcesses(procs, "chrome")
	if len(res) != 1 || res[0].Name != "chrome.exe" {
		t.Errorf("Filtro incorreto para 'chrome': %+v", res)
	}
}
