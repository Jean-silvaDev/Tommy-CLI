package ui

import (
	"testing"
)

func TestRenderTable(t *testing.T) {
	headers := []string{"ID", "NOME", "STACK", "IMAGEM", "STATUS", "ESTADO"}
	rows := [][]string{
		{"68a54b697fa5", "komvos-supabase-auth-1", "komvos-supabase", "supabase/gotrue:v2.186.0", "Up 7 seconds (health: starting)", "RUNNING"},
		{"cee7ba19b6d4", "komvos-supabase-studio-1", "komvos-supabase", "supabase/studio:2026.06.03-sha-0bca601", "Up 7 seconds", "RUNNING"},
	}

	// Verificar se o renderizador executa sem pânicos
	RenderTable(headers, rows)
}

func TestRenderGitLog(t *testing.T) {
	logs := []GitLogItem{
		{
			Hash:    "edae38a",
			Author:  "Jean-silvaDev",
			Date:    "88 seconds ago",
			Subject: "feat: initial commit",
		},
	}

	// Verificar se o renderizador do Git Log executa sem pânicos
	RenderGitLog(logs)
}
