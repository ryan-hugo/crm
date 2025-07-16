# Arquitetura do Sistema CRM Backend

## Visão Geral da Arquitetura

O CRM Backend foi projetado seguindo os princípios de Clean Architecture e Domain-Driven Design (DDD), adaptados para as convenções e idiomas da linguagem Go. A arquitetura é organizada em camadas bem definidas que promovem separação de responsabilidades, testabilidade e manutenibilidade.

## Princípios Arquiteturais

### 1. Separação de Responsabilidades
Cada camada tem uma responsabilidade específica e bem definida:
- **Handlers**: Gerenciam requisições HTTP e respostas
- **Services**: Implementam regras de negócio
- **Repositories**: Abstraem acesso a dados
- **Models**: Definem estruturas de dados

### 2. Inversão de Dependências
As camadas superiores não dependem diretamente das inferiores. Interfaces são utilizadas para inverter dependências e facilitar testes.

### 3. Isolamento de Domínio
A lógica de negócio está isolada em services, independente de frameworks e tecnologias específicas.

## Estrutura de Camadas

### Camada de Apresentação (Handlers)
**Localização**: `internal/handlers/`

Responsável por:
- Receber e validar requisições HTTP
- Converter dados de entrada para DTOs
- Chamar services apropriados
- Formatar respostas HTTP
- Tratamento de erros específicos da API

**Exemplo de Handler**:
```go
type ContactHandler struct {
    contactService services.ContactService
}

func (h *ContactHandler) Create(c *gin.Context) {
    var req models.ContactCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(errors.NewBadRequestError(err.Error()))
        return
    }
    
    userID := c.GetUint("user_id")
    contact, err := h.contactService.Create(userID, &req)
    if err != nil {
        c.Error(err)
        return
    }
    
    c.JSON(http.StatusCreated, contact)
}
```

### Camada de Serviços (Services)
**Localização**: `internal/services/`

Responsável por:
- Implementar regras de negócio
- Orquestrar operações entre diferentes repositories
- Validar dados de negócio
- Aplicar políticas de autorização
- Coordenar transações

**Exemplo de Service**:
```go
type ContactService interface {
    Create(userID uint, req *models.ContactCreateRequest) (*models.Contact, error)
    GetByID(userID, contactID uint) (*models.Contact, error)
    List(userID uint, filter *models.ContactListFilter) ([]models.Contact, error)
    Update(userID, contactID uint, req *models.ContactUpdateRequest) (*models.Contact, error)
    Delete(userID, contactID uint) error
}
```

### Camada de Repositórios (Repositories)
**Localização**: `internal/repositories/`

Responsável por:
- Abstrair acesso ao banco de dados
- Implementar operações CRUD
- Executar consultas específicas
- Gerenciar relacionamentos entre entidades

**Exemplo de Repository**:
```go
type ContactRepository interface {
    Create(contact *models.Contact) error
    GetByID(id uint) (*models.Contact, error)
    GetByUserID(userID uint, filter *models.ContactListFilter) ([]models.Contact, error)
    Update(contact *models.Contact) error
    Delete(id uint) error
}
```

### Camada de Modelos (Models)
**Localização**: `internal/models/`

Responsável por:
- Definir estruturas de dados do domínio
- Especificar relacionamentos entre entidades
- Definir DTOs para entrada e saída de dados
- Implementar validações de dados

## Fluxo de Dados

### Requisição HTTP
1. **Middleware**: Processa autenticação, CORS, logging
2. **Handler**: Recebe requisição, valida entrada
3. **Service**: Aplica regras de negócio
4. **Repository**: Acessa banco de dados
5. **Response**: Retorna dados formatados

### Exemplo de Fluxo Completo
```
POST /api/contacts
    ↓
[CORS Middleware] → [Auth Middleware] → [Logger Middleware]
    ↓
[ContactHandler.Create]
    ↓ (ContactCreateRequest)
[ContactService.Create]
    ↓ (Business Logic + Validation)
[ContactRepository.Create]
    ↓ (Database Operation)
[GORM] → [PostgreSQL]
    ↑ (Contact Entity)
[ContactHandler.Create]
    ↑ (JSON Response)
HTTP 201 Created
```

## Padrões de Design Utilizados

### 1. Repository Pattern
Abstrai a camada de persistência, permitindo trocar implementações sem afetar a lógica de negócio.

```go
type UserRepository interface {
    Create(user *models.User) error
    GetByEmail(email string) (*models.User, error)
    GetByID(id uint) (*models.User, error)
    Update(user *models.User) error
}
```

### 2. Service Layer Pattern
Encapsula a lógica de negócio em uma camada dedicada.

```go
type AuthService interface {
    Register(req *models.UserCreateRequest) (*models.User, error)
    Login(email, password string) (string, error)
    ValidateToken(token string) (*models.User, error)
}
```

### 3. DTO (Data Transfer Object) Pattern
Separa modelos de domínio dos modelos de API.

```go
// Modelo de domínio
type Contact struct {
    ID       uint   `gorm:"primaryKey"`
    Name     string `gorm:"not null"`
    Email    string `gorm:"not null"`
    // ... outros campos
}

// DTO para criação
type ContactCreateRequest struct {
    Name  string      `json:"name" validate:"required"`
    Email string      `json:"email" validate:"required,email"`
    Type  ContactType `json:"type" validate:"required"`
}
```

### 4. Middleware Pattern
Implementa funcionalidades transversais como autenticação, logging e tratamento de erros.

```go
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Validar token JWT
        // Adicionar user_id ao contexto
        c.Next()
    }
}
```

## Gerenciamento de Dependências

### Injeção de Dependências
As dependências são injetadas através de construtores, facilitando testes e manutenção.

```go
// main.go
userRepo := repositories.NewUserRepository(db)
authService := services.NewAuthService(userRepo, cfg.JWTSecret)
authHandler := handlers.NewAuthHandler(authService)
```

### Interfaces
Todas as dependências são definidas através de interfaces, permitindo mock em testes.

```go
type ContactService interface {
    Create(userID uint, req *models.ContactCreateRequest) (*models.Contact, error)
    // ... outros métodos
}

type contactService struct {
    contactRepo repositories.ContactRepository
}
```

## Tratamento de Erros

### Erros Customizados
Sistema de erros tipados para diferentes cenários.

```go
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

var (
    ErrNotFound     = NewAppError(404, "Recurso não encontrado", "")
    ErrUnauthorized = NewAppError(401, "Não autorizado", "")
    ErrBadRequest   = NewAppError(400, "Requisição inválida", "")
)
```

### Middleware de Tratamento
Captura e formata erros de forma consistente.

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            err := c.Errors.Last()
            if appErr, ok := err.Err.(*errors.AppError); ok {
                c.JSON(appErr.Code, appErr)
                return
            }
            c.JSON(500, gin.H{"error": "Erro interno"})
        }
    }
}
```

## Segurança

### Autenticação JWT
Tokens JWT são utilizados para autenticação stateless.

```go
func GenerateToken(userID uint, secret string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

### Autorização
Cada operação verifica se o usuário tem permissão para acessar o recurso.

```go
func (s *contactService) GetByID(userID, contactID uint) (*models.Contact, error) {
    contact, err := s.contactRepo.GetByID(contactID)
    if err != nil {
        return nil, err
    }
    
    // Verificar se o contato pertence ao usuário
    if contact.UserID != userID {
        return nil, errors.ErrForbidden
    }
    
    return contact, nil
}
```

## Configuração

### Variáveis de Ambiente
Configurações são carregadas de variáveis de ambiente com valores padrão.

```go
type Config struct {
    DatabaseURL string
    JWTSecret   string
    Port        string
    Environment string
}

func Load() *Config {
    return &Config{
        DatabaseURL: getEnv("DATABASE_URL", "postgres://localhost/crm_db"),
        JWTSecret:   getEnv("JWT_SECRET", "default-secret"),
        Port:        getEnv("PORT", "8080"),
        Environment: getEnv("ENVIRONMENT", "development"),
    }
}
```

## Banco de Dados

### ORM e Migrações
GORM é utilizado para mapeamento objeto-relacional e migrações automáticas.

```go
func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
        &models.Contact{},
        &models.Interaction{},
        &models.Task{},
        &models.Project{},
    )
}
```

### Relacionamentos
Relacionamentos são definidos através de tags GORM.

```go
type Contact struct {
    ID           uint          `gorm:"primaryKey"`
    UserID       uint          `gorm:"not null"`
    User         User          `gorm:"foreignKey:UserID"`
    Interactions []Interaction `gorm:"foreignKey:ContactID"`
    Tasks        []Task        `gorm:"foreignKey:ContactID"`
}
```

## Logging

### Sistema de Logs
Logs estruturados para diferentes níveis de severidade.

```go
func Info(v ...interface{}) {
    InfoLogger.Println(v...)
}

func Error(v ...interface{}) {
    ErrorLogger.Println(v...)
}
```

## Testabilidade

### Interfaces para Mocking
Todas as dependências são interfaces, facilitando criação de mocks.

```go
type MockContactRepository struct {
    contacts map[uint]*models.Contact
}

func (m *MockContactRepository) Create(contact *models.Contact) error {
    contact.ID = uint(len(m.contacts) + 1)
    m.contacts[contact.ID] = contact
    return nil
}
```

### Testes de Unidade
Cada camada pode ser testada independentemente.

```go
func TestContactService_Create(t *testing.T) {
    mockRepo := &MockContactRepository{contacts: make(map[uint]*models.Contact)}
    service := services.NewContactService(mockRepo)
    
    req := &models.ContactCreateRequest{
        Name:  "João Silva",
        Email: "joao@example.com",
        Type:  models.ContactTypeClient,
    }
    
    contact, err := service.Create(1, req)
    assert.NoError(t, err)
    assert.Equal(t, "João Silva", contact.Name)
}
```

## Escalabilidade

### Preparação para Microserviços
A arquitetura em camadas facilita a divisão em microserviços no futuro.

### Cache
Estrutura preparada para implementação de cache em diferentes camadas.

### Monitoramento
Logs estruturados facilitam implementação de monitoramento e observabilidade.

## Conclusão

Esta arquitetura fornece uma base sólida para o desenvolvimento do CRM Backend, seguindo as melhores práticas da linguagem Go e padrões de arquitetura modernos. A estrutura é flexível o suficiente para evoluir conforme as necessidades do projeto, mantendo a qualidade e manutenibilidade do código.

