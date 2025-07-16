# Sistema de Logging - CRM Backend

## Visão Geral

O sistema de logging do CRM Backend foi projetado para fornecer informações detalhadas sobre o funcionamento da aplicação, facilitando debugging, monitoramento e auditoria.

## Funcionalidades

### 1. Níveis de Log

- **DEBUG**: Informações detalhadas para debugging
- **INFO**: Informações gerais sobre o funcionamento da aplicação
- **WARNING**: Situações que merecem atenção mas não impedem o funcionamento
- **ERROR**: Erros que afetam o funcionamento da aplicação

### 2. Formatos de Log

#### Formato Texto (Desenvolvimento)

```
INFO: 2025-07-16 10:30:45.123456 HTTP Request method=POST path=/api/contacts status_code=201 duration=45ms client_ip=192.168.1.100
```

#### Formato JSON (Produção)

```json
{
  "timestamp": "2025-07-16T10:30:45.123456Z",
  "level": "INFO",
  "message": "HTTP Request",
  "fields": {
    "method": "POST",
    "path": "/api/contacts",
    "status_code": 201,
    "duration": "45ms",
    "client_ip": "192.168.1.100"
  },
  "source": "middleware.go:45"
}
```

## Configuração

### Variáveis de Ambiente

```bash
# Nível de log (DEBUG, INFO, WARNING, ERROR)
LOG_LEVEL=INFO

# Formato de log (text, json)
LOG_FORMAT=text

# Saída de log (stdout, file, both)
LOG_OUTPUT=stdout

# Ativar modo debug
DEBUG=true

# Configurações de arquivo (quando LOG_OUTPUT=file ou both)
LOG_MAX_SIZE=100        # Tamanho máximo em MB
LOG_MAX_BACKUPS=10      # Número de backups
LOG_MAX_AGE=30          # Idade máxima em dias
LOG_COMPRESS=true       # Compressão dos backups
```

## Como Usar

### 1. Logging Básico

```go
import "crm-backend/pkg/logger"

// Logs simples
logger.Info("Usuário logado com sucesso")
logger.Warning("Tentativa de acesso negada")
logger.Error("Falha ao conectar com banco de dados")
logger.Debug("Variável x =", x)

// Logs formatados
logger.Infof("Usuário %s logado às %s", username, time.Now())
logger.Errorf("Erro ao processar ID %d: %v", id, err)
```

### 2. Logging Estruturado

```go
// Com campos estruturados
logger.WithFields("INFO", "User Login", map[string]interface{}{
    "user_id": 123,
    "email": "user@example.com",
    "ip": "192.168.1.100",
    "timestamp": time.Now(),
})

// Log de erro com contexto
logger.LogError(err, "Database Operation", map[string]interface{}{
    "operation": "SELECT",
    "table": "users",
    "user_id": 123,
})
```

### 3. Logging de Requisições HTTP

```go
// Automático via middleware
// Registra automaticamente todas as requisições HTTP

// Manual em handlers
logger.LogRequest(
    method,
    path,
    statusCode,
    duration,
    clientIP,
    userAgent,
)
```

### 4. Logging de Chamadas de Serviço

```go
start := time.Now()
// ... operação do serviço ...
logger.LogServiceCall("UserService", "CreateUser", time.Since(start), true)
```

### 5. Logging Estruturado em JSON

```go
// Inicializar logger estruturado
logger.InitStructuredLogger()

// Usar logger estruturado
logger.StructuredLog.Info("User created", map[string]interface{}{
    "user_id": 123,
    "email": "user@example.com",
})

// Funções específicas para diferentes tipos de eventos
logger.LogDatabaseOperation("INSERT", "users", duration, true, nil)
logger.LogAPICall("POST", "/api/users", 201, duration, userID)
logger.LogBusinessEvent("user_registered", "user", userID, userID, details)
```

## Exemplos Práticos

### 1. Em Handlers

```go
func (h *ContactHandler) Create(c *gin.Context) {
    start := time.Now()
    userID := c.GetUint("user_id")

    logger.Debugf("Criando novo contato para usuário %d", userID)

    // ... lógica do handler ...

    if err != nil {
        logger.LogError(err, "Contact Creation", map[string]interface{}{
            "user_id": userID,
            "request": req,
        })
        c.Error(err)
        return
    }

    logger.WithFields("INFO", "Contact Created", map[string]interface{}{
        "user_id": userID,
        "contact_id": contact.ID,
        "duration": time.Since(start),
    })

    c.JSON(http.StatusCreated, contact)
}
```

### 2. Em Services

```go
func (s *ContactService) Create(userID uint, req *models.ContactCreateRequest) (*models.Contact, error) {
    start := time.Now()

    logger.Debugf("Criando contato para usuário %d", userID)

    // ... lógica do serviço ...

    if err != nil {
        logger.LogError(err, "Contact Service Create", map[string]interface{}{
            "user_id": userID,
            "email": req.Email,
        })
        return nil, err
    }

    logger.LogServiceCall("ContactService", "Create", time.Since(start), true)
    return contact, nil
}
```

### 3. Em Repositories

```go
func (r *ContactRepository) Create(contact *models.Contact) error {
    start := time.Now()

    err := r.db.Create(contact).Error
    success := err == nil

    logger.LogDatabaseOperation("INSERT", "contacts", time.Since(start), success, err)

    return err
}
```

## Monitoramento e Alertas

### 1. Logs de Erro

- Todos os erros são automaticamente registrados com contexto
- Incluem stack trace quando disponível
- Contêm informações sobre usuário e operação

### 2. Métricas de Performance

- Duração de operações
- Chamadas de API
- Operações de banco de dados
- Uso de memória

### 3. Auditoria

- Todas as operações de usuário são registradas
- Eventos de negócio são rastreados
- Informações de segurança são capturadas

## Boas Práticas

1. **Use níveis apropriados**: DEBUG para desenvolvimento, INFO para produção
2. **Inclua contexto**: Sempre adicione informações relevantes (user_id, request_id)
3. **Não registre informações sensíveis**: Senhas, tokens, dados pessoais
4. **Use logging estruturado**: Facilita análise e busca
5. **Monitore performance**: Registre duração de operações importantes
6. **Trate erros**: Sempre registre erros com contexto adequado

## Troubleshooting

### Logs não aparecem

1. Verifique o nível de log configurado
2. Confirme as variáveis de ambiente
3. Verifique se o logger foi inicializado

### Performance impacto

1. Reduza o nível de log em produção
2. Use logging assíncrono para alta carga
3. Configure rotação de logs adequada

### Logs muito verbosos

1. Ajuste o nível de log
2. Filtre logs desnecessários
3. Use campos estruturados em vez de strings longas
