//go:build windows

package system

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

type memoryStatusEx struct {
	cbSize                  uint32
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

type processMemoryCounters struct {
	cb                         uint32
	PageFaultCount             uint32
	PeakWorkingSetSize         uintptr
	WorkingSetSize             uintptr
	QuotaPeakPagedPoolUsage    uintptr
	QuotaPagedPoolUsage        uintptr
	QuotaPeakNonPagedPoolUsage uintptr
	QuotaNonPagedPoolUsage     uintptr
	PagefileUsage              uintptr
	PeakPagefileUsage          uintptr
}

type WindowsSystemProvider struct{}

func NewSystemProvider() SystemProvider {
	return &WindowsSystemProvider{}
}

func (p *WindowsSystemProvider) GetMemoryInfo() (MemoryInfo, error) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	globalMemoryStatusEx := kernel32.NewProc("GlobalMemoryStatusEx")

	var memStatus memoryStatusEx
	memStatus.cbSize = uint32(unsafe.Sizeof(memStatus))

	ret, _, err := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	if ret == 0 {
		return MemoryInfo{}, fmt.Errorf("falha ao consultar GlobalMemoryStatusEx: %w", err)
	}

	total := memStatus.ullTotalPhys
	avail := memStatus.ullAvailPhys
	used := total - avail
	pct := CalculatePercent(used, total)

	return MemoryInfo{
		TotalBytes:     total,
		UsedBytes:      used,
		AvailableBytes: avail,
		UsedPercent:    pct,
	}, nil
}

func (p *WindowsSystemProvider) ListVolumes() ([]VolumeInfo, error) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	getLogicalDrives := kernel32.NewProc("GetLogicalDrives")
	getDiskFreeSpaceEx := kernel32.NewProc("GetDiskFreeSpaceExW")
	getDriveType := kernel32.NewProc("GetDriveTypeW")
	getVolumeInformation := kernel32.NewProc("GetVolumeInformationW")

	drivesBitmask, _, _ := getLogicalDrives.Call()
	if drivesBitmask == 0 {
		return nil, fmt.Errorf("nenhuma unidade de disco encontrada no sistema")
	}

	var volumes []VolumeInfo

	for i := 0; i < 26; i++ {
		if (drivesBitmask & (1 << i)) != 0 {
			driveLetter := fmt.Sprintf("%c:", 'A'+i)
			driveRoot := driveLetter + `\`
			driveRootUTF16, err := syscall.UTF16PtrFromString(driveRoot)
			if err != nil {
				continue
			}

			dtRet, _, _ := getDriveType.Call(uintptr(unsafe.Pointer(driveRootUTF16)))
			driveTypeStr := "Desconhecido"
			diskTypeStr := "SSD/HDD"

			switch dtRet {
			case 2:
				driveTypeStr = "Removível"
				diskTypeStr = "Pendrive/Removível"
			case 3:
				driveTypeStr = "Fixo"
				diskTypeStr = "SSD/HDD"
			case 4:
				driveTypeStr = "Rede"
				diskTypeStr = "Unidade de Rede"
			case 5:
				driveTypeStr = "CD/DVD"
				diskTypeStr = "Mídia Óptica"
			case 6:
				driveTypeStr = "RAM Disk"
				diskTypeStr = "Memória RAM"
			}

			var freeBytesAvail, totalBytes, totalFreeBytes uint64
			retSpace, _, _ := getDiskFreeSpaceEx.Call(
				uintptr(unsafe.Pointer(driveRootUTF16)),
				uintptr(unsafe.Pointer(&freeBytesAvail)),
				uintptr(unsafe.Pointer(&totalBytes)),
				uintptr(unsafe.Pointer(&totalFreeBytes)),
			)

			if retSpace == 0 || totalBytes == 0 {
				continue
			}

			var volumeNameBuf [256]uint16
			var fileSystemBuf [256]uint16
			var serialNumber, maxComponentLen, flags uint32

			_, _, _ = getVolumeInformation.Call(
				uintptr(unsafe.Pointer(driveRootUTF16)),
				uintptr(unsafe.Pointer(&volumeNameBuf[0])),
				uintptr(256),
				uintptr(unsafe.Pointer(&serialNumber)),
				uintptr(unsafe.Pointer(&maxComponentLen)),
				uintptr(unsafe.Pointer(&flags)),
				uintptr(unsafe.Pointer(&fileSystemBuf[0])),
				uintptr(256),
			)

			volName := syscall.UTF16ToString(volumeNameBuf[:])
			fileSys := syscall.UTF16ToString(fileSystemBuf[:])
			if volName == "" {
				volName = "Disco Local"
			}
			if fileSys == "" {
				fileSys = "NTFS"
			}

			usedBytes := totalBytes - totalFreeBytes
			pct := CalculatePercent(usedBytes, totalBytes)

			volumes = append(volumes, VolumeInfo{
				Device:       driveLetter,
				Label:        volName,
				DriveType:    driveTypeStr,
				DiskType:     diskTypeStr,
				TotalBytes:   totalBytes,
				UsedBytes:    usedBytes,
				FreeBytes:    totalFreeBytes,
				UsedPercent:  pct,
				FileSystem:   fileSys,
				Status:       "OK",
				HealthStatus: "OK (Saudável)",
			})
		}
	}

	return volumes, nil
}

func (p *WindowsSystemProvider) ListProcesses() ([]ProcessInfo, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar snapshot de processos: %w", err)
	}
	defer func() { _ = windows.CloseHandle(snapshot) }()

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))

	err = windows.Process32First(snapshot, &pe)
	if err != nil {
		return nil, fmt.Errorf("falha ao obter primeiro processo: %w", err)
	}

	memInfo, _ := p.GetMemoryInfo()
	psapi := windows.NewLazySystemDLL("psapi.dll")
	getProcessMemoryInfo := psapi.NewProc("GetProcessMemoryInfo")

	var processes []ProcessInfo

	for {
		pid := int(pe.ProcessID)
		name := syscall.UTF16ToString(pe.ExeFile[:])

		var memoryBytes uint64
		var path string

		hProc, errOpen := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION|windows.PROCESS_VM_READ, false, pe.ProcessID)
		if errOpen == nil {
			var counters processMemoryCounters
			counters.cb = uint32(unsafe.Sizeof(counters))
			ret, _, _ := getProcessMemoryInfo.Call(
				uintptr(hProc),
				uintptr(unsafe.Pointer(&counters)),
				uintptr(counters.cb),
			)
			if ret != 0 {
				memoryBytes = uint64(counters.WorkingSetSize)
			}

			var buf [1024]uint16
			size := uint32(len(buf))
			if errPath := windows.QueryFullProcessImageName(hProc, 0, &buf[0], &size); errPath == nil {
				path = syscall.UTF16ToString(buf[:size])
			}
			_ = windows.CloseHandle(hProc)
		}

		memPct := CalculatePercent(memoryBytes, memInfo.TotalBytes)
		isProt := IsProtectedProcess(pid, name)

		processes = append(processes, ProcessInfo{
			PID:           pid,
			Name:          name,
			Path:          path,
			User:          "Sistema / Usuário",
			MemoryBytes:   memoryBytes,
			MemoryPercent: memPct,
			Threads:       int(pe.Threads),
			Status:        "Em Execução",
			IsProtected:   isProt,
		})

		err = windows.Process32Next(snapshot, &pe)
		if err != nil {
			break
		}
	}

	return processes, nil
}

func (p *WindowsSystemProvider) KillProcess(pid int) error {
	if pid <= 4 {
		return fmt.Errorf("o processo (PID %d) é um processo vital do sistema operacional e está protegido", pid)
	}

	procs, err := p.ListProcesses()
	if err == nil {
		for _, pr := range procs {
			if pr.PID == pid && pr.IsProtected {
				return fmt.Errorf("o processo '%s' (PID %d) é vital para a estabilidade do sistema e não pode ser encerrado", pr.Name, pid)
			}
		}
	}

	hProc, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return fmt.Errorf("falha ao abrir o processo (PID %d) para encerramento (verifique privilégios): %w", pid, err)
	}
	defer func() { _ = windows.CloseHandle(hProc) }()

	err = windows.TerminateProcess(hProc, 1)
	if err != nil {
		return fmt.Errorf("falha ao encerrar o processo (PID %d): %w", pid, err)
	}

	return nil
}

func (p *WindowsSystemProvider) ListApplications() ([]AppInfo, error) {
	registryRoots := []struct {
		Key  registry.Key
		Path string
		Src  string
	}{
		{registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, "Registry (HKLM)"},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`, "Registry (HKLM32)"},
		{registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`, "Registry (HKCU)"},
	}

	var apps []AppInfo
	seenNames := make(map[string]bool)

	for _, root := range registryRoots {
		key, err := registry.OpenKey(root.Key, root.Path, registry.READ)
		if err != nil {
			continue
		}

		subkeys, errSub := key.ReadSubKeyNames(-1)
		_ = key.Close()
		if errSub != nil {
			continue
		}

		for _, sk := range subkeys {
			subKeyPath := root.Path + `\` + sk
			appKey, errApp := registry.OpenKey(root.Key, subKeyPath, registry.READ)
			if errApp != nil {
				continue
			}

			dispName, _, _ := appKey.GetStringValue("DisplayName")
			dispName = strings.TrimSpace(dispName)
			if dispName == "" || seenNames[dispName] {
				_ = appKey.Close()
				continue
			}

			seenNames[dispName] = true
			dispVersion, _, _ := appKey.GetStringValue("DisplayVersion")
			publisher, _, _ := appKey.GetStringValue("Publisher")
			installLoc, _, _ := appKey.GetStringValue("InstallLocation")
			uninstStr, _, _ := appKey.GetStringValue("UninstallString")
			estSize, _, errEst := appKey.GetIntegerValue("EstimatedSize")

			var sizeBytes uint64
			if errEst == nil {
				sizeBytes = estSize * 1024 // EstimatedSize é registrado em KB no Registro do Windows
			}

			_ = appKey.Close()

			apps = append(apps, AppInfo{
				ID:              sk,
				Name:            dispName,
				Version:         dispVersion,
				Publisher:       publisher,
				SizeBytes:       sizeBytes,
				InstallLocation: installLoc,
				UninstallString: uninstStr,
				Source:          root.Src,
			})
		}
	}

	return apps, nil
}

func (p *WindowsSystemProvider) UninstallApplication(app AppInfo) error {
	if strings.TrimSpace(app.UninstallString) == "" {
		return fmt.Errorf("o aplicativo '%s' não possui uma instrução registrada de desinstalação automática", app.Name)
	}

	uninstCmd := strings.TrimSpace(app.UninstallString)
	cmd := exec.Command("cmd", "/c", uninstCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("falha ao executar a desinstalação de '%s': %w", app.Name, err)
	}

	return nil
}

func (p *WindowsSystemProvider) AnalyzePath(path string) (*FileNode, error) {
	return AnalyzePath(path)
}

func (p *WindowsSystemProvider) GetLargestFiles(path string, limit int) ([]FileNode, error) {
	return GetLargestFiles(path, limit)
}

func (p *WindowsSystemProvider) GetCPUUsage() (float64, error) {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	getSystemTimes := kernel32.NewProc("GetSystemTimes")

	var idleTime, kernelTime, userTime syscall.Filetime
	ret, _, _ := getSystemTimes.Call(
		uintptr(unsafe.Pointer(&idleTime)),
		uintptr(unsafe.Pointer(&kernelTime)),
		uintptr(unsafe.Pointer(&userTime)),
	)

	if ret == 0 {
		return 0.0, fmt.Errorf("falha ao obter tempos do sistema")
	}

	idle := filetimeToUint64(idleTime)
	kernel := filetimeToUint64(kernelTime)
	user := filetimeToUint64(userTime)

	total := kernel + user
	if total == 0 {
		return 0.0, nil
	}

	busy := total - idle
	pct := (float64(busy) / float64(total)) * 100.0
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	return pct, nil
}

func filetimeToUint64(ft syscall.Filetime) uint64 {
	return uint64(ft.HighDateTime)<<32 | uint64(ft.LowDateTime)
}

func (p *WindowsSystemProvider) GetSystemSummary() (SystemSummary, error) {
	mem, errMem := p.GetMemoryInfo()
	if errMem != nil {
		return SystemSummary{}, errMem
	}

	vols, _ := p.ListVolumes()
	cpu, _ := p.GetCPUUsage()

	summary := SystemSummary{
		Memory:  mem,
		CPU:     CPUInfo{UsagePercent: cpu},
		Volumes: vols,
	}

	summary.Alerts = EvaluateAlerts(summary)
	return summary, nil
}
