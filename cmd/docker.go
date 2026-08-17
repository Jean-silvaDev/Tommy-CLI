package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"tommy/pkg/docker"
	"tommy/pkg/ui"

	"github.com/spf13/cobra"
)

var dockerCmd = &cobra.Command{
	Use:     "docker",
	Aliases: []string{"dock", "containers"},
	Short:   "Hub completo de gerenciamento e limpeza do Docker",
	Long:    `Permite gerenciar Containers, Stacks Compose, Imagens, Volumes, Redes e executar limpezas seletivas ou profundas no Docker.`,
	Run: func(cmd *cobra.Command, args []string) {
		runDockerDashboard()
	},
}

var dockerListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lista todos os containers do Docker",
	Run: func(cmd *cobra.Command, args []string) {
		ui.PrintBanner("Containers Docker")
		displayContainersTable()
	},
}

var dockerStartCmd = &cobra.Command{
	Use:   "start [container_id_ou_nome]",
	Short: "Inicia um container parado",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente INICIAR o container '%s'?", target)) {
			ui.PrintInfo("Operação cancelada.")
			return
		}
		ui.PrintInfo("Iniciando container: %s...", target)
		err := docker.StartContainer(target)
		if err != nil {
			ui.PrintError("%v", err)
			return
		}
		ui.PrintSuccess("Container '%s' iniciado com sucesso!", target)
	},
}

var dockerStopCmd = &cobra.Command{
	Use:   "stop [container_id_ou_nome]",
	Short: "Para um container rodando",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente PARAR o container '%s'?", target)) {
			ui.PrintInfo("Operação cancelada.")
			return
		}
		ui.PrintInfo("Parando container: %s...", target)
		err := docker.StopContainer(target)
		if err != nil {
			ui.PrintError("%v", err)
			return
		}
		ui.PrintSuccess("Container '%s' parado com sucesso!", target)
	},
}

func runDockerDashboard() {
	for {
		ui.PrintBanner("Hub Completo Docker")

		// 1. Checar se o Daemon do Docker está rodando
		if !docker.IsDockerRunning() {
			ui.PrintWarning("O serviço do Docker (Daemon/Desktop) NÃO está rodando na sua máquina.")
			fmt.Println()

			options := []string{
				"1. 🚀 Iniciar Serviço / App do Docker (Docker Desktop)",
				"2. ⬅️  Voltar ao Menu Principal",
				"3. 🚪 Sair da Aplicação",
			}

			idx, _, err := ui.SelectOption("Escolha uma ação para configurar o Docker", options)
			if err != nil || idx == 1 {
				break
			}
			if idx == 2 {
				ui.PrintExitMessage()
				os.Exit(0)
			}

			if idx == 0 {
				ui.PrintInfo("Inicializando o Docker Desktop...")
				err := docker.StartDockerEngine()
				if err != nil {
					ui.PrintError("%v", err)
					ui.WaitForEnter("")
					continue
				}

				ui.PrintInfo("Aguardando o serviço do Docker responder (pode levar alguns segundos)...")
				started := false
				for i := 0; i < 12; i++ {
					time.Sleep(2 * time.Second)
					if docker.IsDockerRunning() {
						started = true
						ui.PrintSuccess("Serviço do Docker iniciado com sucesso!")
						break
					}
				}
				if !started {
					ui.PrintError("O Docker Desktop está demorando para iniciar. Verifique a janela do aplicativo.")
					ui.WaitForEnter("")
				}
			}
			continue
		}

		options := []string{
			" 1. 🐳 Gerenciador de Containers & Stacks Compose",
			" 2. 🖼️  Gerenciador de Imagens Docker (Em uso vs Órfãs)",
			" 3. 💾 Gerenciador de Volumes Docker (Em uso vs Órfãos)",
			" 4. 🌐 Gerenciador de Redes Docker",
			" 5. 🧹 Central de Limpeza & Pruning Seletivo Docker",
			" 6. ⏹️  Parar Serviço do Docker (Docker Desktop)",
			" 7. ⬅️  Voltar ao Menu Principal",
			" 8. 🚪 Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Selecione um recurso do Docker para gerenciar", options)
		if err != nil || idx == 6 {
			break
		}

		if idx == 7 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			runContainerHub()
		case 1:
			runImageHub()
		case 2:
			runVolumeHub()
		case 3:
			runNetworkHub()
		case 4:
			runDockerCleanCenter()
		case 5:
			if ui.ConfirmPrompt("Deseja realmente PARAR o serviço do Docker (Docker Desktop)?") {
				ui.PrintInfo("Encerrando o serviço do Docker...")
				err := docker.StopDockerEngine()
				if err != nil {
					ui.PrintError("%v", err)
				} else {
					ui.PrintSuccess("Serviço do Docker parado com sucesso!")
				}
				ui.WaitForEnter("")
			}
		}
	}
}

func runContainerHub() {
	for {
		ui.PrintBanner("Containers & Stacks Docker")

		containers, err := displayContainersTable()
		if err != nil {
			ui.WaitForEnter("")
			break
		}

		options := []string{
			" 1. 📦 Gerenciar Stacks / Grupos Docker (Compose)",
			" 2. ▶️  Iniciar um Container Parado",
			" 3. ⏹️  Parar um Container Rodando",
			" 4. 🔄 Reiniciar um Container",
			" 5. 📜 Ver Logs do Container (docker logs)",
			" 6. 🔍 Inspecionar Detalhes do Container (docker inspect)",
			" 7. 🔍 Atualizar Listagem",
			" 8. ⬅️  Voltar ao Hub Docker",
			" 9. 🚪 Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Ação para Containers", options)
		if err != nil || idx == 7 {
			break
		}
		if idx == 8 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			runStackWizard(containers)
		case 1:
			selectAndActionContainer(containers, "running", func(target string) error {
				return docker.StartContainer(target)
			}, "INICIAR", "iniciado")
		case 2:
			selectAndActionContainer(containers, "", func(target string) error {
				return docker.StopContainer(target)
			}, "PARAR", "parado")
		case 3:
			selectAndActionContainer(containers, "all", func(target string) error {
				return docker.RestartContainer(target)
			}, "REINICIAR", "reiniciado")
		case 4:
			selectAndShowLogs(containers)
		case 5:
			selectAndInspectContainer(containers)
		case 6:
			continue
		}
	}
}

func selectAndShowLogs(containers []docker.ContainerInfo) {
	if len(containers) == 0 {
		ui.PrintWarning("Nenhum container disponível para leitura de logs.")
		ui.WaitForEnter("")
		return
	}

	var opts []string
	for _, c := range containers {
		opts = append(opts, fmt.Sprintf("%s (%s) [%s]", c.Name, c.ID, strings.ToUpper(c.State)))
	}
	opts = append(opts, "Cancelar")

	idx, _, err := ui.SelectOption("Escolha o container para visualizar os logs (últimas 100 linhas)", opts)
	if err != nil || idx == len(opts)-1 {
		return
	}

	target := containers[idx]
	ui.PrintInfo("Buscando logs de '%s' (%s)...", target.Name, target.ID)
	logs, err := docker.GetContainerLogs(target.ID, 100)
	if err != nil {
		ui.PrintError("%v", err)
	} else {
		fmt.Println()
		fmt.Println(ui.SubtitleStyle.Render(fmt.Sprintf("Logs do Container: %s", target.Name)))
		fmt.Println(logs)
		fmt.Println()
	}
	ui.WaitForEnter("")
}

func selectAndInspectContainer(containers []docker.ContainerInfo) {
	if len(containers) == 0 {
		ui.PrintWarning("Nenhum container disponível para inspeção.")
		ui.WaitForEnter("")
		return
	}

	var opts []string
	for _, c := range containers {
		opts = append(opts, fmt.Sprintf("%s (%s) [%s]", c.Name, c.ID, strings.ToUpper(c.State)))
	}
	opts = append(opts, "Cancelar")

	idx, _, err := ui.SelectOption("Escolha o container para inspecionar", opts)
	if err != nil || idx == len(opts)-1 {
		return
	}

	target := containers[idx]
	ui.PrintInfo("Inspecionando detalhes de '%s' (%s)...", target.Name, target.ID)
	inspectData, err := docker.InspectContainer(target.ID)
	if err != nil {
		ui.PrintError("%v", err)
	} else {
		fmt.Println()
		fmt.Println(ui.SubtitleStyle.Render(fmt.Sprintf("Inspeção Docker: %s", target.Name)))
		fmt.Println(inspectData)
		fmt.Println()
	}
	ui.WaitForEnter("")
}

func displayContainersTable() ([]docker.ContainerInfo, error) {
	containers, err := docker.ListContainers()
	if err != nil {
		ui.PrintError("%v", err)
		return nil, err
	}

	if len(containers) == 0 {
		ui.PrintWarning("Nenhum container encontrado no ambiente Docker.")
		return containers, nil
	}

	headers := []string{"ID", "NOME", "STACK", "IMAGEM", "STATUS", "ESTADO"}
	var rows [][]string

	for _, c := range containers {
		stateBadge := ui.BadgeStopped.Render("STOPPED")
		if strings.ToLower(c.State) == "running" {
			stateBadge = ui.BadgeRunning.Render("RUNNING")
		}

		rows = append(rows, []string{
			c.ID,
			c.Name,
			c.Project,
			c.Image,
			c.Status,
			stateBadge,
		})
	}

	ui.RenderTable(headers, rows)
	return containers, nil
}

func runImageHub() {
	for {
		ui.PrintBanner("Gerenciador de Imagens Docker")

		images, err := docker.ListImages()
		if err != nil {
			ui.PrintError("%v", err)
			ui.WaitForEnter("")
			break
		}

		if len(images) == 0 {
			ui.PrintWarning("Nenhuma imagem Docker encontrada localmente.")
		} else {
			headers := []string{"REPOSITÓRIO", "TAG", "ID IMAGEM", "TAMANHO", "STATUS DA IMAGEM"}
			var rows [][]string

			for _, img := range images {
				statusBadge := ui.WarningStyle.Render("⚠️ ÓRFÃ (SEM USO)")
				if img.InUse {
					statusBadge = ui.SuccessStyle.Render("✔ EM USO")
				}

				rows = append(rows, []string{
					img.Repository,
					img.Tag,
					img.ID,
					img.Size,
					statusBadge,
				})
			}
			ui.RenderTable(headers, rows)
		}

		options := []string{
			" 1. 📥 Baixar / Pull de Nova Imagem Docker (docker pull)",
			" 2. 🗑️  Excluir uma Imagem Específica",
			" 3. 🧹 Limpar Todas as Imagens Órfãs / Sem Uso (image prune -a)",
			" 4. 🔍 Atualizar Listagem",
			" 5. ⬅️  Voltar ao Hub Docker",
			" 6. 🚪 Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Escolha uma ação para Imagens", options)
		if err != nil || idx == 4 {
			break
		}
		if idx == 5 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			imageRef := ui.PromptInput("Informe a referência da imagem no Docker Hub (ex: postgres:16-alpine, redis:alpine, nginx:latest):")
			if strings.TrimSpace(imageRef) == "" {
				ui.PrintInfo("Operação de pull cancelada.")
				ui.WaitForEnter("")
				continue
			}

			ui.PrintInfo("Baixando a imagem '%s' do Docker Hub...", imageRef)
			errPull := docker.PullImage(imageRef)
			if errPull != nil {
				ui.PrintError("%v", errPull)
			} else {
				ui.PrintSuccess("Imagem '%s' baixada com sucesso!", imageRef)
			}
			ui.WaitForEnter("")

		case 1:
			if len(images) == 0 {
				ui.PrintWarning("Não há imagens para remover.")
				ui.WaitForEnter("")
				continue
			}

			var imgOpts []string
			for _, img := range images {
				statusStr := "SEM USO"
				if img.InUse {
					statusStr = "EM USO"
				}
				imgOpts = append(imgOpts, fmt.Sprintf("%s:%s (%s) - %s [%s]", img.Repository, img.Tag, img.ID, img.Size, statusStr))
			}
			imgOpts = append(imgOpts, "Cancelar")

			selectedIdx, _, sErr := ui.SelectOption("Selecione a imagem para excluir", imgOpts)
			if sErr != nil || selectedIdx == len(imgOpts)-1 {
				continue
			}

			targetImg := images[selectedIdx]

			if targetImg.InUse {
				ui.PrintWarning("Atenção: A imagem '%s:%s' (%s) está marcada como EM USO por containers existentes.", targetImg.Repository, targetImg.Tag, targetImg.ID)
				if !ui.ConfirmPrompt("Deseja FORÇAR a exclusão desta imagem mesmo em uso (-f)?") {
					ui.PrintInfo("Operação cancelada.")
					ui.WaitForEnter("")
					continue
				}
				ui.PrintInfo("Forçando exclusão da imagem %s...", targetImg.ID)
				errRmi := docker.RemoveImage(targetImg.ID, true)
				if errRmi != nil {
					ui.PrintError("%v", errRmi)
				} else {
					ui.PrintSuccess("Imagem %s excluída com sucesso!", targetImg.Repository)
				}
				ui.WaitForEnter("")
			} else {
				if ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente EXCLUIR a imagem '%s:%s'?", targetImg.Repository, targetImg.Tag)) {
					ui.PrintInfo("Excluindo imagem %s...", targetImg.Repository)
					errRmi := docker.RemoveImage(targetImg.ID, false)
					if errRmi != nil {
						ui.PrintError("%v", errRmi)
					} else {
						ui.PrintSuccess("Imagem %s excluída com sucesso!", targetImg.Repository)
					}
					ui.WaitForEnter("")
				}
			}

		case 2:
			var unusedCount int
			for _, img := range images {
				if !img.InUse {
					unusedCount++
				}
			}
			if unusedCount == 0 {
				ui.PrintWarning("Não há imagens órfãs (sem uso) para remover no momento.")
				ui.WaitForEnter("")
				continue
			}

			if ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente remover %d imagem(ns) Docker órfã(s) não utilizada(s)?", unusedCount)) {
				ui.PrintInfo("Executando docker image prune -a...")
				errP := docker.PruneImages(true)
				if errP != nil {
					ui.PrintError("%v", errP)
				} else {
					ui.PrintSuccess("Imagens não utilizadas removidas com sucesso!")
				}
				ui.WaitForEnter("")
			}
		case 3:
			continue
		}
	}
}

func runVolumeHub() {
	for {
		ui.PrintBanner("Gerenciador de Volumes Docker")

		volumes, err := docker.ListVolumes()
		if err != nil {
			ui.PrintError("%v", err)
			ui.WaitForEnter("")
			break
		}

		if len(volumes) == 0 {
			ui.PrintWarning("Nenhum volume Docker encontrado.")
		} else {
			headers := []string{"NOME DO VOLUME", "DRIVER", "ESCOPO", "STATUS DO VOLUME"}
			var rows [][]string

			for _, v := range volumes {
				statusBadge := ui.WarningStyle.Render("⚠️ ÓRFÃO (SEM USO)")
				if v.InUse {
					statusBadge = ui.SuccessStyle.Render("✔ EM USO")
				}

				rows = append(rows, []string{
					v.Name,
					v.Driver,
					v.Scope,
					statusBadge,
				})
			}
			ui.RenderTable(headers, rows)
		}

		options := []string{
			" 1. ➕ Criar Novo Volume Customizado (docker volume create)",
			" 2. 🗑️  Excluir um Volume Específico",
			" 3. 🧹 Limpar Todos os Volumes Órfãos / Sem Uso (volume prune)",
			" 4. 🔍 Atualizar Listagem",
			" 5. ⬅️  Voltar ao Hub Docker",
			" 6. 🚪 Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Escolha uma ação para Volumes", options)
		if err != nil || idx == 4 {
			break
		}
		if idx == 5 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			volName := ui.PromptInput("Informe o Nome para o novo volume Docker:")
			if strings.TrimSpace(volName) == "" {
				ui.PrintInfo("Operação de criação de volume cancelada.")
				ui.WaitForEnter("")
				continue
			}

			ui.PrintInfo("Criando volume '%s'...", volName)
			errCreate := docker.CreateVolume(volName, "local")
			if errCreate != nil {
				ui.PrintError("%v", errCreate)
			} else {
				ui.PrintSuccess("Volume '%s' criado com sucesso!", volName)
			}
			ui.WaitForEnter("")

		case 1:
			if len(volumes) == 0 {
				ui.PrintWarning("Não há volumes para remover.")
				ui.WaitForEnter("")
				continue
			}

			var volOpts []string
			for _, v := range volumes {
				statusStr := "ÓRFÃO"
				if v.InUse {
					statusStr = "EM USO"
				}
				volOpts = append(volOpts, fmt.Sprintf("%s (%s) - [%s]", v.Name, v.Driver, statusStr))
			}
			volOpts = append(volOpts, "Cancelar")

			selectedIdx, _, sErr := ui.SelectOption("Selecione o volume para excluir", volOpts)
			if sErr != nil || selectedIdx == len(volOpts)-1 {
				continue
			}

			targetVol := volumes[selectedIdx]

			if targetVol.InUse {
				ui.PrintWarning("O volume '%s' está em uso por containers ativos.", targetVol.Name)
				if !ui.ConfirmPrompt("Deseja FORÇAR a exclusão deste volume em uso (-f)?") {
					ui.PrintInfo("Operação cancelada.")
					ui.WaitForEnter("")
					continue
				}
				errRm := docker.RemoveVolume(targetVol.Name, true)
				if errRm != nil {
					ui.PrintError("%v", errRm)
				} else {
					ui.PrintSuccess("Volume '%s' excluído com sucesso!", targetVol.Name)
				}
				ui.WaitForEnter("")
			} else {
				if ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente EXCLUIR o volume '%s'?", targetVol.Name)) {
					errRm := docker.RemoveVolume(targetVol.Name, false)
					if errRm != nil {
						ui.PrintError("%v", errRm)
					} else {
						ui.PrintSuccess("Volume '%s' excluído com sucesso!", targetVol.Name)
					}
					ui.WaitForEnter("")
				}
			}

		case 2:
			var orphanedCount int
			for _, v := range volumes {
				if !v.InUse {
					orphanedCount++
				}
			}
			if orphanedCount == 0 {
				ui.PrintWarning("Não há volumes órfãos (sem uso) para remover no momento.")
				ui.WaitForEnter("")
				continue
			}

			if ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente apagar %d volume(s) órfão(s) (desconectados)?", orphanedCount)) {
				ui.PrintInfo("Executando docker volume prune...")
				errP := docker.PruneVolumes()
				if errP != nil {
					ui.PrintError("%v", errP)
				} else {
					ui.PrintSuccess("Volumes órfãos removidos com sucesso!")
				}
				ui.WaitForEnter("")
			}
		case 3:
			continue
		}
	}
}

func runNetworkHub() {
	for {
		ui.PrintBanner("Gerenciador de Redes Docker")

		networks, err := docker.ListNetworks()
		if err != nil {
			ui.PrintError("%v", err)
			ui.WaitForEnter("")
			break
		}

		if len(networks) == 0 {
			ui.PrintWarning("Nenhuma rede Docker encontrada.")
		} else {
			headers := []string{"ID REDE", "NOME DA REDE", "DRIVER", "TIPO / ESCOPO"}
			var rows [][]string

			for _, net := range networks {
				typeBadge := ui.InfoStyle.Render("CUSTOMIZADA")
				if net.IsProtected {
					typeBadge = ui.MutedStyle.Render("SISTEMA (PROTEGIDA)")
				}

				rows = append(rows, []string{
					net.ID,
					net.Name,
					net.Driver,
					typeBadge,
				})
			}
			ui.RenderTable(headers, rows)
		}

		options := []string{
			" 1. ➕ Criar Nova Rede Docker (docker network create)",
			" 2. 🔍 Inspecionar Detalhes da Rede (Subnet, Gateway, Containers)",
			" 3. 🔗 Conectar Container a uma Rede (docker network connect)",
			" 4. 🔌 Desconectar Container de uma Rede (docker network disconnect)",
			" 5. 🗑️  Excluir uma Rede Customizada",
			" 6. 🧹 Limpar Redes não Utilizadas (network prune)",
			" 7. 🔍 Atualizar Listagem",
			" 8. ⬅️  Voltar ao Hub Docker",
			" 9. 🚪 Sair da Aplicação",
		}

		idx, _, err := ui.SelectOption("Escolha uma ação para Redes", options)
		if err != nil || idx == 7 {
			break
		}
		if idx == 8 {
			ui.PrintExitMessage()
			os.Exit(0)
		}

		switch idx {
		case 0:
			netName := ui.PromptInput("Informe o Nome para a nova rede Docker:")
			if strings.TrimSpace(netName) == "" {
				ui.PrintInfo("Operação de criação de rede cancelada.")
				ui.WaitForEnter("")
				continue
			}

			driverOpts := []string{"1. bridge (Padrão para desenvolvimento local)", "2. overlay (Para Docker Swarm / Multi-host)", "3. macvlan (Atribuir MAC físico)"}
			dIdx, _, dErr := ui.SelectOption("Selecione o Driver da Rede", driverOpts)
			if dErr != nil {
				continue
			}

			driver := "bridge"
			if dIdx == 1 {
				driver = "overlay"
			} else if dIdx == 2 {
				driver = "macvlan"
			}

			subnet := ui.PromptInput("Informe a Sub-rede (ex: 172.28.0.0/16 ou aperte Enter para deixar automático):")

			ui.PrintInfo("Criando rede '%s' (Driver: %s)...", netName, driver)
			errCreate := docker.CreateNetwork(netName, driver, strings.TrimSpace(subnet))
			if errCreate != nil {
				ui.PrintError("%v", errCreate)
			} else {
				ui.PrintSuccess("Rede Docker '%s' criada com sucesso!", netName)
			}
			ui.WaitForEnter("")

		case 1:
			if len(networks) == 0 {
				ui.PrintWarning("Nenhuma rede disponível para inspeção.")
				ui.WaitForEnter("")
				continue
			}

			var netOpts []string
			for _, n := range networks {
				netOpts = append(netOpts, fmt.Sprintf("%s (%s) - Driver: %s", n.Name, n.ID, n.Driver))
			}
			netOpts = append(netOpts, "Cancelar")

			nIdx, _, nErr := ui.SelectOption("Selecione a rede para inspecionar", netOpts)
			if nErr != nil || nIdx == len(netOpts)-1 {
				continue
			}

			targetNet := networks[nIdx]
			ui.PrintInfo("Inspecionando rede '%s'...", targetNet.Name)
			details, errInsp := docker.InspectNetwork(targetNet.Name)
			if errInsp != nil {
				ui.PrintError("%v", errInsp)
			} else {
				fmt.Println()
				fmt.Println(ui.SubtitleStyle.Render(fmt.Sprintf("Inspeção da Rede Docker: %s", targetNet.Name)))
				fmt.Println(details)
				fmt.Println()
			}
			ui.WaitForEnter("")

		case 2:
			containers, errC := docker.ListContainers()
			if errC != nil || len(containers) == 0 {
				ui.PrintWarning("Nenhum container ativo disponível para conexão de rede.")
				ui.WaitForEnter("")
				continue
			}

			var netOpts []string
			for _, n := range networks {
				netOpts = append(netOpts, fmt.Sprintf("%s (%s)", n.Name, n.Driver))
			}
			netOpts = append(netOpts, "Cancelar")

			nIdx, _, nErr := ui.SelectOption("Selecione a REDE de destino", netOpts)
			if nErr != nil || nIdx == len(netOpts)-1 {
				continue
			}
			targetNet := networks[nIdx]

			var cOpts []string
			for _, c := range containers {
				cOpts = append(cOpts, fmt.Sprintf("%s (%s) [%s]", c.Name, c.ID, strings.ToUpper(c.State)))
			}
			cOpts = append(cOpts, "Cancelar")

			cIdx, _, cErr := ui.SelectOption("Selecione o CONTAINER para conectar à rede "+targetNet.Name, cOpts)
			if cErr != nil || cIdx == len(cOpts)-1 {
				continue
			}
			targetC := containers[cIdx]

			ui.PrintInfo("Conectando container '%s' à rede '%s'...", targetC.Name, targetNet.Name)
			errConn := docker.ConnectNetwork(targetNet.Name, targetC.ID)
			if errConn != nil {
				ui.PrintError("%v", errConn)
			} else {
				ui.PrintSuccess("Container '%s' conectado com sucesso à rede '%s'!", targetC.Name, targetNet.Name)
			}
			ui.WaitForEnter("")

		case 3:
			containers, errC := docker.ListContainers()
			if errC != nil || len(containers) == 0 {
				ui.PrintWarning("Nenhum container disponível para desconexão de rede.")
				ui.WaitForEnter("")
				continue
			}

			var netOpts []string
			for _, n := range networks {
				netOpts = append(netOpts, fmt.Sprintf("%s (%s)", n.Name, n.Driver))
			}
			netOpts = append(netOpts, "Cancelar")

			nIdx, _, nErr := ui.SelectOption("Selecione a REDE da qual deseja remover o container", netOpts)
			if nErr != nil || nIdx == len(netOpts)-1 {
				continue
			}
			targetNet := networks[nIdx]

			var cOpts []string
			for _, c := range containers {
				cOpts = append(cOpts, fmt.Sprintf("%s (%s)", c.Name, c.ID))
			}
			cOpts = append(cOpts, "Cancelar")

			cIdx, _, cErr := ui.SelectOption("Selecione o CONTAINER para desconectar", cOpts)
			if cErr != nil || cIdx == len(cOpts)-1 {
				continue
			}
			targetC := containers[cIdx]

			ui.PrintInfo("Desconectando container '%s' da rede '%s'...", targetC.Name, targetNet.Name)
			errDisc := docker.DisconnectNetwork(targetNet.Name, targetC.ID)
			if errDisc != nil {
				ui.PrintError("%v", errDisc)
			} else {
				ui.PrintSuccess("Container '%s' desconectado com sucesso da rede '%s'!", targetC.Name, targetNet.Name)
			}
			ui.WaitForEnter("")

		case 4:
			var customNets []docker.NetworkInfo
			var netOpts []string

			for _, net := range networks {
				if !net.IsProtected {
					customNets = append(customNets, net)
					netOpts = append(netOpts, fmt.Sprintf("%s (%s) - Driver: %s", net.Name, net.ID, net.Driver))
				}
			}

			if len(customNets) == 0 {
				ui.PrintWarning("Não há redes customizadas elegíveis para exclusão (redes padrão do sistema são protegidas).")
				ui.WaitForEnter("")
				continue
			}

			netOpts = append(netOpts, "Cancelar")
			selectedIdx, _, sErr := ui.SelectOption("Selecione a rede para remover", netOpts)
			if sErr != nil || selectedIdx == len(netOpts)-1 {
				continue
			}

			targetNet := customNets[selectedIdx]
			if ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente REMOVER a rede '%s'?", targetNet.Name)) {
				errRm := docker.RemoveNetwork(targetNet.Name)
				if errRm != nil {
					ui.PrintError("%v", errRm)
				} else {
					ui.PrintSuccess("Rede '%s' removida com sucesso!", targetNet.Name)
				}
				ui.WaitForEnter("")
			}

		case 5:
			var customUnusedCount int
			for _, net := range networks {
				if !net.IsProtected {
					customUnusedCount++
				}
			}
			if customUnusedCount == 0 {
				ui.PrintWarning("Não há redes customizadas sem uso para remover no momento.")
				ui.WaitForEnter("")
				continue
			}

			if ui.ConfirmPrompt(fmt.Sprintf("Deseja limpar %d rede(s) customizada(s) não utilizada(s)?", customUnusedCount)) {
				ui.PrintInfo("Executando docker network prune...")
				errP := docker.PruneNetworks()
				if errP != nil {
					ui.PrintError("%v", errP)
				} else {
					ui.PrintSuccess("Redes não utilizadas removidas com sucesso!")
				}
				ui.WaitForEnter("")
			}
		case 6:
			continue
		}
	}
}

func runDockerCleanCenter() {
	for {
		ui.PrintBanner("Central de Limpeza & Pruning Docker")

		dfSummary, errDf := docker.GetDockerUsageSummary()
		if errDf == nil && dfSummary != "" {
			ui.PrintInfo("Resumo de Uso de Disco do Docker (docker system df):")
			fmt.Println(dfSummary)
			fmt.Println()
		}

		options := []string{
			" 1. 🧹 Limpar Apenas Containers Parados (container prune)",
			" 2. 🧹 Limpar Apenas Imagens sem Uso (image prune -a)",
			" 3. 🧹 Limpar Apenas Volumes Órfãos (volume prune)",
			" 4. 🧹 Limpar Apenas Redes não Utilizadas (network prune)",
			" 5. 💥 Limpeza COMPLETA do Sistema Docker (Containers, Imagens e Volumes)",
			" 6. ⬅️  Voltar ao Hub Docker",
		}

		idx, _, err := ui.SelectOption("Escolha o tipo de limpeza do Docker", options)
		if err != nil || idx == 5 {
			break
		}

		switch idx {
		case 0:
			containers, _ := docker.ListContainers()
			stoppedCount := 0
			for _, c := range containers {
				if strings.ToLower(c.State) != "running" {
					stoppedCount++
				}
			}
			if stoppedCount == 0 {
				ui.PrintWarning("Nenhum container parado encontrado para limpeza no momento.")
				ui.WaitForEnter("")
				continue
			}

			if ui.ConfirmPrompt(fmt.Sprintf("Deseja remover %d container(s) parado(s)?", stoppedCount)) {
				ui.PrintInfo("Limpando containers parados...")
				errP := docker.PruneContainers()
				if errP != nil {
					ui.PrintError("%v", errP)
				} else {
					ui.PrintSuccess("Containers parados removidos com sucesso!")
				}
				ui.WaitForEnter("")
			}
		case 1:
			images, _ := docker.ListImages()
			unusedCount := 0
			for _, img := range images {
				if !img.InUse {
					unusedCount++
				}
			}
			if unusedCount == 0 {
				ui.PrintWarning("Nenhuma imagem sem uso encontrada para limpeza no momento.")
				ui.WaitForEnter("")
				continue
			}

			if ui.ConfirmPrompt(fmt.Sprintf("Deseja remover %d imagem(ns) sem uso?", unusedCount)) {
				ui.PrintInfo("Limpando imagens...")
				errP := docker.PruneImages(true)
				if errP != nil {
					ui.PrintError("%v", errP)
				} else {
					ui.PrintSuccess("Imagens sem uso removidas com sucesso!")
				}
				ui.WaitForEnter("")
			}
		case 2:
			volumes, _ := docker.ListVolumes()
			orphanedCount := 0
			for _, v := range volumes {
				if !v.InUse {
					orphanedCount++
				}
			}
			if orphanedCount == 0 {
				ui.PrintWarning("Nenhum volume órfão encontrado para limpeza no momento.")
				ui.WaitForEnter("")
				continue
			}

			if ui.ConfirmPrompt(fmt.Sprintf("Deseja remover %d volume(s) órfão(s)?", orphanedCount)) {
				ui.PrintInfo("Limpando volumes...")
				errP := docker.PruneVolumes()
				if errP != nil {
					ui.PrintError("%v", errP)
				} else {
					ui.PrintSuccess("Volumes órfãos removidos com sucesso!")
				}
				ui.WaitForEnter("")
			}
		case 3:
			networks, _ := docker.ListNetworks()
			customCount := 0
			for _, net := range networks {
				if !net.IsProtected {
					customCount++
				}
			}
			if customCount == 0 {
				ui.PrintWarning("Nenhuma rede customizada sem uso encontrada para limpeza no momento.")
				ui.WaitForEnter("")
				continue
			}

			if ui.ConfirmPrompt(fmt.Sprintf("Deseja remover %d rede(s) não utilizada(s)?", customCount)) {
				ui.PrintInfo("Limpando redes...")
				errP := docker.PruneNetworks()
				if errP != nil {
					ui.PrintError("%v", errP)
				} else {
					ui.PrintSuccess("Redes não utilizadas removidas com sucesso!")
				}
				ui.WaitForEnter("")
			}
		case 4:
			ui.PrintWarning("A limpeza completa removerá containers parados, imagens sem uso e volumes órfãos de uma só vez.")
			if ui.ConfirmPrompt("CONFIRMA A LIMPEZA PROFUNDA E COMPLETA DO DOCKER?") {
				ui.PrintInfo("Executando docker system prune -a --volumes...")
				errP := docker.PruneSystem(true, true)
				if errP != nil {
					ui.PrintError("%v", errP)
				} else {
					ui.PrintSuccess("Limpeza profunda do Docker concluída com sucesso!")
				}
				ui.WaitForEnter("")
			}
		}
	}
}

func runStackWizard(containers []docker.ContainerInfo) {
	stacks := docker.GetDockerStacks(containers)
	if len(stacks) == 0 {
		ui.PrintWarning("Nenhuma Stack / Grupo Compose detectado.")
		ui.WaitForEnter("")
		return
	}

	var options []string
	for _, s := range stacks {
		options = append(options, fmt.Sprintf("📦 Stack: %s (%d containers - %d RUNNING)", s.Name, s.TotalCount, s.RunningCount))
	}
	options = append(options, "⬅️  Voltar")

	idx, _, err := ui.SelectOption("Selecione uma Stack do Docker Compose para gerenciar", options)
	if err != nil || idx == len(options)-1 {
		return
	}

	targetStack := stacks[idx]

	for {
		fmt.Println()
		ui.PrintBanner(fmt.Sprintf("Stack: %s", targetStack.Name))
		ui.PrintInfo("Containers desta Stack:")
		for _, c := range targetStack.Containers {
			badge := ui.BadgeStopped.Render("STOPPED")
			if strings.ToLower(c.State) == "running" {
				badge = ui.BadgeRunning.Render("RUNNING")
			}
			fmt.Printf("  • %-32s (%s) [%s]\n", c.Name, c.ID, badge)
		}
		fmt.Println()

		stackOptions := []string{
			fmt.Sprintf("1. ▶️  Iniciar TODOS os containers da stack '%s'", targetStack.Name),
			fmt.Sprintf("2. ⏹️  Parar TODOS os containers da stack '%s'", targetStack.Name),
			fmt.Sprintf("3. 🔄 Reiniciar TODOS os containers da stack '%s'", targetStack.Name),
			"4. ⬅️  Voltar ao Menu de Stacks",
		}

		sIdx, _, sErr := ui.SelectOption("Ação para a Stack "+targetStack.Name, stackOptions)
		if sErr != nil || sIdx == 3 {
			break
		}

		var ids []string
		for _, c := range targetStack.Containers {
			ids = append(ids, c.ID)
		}

		switch sIdx {
		case 0:
			if ui.ConfirmPrompt(fmt.Sprintf("Deseja INICIAR todos os %d containers da stack '%s'?", len(ids), targetStack.Name)) {
				ui.PrintInfo("Iniciando containers da stack '%s'...", targetStack.Name)
				err := docker.StartStack(ids)
				if err != nil {
					ui.PrintError("%v", err)
				} else {
					ui.PrintSuccess("Todos os containers da stack '%s' foram iniciados!", targetStack.Name)
				}
				ui.WaitForEnter("")
				return
			}
		case 1:
			if ui.ConfirmPrompt(fmt.Sprintf("Deseja PARAR todos os %d containers da stack '%s'?", len(ids), targetStack.Name)) {
				ui.PrintInfo("Parando containers da stack '%s'...", targetStack.Name)
				err := docker.StopStack(ids)
				if err != nil {
					ui.PrintError("%v", err)
				} else {
					ui.PrintSuccess("Todos os containers da stack '%s' foram parados!", targetStack.Name)
				}
				ui.WaitForEnter("")
				return
			}
		case 2:
			if ui.ConfirmPrompt(fmt.Sprintf("Deseja REINICIAR todos os %d containers da stack '%s'?", len(ids), targetStack.Name)) {
				ui.PrintInfo("Reiniciando containers da stack '%s'...", targetStack.Name)
				err := docker.RestartStack(ids)
				if err != nil {
					ui.PrintError("%v", err)
				} else {
					ui.PrintSuccess("Todos os containers da stack '%s' foram reiniciados!", targetStack.Name)
				}
				ui.WaitForEnter("")
				return
			}
		}
	}
}

func selectAndActionContainer(containers []docker.ContainerInfo, filterState string, action func(string) error, actionVerb string, actionPast string) {
	var filtered []docker.ContainerInfo
	var options []string

	for _, c := range containers {
		if filterState == "running" && strings.ToLower(c.State) == "running" {
			continue
		}
		filtered = append(filtered, c)
		options = append(options, fmt.Sprintf("%s (%s) - Stack: %s [%s]", c.Name, c.ID, c.Project, strings.ToUpper(c.State)))
	}

	if len(options) == 0 {
		ui.PrintWarning("Nenhum container elegível para esta ação.")
		ui.WaitForEnter("")
		return
	}

	options = append(options, "Cancelar")
	idx, _, err := ui.SelectOption("Escolha o container", options)
	if err != nil || idx == len(options)-1 {
		return
	}

	target := filtered[idx]

	fmt.Println()
	ui.PrintWarning("Container Selecionado: Nome: '%s' | ID: %s | Stack: %s | Estado Atual: %s",
		target.Name, target.ID, target.Project, strings.ToUpper(target.State))

	if !ui.ConfirmPrompt(fmt.Sprintf("Deseja realmente %s o container '%s'?", actionVerb, target.Name)) {
		ui.PrintInfo("Operação cancelada pelo usuário.")
		ui.WaitForEnter("")
		return
	}

	ui.PrintInfo("Executando operação em '%s'...", target.Name)
	err = action(target.ID)
	if err != nil {
		ui.PrintError("%v", err)
	} else {
		ui.PrintSuccess("Container '%s' %s com sucesso!", target.Name, actionPast)
	}
	ui.WaitForEnter("")
}

func init() {
	dockerCmd.AddCommand(dockerListCmd)
	dockerCmd.AddCommand(dockerStartCmd)
	dockerCmd.AddCommand(dockerStopCmd)
	rootCmd.AddCommand(dockerCmd)
}
