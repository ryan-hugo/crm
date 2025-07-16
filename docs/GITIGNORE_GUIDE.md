# Documentação do .gitignore

## 📁 Estrutura de Pastas Go - O que NÃO é ignorado

### ✅ Pastas IMPORTANTES que são commitadas:

- **`internal/`**: Código interno da aplicação (handlers, services, models, etc.)
- **`pkg/`**: Pacotes reutilizáveis (logger, errors, utils, etc.)
- **`cmd/`**: Pontos de entrada da aplicação (main.go)
- **`api/`**: Definições de API (rotas, documentação)
- **`docs/`**: Documentação do projeto
- **`examples/`**: Exemplos de uso
- **`scripts/`**: Scripts de build/deploy
- **`migrations/`**: Migrações de banco de dados

### ❌ Pastas que são ignoradas:

- **`build/`**: Arquivos de build/compilação
- **`dist/`**: Distribuição compilada
- **`bin/`**: Binários compilados
- **`vendor/`**: Dependências Go (use go mod)
- **`tmp/`**: Arquivos temporários
- **`logs/`**: Arquivos de log
- **`coverage/`**: Relatórios de cobertura

## 🔒 Segurança - Arquivos Sensíveis Ignorados

### Configurações:

- `.env*` (todas as variações)
- `config.json`, `secrets.yaml`
- `*.key`, `*.pem`, `*.crt`

### Credenciais:

- `*token*`, `*secret*`, `*password*`
- `credentials.json`, `auth.json`
- `service-account*.json`

### Banco de Dados:

- `*.db`, `*.sqlite`, `*.sqlite3`
- `dump.sql`, `*.backup`

## 🚨 Armadilhas Comuns

### ❌ Erros que foram corrigidos:

1. **`pkg/` estava sendo ignorado**: Corrigido para `build/pkg/`
2. **`internal*` muito genérico**: Agora só ignora `internal-config*`, `internal-secret*`
3. **`.internal/` desnecessário**: Removido

### ✅ Padrões seguros:

```gitignore
# ✅ CORRETO - Ignora apenas arquivos específicos
internal-config*
internal-secret*

# ❌ ERRADO - Ignoraria toda a pasta internal/
internal*

# ✅ CORRETO - Ignora apenas build/pkg/
build/pkg/

# ❌ ERRADO - Ignoraria a pasta pkg/ raiz
pkg/
```

## 📝 Validação

Para verificar se o .gitignore está correto:

```bash
# Verificar se pastas importantes não estão ignoradas
git check-ignore internal/
git check-ignore pkg/
git check-ignore cmd/

# Não devem retornar nada (não estão ignoradas)

# Verificar se arquivos sensíveis estão ignorados
git check-ignore .env
git check-ignore secrets.json
git check-ignore *.key

# Devem retornar o caminho (estão ignoradas)
```

## 🔧 Manutenção

### Adicionar novos padrões:

1. Sempre teste com `git check-ignore`
2. Use padrões específicos, não genéricos
3. Adicione comentários explicativos
4. Documente mudanças importantes

### Revisão periódica:

- Verificar se novas ferramentas geram arquivos que devem ser ignorados
- Validar que pastas importantes não foram ignoradas por engano
- Atualizar documentação quando necessário

## 🎯 Checklist de Verificação

Antes de fazer commit do .gitignore:

- [ ] `internal/` não está ignorado
- [ ] `pkg/` não está ignorado
- [ ] `cmd/` não está ignorado
- [ ] `.env` está ignorado
- [ ] `*.key` está ignorado
- [ ] `*.log` está ignorado
- [ ] `build/` está ignorado
- [ ] `vendor/` está ignorado

## 📞 Problemas Comuns

### "Minha pasta internal sumiu!"

- Verifique se `internal*` não está no .gitignore
- Use `git status` para ver se está sendo ignorada
- Use `git check-ignore internal/` para confirmar

### "Arquivos sensíveis foram commitados!"

- Execute os hooks de segurança
- Use `git secrets --scan` se disponível
- Remova do histórico se necessário

### "Build não funciona!"

- Verifique se `pkg/` não está sendo ignorado
- Confirme que `go.mod` e `go.sum` estão sendo commitados
- Verifique se código fonte não está sendo ignorado

## 📚 Referências

- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Git Ignore Patterns](https://git-scm.com/docs/gitignore)
- [Security Best Practices](./SECURITY_GUIDE.md)
