---
name: network-diagnostics
description: >-
  Use esta skill quando for necessário inspecionar interfaces de rede, liberar portas em uso (kill process),
  medir latência (ping), escanear portas locais de dev, calcular CIDR ou limpar cache DNS.
---

# Skill: Diagnósticos de Redes & Process Kill por Porta

Runbook procedimental para diagnóstico de conectividade e gerenciamento de portas em `pkg/network/network.go`.

---

## Passos Procedimentais

1. **Inspeção de Interfaces & WAN**:
   - Utilize as funções de diagnóstico em `pkg/network/network.go`:
     - `GetNetworkInterfaces()` para obter IPv4, IPv6 e MAC.
     - `GetPublicIP()` para obter o IP público WAN via serviço remoto.

2. **Liberar Porta em Uso (`KillProcessOnPort`)**:
   - Para identificar e encerrar processos que estejam travando portas locais (ex: `3000`, `8080`, `5432`):
     - Mapear a porta via `GetProcessOnPort(port)`.
     - Confirmar com o usuário via `ui.ConfirmPrompt(...)`.
     - Encerrar o processo via `KillProcess(pid)`.

3. **Cálculo CIDR**:
   - Utilize `CalculateCIDR(cidrStr)` para calcular a máscara de rede, IP de rede, IP de broadcast e faixa de IPs válidos sem depender de utilitários externos de terminal.

4. **Flush DNS Cross-Platform**:
   - Executar o comando nativo do SO (`ipconfig /flushdns` no Windows, `systemd-resolve` ou `killall -HUP mDNSResponder` no Linux/macOS).
