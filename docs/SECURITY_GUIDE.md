# Guia de Segurança - Git e Versionamento

## 🔐 Visão Geral

Este documento descreve as práticas de segurança implementadas no projeto CRM Backend para proteger informações sensíveis e garantir a integridade do código.

## 📋 Arquivos de Segurança

### 1. `.gitignore`

Arquivo abrangente que protege contra commit acidental de:

- **Arquivos de configuração sensíveis**: `.env`, `config.json`, `secrets.yaml`
- **Chaves e certificados**: `*.key`, `*.pem`, `*.crt`
- **Tokens e credenciais**: `*token*`, `*secret*`, `*auth*`
- **Arquivos de banco de dados**: `*.db`, `*.sqlite`
- **Logs e arquivos temporários**: `*.log`, `tmp/`, `temp/`
- **Binários e builds**: `*.exe`, `build/`, `dist/`

### 2. `.gitattributes`

Configurações para:

- **Normalização de linha**: Garante consistência entre sistemas
- **Filtros de limpeza**: Remove conteúdo sensível automaticamente
- **Controle de exports**: Exclui arquivos sensíveis de archives
- **Tratamento de binários**: Identifica corretamente tipos de arquivo

### 3. Scripts de Segurança

- **`setup-git-security.sh`**: Configura hooks e filtros de segurança
- **Hooks pre-commit**: Verifica arquivos sensíveis antes do commit
- **Hooks pre-push**: Verifica histórico antes do push

## 🚀 Configuração Inicial

### 1. Executar Script de Segurança

```bash
# Linux/Mac
chmod +x scripts/setup-git-security.sh
./scripts/setup-git-security.sh

# Windows PowerShell
.git/hooks/setup-security.ps1
```

### 2. Configurar GPG (Recomendado)

```bash
# Gerar chave GPG
gpg --full-generate-key

# Configurar Git para usar GPG
git config user.signingkey YOUR_GPG_KEY_ID
git config commit.gpgsign true
git config tag.gpgsign true
```

### 3. Instalar Ferramentas Adicionais

```bash
# git-secrets (detecta secrets em commits)
brew install git-secrets  # Mac
apt-get install git-secrets  # Linux

# Configurar git-secrets
git secrets --register-aws
git secrets --install
```

## 🔍 Verificações de Segurança

### 1. Hooks Pre-commit

O hook pre-commit verifica:

- Padrões sensíveis em arquivos
- Tipos de arquivo não permitidos
- Tamanho de arquivos (>10MB)
- Conteúdo potencialmente sensível

### 2. Hooks Pre-push

O hook pre-push verifica:

- Histórico de commits por arquivos sensíveis
- Secrets hardcoded nos últimos commits
- Integridade do repositório

### 3. Filtros de Limpeza

Filtros automáticos que:

- Removem passwords/secrets de arquivos
- Limpam tokens antes do commit
- Normalizam conteúdo sensível

## ⚠️ Padrões Sensíveis Detectados

### Palavras-chave Monitoradas

- `password`, `secret`, `token`, `api_key`
- `private_key`, `credential`, `auth`
- `jwt`, `database_url`, `connection_string`
- `smtp`, `email_password`, `redis_password`
- `aws_secret`, `gcp_key`, `azure_key`

### Tipos de Arquivo Bloqueados

- `.env*` (arquivos de ambiente)
- `*.key` (chaves privadas)
- `*.pem` (certificados)
- `*.crt` (certificados)
- `*.p12`, `*.pfx` (keystores)
- `*secret*` (qualquer arquivo com "secret")

## 🛠️ Práticas Recomendadas

### 1. Gerenciamento de Secrets

```bash
# ✅ CORRETO - Usar variáveis de ambiente
DATABASE_URL=postgres://user:pass@localhost/db

# ❌ ERRADO - Hardcoded no código
const dbURL = "postgres://user:pass@localhost/db"
```

### 2. Configuração por Ambiente

```bash
# Desenvolvimento
cp .env.example .env.development

# Produção (nunca commitar)
cp .env.example .env.production
```

### 3. Rotação de Secrets

- Alterar secrets regularmente
- Usar gerenciadores de secrets (Vault, AWS Secrets Manager)
- Implementar rotação automática

### 4. Auditoria Regular

```bash
# Verificar histórico por secrets
git log --all --grep="password\|secret\|token" --oneline

# Verificar arquivos sensíveis no histórico
git log --name-only --pretty=format: | grep -E '\.(env|key|pem)$'

# Escanear todo o repositório
git secrets --scan-history
```

## 🚨 Resposta a Incidentes

### 1. Secret Commitado Acidentalmente

```bash
# Remover do último commit
git reset --soft HEAD~1
git reset HEAD arquivo-sensivel
git commit -m "Remove sensitive file"

# Remover do histórico (CUIDADO!)
git filter-branch --tree-filter 'rm -f arquivo-sensivel' HEAD
git push --force-with-lease origin main
```

### 2. Limpar Histórico Completamente

```bash
# Usar BFG Repo-Cleaner
java -jar bfg.jar --replace-text passwords.txt repo.git
git reflog expire --expire=now --all
git gc --prune=now --aggressive
```

### 3. Invalidar Credenciais

- Alterar imediatamente todas as credenciais comprometidas
- Revogar tokens e chaves de API
- Notificar equipe de segurança
- Documentar incidente

## 🔒 Configurações de Segurança do Git

### 1. Configurações Básicas

```bash
git config --global user.name "Seu Nome"
git config --global user.email "seu.email@empresa.com"
git config --global init.defaultBranch main
git config --global pull.rebase true
git config --global push.default simple
```

### 2. Configurações de Segurança

```bash
git config --global core.autocrlf false
git config --global core.safecrlf true
git config --global commit.gpgsign true
git config --global tag.gpgsign true
```

### 3. Aliases Úteis

```bash
git config --global alias.check-secrets "!git secrets --scan"
git config --global alias.scan-history "!git secrets --scan-history"
git config --global alias.clean-history "!git reflog expire --expire=now --all && git gc --prune=now --aggressive"
```

## 📊 Monitoramento e Alertas

### 1. CI/CD Integration

```yaml
# .github/workflows/security.yml
name: Security Scan
on: [push, pull_request]
jobs:
  security:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run git-secrets
        run: |
          git secrets --install
          git secrets --register-aws
          git secrets --scan-history
```

### 2. Pre-commit Hooks

```yaml
# .pre-commit-config.yaml
repos:
  - repo: https://github.com/Yelp/detect-secrets
    rev: v1.2.0
    hooks:
      - id: detect-secrets
        args: ["--baseline", ".secrets.baseline"]
```

## 📚 Recursos Adicionais

### Ferramentas Recomendadas

- **git-secrets**: Detecta secrets em repositórios
- **detect-secrets**: Scanner avançado de secrets
- **truffleHog**: Busca por secrets em histórico Git
- **GitGuardian**: Monitoramento contínuo de secrets
- **SOPS**: Criptografia de arquivos de configuração

### Documentação

- [Git Security Best Practices](https://git-scm.com/book/en/v2/Git-Tools-Signing-Your-Work)
- [GitHub Security Features](https://docs.github.com/en/code-security)
- [OWASP Secure Coding Practices](https://owasp.org/www-project-secure-coding-practices-quick-reference-guide/)

### Treinamento

- Workshops sobre segurança em Git
- Simulações de resposta a incidentes
- Revisão regular de práticas de segurança

## 🎯 Checklist de Segurança

### Antes de Cada Commit

- [ ] Revisar arquivos a serem commitados
- [ ] Verificar se não há secrets hardcoded
- [ ] Confirmar que .env não está sendo commitado
- [ ] Executar `git secrets --scan`

### Antes de Cada Push

- [ ] Verificar histórico recente
- [ ] Confirmar que branch está atualizada
- [ ] Validar que não há conflicts
- [ ] Executar testes de segurança

### Periodicamente

- [ ] Auditar histórico do repositório
- [ ] Rotar secrets e credenciais
- [ ] Atualizar ferramentas de segurança
- [ ] Revisar configurações de segurança
- [ ] Treinar equipe sobre novas práticas

## 🆘 Contatos de Emergência

Em caso de incidente de segurança:

1. **Equipe de Segurança**: security@empresa.com
2. **DevOps**: devops@empresa.com
3. **Gerente de Projeto**: manager@empresa.com

**Lembre-se**: A segurança é responsabilidade de todos!
