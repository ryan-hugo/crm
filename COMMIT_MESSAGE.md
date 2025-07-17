fix: Corrigir autenticação JWT na rota /api/users/profile

- Corrigir conversão de user_id de float64 para uint no middleware
- Adicionar validação adicional de user_id > 0
- Implementar logging detalhado para debug de autenticação
- Melhorar tratamento de erros no handler GetProfile
- Adicionar scripts de teste automatizados (bash e PowerShell)

Problema:
- Middleware estava passando user_id como interface{} diretamente
- JWT armazena números como float64 por padrão
- c.GetUint("user_id") falhava na conversão

Solução:
- Converter explicitamente claims["user_id"].(float64) para uint
- Validar user_id antes de adicionar ao contexto
- Adicionar logs estruturados para debugging
- Criar testes automatizados para validação

Arquivos modificados:
- internal/middleware/auth.go
- internal/handlers/user_handler.go
- scripts/test-auth.sh (novo)
- scripts/test-auth.ps1 (novo)
- docs/JWT_AUTH_FIX.md (novo)
- AUTH_FIX_SUMMARY.md (novo)

Testes:
- Login com credenciais válidas
- Acesso a rota protegida com token válido
- Rejeição de token inválido
- Rejeição de requisição sem token

Co-authored-by: GitHub Copilot <copilot@github.com>
