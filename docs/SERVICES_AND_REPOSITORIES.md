# Services e Repositories - CRM Backend

## Visão Geral

Este documento detalha a implementação dos Services e Repositories do backend GoLang do CRM. Estas camadas são fundamentais para a arquitetura do sistema, implementando respectivamente a lógica de negócio e o acesso aos dados.

## Repositories

Os repositories implementam o padrão Repository, abstraindo o acesso ao banco de dados e fornecendo uma interface limpa para operações CRUD e consultas específicas.

### UserRepository

**Localização**: `internal/repositories/user_repository.go`

O UserRepository gerencia todas as operações relacionadas aos usuários no banco de dados.

#### Interface
```go
type UserRepository interface {
    Create(user *models.User) error
    GetByID(id uint) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
    EmailExists(email string) (bool, error)
}
```

#### Funcionalidades Principais
- **Create**: Cria um novo usuário no banco de dados
- **GetByID**: Busca usuário por ID único
- **GetByEmail**: Busca usuário por email (usado na autenticação)
- **Update**: Atualiza dados de um usuário existente
- **Delete**: Remove usuário (soft delete via GORM)
- **EmailExists**: Verifica se um email já está em uso

### InteractionRepository

**Localização**: `internal/repositories/interaction_repository.go`

Gerencia as interações entre usuários e seus contatos.

#### Interface
```go
type InteractionRepository interface {
    Create(interaction *models.Interaction) error
    GetByID(id uint) (*models.Interaction, error)
    GetByContactID(contactID uint, filter *models.InteractionListFilter) ([]models.Interaction, error)
    Update(interaction *models.Interaction) error
    Delete(id uint) error
    GetByUserID(userID uint, filter *models.InteractionListFilter) ([]models.Interaction, error)
    CountByContactID(contactID uint) (int64, error)
}
```

#### Funcionalidades Principais
- **Filtros Avançados**: Suporte a filtros por tipo, data, contato
- **Paginação**: Implementa limit e offset para grandes volumes de dados
- **Ordenação**: Ordena por data (mais recente primeiro)
- **Relacionamentos**: Carrega automaticamente dados do contato associado
- **Agregações**: Conta interações por contato

### TaskRepository

**Localização**: `internal/repositories/task_repository.go`

Gerencia tarefas dos usuários com funcionalidades avançadas de filtro e busca.

#### Interface
```go
type TaskRepository interface {
    Create(task *models.Task) error
    GetByID(id uint) (*models.Task, error)
    GetByUserID(userID uint, filter *models.TaskListFilter) ([]models.Task, error)
    Update(task *models.Task) error
    Delete(id uint) error
    GetByContactID(contactID uint) ([]models.Task, error)
    GetByProjectID(projectID uint) ([]models.Task, error)
    CountByUserID(userID uint) (int64, error)
    CountPendingByUserID(userID uint) (int64, error)
    GetOverdueTasks(userID uint) ([]models.Task, error)
}
```

#### Funcionalidades Principais
- **Filtros Múltiplos**: Status, prioridade, contato, projeto, datas
- **Ordenação Inteligente**: Por prioridade e data de vencimento
- **Tarefas em Atraso**: Busca específica para tarefas vencidas
- **Estatísticas**: Contadores para dashboard e relatórios
- **Associações**: Suporte a tarefas vinculadas a contatos e projetos

### ProjectRepository

**Localização**: `internal/repositories/project_repository.go`

Gerencia projetos e suas relações com clientes e tarefas.

#### Interface
```go
type ProjectRepository interface {
    Create(project *models.Project) error
    GetByID(id uint) (*models.Project, error)
    GetByUserID(userID uint, filter *models.ProjectListFilter) ([]models.Project, error)
    Update(project *models.Project) error
    Delete(id uint) error
    GetByClientID(clientID uint) ([]models.Project, error)
    CountByUserID(userID uint) (int64, error)
    CountByStatus(userID uint, status models.ProjectStatus) (int64, error)
    GetWithTasks(id uint) (*models.Project, error)
}
```

#### Funcionalidades Principais
- **Relacionamentos Complexos**: Carrega cliente, usuário e tarefas
- **Filtros por Status**: IN_PROGRESS, COMPLETED, CANCELLED
- **Estatísticas por Status**: Para dashboards e relatórios
- **Busca por Cliente**: Lista todos os projetos de um cliente específico

## Services

Os services implementam a lógica de negócio, orquestrando operações entre diferentes repositories e aplicando regras de validação e autorização.

### AuthService

**Localização**: `internal/services/auth_service.go`

Responsável por toda a lógica de autenticação e autorização do sistema.

#### Interface
```go
type AuthService interface {
    Register(req *models.UserCreateRequest) (*models.UserResponse, error)
    Login(email, password string) (string, *models.UserResponse, error)
    ValidateToken(tokenString string) (uint, error)
    GenerateToken(userID uint) (string, error)
}
```

#### Funcionalidades Principais

**Registro de Usuários**
- Validação de email único
- Hash seguro de senhas usando bcrypt
- Criação de conta com dados básicos

**Login e Autenticação**
- Verificação de credenciais
- Geração de tokens JWT com expiração de 24 horas
- Retorno de dados do usuário (sem senha)

**Validação de Tokens**
- Parse e validação de tokens JWT
- Verificação de expiração
- Extração de claims do usuário
- Validação de existência do usuário

**Segurança**
- Uso de bcrypt para hash de senhas
- Tokens JWT com assinatura HMAC
- Validação rigorosa de tokens
- Proteção contra ataques de timing

### UserService

**Localização**: `internal/services/user_service.go`

Gerencia operações relacionadas ao perfil e dados do usuário.

#### Interface
```go
type UserService interface {
    GetProfile(userID uint) (*models.UserResponse, error)
    UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
    ChangePassword(userID uint, currentPassword, newPassword string) error
    DeleteAccount(userID uint, password string) error
    GetUserStats(userID uint) (*UserStats, error)
}
```

#### Funcionalidades Principais

**Gestão de Perfil**
- Visualização de dados do usuário
- Atualização de nome e email
- Validação de email único em atualizações

**Gestão de Senha**
- Alteração segura de senha
- Verificação da senha atual
- Hash da nova senha

**Exclusão de Conta**
- Verificação de senha para confirmação
- Soft delete preservando integridade referencial

**Estatísticas do Usuário**
- Dashboard com métricas consolidadas
- Contadores de contatos, tarefas, projetos
- Estatísticas por status e tipo

### InteractionService

**Localização**: `internal/services/interaction_service.go`

Gerencia todas as operações relacionadas às interações com contatos.

#### Interface
```go
type InteractionService interface {
    Create(userID, contactID uint, req *models.InteractionCreateRequest) (*models.Interaction, error)
    GetByID(userID, interactionID uint) (*models.Interaction, error)
    GetByContactID(userID, contactID uint, filter *models.InteractionListFilter) ([]models.Interaction, error)
    GetByUserID(userID uint, filter *models.InteractionListFilter) ([]models.Interaction, error)
    Update(userID, interactionID uint, req *models.InteractionUpdateRequest) (*models.Interaction, error)
    Delete(userID, interactionID uint) error
    GetRecentInteractions(userID uint, limit int) ([]models.Interaction, error)
}
```

#### Funcionalidades Principais

**Criação e Gestão**
- Validação de propriedade do contato
- Criação com tipos: EMAIL, CALL, MEETING, OTHER
- Atualização de dados existentes

**Consultas e Filtros**
- Busca por contato específico
- Filtros por tipo, data, contato
- Paginação para grandes volumes
- Interações recentes para dashboard

**Autorização**
- Verificação de propriedade em todas as operações
- Proteção contra acesso não autorizado
- Validação de relacionamentos

### TaskService

**Localização**: `internal/services/task_service.go`

Implementa toda a lógica de negócio para gestão de tarefas.

#### Interface
```go
type TaskService interface {
    Create(userID uint, req *models.TaskCreateRequest) (*models.Task, error)
    GetByID(userID, taskID uint) (*models.Task, error)
    GetByUserID(userID uint, filter *models.TaskListFilter) ([]models.Task, error)
    Update(userID, taskID uint, req *models.TaskUpdateRequest) (*models.Task, error)
    Delete(userID, taskID uint) error
    MarkAsCompleted(userID, taskID uint) (*models.Task, error)
    MarkAsPending(userID, taskID uint) (*models.Task, error)
    GetByContactID(userID, contactID uint) ([]models.Task, error)
    GetByProjectID(userID, projectID uint) ([]models.Task, error)
    GetOverdueTasks(userID uint) ([]models.Task, error)
    GetUpcomingTasks(userID uint, days int) ([]models.Task, error)
}
```

#### Funcionalidades Principais

**Gestão Completa de Tarefas**
- Criação com prioridades: LOW, MEDIUM, HIGH
- Status: PENDING, COMPLETED
- Associação opcional a contatos e projetos
- Datas de vencimento

**Operações de Status**
- Marcação rápida como concluída/pendente
- Histórico de alterações via timestamps

**Consultas Especializadas**
- Tarefas em atraso
- Tarefas próximas do vencimento
- Tarefas por contato/projeto
- Filtros avançados

**Validações de Negócio**
- Verificação de propriedade de associações
- Validação de relacionamentos
- Autorização em todas as operações

### ProjectService

**Localização**: `internal/services/project_service.go`

Gerencia projetos e suas complexas relações com clientes e tarefas.

#### Interface
```go
type ProjectService interface {
    Create(userID uint, req *models.ProjectCreateRequest) (*models.Project, error)
    GetByID(userID, projectID uint) (*models.Project, error)
    GetWithTasks(userID, projectID uint) (*models.Project, error)
    GetByUserID(userID uint, filter *models.ProjectListFilter) ([]models.Project, error)
    Update(userID, projectID uint, req *models.ProjectUpdateRequest) (*models.Project, error)
    Delete(userID, projectID uint) error
    GetByClientID(userID, clientID uint) ([]models.Project, error)
    ChangeStatus(userID, projectID uint, status models.ProjectStatus) (*models.Project, error)
    GetProjectSummary(userID, projectID uint) (*ProjectSummary, error)
}
```

#### Funcionalidades Principais

**Gestão de Projetos**
- Criação vinculada a clientes (tipo CLIENT)
- Status: IN_PROGRESS, COMPLETED, CANCELLED
- Validação de tipo de contato

**Relacionamentos Complexos**
- Associação obrigatória com cliente
- Carregamento de tarefas relacionadas
- Histórico completo do projeto

**Resumos e Relatórios**
- ProjectSummary com estatísticas
- Progresso baseado em tarefas
- Métricas de conclusão

**Regras de Negócio**
- Proteção contra exclusão com tarefas ativas
- Validação de mudanças de cliente
- Controle de status do projeto

## Padrões de Implementação

### Tratamento de Erros

Todos os services utilizam o sistema de erros customizado:

```go
// Erro de recurso não encontrado
if err == gorm.ErrRecordNotFound {
    return nil, errors.NewNotFoundError("Recurso")
}

// Erro de autorização
if resource.UserID != userID {
    return nil, errors.ErrForbidden
}

// Erro de conflito
if exists {
    return nil, errors.NewConflictError("Email já está em uso")
}
```

### Validação de Autorização

Padrão consistente de verificação de propriedade:

```go
// Verificar se o recurso pertence ao usuário
resource, err := s.repository.GetByID(resourceID)
if err != nil {
    return handleRepositoryError(err)
}

if resource.UserID != userID {
    return nil, errors.ErrForbidden
}
```

### Paginação e Filtros

Implementação padronizada de filtros:

```go
// Aplicar valores padrão
if filter == nil {
    filter = &models.FilterType{}
}
if filter.Limit == 0 {
    filter.Limit = 50 // Limite padrão
}
```

### Relacionamentos

Carregamento consistente de relacionamentos:

```go
// Buscar com relacionamentos
resource, err := s.repository.GetByID(id)
if err != nil {
    return nil, errors.ErrInternalServer
}

// Retornar com dados relacionados carregados
return resource, nil
```

## Integração entre Camadas

### Fluxo de Dados

1. **Handler** recebe requisição HTTP
2. **Handler** valida entrada e extrai contexto
3. **Service** aplica lógica de negócio
4. **Repository** executa operação no banco
5. **Service** processa resultado
6. **Handler** formata resposta

### Injeção de Dependências

```go
// main.go
userRepo := repositories.NewUserRepository(db)
authService := services.NewAuthService(userRepo, cfg.JWTSecret)
authHandler := handlers.NewAuthHandler(authService)
```

### Testes

Cada camada pode ser testada independentemente:

```go
// Mock do repository para testar service
mockRepo := &MockUserRepository{}
service := services.NewUserService(mockRepo)

// Teste do service
result, err := service.GetProfile(1)
assert.NoError(t, err)
assert.NotNil(t, result)
```

## Próximos Passos

Para completar a implementação:

1. **Handlers**: Criar controladores HTTP para cada service
2. **Middleware**: Integrar autenticação nos handlers
3. **Validação**: Implementar validações de entrada
4. **Testes**: Criar testes unitários e de integração
5. **Documentação**: Gerar documentação da API

## Considerações de Performance

### Otimizações Implementadas

- **Preload de Relacionamentos**: Evita N+1 queries
- **Índices de Banco**: Definidos nos models
- **Paginação**: Limita resultados grandes
- **Filtros no Banco**: Reduz transferência de dados

### Melhorias Futuras

- **Cache**: Implementar cache para consultas frequentes
- **Conexão Pool**: Otimizar pool de conexões do banco
- **Queries Otimizadas**: Implementar queries específicas para relatórios
- **Compressão**: Comprimir respostas grandes

Este conjunto de repositories e services fornece uma base sólida e escalável para o CRM, seguindo as melhores práticas de desenvolvimento em Go e arquitetura limpa.



### ContactRepository

**Localização**: `internal/repositories/contact_repository.go`

O ContactRepository é responsável por todas as operações relacionadas aos contatos no banco de dados, incluindo clientes e leads.

#### Interface
```go
type ContactRepository interface {
    Create(contact *models.Contact) error
    GetByID(id uint) (*models.Contact, error)
    GetByUserID(userID uint, filter *models.ContactListFilter) ([]models.Contact, error)
    Update(contact *models.Contact) error
    Delete(id uint) error
    GetByEmail(email string) (*models.Contact, error)
    CountByUserID(userID uint) (int64, error)
    CountByType(userID uint, contactType models.ContactType) (int64, error)
    SearchByName(userID uint, name string) ([]models.Contact, error)
    GetWithInteractions(id uint) (*models.Contact, error)
    GetWithTasks(id uint) (*models.Contact, error)
    GetWithProjects(id uint) (*models.Contact, error)
}
```

#### Funcionalidades Principais

**Operações CRUD Básicas**
- **Create**: Cria um novo contato no banco de dados
- **GetByID**: Busca contato por ID com preload do usuário
- **Update**: Atualiza dados de um contato existente
- **Delete**: Remove contato (soft delete via GORM)

**Consultas Especializadas**
- **GetByEmail**: Busca contato por email único
- **SearchByName**: Busca parcial por nome com ILIKE
- **GetByUserID**: Lista contatos do usuário com filtros avançados

**Relacionamentos Complexos**
- **GetWithInteractions**: Carrega contato com histórico de interações ordenado por data
- **GetWithTasks**: Carrega contato com tarefas ordenadas por vencimento
- **GetWithProjects**: Carrega contato com projetos ordenados por criação

**Filtros e Busca Avançada**
- Filtro por tipo (CLIENT, LEAD)
- Busca textual em nome, email e empresa
- Paginação com limit e offset
- Ordenação alfabética por nome

**Estatísticas e Agregações**
- **CountByUserID**: Conta total de contatos do usuário
- **CountByType**: Conta contatos por tipo específico

### ContactService

**Localização**: `internal/services/contact_service.go`

O ContactService implementa toda a lógica de negócio relacionada aos contatos, incluindo validações, conversões e operações complexas.

#### Interface
```go
type ContactService interface {
    Create(userID uint, req *models.ContactCreateRequest) (*models.Contact, error)
    GetByID(userID, contactID uint) (*models.Contact, error)
    GetWithDetails(userID, contactID uint) (*ContactDetails, error)
    GetByUserID(userID uint, filter *models.ContactListFilter) ([]models.Contact, error)
    Update(userID, contactID uint, req *models.ContactUpdateRequest) (*models.Contact, error)
    Delete(userID, contactID uint) error
    SearchByName(userID uint, name string) ([]models.Contact, error)
    GetContactSummary(userID, contactID uint) (*ContactSummary, error)
    ConvertLeadToClient(userID, contactID uint) (*models.Contact, error)
}
```

#### Funcionalidades Principais

**Gestão Completa de Contatos**
- Criação com validação de email único por usuário
- Atualização com verificação de conflitos
- Exclusão com proteção para clientes com projetos
- Busca e listagem com filtros

**Detalhes e Relacionamentos**
- **ContactDetails**: Estrutura completa com interações, tarefas e projetos
- **ContactSummary**: Estatísticas consolidadas do contato
- Carregamento otimizado de relacionamentos

**Conversão de Leads**
- **ConvertLeadToClient**: Converte lead em cliente
- Validação de tipo antes da conversão
- Atualização automática do status

**Validações de Negócio**
- Email único por usuário
- Verificação de propriedade em todas as operações
- Proteção contra exclusão de clientes com projetos ativos
- Validação de tipos de contato

**Estatísticas Avançadas**
- Contadores de interações, tarefas e projetos
- Progresso de tarefas (concluídas vs pendentes)
- Status de projetos (ativos vs concluídos)
- Data da última interação

## Handlers

Os handlers implementam a camada de apresentação da API REST, gerenciando requisições HTTP, validações de entrada e formatação de respostas.

### AuthHandler

**Localização**: `internal/handlers/auth_handler.go`

Gerencia todas as operações de autenticação e autorização do sistema.

#### Endpoints Implementados

**POST /api/auth/register**
- Registra novo usuário no sistema
- Validação de dados de entrada
- Verificação de email único
- Hash seguro da senha
- Retorna dados do usuário (sem senha)

**POST /api/auth/login**
- Autentica usuário com email e senha
- Validação de credenciais
- Geração de token JWT
- Retorna token e dados do usuário

**POST /api/auth/logout**
- Logout do usuário (stateless)
- Instrução para remoção do token no frontend

**GET /api/auth/validate**
- Valida token JWT atual
- Protegido por middleware de autenticação
- Retorna status de validade

#### Estruturas de Dados

```go
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
    Token  string `json:"token"`
    UserID uint   `json:"user_id"`
    Email  string `json:"email"`
    Name   string `json:"name"`
}
```

### UserHandler

**Localização**: `internal/handlers/user_handler.go`

Gerencia operações relacionadas ao perfil e dados do usuário autenticado.

#### Endpoints Implementados

**GET /api/users/profile**
- Obtém perfil do usuário autenticado
- Retorna dados sem senha
- Protegido por autenticação

**PUT /api/users/profile**
- Atualiza dados do perfil
- Validação de email único
- Campos opcionais para atualização parcial

**PUT /api/users/change-password**
- Altera senha do usuário
- Verificação da senha atual
- Validação de confirmação
- Hash seguro da nova senha

**DELETE /api/users/delete-account**
- Exclui conta do usuário
- Confirmação obrigatória com senha
- Soft delete preservando integridade

**GET /api/users/stats**
- Estatísticas consolidadas do usuário
- Contadores de contatos, tarefas, projetos
- Métricas por status e tipo

### ContactHandler

**Localização**: `internal/handlers/contact_handler.go`

Gerencia todas as operações CRUD e funcionalidades especiais dos contatos.

#### Endpoints Implementados

**POST /api/contacts**
- Cria novo contato (cliente ou lead)
- Validação de dados obrigatórios
- Verificação de email único

**GET /api/contacts**
- Lista contatos do usuário
- Filtros por tipo, busca textual
- Paginação com limit/offset

**GET /api/contacts/{id}**
- Obtém contato específico
- Verificação de propriedade
- Dados básicos do contato

**GET /api/contacts/{id}/details**
- Detalhes completos do contato
- Inclui interações, tarefas e projetos
- Carregamento otimizado de relacionamentos

**PUT /api/contacts/{id}**
- Atualiza contato existente
- Validação de conflitos de email
- Atualização parcial de campos

**DELETE /api/contacts/{id}**
- Exclui contato
- Proteção para clientes com projetos
- Soft delete com cascata

**GET /api/contacts/search**
- Busca contatos por nome
- Busca parcial case-insensitive
- Parâmetro obrigatório 'q'

**GET /api/contacts/{id}/summary**
- Resumo estatístico do contato
- Contadores de atividades
- Progresso e métricas

**PUT /api/contacts/{id}/convert-to-client**
- Converte lead em cliente
- Validação de tipo atual
- Atualização automática

### InteractionHandler

**Localização**: `internal/handlers/interaction_handler.go`

Gerencia o histórico de comunicações e interações com contatos.

#### Endpoints Implementados

**POST /api/contacts/{contactId}/interactions**
- Cria nova interação para contato
- Validação de propriedade do contato
- Tipos: EMAIL, CALL, MEETING, OTHER

**GET /api/contacts/{contactId}/interactions**
- Lista interações de um contato
- Filtros por tipo e data
- Ordenação cronológica reversa

**GET /api/interactions**
- Lista todas as interações do usuário
- Filtros avançados por contato, tipo, data
- Paginação e ordenação

**GET /api/interactions/{id}**
- Obtém interação específica
- Verificação de propriedade
- Dados completos com relacionamentos

**PUT /api/interactions/{id}**
- Atualiza interação existente
- Campos opcionais para atualização
- Preservação de relacionamentos

**DELETE /api/interactions/{id}**
- Exclui interação
- Verificação de propriedade
- Remoção permanente

**GET /api/interactions/recent**
- Interações mais recentes
- Limite configurável (padrão: 10)
- Para dashboard e resumos

### TaskHandler

**Localização**: `internal/handlers/task_handler.go`

Gerencia o sistema completo de tarefas com prioridades, status e associações.

#### Endpoints Implementados

**POST /api/tasks**
- Cria nova tarefa
- Associação opcional a contatos/projetos
- Prioridades: LOW, MEDIUM, HIGH
- Status: PENDING, COMPLETED

**GET /api/tasks**
- Lista tarefas do usuário
- Filtros por status, prioridade, associações
- Ordenação por prioridade e vencimento

**GET /api/tasks/{id}**
- Obtém tarefa específica
- Dados completos com relacionamentos
- Verificação de propriedade

**PUT /api/tasks/{id}**
- Atualiza tarefa existente
- Validação de associações
- Atualização parcial de campos

**DELETE /api/tasks/{id}**
- Exclui tarefa
- Verificação de propriedade
- Remoção permanente

**PUT /api/tasks/{id}/complete**
- Marca tarefa como concluída
- Atualização automática de status
- Timestamp de conclusão

**PUT /api/tasks/{id}/pending**
- Marca tarefa como pendente
- Reversão de conclusão
- Atualização de status

**GET /api/contacts/{contactId}/tasks**
- Tarefas de um contato específico
- Verificação de propriedade do contato
- Ordenação por vencimento

**GET /api/projects/{projectId}/tasks**
- Tarefas de um projeto específico
- Verificação de propriedade do projeto
- Ordenação por vencimento

**GET /api/tasks/overdue**
- Tarefas em atraso
- Filtro automático por data atual
- Ordenação por vencimento

**GET /api/tasks/upcoming**
- Tarefas próximas do vencimento
- Parâmetro 'days' configurável (padrão: 7)
- Para alertas e planejamento

### ProjectHandler

**Localização**: `internal/handlers/project_handler.go`

Gerencia projetos e suas complexas relações com clientes e tarefas.

#### Endpoints Implementados

**POST /api/projects**
- Cria novo projeto
- Associação obrigatória com cliente
- Validação de tipo de contato (CLIENT)
- Status inicial configurável

**GET /api/projects**
- Lista projetos do usuário
- Filtros por status e cliente
- Paginação e ordenação

**GET /api/projects/{id}**
- Obtém projeto específico
- Dados básicos com relacionamentos
- Verificação de propriedade

**GET /api/projects/{id}/with-tasks**
- Projeto com todas as tarefas
- Carregamento otimizado
- Ordenação de tarefas por vencimento

**PUT /api/projects/{id}**
- Atualiza projeto existente
- Validação de novo cliente
- Atualização parcial de campos

**DELETE /api/projects/{id}**
- Exclui projeto
- Proteção contra exclusão com tarefas ativas
- Verificação de dependências

**GET /api/clients/{clientId}/projects**
- Projetos de um cliente específico
- Verificação de propriedade do cliente
- Ordenação cronológica

**PUT /api/projects/{id}/status**
- Altera status do projeto
- Status: IN_PROGRESS, COMPLETED, CANCELLED
- Validação de transições

**GET /api/projects/{id}/summary**
- Resumo detalhado do projeto
- Estatísticas de tarefas
- Progresso percentual
- Métricas de conclusão

## Padrões de Implementação dos Handlers

### Validação de Entrada

Todos os handlers seguem padrões consistentes de validação:

```go
// Validação de JSON
if err := c.ShouldBindJSON(&req); err != nil {
    c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
    return
}

// Validação de parâmetros de URL
contactIDStr := c.Param("id")
contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
if err != nil {
    c.Error(errors.NewBadRequestError("ID do contato inválido"))
    return
}
```

### Tratamento de Erros

Sistema unificado de tratamento de erros:

```go
// Chamar service e propagar erro
result, err := h.service.Operation(userID, params)
if err != nil {
    c.Error(err) // Middleware de erro processa
    return
}

// Resposta de sucesso
c.JSON(http.StatusOK, result)
```

### Autenticação e Autorização

Extração consistente do usuário autenticado:

```go
// Obtido do middleware de autenticação
userID := c.GetUint("user_id")

// Usado em todas as operações do service
result, err := h.service.GetByUserID(userID, filter)
```

### Documentação Swagger

Todos os endpoints incluem documentação completa:

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

### Respostas Padronizadas

Formato consistente de respostas:

```go
// Sucesso simples
c.JSON(http.StatusOK, result)

// Sucesso com mensagem
c.JSON(http.StatusOK, gin.H{
    "message": "Operação realizada com sucesso",
    "data":    result,
})

// Exclusão bem-sucedida
c.Status(http.StatusNoContent)
```

## Integração Completa

### Fluxo de Requisição

1. **Middleware CORS**: Permite requisições cross-origin
2. **Middleware Logger**: Registra todas as requisições
3. **Middleware Auth**: Valida token JWT (rotas protegidas)
4. **Handler**: Processa requisição e valida entrada
5. **Service**: Aplica lógica de negócio e validações
6. **Repository**: Executa operações no banco de dados
7. **Middleware Error**: Formata e retorna erros
8. **Response**: JSON estruturado para o cliente

### Estrutura de Rotas

```go
// Rotas públicas
auth := router.Group("/api/auth")
{
    auth.POST("/register", authHandler.Register)
    auth.POST("/login", authHandler.Login)
}

// Rotas protegidas
api := router.Group("/api")
api.Use(middleware.AuthMiddleware())
{
    // Usuários
    users := api.Group("/users")
    users.GET("/profile", userHandler.GetProfile)
    users.PUT("/profile", userHandler.UpdateProfile)
    
    // Contatos
    contacts := api.Group("/contacts")
    contacts.POST("", contactHandler.Create)
    contacts.GET("", contactHandler.List)
    contacts.GET("/:id", contactHandler.GetByID)
    
    // Interações aninhadas
    contacts.POST("/:contactId/interactions", interactionHandler.Create)
    contacts.GET("/:contactId/interactions", interactionHandler.ListByContact)
}
```

### Middleware de Erro

Processamento centralizado de erros:

```go
func ErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            
            switch e := err.(type) {
            case *errors.BadRequestError:
                c.JSON(http.StatusBadRequest, gin.H{"error": e.Message})
            case *errors.NotFoundError:
                c.JSON(http.StatusNotFound, gin.H{"error": e.Message})
            case *errors.UnauthorizedError:
                c.JSON(http.StatusUnauthorized, gin.H{"error": e.Message})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
            }
        }
    }
}
```

Este conjunto completo de repositories, services e handlers fornece uma API REST robusta e bem estruturada para o sistema CRM, seguindo as melhores práticas de desenvolvimento em Go e arquitetura limpa.

