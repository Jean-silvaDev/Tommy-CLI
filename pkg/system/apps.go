package system

import (
	"sort"
	"strings"
)

// FilterApps filtra uma lista de aplicativos pelo nome, fabricante ou ID
func FilterApps(apps []AppInfo, query string) []AppInfo {
	if strings.TrimSpace(query) == "" {
		return apps
	}
	q := strings.ToLower(strings.TrimSpace(query))

	var filtered []AppInfo
	for _, app := range apps {
		if strings.Contains(strings.ToLower(app.Name), q) ||
			strings.Contains(strings.ToLower(app.Publisher), q) ||
			strings.Contains(strings.ToLower(app.ID), q) {
			filtered = append(filtered, app)
		}
	}
	return filtered
}

// SortAppsBySize ordena os aplicativos por tamanho estimado no disco (decrescente)
func SortAppsBySize(apps []AppInfo) []AppInfo {
	sorted := make([]AppInfo, len(apps))
	copy(sorted, apps)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].SizeBytes > sorted[j].SizeBytes
	})
	return sorted
}

// SortAppsByName ordena os aplicativos por nome (alfabético)
func SortAppsByName(apps []AppInfo) []AppInfo {
	sorted := make([]AppInfo, len(apps))
	copy(sorted, apps)

	sort.Slice(sorted, func(i, j int) bool {
		return strings.ToLower(sorted[i].Name) < strings.ToLower(sorted[j].Name)
	})
	return sorted
}

// FilterAppsByDrive filtra os aplicativos instalados em um volume específico (ex: "C:")
func FilterAppsByDrive(apps []AppInfo, driveLetter string) []AppInfo {
	if strings.TrimSpace(driveLetter) == "" {
		return apps
	}
	driveUpper := strings.ToUpper(strings.TrimSpace(driveLetter))
	if !strings.HasSuffix(driveUpper, ":") {
		driveUpper += ":"
	}

	var filtered []AppInfo
	for _, app := range apps {
		if strings.HasPrefix(strings.ToUpper(app.InstallLocation), driveUpper) {
			filtered = append(filtered, app)
		}
	}
	return filtered
}
