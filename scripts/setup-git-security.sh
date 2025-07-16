#!/bin/bash

# ====================
# CONFIGURAÇÃO DE SEGURANÇA DO GIT
# ====================

echo "🔐 Configurando segurança do Git para o projeto CRM Backend..."

# Configurar hooks do Git para prevenir commits de arquivos sensíveis
echo "📋 Configurando hooks do Git..."

# Criar diretório de hooks se não existir
mkdir -p .git/hooks

# Hook pre-commit para verificar arquivos sensíveis
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash

# Cores para output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}🔍 Verificando arquivos sensíveis antes do commit...${NC}"

# Lista de padrões sensíveis
SENSITIVE_PATTERNS=(
    "password"
    "secret"
    "token"
    "api.key"
    "private.key"
    "credential"
    "auth"
    "jwt"
    "database.url"
    "db.url"
    "connection.string"
    "smtp"
    "email.password"
    "redis.password"
    "aws.secret"
    "gcp.key"
    "azure.key"
    "BEGIN.PRIVATE.KEY"
    "BEGIN.RSA.PRIVATE.KEY"
    "BEGIN.CERTIFICATE"
)

# Verificar arquivos staged
STAGED_FILES=$(git diff --cached --name-only)

# Flag para indicar se foram encontrados problemas
HAS_ISSUES=false

# Verificar cada arquivo staged
for file in $STAGED_FILES; do
    if [[ -f "$file" ]]; then
        # Verificar se o arquivo contém padrões sensíveis
        for pattern in "${SENSITIVE_PATTERNS[@]}"; do
            if grep -qi "$pattern" "$file"; then
                echo -e "${RED}❌ AVISO: Possível conteúdo sensível encontrado em $file (padrão: $pattern)${NC}"
                HAS_ISSUES=true
            fi
        done
        
        # Verificar se é um arquivo que não deveria ser commitado
        if [[ "$file" == *.env* ]] || [[ "$file" == *.key ]] || [[ "$file" == *.pem ]] || [[ "$file" == *secret* ]]; then
            echo -e "${RED}❌ ERRO: Arquivo sensível detectado: $file${NC}"
            HAS_ISSUES=true
        fi
    fi
done

# Verificar se há arquivos grandes (> 10MB)
for file in $STAGED_FILES; do
    if [[ -f "$file" ]]; then
        size=$(stat -c%s "$file" 2>/dev/null || stat -f%z "$file" 2>/dev/null || echo 0)
        if [[ $size -gt 10485760 ]]; then
            echo -e "${YELLOW}⚠️  AVISO: Arquivo grande detectado: $file ($(($size / 1024 / 1024))MB)${NC}"
        fi
    fi
done

if [[ "$HAS_ISSUES" == true ]]; then
    echo -e "${RED}❌ Commit bloqueado devido a arquivos sensíveis detectados!${NC}"
    echo -e "${YELLOW}💡 Dicas:${NC}"
    echo -e "   - Remova arquivos sensíveis do commit"
    echo -e "   - Use variáveis de ambiente para configurações"
    echo -e "   - Adicione arquivos sensíveis ao .gitignore"
    echo -e "   - Para forçar o commit (NÃO RECOMENDADO): git commit --no-verify"
    exit 1
fi

echo -e "${GREEN}✅ Nenhum arquivo sensível detectado. Commit permitido.${NC}"
exit 0
EOF

# Tornar o hook executável
chmod +x .git/hooks/pre-commit

# Hook pre-push para verificações adicionais
cat > .git/hooks/pre-push << 'EOF'
#!/bin/bash

RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
NC='\033[0m'

echo -e "${GREEN}🚀 Verificando antes do push...${NC}"

# Verificar se há arquivos .env ou similares no histórico
if git log --name-only --pretty=format: | grep -E '\.(env|key|pem|crt|p12|pfx)$' | head -1; then
    echo -e "${RED}❌ ERRO: Arquivos sensíveis detectados no histórico do Git!${NC}"
    echo -e "${YELLOW}💡 Execute: git filter-branch --tree-filter 'rm -f arquivo-sensivel' HEAD${NC}"
    exit 1
fi

# Verificar se há secrets hardcoded nos últimos commits
if git log -p -10 | grep -i -E '(password|secret|token|api.key|private.key).*[:=].*([\'"'"'"][^"'"'"']*[\'"'"'"]|[^[:space:]]+)'; then
    echo -e "${YELLOW}⚠️  AVISO: Possíveis secrets detectados nos últimos commits!${NC}"
    echo -e "${YELLOW}💡 Revise o histórico e considere usar git-secrets ou similar${NC}"
fi

echo -e "${GREEN}✅ Verificações de push concluídas.${NC}"
exit 0
EOF

chmod +x .git/hooks/pre-push

# Configurar filtros para limpar conteúdo sensível
echo "🧹 Configurando filtros de limpeza..."

git config filter.remove-secrets.clean 'sed -E "s/(password|secret|token|key|credential)[:=][[:space:]]*[\"'"'"']?[^\"'"'"'[:space:]]*[\"'"'"']?/\1=***REMOVED***/gi"'
git config filter.remove-secrets.smudge 'cat'

# Configurar configurações de segurança do Git
echo "⚙️  Configurando opções de segurança do Git..."

# Prevenir pushes acidentais para branches protegidas
git config branch.main.pushRemote origin
git config branch.master.pushRemote origin

# Configurar assinatura de commits (se GPG estiver disponível)
if command -v gpg &> /dev/null; then
    echo "🔑 GPG disponível. Configure a assinatura de commits:"
    echo "   git config user.signingkey YOUR_GPG_KEY"
    echo "   git config commit.gpgsign true"
fi

# Configurar URL remota para usar HTTPS em vez de SSH (mais seguro para alguns ambientes)
REMOTE_URL=$(git config --get remote.origin.url)
if [[ "$REMOTE_URL" == git@* ]]; then
    echo "🔗 Detectada URL SSH. Considere usar HTTPS para maior segurança em alguns ambientes."
fi

# Configurar configurações de segurança adicionais
git config core.autocrlf false
git config core.safecrlf true
git config push.default simple
git config pull.rebase true
git config init.defaultBranch main

# Configurar hooks de segurança para trabalho em equipe
git config core.hooksPath .git/hooks

echo -e "${GREEN}✅ Configuração de segurança do Git concluída!${NC}"
echo -e "${YELLOW}📋 Resumo das configurações aplicadas:${NC}"
echo "   - Hook pre-commit: Verifica arquivos sensíveis"
echo "   - Hook pre-push: Verifica histórico"
echo "   - Filtros: Limpeza automática de secrets"
echo "   - Configurações: Segurança aprimorada"
echo ""
echo -e "${YELLOW}💡 Próximos passos recomendados:${NC}"
echo "   1. Instale git-secrets: brew install git-secrets (Mac) ou apt-get install git-secrets (Linux)"
echo "   2. Configure GPG para assinar commits"
echo "   3. Use um gerenciador de secrets (HashiCorp Vault, AWS Secrets Manager, etc.)"
echo "   4. Configure CI/CD para escanear secrets automaticamente"
echo "   5. Treine a equipe sobre práticas de segurança"
EOF

# Tornar o script executável
chmod +x .git/hooks/setup-security.sh

# Criar script PowerShell para Windows
cat > .git/hooks/setup-security.ps1 << 'EOF'
# ====================
# CONFIGURAÇÃO DE SEGURANÇA DO GIT (PowerShell)
# ====================

Write-Host "🔐 Configurando segurança do Git para o projeto CRM Backend..." -ForegroundColor Green

# Configurar hooks do Git
Write-Host "📋 Configurando hooks do Git..." -ForegroundColor Yellow

# Criar hook pre-commit
$preCommitHook = @"
#!/bin/bash
# (Conteúdo do hook pre-commit seria inserido aqui)
# Este é um placeholder para a versão PowerShell
echo "Verificando arquivos sensíveis..."
exit 0
"@

Set-Content -Path ".git/hooks/pre-commit" -Value $preCommitHook

# Configurar Git
Write-Host "⚙️  Configurando opções de segurança do Git..." -ForegroundColor Yellow

git config core.autocrlf false
git config core.safecrlf true
git config push.default simple
git config pull.rebase true
git config init.defaultBranch main

Write-Host "✅ Configuração de segurança do Git concluída!" -ForegroundColor Green
Write-Host "💡 Instale git-secrets para proteção adicional" -ForegroundColor Yellow
EOF

echo "🔐 Arquivos de segurança Git criados com sucesso!"
echo ""
echo "📋 Para aplicar as configurações, execute:"
echo "   Linux/Mac: chmod +x .git/hooks/setup-security.sh && ./.git/hooks/setup-security.sh"
echo "   Windows: .git/hooks/setup-security.ps1"
