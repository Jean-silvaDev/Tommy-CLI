# Regras de Arquitetura e Estilo Go (Tommy CLI)

Diretrizes obrigatórias para escrita e manutenção de código em Go no projeto **Tommy CLI**.

---

## 📐 Padrões de Código e Manipulação de Erros

1. **Empacotamento de Erros (`fmt.Errorf`)**:
   - Sempre envolva erros internos usando o verbo `%w` para preservar a cadeia de exceções (*error wrapping*):
     ```go
     if err != nil {
         return fmt.Errorf("falha ao processar operação: %w", err)
     }
     ```

2. **Formatação de Mensagens de Erro**:
   - Mensagens de erro em Go não devem começar com letra maiúscula nem terminar com ponto final (convenção padrão de Go):
     ```go
     // Correto:
     return fmt.Errorf("arquivo de configuração não encontrado")
     ```

3. **Uso do Pacote de UI (`pkg/ui`)**:
   - Para exibição no terminal interativo, utilize os helpers padronizados do pacote `pkg/ui`:
     - `ui.PrintInfo(...)`
     - `ui.PrintSuccess(...)`
     - `ui.PrintWarning(...)`
     - `ui.PrintError(...)`
     - `ui.ConfirmPrompt(...)`
     - `ui.SelectOption(...)`

4. **Tratamento de Permissões no Windows**:
   - Para remoção de arquivos e pastas no Windows (ex: arquivos temporários), desmarque o atributo Read-Only aplicando `os.Chmod(path, 0666)` antes de chamar `os.RemoveAll(path)`.
