package cmd

import (
	"fmt"
	"os"
	"strings"

	"tommy/pkg/git"
	"tommy/pkg/ui"

	"github.com/spf13/cobra"
)

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Assistente de fluxo de trabalho Git (commits, logs, branches, tags, push)",
	Long:  `Subcomando interativo para automação de tarefas do Git: commits padronizados em estilo Conventional Commits, histórico de logs, gerenciamento de branches, tags e push.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGitDashboard()
	},
}

var gitCommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Cria um commit padronizado via assistente interativo",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Novo Commit (Conventional Commits)")
		runInteractiveCommitWizard()
	},
}

var gitLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Visualiza o histórico recente de commits",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Histórico de Commits Recente")
		displayGitLog(10)
	},
}

var gitBranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Gerencia branches do repositório",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Gerenciador de Branches")
		runBranchWizard()
	},
}

var gitTagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Gerencia tags do repositório",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Gerenciador de Tags")
		runTagWizard()
	},
}

var gitPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Baixa e aplica as alterações do repositório remoto",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Git Pull")
		if !ui.ConfirmPrompt("Deseja realmente atualizar (PULL) as alterações do repositório remoto?") {
			ui.PrintInfo("Operação de Pull cancelada.")
			return
		}
		ui.PrintInfo("Buscando e aplicando alterações do repositório remoto...")
		err := git.Pull()
		if err != nil {
			ui.PrintError("%v", err)
			conflicts, _ := git.GetConflictedFiles()
			if len(conflicts) > 0 {
				fmt.Println()
				ui.PrintWarning("Foram detectados CONFLITOS DE MERGE em %d arquivo(s)!", len(conflicts))
				if ui.ConfirmPrompt("Deseja abrir o Assistente de Resolução de Conflitos agora?") {
					runConflictWizard(conflicts)
				}
			}
			return
		}
		ui.PrintSuccess("Pull realizado com sucesso do repositório remoto!")
	},
}

var gitPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Envia os commits locais para o repositório remoto",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Git Push")
		if !ui.ConfirmPrompt("Deseja realmente enviar (PUSH) as alterações locais para o repositório remoto?") {
			ui.PrintInfo("Operação de Push cancelada.")
			return
		}
		ui.PrintInfo("Enviando alterações para o repositório remoto...")
		err := git.Push(true)
		if err != nil {
			ui.PrintError("%v", err)
			return
		}
		ui.PrintSuccess("Push realizado com sucesso para o remoto!")
	},
}

func getCurrentWorkingDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir
}

func runGitDashboard() {
	for {
		ui.PrintBanner("Assistente Git")

		if !git.IsGitRepo() {
			ui.PrintWarning("O diretório atual '%s' NÃO é um repositório Git.", getCurrentWorkingDir())
			ui.PrintInfo("Inicialize um repositório ou vincule a um repositório remoto para utilizar o assistente.\n")

			options := []string{
				"1. 🚀 Inicializar Repositório Git (git init)",
				"2. 🔗 Vincular Repositório Remoto (git remote origin)",
				"3. ⬅️  Voltar ao Menu Principal",
				"4. 🚪 Sair da Aplicação",
			}

			idx, _, err := ui.SelectOption("Escolha uma ação para configurar o Git nesta pasta", options)
			if err != nil || idx == 2 {
				break // Voltar ao Menu Principal
			}

			if idx == 3 {
				ui.PrintExitMessage()
				os.Exit(0)
			}

			switch idx {
			case 0:
				if ui.ConfirmPrompt("Deseja inicializar um repositório Git nesta pasta?") {
					err := git.InitRepo()
					if err != nil {
						ui.PrintError("%v", err)
					} else {
						ui.PrintSuccess("Repositório Git inicializado com sucesso!")
					}
					ui.WaitForEnter("")
				}
			case 1:
				runRemoteWizard(true)
			}
			continue
		}

		// Diretório já é um repositório Git
		remoteURL, errRemote := git.GetRemoteURL()
		remoteStr := "Nenhum"
		if errRemote == nil {
			remoteStr = remoteURL
		}

		_, currentBranch, _ := git.GetBranches()
		if currentBranch == "" {
			currentBranch = "Sem branch (repositório recente)"
		}

		conflictedFiles, _ := git.GetConflictedFiles()
		hasConflicts := len(conflictedFiles) > 0

		if hasConflicts {
			ui.PrintError("🚨 ATENÇÃO: CONFLITOS DE MERGE EM ANDAMENTO (%d arquivo(s) afetados)!", len(conflictedFiles))
			ui.PrintInfo("Utilize o Assistente de Conflitos para visualizar, resolver ou abortar o merge.\n")
		} else {
			ui.PrintInfo("Branch Ativa: %s | 🔗 Remoto Origin: %s\n",
				ui.SuccessStyle.Render(currentBranch),
				ui.InfoStyle.Render(remoteStr))
		}

		var options []string
		if hasConflicts {
			options = []string{
				"🚨 1. RESOLVER CONFLITOS DE MERGE (Assistente Interativo)",
				"2. 📜 Ver Histórico de Commits (Log)",
				"3. 🌿 Gerenciar Branches (Trocar / Criar / Renomear / Excluir)",
				"4. ⬅️  Voltar ao Menu Principal",
				"5. 🚪 Sair da Aplicação",
			}
		} else {
			options = []string{
				"1. 📝 Criar Commit Interativo (Conventional Commits)",
				"2. 📜 Ver Histórico de Commits (Log)",
				"3. 🌿 Gerenciar Branches (Trocar / Criar / Renomear / Excluir)",
				"4. 🏷️  Gerenciar Tags (Listar / Criar)",
				"5. 📥 Dar Pull no Repositório Remoto (git pull)",
				"6. 🚀 Dar Push no Repositório Remoto (git push)",
				"7. 🔗 Vincular / Alterar Repositório Remoto",
				"8. ⬅️  Voltar ao Menu Principal",
				"9. 🚪 Sair da Aplicação",
			}
		}

		idx, _, err := ui.SelectOption("Escolha uma ação do Git", options)

		if hasConflicts {
			if err != nil || idx == 3 {
				break
			}
			if idx == 4 {
				ui.PrintExitMessage()
				os.Exit(0)
			}
			switch idx {
			case 0:
				runConflictWizard(conflictedFiles)
			case 1:
				displayGitLog(10)
			case 2:
				runBranchWizard()
			}
			continue
		}

		// Sem conflitos - fluxo normal
		if err != nil || idx == 7 {
			break // Voltar ao Menu Principal
		}

		if idx == 8 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			runInteractiveCommitWizard()
		case 1:
			displayGitLog(10)
		case 2:
			runBranchWizard()
		case 3:
			runTagWizard()
		case 4:
			gitPullCmd.Run(nil, nil)
			ui.WaitForEnter("")
		case 5:
			gitPushCmd.Run(nil, nil)
			ui.WaitForEnter("")
		case 6:
			runRemoteWizard(false)
		}
	}
}

func runConflictWizard(conflictedFiles []string) {
	for {
		ui.PrintBanner("Assistente de Conflitos de Merge")

		ui.PrintError("🚨 CONFLITO DE MERGE DETECTADO em %d arquivo(s)!", len(conflictedFiles))
		ui.PrintInfo("Arquivos afetados com marcações de conflito:")
		for _, f := range conflictedFiles {
			fmt.Printf("  • %s\n", ui.ErrorStyle.Render(f))
		}
		fmt.Println()

		options := []string{
			"1. 📜 Ver Guia de Resolução de Conflitos",
			"2. ➕ Marcar Arquivos Resolvidos no Stage (git add -A)",
			"3. ✅ Concluir Merge com Commit (git commit)",
			"4. 🚫 Abortar Merge e Restaurar Estado Anterior (git merge --abort)",
			"⬅️  5. Voltar ao Assistente Git",
		}

		idx, _, err := ui.SelectOption("Escolha uma ação para gerenciar os conflitos", options)
		if err != nil || idx == 4 {
			return
		}

		switch idx {
		case 0:
			fmt.Println(ui.SubtitleStyle.Render("Guia Passo a Passo para Resolver Conflitos:"))
			ui.PrintInfo("1. Abra os arquivos listados no seu editor (ex: VS Code).")
			ui.PrintInfo("2. Procure por marcações '<<<<<<<', '=======' e '>>>>>>>'.")
			ui.PrintInfo("3. Escolha quais linhas manter, remova as marcações e salve os arquivos.")
			ui.PrintInfo("4. Volte aqui e selecione a Opção 2 (Marcar Arquivos Resolvidos).")
			fmt.Println()
			ui.WaitForEnter("")
		case 1:
			if ui.ConfirmPrompt("Deseja adicionar todos os arquivos resolvidos ao Stage (git add -A)?") {
				err := git.StageAll()
				if err != nil {
					ui.PrintError("Erro ao fazer stage: %v", err)
				} else {
					ui.PrintSuccess("Arquivos adicionados ao Stage com sucesso!")
				}
				ui.WaitForEnter("")
			}
		case 2:
			remaining, _ := git.GetConflictedFiles()
			if len(remaining) > 0 {
				ui.PrintError("Ainda existem %d arquivo(s) com conflitos não resolvidos!", len(remaining))
				ui.WaitForEnter("")
				continue
			}
			commitMsg, _ := ui.PromptText("Mensagem do commit de merge", "merge: resolve merge conflicts")
			err := git.CreateCommit(commitMsg)
			if err != nil {
				ui.PrintError("Erro ao concluir commit de merge: %v", err)
			} else {
				ui.PrintSuccess("Merge concluído com sucesso! Conflitos resolvidos.")
				ui.WaitForEnter("")
				return
			}
			ui.WaitForEnter("")
		case 3:
			if ui.ConfirmPrompt("Deseja realmente ABORTAR o merge? Todas as alterações não resolvidas do merge serão desfeitas.") {
				err := git.AbortMerge()
				if err != nil {
					ui.PrintError("%v", err)
				} else {
					ui.PrintSuccess("Merge abortado com sucesso! Repositório restaurado ao estado anterior.")
					ui.WaitForEnter("")
					return
				}
				ui.WaitForEnter("")
			}
		}

		conflictedFiles, _ = git.GetConflictedFiles()
		if len(conflictedFiles) == 0 {
			ui.PrintSuccess("Todos os conflitos foram resolvidos!")
			ui.WaitForEnter("")
			return
		}
	}
}

func runRemoteWizard(autoInit bool) {
	if autoInit && !git.IsGitRepo() {
		if ui.ConfirmPrompt("O diretório atual não é um repositório Git. Deseja inicializar (git init) primeiro?") {
			err := git.InitRepo()
			if err != nil {
				ui.PrintError("%v", err)
				ui.WaitForEnter("")
				return
			}
			ui.PrintSuccess("Repositório Git inicializado com sucesso!")
		} else {
			ui.PrintInfo("Operação cancelada.")
			ui.WaitForEnter("")
			return
		}
	}

	currentRemote, _ := git.GetRemoteURL()
	if currentRemote != "" {
		ui.PrintInfo("Remoto 'origin' atual: %s", currentRemote)
	}

	remoteURL, err := ui.PromptText("Informe a URL do repositório remoto (ex: https://github.com/usuario/repositorio.git)", currentRemote)
	if err != nil || strings.TrimSpace(remoteURL) == "" {
		ui.PrintError("URL remota é obrigatória!")
		ui.WaitForEnter("")
		return
	}
	remoteURL = strings.TrimSpace(remoteURL)

	if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente configurar '%s' como o remoto 'origin'?", remoteURL)) {
		ui.PrintInfo("Configuração do remoto cancelada.")
		ui.WaitForEnter("")
		return
	}

	err = git.AddOrSetRemote(remoteURL)
	if err != nil {
		ui.PrintError("%v", err)
	} else {
		ui.PrintSuccess("Repositório remoto 'origin' configurado para '%s'!", remoteURL)
	}
	ui.WaitForEnter("")
}

func runInteractiveCommitWizard() {
	status, err := git.GetStatus()
	if err != nil {
		ui.PrintError("%v", err)
		ui.WaitForEnter("")
		return
	}

	if status == "" {
		ui.PrintWarning("Nenhuma alteração detectada no repositório Git.")
		ui.WaitForEnter("")
		return
	}

	fmt.Println(ui.SubtitleStyle.Render("Alterações Detectadas:"))
	fmt.Println(status)
	fmt.Println()

	if ui.ConfirmPrompt("Deseja adicionar TODAS as alterações ao Stage (git add -A)?") {
		err := git.StageAll()
		if err != nil {
			ui.PrintError("Erro ao fazer stage: %v", err)
			ui.WaitForEnter("")
			return
		}
		ui.PrintSuccess("Arquivos adicionados ao Stage!")
	}

	commitTypes := []string{
		"feat: Nova funcionalidade para o usuário",
		"fix: Correção de bug ou erro",
		"chore: Alterações de build, dependências ou tarefas secundárias",
		"docs: Alterações apenas em documentação",
		"refactor: Mudança de código que não corrige bug nem adiciona funcionalidade",
		"test: Adição ou correção de testes",
		"style: Mudanças de formatação de código (espaços, vírgulas, etc)",
	}

	typeIdx, _, err := ui.SelectOption("Selecione o TIPO do commit", commitTypes)
	if err != nil {
		return
	}

	typePrefix := strings.Split(commitTypes[typeIdx], ":")[0]

	scope, _ := ui.PromptText("Informe o ESCOPO (opcional, ex: cli, cleaner, docker)", "")
	message, err := ui.PromptText("Informe o TÍTULO/MENSAGEM do commit (ou cole o texto completo)", "")
	if err != nil || message == "" {
		ui.PrintError("A mensagem do commit é obrigatória!")
		ui.WaitForEnter("")
		return
	}

	var body string
	if !strings.Contains(message, " - ") && !strings.Contains(message, "\n") {
		body, _ = ui.PromptText("Informe a DESCRIÇÃO/CORPO detalhado (opcional, use ' - ' para tópicos)", "")
	}

	formattedMsg := FormatCommitMessage(typePrefix, scope, message, body)

	fmt.Println()
	ui.PrintInfo("Mensagem final do commit:\n%s\n", formattedMsg)

	if ui.ConfirmPrompt("Confirmar e realizar este commit?") {
		err := git.CreateCommit(formattedMsg)
		if err != nil {
			ui.PrintError("%v", err)
			ui.WaitForEnter("")
			return
		}
		ui.PrintSuccess("Commit criado com sucesso!")

		if ui.ConfirmPrompt("Deseja fazer PUSH para o remoto agora?") {
			gitPushCmd.Run(nil, nil)
		}
		ui.WaitForEnter("")
	}
}

func FormatCommitMessage(typePrefix, scope, message, body string) string {
	message = strings.TrimSpace(message)
	var subject string
	var bullets []string

	if strings.Contains(message, " - ") {
		parts := strings.Split(message, " - ")
		subject = strings.TrimSpace(parts[0])
		for _, p := range parts[1:] {
			p = strings.TrimSpace(p)
			if p != "" {
				bullets = append(bullets, p)
			}
		}
	} else if strings.Contains(message, "\n") {
		lines := strings.Split(message, "\n")
		subject = strings.TrimSpace(lines[0])
		for _, l := range lines[1:] {
			l = strings.TrimSpace(l)
			if l != "" {
				bullets = append(bullets, l)
			}
		}
	} else {
		subject = message
	}

	if strings.TrimSpace(body) != "" {
		bodyClean := strings.TrimSpace(body)
		if strings.Contains(bodyClean, " - ") {
			parts := strings.Split(bodyClean, " - ")
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p != "" {
					bullets = append(bullets, p)
				}
			}
		} else if strings.Contains(bodyClean, "\n") {
			lines := strings.Split(bodyClean, "\n")
			for _, l := range lines {
				l = strings.TrimSpace(l)
				if l != "" {
					bullets = append(bullets, l)
				}
			}
		} else {
			bullets = append(bullets, bodyClean)
		}
	}

	var formattedSubject string
	if scope != "" {
		formattedSubject = fmt.Sprintf("%s(%s): %s", typePrefix, scope, subject)
	} else {
		formattedSubject = fmt.Sprintf("%s: %s", typePrefix, subject)
	}

	if len(bullets) == 0 {
		return formattedSubject
	}

	var bodyLines []string
	for _, b := range bullets {
		bClean := strings.TrimPrefix(b, "-")
		bClean = strings.TrimPrefix(bClean, "•")
		bClean = strings.TrimSpace(bClean)
		bodyLines = append(bodyLines, fmt.Sprintf("• %s", bClean))
	}

	return fmt.Sprintf("%s\n\n%s", formattedSubject, strings.Join(bodyLines, "\n"))
}

func displayGitLog(count int) {
	logs, err := git.GetRecentLogs(count)
	if err != nil {
		ui.PrintError("%v", err)
		ui.WaitForEnter("")
		return
	}

	if len(logs) == 0 {
		ui.PrintWarning("Nenhum commit encontrado no repositório.")
		ui.WaitForEnter("")
		return
	}

	var gitLogs []ui.GitLogItem
	for _, l := range logs {
		gitLogs = append(gitLogs, ui.GitLogItem{
			Hash:    l.Hash,
			Author:  l.Author,
			Date:    l.Date,
			Subject: l.Subject,
		})
	}

	ui.RenderGitLog(gitLogs)
	ui.WaitForEnter("")
}

func runBranchWizard() {
	for {
		branches, current, err := git.GetBranches()
		if err != nil {
			ui.PrintError("%v", err)
			return
		}

		ui.PrintBanner("Gerenciador de Branches")
		ui.PrintInfo("Branch ativa atual: %s", ui.SuccessStyle.Render(current))
		fmt.Println()

		var options []string
		for _, b := range branches {
			if b == current {
				options = append(options, fmt.Sprintf("⭐ %s (Ativa)", b))
			} else {
				options = append(options, fmt.Sprintf("🌿 %s", b))
			}
		}

		options = append(options, "➕ Criar Nova Branch")
		options = append(options, "⬅️  Voltar ao Assistente Git")

		idx, _, err := ui.SelectOption("Selecione uma branch para gerenciar ou crie uma nova", options)
		if err != nil || idx == len(options)-1 {
			return
		}

		if idx == len(options)-2 { // Criar Nova Branch
			newName, err := ui.PromptText("Nome da nova branch", "")
			if err != nil || newName == "" {
				ui.PrintError("Nome inválido para a branch!")
				continue
			}
			if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente criar a branch '%s'?", newName)) {
				ui.PrintInfo("Criação de branch cancelada.")
				continue
			}
			err = git.CreateBranch(newName)
			if err != nil {
				ui.PrintError("%v", err)
			} else {
				ui.PrintSuccess("Branch '%s' criada e ativada!", newName)
			}
			continue
		}

		selectedBranch := branches[idx]
		isCurrent := (selectedBranch == current)

		handleBranchAction(selectedBranch, isCurrent)
	}
}

func handleBranchAction(targetBranch string, isCurrent bool) {
	var actionOptions []string
	if isCurrent {
		actionOptions = []string{
			"✏️  Renomear esta branch",
			"⬅️  Voltar",
		}
	} else {
		actionOptions = []string{
			"🔀 Alternar (Checkout) para esta branch",
			"✏️  Renomear esta branch",
			"🗑️  Excluir esta branch",
			"⬅️  Voltar",
		}
	}

	header := fmt.Sprintf("Ações para a branch '%s'", targetBranch)
	if isCurrent {
		header += " (Ativa)"
	}

	idx, _, err := ui.SelectOption(header, actionOptions)
	if err != nil {
		return
	}

	if isCurrent {
		switch idx {
		case 0:
			renameBranchFlow(targetBranch)
		case 1:
			return
		}
	} else {
		switch idx {
		case 0:
			switchBranchFlow(targetBranch)
		case 1:
			renameBranchFlow(targetBranch)
		case 2:
			deleteBranchFlow(targetBranch)
		case 3:
			return
		}
	}
}

func switchBranchFlow(targetBranch string) {
	if !ui.ConfirmPrompt(fmt.Sprintf("Deseja alternar para a branch '%s'?", targetBranch)) {
		ui.PrintInfo("Alternância de branch cancelada.")
		return
	}
	err := git.SwitchBranch(targetBranch)
	if err != nil {
		ui.PrintError("%v", err)
	} else {
		ui.PrintSuccess("Alternado com sucesso para a branch '%s'!", targetBranch)
	}
}

func renameBranchFlow(targetBranch string) {
	newName, err := ui.PromptText(fmt.Sprintf("Novo nome para a branch '%s'", targetBranch), targetBranch)
	if err != nil || strings.TrimSpace(newName) == "" {
		ui.PrintError("Nome inválido para a branch!")
		return
	}
	newName = strings.TrimSpace(newName)
	if newName == targetBranch {
		ui.PrintInfo("O novo nome é idêntico ao nome atual.")
		return
	}

	if !ui.ConfirmPrompt(fmt.Sprintf("Deseja renomear a branch '%s' para '%s'?", targetBranch, newName)) {
		ui.PrintInfo("Renomeação cancelada.")
		return
	}

	err = git.RenameBranch(targetBranch, newName)
	if err != nil {
		ui.PrintError("%v", err)
	} else {
		ui.PrintSuccess("Branch '%s' renomeada para '%s' com sucesso!", targetBranch, newName)
	}
}

func deleteBranchFlow(targetBranch string) {
	if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente EXCLUIR a branch '%s'?", targetBranch)) {
		ui.PrintInfo("Exclusão de branch cancelada.")
		return
	}

	err := git.DeleteBranch(targetBranch, false)
	if err != nil {
		ui.PrintWarning("Não foi possível excluir suavemente: %v", err)
		if ui.ConfirmPrompt(fmt.Sprintf("Deseja FORÇAR a exclusão (-D) da branch '%s'?", targetBranch)) {
			errForce := git.DeleteBranch(targetBranch, true)
			if errForce != nil {
				ui.PrintError("%v", errForce)
			} else {
				ui.PrintSuccess("Branch '%s' excluída forçadamente com sucesso!", targetBranch)
			}
		}
	} else {
		ui.PrintSuccess("Branch '%s' excluída com sucesso!", targetBranch)
	}
}

func runTagWizard() {
	for {
		tags, err := git.GetTags()
		if err != nil {
			ui.PrintError("%v", err)
			return
		}

		ui.PrintBanner("Gerenciador de Tags")
		if len(tags) == 0 {
			ui.PrintWarning("Nenhuma tag existente no repositório.\n")
		} else {
			ui.PrintInfo("Tags encontradas: %s\n", strings.Join(tags, ", "))
		}

		options := []string{
			"🏷️  Criar Nova Tag",
			"⬅️  Voltar ao Assistente Git",
		}

		idx, _, err := ui.SelectOption("Ação desejada com Tags", options)
		if err != nil || idx == 1 {
			return
		}

		tagName, err := ui.PromptText("Nome da tag (ex: v1.0.0)", "")
		if err != nil || tagName == "" {
			ui.PrintError("Nome de tag é obrigatório!")
			continue
		}

		tagMsg, _ := ui.PromptText("Descrição/Mensagem da tag (opcional)", fmt.Sprintf("Release %s", tagName))

		if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente criar a tag '%s'?", tagName)) {
			ui.PrintInfo("Criação de tag cancelada.")
			continue
		}

		err = git.CreateTag(tagName, tagMsg)
		if err != nil {
			ui.PrintError("%v", err)
		} else {
			ui.PrintSuccess("Tag '%s' criada com sucesso!", tagName)
		}
	}
}

func init() {
	gitCmd.AddCommand(gitCommitCmd)
	gitCmd.AddCommand(gitLogCmd)
	gitCmd.AddCommand(gitBranchCmd)
	gitCmd.AddCommand(gitTagCmd)
	gitCmd.AddCommand(gitPullCmd)
	gitCmd.AddCommand(gitPushCmd)
	rootCmd.AddCommand(gitCmd)
}
