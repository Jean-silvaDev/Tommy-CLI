package system

import (
	"testing"
)

func TestFilterAndSortApps(t *testing.T) {
	apps := []AppInfo{
		{Name: "Google Chrome", Publisher: "Google LLC", SizeBytes: 500000, InstallLocation: `C:\Program Files\Google`},
		{Name: "Visual Studio Code", Publisher: "Microsoft", SizeBytes: 1500000, InstallLocation: `C:\Users\Jean\AppData`},
		{Name: "Docker Desktop", Publisher: "Docker Inc", SizeBytes: 3000000, InstallLocation: `D:\Docker`},
	}

	// Teste filtro por nome
	chrome := FilterApps(apps, "chrome")
	if len(chrome) != 1 || chrome[0].Name != "Google Chrome" {
		t.Errorf("Filtro por nome falhou: %+v", chrome)
	}

	// Teste ordenação por tamanho
	sortedSize := SortAppsBySize(apps)
	if sortedSize[0].Name != "Docker Desktop" {
		t.Errorf("Ordenação por tamanho falhou: %+v", sortedSize)
	}

	// Teste filtro por disco
	dApps := FilterAppsByDrive(apps, "D:")
	if len(dApps) != 1 || dApps[0].Name != "Docker Desktop" {
		t.Errorf("Filtro por disco falhou: %+v", dApps)
	}
}
