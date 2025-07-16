# Guia de Desenvolvimento - CRM Backend

## Introdução

Este guia fornece instruções detalhadas para desenvolvedores que trabalharão no projeto CRM Backend. Ele cobre desde a configuração do ambiente de desenvolvimento até as práticas recomendadas para contribuição com o projeto.

## Configuração do Ambiente de Desenvolvimento

### Pré-requisitos

#### Go
```bash
# Instalar Go 1.21 ou superior
# No Ubuntu/Debian:
sudo apt update
sudo apt install golang-go

# Verificar instalação
go version
```

#### PostgreSQL
```bash
# No Ubuntu/Debian:
sudo apt install postgresql postgresql-contrib

# Iniciar serviço
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Criar usuário e banco
sudo -u postgres createuser --interactive
sudo -u postgres createdb crm_db
```

#### Git
```bash
# No Ubuntu/Debian:
sudo apt install git

# Configurar Git
git config --global user.name "Seu Nome"
git config --global user.email "seu.email@example.com"
```

### Configuração do Projeto

#### 1. Clone do Repositório
```bash
git clone <repository-url>
cd crm-backend
```

#### 2. Configuração de Variáveis de Ambiente
```bash
# Copiar arquivo de exemplo
cp .env.example .env

# Editar configurações
nano .env
```

**Configurações importantes no .env:**
```env
DATABASE_URL=postgres://username:password@localhost:5432/crm_db?sslmode=disable
JWT_SECRET=sua-chave-secreta-muito-segura-aqui
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
```

#### 3. Instalação de Dependências
```bash
# Baixar dependências
go mod download

# Verificar dependências
go mod verify
```

#### 4. Configuração do Banco de Dados
```bash
# Conectar ao PostgreSQL
psql -U username -d crm_db

# Verificar conexão (dentro do psql)
\dt
\q
```

#### 5. Executar a Aplicação
```bash
# Executar em modo desenvolvimento
go run cmd/main.go

# Ou compilar e executar
go build -o bin/crm-backend cmd/main.go
./bin/crm-backend
```

## Estrutura do Projeto Detalhada

### Diretório `cmd/`
Contém o ponto de entrada da aplicação.

```go
// cmd/main.go
package main

func main() {
    // Inicialização da aplicação
    // Configuração de dependências
    // Inicialização do servidor
}
```

### Diretório `internal/`
Código interno da aplicação, não exportável.

#### `internal/config/`
Gerenciamento de configurações.

```go
// internal/config/config.go
type Config struct {
    DatabaseURL string
    JWTSecret   string
    Port        string
    Environment string
}
```

#### `internal/models/`
Definição de modelos de dados e DTOs.

```go
// internal/models/user.go
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" validate:"required"`
    Email     string    `json:"email" validate:"required,email"`
    Password  string    `json:"-" validate:"required,min=6"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

#### `internal/repositories/` (A ser implementado)
Camada de acesso a dados.

```go
// internal/repositories/user_repository.go
type UserRepository interface {
    Create(user *models.User) error
    GetByEmail(email string) (*models.User, error)
    GetByID(id uint) (*models.User, error)
    Update(user *models.User) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{db: db}
}
```

#### `internal/services/` (A ser implementado)
Lógica de negócio.

```go
// internal/services/user_service.go
type UserService interface {
    GetProfile(userID uint) (*models.UserResponse, error)
    UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
}

type userService struct {
    userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
    return &userService{userRepo: userRepo}
}
```

#### `internal/handlers/` (A ser implementado)
Controladores HTTP.

```go
// internal/handlers/user_handler.go
type UserHandler struct {
    userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
    return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
    userID := c.GetUint("user_id")
    profile, err := h.userService.GetProfile(userID)
    if err != nil {
        c.Error(err)
        return
    }
    c.JSON(http.StatusOK, profile)
}
```

### Diretório `pkg/`
Pacotes reutilizáveis.

#### `pkg/errors/`
Sistema de erros customizados.

#### `pkg/logger/`
Sistema de logging.

#### `pkg/validator/` (A ser implementado)
Validações customizadas.

## Padrões de Desenvolvimento

### Nomenclatura

#### Arquivos
- Use snake_case para nomes de arquivos: `user_service.go`
- Agrupe por funcionalidade: `contact_handler.go`, `contact_service.go`

#### Variáveis e Funções
```go
// Variáveis: camelCase
var userID uint
var contactList []models.Contact

// Funções públicas: PascalCase
func CreateUser(req *UserCreateRequest) error

// Funções privadas: camelCase
func validateEmail(email string) bool

// Constantes: PascalCase ou UPPER_CASE
const DefaultPageSize = 20
const MAX_RETRY_ATTEMPTS = 3
```

#### Interfaces
```go
// Interfaces terminam com o nome da funcionalidade
type UserRepository interface {
    Create(user *models.User) error
}

type ContactService interface {
    Create(userID uint, req *ContactCreateRequest) error
}
```

### Estrutura de Handlers

```go
type ContactHandler struct {
    contactService services.ContactService
}

func NewContactHandler(contactService services.ContactService) *ContactHandler {
    return &ContactHandler{contactService: contactService}
}

func (h *ContactHandler) Create(c *gin.Context) {
    // 1. Validar entrada
    var req models.ContactCreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(errors.NewBadRequestError(err.Error()))
        return
    }

    // 2. Obter contexto do usuário
    userID := c.GetUint("user_id")

    // 3. Chamar service
    contact, err := h.contactService.Create(userID, &req)
    if err != nil {
        c.Error(err)
        return
    }

    // 4. Retornar resposta
    c.JSON(http.StatusCreated, contact)
}
```

### Estrutura de Services

```go
type contactService struct {
    contactRepo repositories.ContactRepository
}

func NewContactService(contactRepo repositories.ContactRepository) ContactService {
    return &contactService{contactRepo: contactRepo}
}

func (s *contactService) Create(userID uint, req *models.ContactCreateRequest) (*models.Contact, error) {
    // 1. Validações de negócio
    if err := s.validateContactData(req); err != nil {
        return nil, err
    }

    // 2. Criar modelo
    contact := &models.Contact{
        Name:   req.Name,
        Email:  req.Email,
        Type:   req.Type,
        UserID: userID,
    }

    // 3. Persistir
    if err := s.contactRepo.Create(contact); err != nil {
        return nil, errors.ErrInternalServer
    }

    return contact, nil
}
```

### Estrutura de Repositories

```go
type contactRepository struct {
    db *gorm.DB
}

func NewContactRepository(db *gorm.DB) ContactRepository {
    return &contactRepository{db: db}
}

func (r *contactRepository) Create(contact *models.Contact) error {
    if err := r.db.Create(contact).Error; err != nil {
        return err
    }
    return nil
}

func (r *contactRepository) GetByUserID(userID uint, filter *models.ContactListFilter) ([]models.Contact, error) {
    var contacts []models.Contact
    query := r.db.Where("user_id = ?", userID)

    // Aplicar filtros
    if filter.Type != "" {
        query = query.Where("type = ?", filter.Type)
    }
    if filter.Search != "" {
        query = query.Where("name ILIKE ? OR email ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
    }

    // Paginação
    if filter.Limit > 0 {
        query = query.Limit(filter.Limit)
    }
    if filter.Offset > 0 {
        query = query.Offset(filter.Offset)
    }

    if err := query.Find(&contacts).Error; err != nil {
        return nil, err
    }

    return contacts, nil
}
```

## Tratamento de Erros

### Criação de Erros
```go
// Erro genérico
return errors.ErrInternalServer

// Erro específico
return errors.NewNotFoundError("Contato")

// Erro com detalhes
return errors.NewBadRequestError("Email já está em uso")
```

### Propagação de Erros
```go
// Em services
contact, err := s.contactRepo.GetByID(contactID)
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, errors.NewNotFoundError("Contato")
    }
    return nil, errors.ErrInternalServer
}
```

## Validações

### Validações de Modelo
```go
type ContactCreateRequest struct {
    Name  string      `json:"name" validate:"required,min=2,max=255"`
    Email string      `json:"email" validate:"required,email"`
    Type  ContactType `json:"type" validate:"required,oneof=CLIENT LEAD"`
}
```

### Validações Customizadas
```go
func (s *contactService) validateContactData(req *models.ContactCreateRequest) error {
    // Verificar se email já existe
    existing, err := s.contactRepo.GetByEmail(req.Email)
    if err == nil && existing != nil {
        return errors.NewConflictError("Email já está em uso")
    }
    
    return nil
}
```

## Testes

### Estrutura de Testes
```
internal/
├── handlers/
│   ├── contact_handler.go
│   └── contact_handler_test.go
├── services/
│   ├── contact_service.go
│   └── contact_service_test.go
└── repositories/
    ├── contact_repository.go
    └── contact_repository_test.go
```

### Testes de Repository
```go
func TestContactRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := repositories.NewContactRepository(db)

    contact := &models.Contact{
        Name:   "João Silva",
        Email:  "joao@example.com",
        Type:   models.ContactTypeClient,
        UserID: 1,
    }

    err := repo.Create(contact)
    assert.NoError(t, err)
    assert.NotZero(t, contact.ID)
}
```

### Testes de Service
```go
func TestContactService_Create(t *testing.T) {
    mockRepo := &MockContactRepository{}
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

### Testes de Handler
```go
func TestContactHandler_Create(t *testing.T) {
    mockService := &MockContactService{}
    handler := handlers.NewContactHandler(mockService)

    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/contacts", handler.Create)

    body := `{"name":"João Silva","email":"joao@example.com","type":"CLIENT"}`
    req := httptest.NewRequest("POST", "/contacts", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## Comandos Úteis

### Desenvolvimento
```bash
# Executar aplicação
go run cmd/main.go

# Executar com hot reload (usando air)
air

# Formatar código
go fmt ./...

# Verificar código
go vet ./...

# Executar testes
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar testes específicos
go test -run TestContactService_Create ./internal/services/
```

### Build e Deploy
```bash
# Build para produção
go build -o bin/crm-backend cmd/main.go

# Build com otimizações
go build -ldflags="-w -s" -o bin/crm-backend cmd/main.go

# Cross-compilation para Linux
GOOS=linux GOARCH=amd64 go build -o bin/crm-backend-linux cmd/main.go
```

### Banco de Dados
```bash
# Conectar ao banco
psql -U username -d crm_db

# Backup do banco
pg_dump crm_db > backup.sql

# Restaurar backup
psql -U username -d crm_db < backup.sql
```

## Ferramentas Recomendadas

### IDEs e Editores
- **VS Code** com extensão Go
- **GoLand** (JetBrains)
- **Vim/Neovim** com vim-go

### Ferramentas de Desenvolvimento
```bash
# Air (hot reload)
go install github.com/cosmtrek/air@latest

# golangci-lint (linter)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# godoc (documentação)
go install golang.org/x/tools/cmd/godoc@latest
```

### Extensões VS Code Recomendadas
- Go (oficial)
- REST Client
- GitLens
- Thunder Client
- PostgreSQL

## Debugging

### Logs de Debug
```go
import "crm-backend/pkg/logger"

func (s *contactService) Create(userID uint, req *models.ContactCreateRequest) (*models.Contact, error) {
    logger.Info("Creating contact for user", userID, "with email", req.Email)
    
    // ... lógica
    
    logger.Info("Contact created successfully with ID", contact.ID)
    return contact, nil
}
```

### Debugging com Delve
```bash
# Instalar delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug da aplicação
dlv debug cmd/main.go

# Debug de testes
dlv test ./internal/services/
```

## Contribuição

### Fluxo de Trabalho
1. **Fork** do repositório
2. **Clone** do fork
3. **Branch** para feature: `git checkout -b feature/nova-funcionalidade`
4. **Desenvolvimento** seguindo os padrões
5. **Testes** para garantir qualidade
6. **Commit** com mensagens descritivas
7. **Push** para o fork
8. **Pull Request** para o repositório principal

### Padrões de Commit
```bash
# Formato: tipo(escopo): descrição

# Exemplos:
git commit -m "feat(contacts): adicionar endpoint de criação de contatos"
git commit -m "fix(auth): corrigir validação de token JWT"
git commit -m "docs(readme): atualizar instruções de instalação"
git commit -m "test(services): adicionar testes para ContactService"
```

### Checklist antes do PR
- [ ] Código formatado (`go fmt`)
- [ ] Sem warnings (`go vet`)
- [ ] Testes passando (`go test`)
- [ ] Documentação atualizada
- [ ] Variáveis de ambiente documentadas
- [ ] Logs apropriados adicionados

## Troubleshooting

### Problemas Comuns

#### Erro de Conexão com Banco
```bash
# Verificar se PostgreSQL está rodando
sudo systemctl status postgresql

# Verificar configurações de conexão
psql -U username -d crm_db -h localhost -p 5432
```

#### Problemas com Dependências
```bash
# Limpar cache de módulos
go clean -modcache

# Redownload dependências
go mod download
```

#### Problemas de Permissão
```bash
# Verificar permissões do usuário PostgreSQL
sudo -u postgres psql -c "\du"

# Dar permissões ao usuário
sudo -u postgres psql -c "ALTER USER username CREATEDB;"
```

## Próximos Passos

Para completar a implementação do projeto, siga esta ordem:

1. **Implementar Repositories**
   - UserRepository
   - ContactRepository
   - InteractionRepository
   - TaskRepository
   - ProjectRepository

2. **Implementar Services**
   - AuthService
   - UserService
   - ContactService
   - InteractionService
   - TaskService
   - ProjectService

3. **Implementar Handlers**
   - AuthHandler
   - UserHandler
   - ContactHandler
   - InteractionHandler
   - TaskHandler
   - ProjectHandler

4. **Adicionar Testes**
   - Testes unitários para cada camada
   - Testes de integração
   - Testes de API

5. **Documentação da API**
   - Swagger/OpenAPI
   - Exemplos de uso
   - Postman Collection

6. **Deploy e Produção**
   - Docker
   - CI/CD
   - Monitoramento

Este guia fornece uma base sólida para o desenvolvimento do projeto. Mantenha-o atualizado conforme o projeto evolui.

