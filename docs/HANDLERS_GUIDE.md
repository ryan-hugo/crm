# Guia Completo dos Handlers - CRM Backend

## Visão Geral

Este documento fornece um guia completo dos handlers implementados no backend GoLang do CRM. Os handlers constituem a camada de apresentação da API REST, responsáveis por processar requisições HTTP, validar entradas e formatar respostas.

## Estrutura Geral dos Handlers

Todos os handlers seguem um padrão consistente de implementação:

### Padrão de Estrutura

```go
type XHandler struct {
    xService services.XService
}

func NewXHandler(xService services.XService) *XHandler {
    return &XHandler{xService: xService}
}
```

### Padrão de Métodos

```go
func (h *XHandler) Create(c *gin.Context) {
    userID := c.GetUint("user_id")
    var req models.XCreateRequest
    
    // Validação de entrada
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(errors.NewBadRequestError("Dados inválidos: " + err.Error()))
        return
    }
    
    // Chamada do service
    result, err := h.xService.Create(userID, &req)
    if err != nil {
        c.Error(err)
        return
    }
    
    // Resposta de sucesso
    c.JSON(http.StatusCreated, result)
}
```

## AuthHandler

### Responsabilidades
- Registro de novos usuários
- Autenticação e geração de tokens JWT
- Validação de tokens
- Logout (stateless)

### Endpoints

#### POST /api/auth/register
**Descrição**: Registra um novo usuário no sistema

**Request Body**:
```json
{
    "name": "João Silva",
    "email": "joao@example.com",
    "password": "senha123"
}
```

**Response (201)**:
```json
{
    "message": "Usuário registrado com sucesso",
    "user": {
        "id": 1,
        "name": "João Silva",
        "email": "joao@example.com",
        "created_at": "2024-01-01T10:00:00Z"
    }
}
```

**Validações**:
- Email deve ser único
- Senha mínima de 6 caracteres
- Nome é obrigatório

#### POST /api/auth/login
**Descrição**: Autentica usuário e retorna token JWT

**Request Body**:
```json
{
    "email": "joao@example.com",
    "password": "senha123"
}
```

**Response (200)**:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user_id": 1,
    "email": "joao@example.com",
    "name": "João Silva"
}
```

**Validações**:
- Email e senha obrigatórios
- Credenciais devem ser válidas
- Token expira em 24 horas

#### GET /api/auth/validate
**Descrição**: Valida token JWT atual

**Headers**: `Authorization: Bearer <token>`

**Response (200)**:
```json
{
    "valid": true,
    "user_id": 1,
    "message": "Token válido"
}
```

#### POST /api/auth/logout
**Descrição**: Logout do usuário (stateless)

**Response (200)**:
```json
{
    "message": "Logout realizado com sucesso"
}
```

## UserHandler

### Responsabilidades
- Gestão de perfil do usuário
- Alteração de senha
- Exclusão de conta
- Estatísticas do usuário

### Endpoints

#### GET /api/users/profile
**Descrição**: Obtém perfil do usuário autenticado

**Headers**: `Authorization: Bearer <token>`

**Response (200)**:
```json
{
    "id": 1,
    "name": "João Silva",
    "email": "joao@example.com",
    "created_at": "2024-01-01T10:00:00Z",
    "updated_at": "2024-01-01T10:00:00Z"
}
```

#### PUT /api/users/profile
**Descrição**: Atualiza perfil do usuário

**Request Body**:
```json
{
    "name": "João Santos",
    "email": "joao.santos@example.com"
}
```

**Response (200)**:
```json
{
    "message": "Perfil atualizado com sucesso",
    "user": {
        "id": 1,
        "name": "João Santos",
        "email": "joao.santos@example.com",
        "updated_at": "2024-01-01T11:00:00Z"
    }
}
```

#### PUT /api/users/change-password
**Descrição**: Altera senha do usuário

**Request Body**:
```json
{
    "current_password": "senhaAtual123",
    "new_password": "novaSenha456",
    "confirm_password": "novaSenha456"
}
```

**Response (200)**:
```json
{
    "message": "Senha alterada com sucesso"
}
```

#### GET /api/users/stats
**Descrição**: Estatísticas consolidadas do usuário

**Response (200)**:
```json
{
    "total_contacts": 25,
    "total_clients": 15,
    "total_leads": 10,
    "total_tasks": 45,
    "pending_tasks": 12,
    "completed_tasks": 33,
    "total_projects": 8,
    "active_projects": 3,
    "completed_projects": 5,
    "total_interactions": 120
}
```

## ContactHandler

### Responsabilidades
- CRUD completo de contatos
- Busca e filtros avançados
- Detalhes com relacionamentos
- Conversão de leads em clientes
- Resumos e estatísticas

### Endpoints

#### POST /api/contacts
**Descrição**: Cria novo contato

**Request Body**:
```json
{
    "name": "Maria Silva",
    "email": "maria@empresa.com",
    "phone": "(11) 99999-9999",
    "company": "Empresa ABC",
    "position": "Gerente",
    "type": "LEAD",
    "notes": "Interessada em nossos serviços"
}
```

**Response (201)**:
```json
{
    "id": 1,
    "name": "Maria Silva",
    "email": "maria@empresa.com",
    "phone": "(11) 99999-9999",
    "company": "Empresa ABC",
    "position": "Gerente",
    "type": "LEAD",
    "notes": "Interessada em nossos serviços",
    "user_id": 1,
    "created_at": "2024-01-01T10:00:00Z"
}
```

#### GET /api/contacts
**Descrição**: Lista contatos com filtros

**Query Parameters**:
- `type`: CLIENT ou LEAD
- `search`: busca em nome, email, empresa
- `limit`: limite de resultados (padrão: 50)
- `offset`: offset para paginação

**Response (200)**:
```json
[
    {
        "id": 1,
        "name": "Maria Silva",
        "email": "maria@empresa.com",
        "type": "LEAD",
        "company": "Empresa ABC"
    }
]
```

#### GET /api/contacts/{id}
**Descrição**: Obtém contato específico

**Response (200)**:
```json
{
    "id": 1,
    "name": "Maria Silva",
    "email": "maria@empresa.com",
    "phone": "(11) 99999-9999",
    "company": "Empresa ABC",
    "position": "Gerente",
    "type": "LEAD",
    "notes": "Interessada em nossos serviços",
    "user_id": 1,
    "created_at": "2024-01-01T10:00:00Z"
}
```

#### GET /api/contacts/{id}/details
**Descrição**: Detalhes completos com relacionamentos

**Response (200)**:
```json
{
    "contact": { /* dados do contato */ },
    "interactions": [
        {
            "id": 1,
            "type": "EMAIL",
            "subject": "Proposta comercial",
            "date": "2024-01-01T14:00:00Z"
        }
    ],
    "tasks": [
        {
            "id": 1,
            "title": "Enviar proposta",
            "status": "PENDING",
            "due_date": "2024-01-02T10:00:00Z"
        }
    ],
    "projects": []
}
```

#### PUT /api/contacts/{id}/convert-to-client
**Descrição**: Converte lead em cliente

**Response (200)**:
```json
{
    "message": "Lead convertido em cliente com sucesso",
    "contact": {
        "id": 1,
        "name": "Maria Silva",
        "type": "CLIENT",
        "updated_at": "2024-01-01T15:00:00Z"
    }
}
```

## InteractionHandler

### Responsabilidades
- Histórico de comunicações
- Interações por contato
- Filtros por tipo e data
- Interações recentes

### Endpoints

#### POST /api/contacts/{contactId}/interactions
**Descrição**: Cria nova interação para contato

**Request Body**:
```json
{
    "type": "EMAIL",
    "date": "2024-01-01T14:00:00Z",
    "subject": "Proposta comercial",
    "description": "Enviada proposta detalhada por email"
}
```

**Response (201)**:
```json
{
    "id": 1,
    "type": "EMAIL",
    "date": "2024-01-01T14:00:00Z",
    "subject": "Proposta comercial",
    "description": "Enviada proposta detalhada por email",
    "contact_id": 1,
    "contact": {
        "id": 1,
        "name": "Maria Silva"
    }
}
```

#### GET /api/contacts/{contactId}/interactions
**Descrição**: Lista interações de um contato

**Query Parameters**:
- `type`: EMAIL, CALL, MEETING, OTHER
- `date_from`: data inicial
- `date_to`: data final
- `limit`: limite de resultados

**Response (200)**:
```json
[
    {
        "id": 1,
        "type": "EMAIL",
        "subject": "Proposta comercial",
        "date": "2024-01-01T14:00:00Z",
        "contact": {
            "id": 1,
            "name": "Maria Silva"
        }
    }
]
```

#### GET /api/interactions/recent
**Descrição**: Interações mais recentes

**Query Parameters**:
- `limit`: número de interações (padrão: 10)

**Response (200)**:
```json
[
    {
        "id": 1,
        "type": "EMAIL",
        "subject": "Proposta comercial",
        "date": "2024-01-01T14:00:00Z",
        "contact": {
            "id": 1,
            "name": "Maria Silva"
        }
    }
]
```

## TaskHandler

### Responsabilidades
- Gestão completa de tarefas
- Prioridades e status
- Associações com contatos/projetos
- Tarefas em atraso e próximas

### Endpoints

#### POST /api/tasks
**Descrição**: Cria nova tarefa

**Request Body**:
```json
{
    "title": "Enviar proposta",
    "description": "Preparar e enviar proposta comercial",
    "due_date": "2024-01-02T10:00:00Z",
    "priority": "HIGH",
    "status": "PENDING",
    "contact_id": 1,
    "project_id": null
}
```

**Response (201)**:
```json
{
    "id": 1,
    "title": "Enviar proposta",
    "description": "Preparar e enviar proposta comercial",
    "due_date": "2024-01-02T10:00:00Z",
    "priority": "HIGH",
    "status": "PENDING",
    "user_id": 1,
    "contact_id": 1,
    "contact": {
        "id": 1,
        "name": "Maria Silva"
    }
}
```

#### GET /api/tasks
**Descrição**: Lista tarefas com filtros

**Query Parameters**:
- `status`: PENDING, COMPLETED
- `priority`: LOW, MEDIUM, HIGH
- `contact_id`: ID do contato
- `project_id`: ID do projeto
- `due_before`: vencimento antes de
- `due_after`: vencimento depois de

**Response (200)**:
```json
[
    {
        "id": 1,
        "title": "Enviar proposta",
        "priority": "HIGH",
        "status": "PENDING",
        "due_date": "2024-01-02T10:00:00Z",
        "contact": {
            "id": 1,
            "name": "Maria Silva"
        }
    }
]
```

#### PUT /api/tasks/{id}/complete
**Descrição**: Marca tarefa como concluída

**Response (200)**:
```json
{
    "message": "Tarefa marcada como concluída",
    "task": {
        "id": 1,
        "title": "Enviar proposta",
        "status": "COMPLETED",
        "updated_at": "2024-01-01T16:00:00Z"
    }
}
```

#### GET /api/tasks/overdue
**Descrição**: Tarefas em atraso

**Response (200)**:
```json
[
    {
        "id": 2,
        "title": "Tarefa atrasada",
        "due_date": "2023-12-30T10:00:00Z",
        "priority": "HIGH",
        "contact": {
            "name": "Cliente ABC"
        }
    }
]
```

## ProjectHandler

### Responsabilidades
- Gestão de projetos
- Associação com clientes
- Status e progresso
- Resumos com estatísticas

### Endpoints

#### POST /api/projects
**Descrição**: Cria novo projeto

**Request Body**:
```json
{
    "name": "Website Corporativo",
    "description": "Desenvolvimento de website institucional",
    "status": "IN_PROGRESS",
    "client_id": 1
}
```

**Response (201)**:
```json
{
    "id": 1,
    "name": "Website Corporativo",
    "description": "Desenvolvimento de website institucional",
    "status": "IN_PROGRESS",
    "user_id": 1,
    "client_id": 1,
    "client": {
        "id": 1,
        "name": "Maria Silva",
        "company": "Empresa ABC"
    },
    "created_at": "2024-01-01T10:00:00Z"
}
```

#### GET /api/projects
**Descrição**: Lista projetos com filtros

**Query Parameters**:
- `status`: IN_PROGRESS, COMPLETED, CANCELLED
- `client_id`: ID do cliente
- `limit`: limite de resultados
- `offset`: offset para paginação

**Response (200)**:
```json
[
    {
        "id": 1,
        "name": "Website Corporativo",
        "status": "IN_PROGRESS",
        "client": {
            "id": 1,
            "name": "Maria Silva",
            "company": "Empresa ABC"
        },
        "created_at": "2024-01-01T10:00:00Z"
    }
]
```

#### GET /api/projects/{id}/summary
**Descrição**: Resumo detalhado do projeto

**Response (200)**:
```json
{
    "project": { /* dados do projeto */ },
    "total_tasks": 10,
    "completed_tasks": 6,
    "pending_tasks": 4,
    "overdue_tasks": 1,
    "tasks_progress": 60.0
}
```

#### PUT /api/projects/{id}/status
**Descrição**: Altera status do projeto

**Request Body**:
```json
{
    "status": "COMPLETED"
}
```

**Response (200)**:
```json
{
    "message": "Status do projeto alterado com sucesso",
    "project": {
        "id": 1,
        "name": "Website Corporativo",
        "status": "COMPLETED",
        "updated_at": "2024-01-01T17:00:00Z"
    }
}
```

## Padrões de Implementação

### Validação de Entrada

```go
// Validação de JSON
if err := c.ShouldBindJSON(&req); err != nil {
    c.Error(errors.NewBadRequestError("Dados inválidos: " + err.Error()))
    return
}

// Validação de parâmetros
if req.Email == "" {
    c.Error(errors.NewBadRequestError("Email é obrigatório"))
    return
}
```

### Extração de Parâmetros

```go
// ID da URL
idStr := c.Param("id")
id, err := strconv.ParseUint(idStr, 10, 32)
if err != nil {
    c.Error(errors.NewBadRequestError("ID inválido"))
    return
}

// Query parameters
var filter models.ListFilter
if err := c.ShouldBindQuery(&filter); err != nil {
    c.Error(errors.NewBadRequestError("Parâmetros inválidos"))
    return
}
```

### Tratamento de Erros

```go
// Propagação de erros do service
result, err := h.service.Operation(params)
if err != nil {
    c.Error(err) // Middleware processa
    return
}

// Resposta de sucesso
c.JSON(http.StatusOK, result)
```

### Respostas Padronizadas

```go
// Criação bem-sucedida
c.JSON(http.StatusCreated, resource)

// Atualização bem-sucedida
c.JSON(http.StatusOK, gin.H{
    "message": "Recurso atualizado com sucesso",
    "data": resource,
})

// Exclusão bem-sucedida
c.Status(http.StatusNoContent)
```

## Middleware de Suporte

### Middleware de Autenticação

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := extractToken(c)
        userID, err := validateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
            c.Abort()
            return
        }
        
        c.Set("user_id", userID)
        c.Next()
    }
}
```

### Middleware de Erro

```go
func ErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            handleError(c, err)
        }
    }
}
```

## Documentação Swagger

Todos os endpoints incluem documentação completa para geração automática da documentação da API:

```go
// @Summary Criar novo contato
// @Description Cria um novo contato (cliente ou lead)
// @Tags contacts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.ContactCreateRequest true "Dados do contato"
// @Success 201 {object} models.Contact
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Router /api/contacts [post]
```

Este conjunto completo de handlers fornece uma API REST robusta e bem documentada para o sistema CRM, seguindo as melhores práticas de desenvolvimento web em Go.

