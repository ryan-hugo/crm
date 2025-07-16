# CRM Backend - Sistema de Gestão de Clientes para Freelancers

Este projeto implementa o backend de um Sistema de Gestão de Clientes (CRM) básico desenvolvido especificamente para freelancers e profissionais liberais. O sistema foi construído em Go utilizando o framework Gin para APIs RESTful e GORM para persistência de dados.

## Visão Geral

O CRM Backend oferece uma API completa para gerenciamento de:
- **Usuários**: Autenticação e perfis de usuário
- **Contatos**: Clientes e leads com informações detalhadas
- **Interações**: Histórico de comunicações com contatos
- **Tarefas**: Gestão de atividades com prazos e prioridades
- **Projetos**: Organização de trabalhos associados a clientes

## Arquitetura

O projeto segue uma arquitetura em camadas bem definida, seguindo as melhores práticas do Go:

### Estrutura de Diretórios

```
crm-backend/
├── cmd/                    # Ponto de entrada da aplicação
│   └── main.go            # Arquivo principal do servidor
├── internal/              # Código interno da aplicação
│   ├── auth/              # Lógica de autenticação
│   ├── config/            # Configurações da aplicação
│   ├── database/          # Conexão e migrações do banco
│   ├── handlers/          # Controladores HTTP (a ser criado)
│   ├── middleware/        # Middlewares HTTP
│   ├── models/            # Modelos de dados e DTOs
│   ├── repositories/      # Camada de acesso a dados (a ser criado)
│   ├── services/          # Lógica de negócio (a ser criado)
│   └── utils/             # Utilitários diversos (a ser criado)
├── pkg/                   # Pacotes reutilizáveis
│   ├── errors/            # Tratamento de erros customizados
│   ├── logger/            # Sistema de logging
│   └── validator/         # Validações customizadas (a ser criado)
├── api/                   # Documentação da API
│   └── v1/                # Versão 1 da API (a ser criado)
├── docs/                  # Documentação adicional
├── scripts/               # Scripts de automação
├── migrations/            # Migrações do banco de dados
├── go.mod                 # Dependências do Go
├── .env.example           # Exemplo de variáveis de ambiente
└── README.md              # Este arquivo
```

## Tecnologias Utilizadas

### Core
- **Go 1.21**: Linguagem de programação principal
- **Gin**: Framework web para APIs RESTful
- **GORM**: ORM para Go com suporte a PostgreSQL
- **PostgreSQL**: Banco de dados relacional

### Autenticação e Segurança
- **JWT (golang-jwt/jwt/v5)**: Tokens de autenticação
- **bcrypt (golang.org/x/crypto)**: Hash de senhas

### Configuração e Ambiente
- **godotenv**: Carregamento de variáveis de ambiente

## Modelos de Dados

### User (Usuário)
Representa os profissionais que utilizam o sistema.

**Campos principais:**
- `ID`: Identificador único
- `Name`: Nome completo
- `Email`: Email único para login
- `Password`: Senha hasheada
- `CreatedAt/UpdatedAt`: Timestamps de auditoria

### Contact (Contato)
Armazena informações de clientes e leads.

**Campos principais:**
- `ID`: Identificador único
- `Name`: Nome do contato
- `Email`: Email do contato
- `Phone`: Telefone (opcional)
- `Company`: Empresa (opcional)
- `Position`: Cargo (opcional)
- `Type`: Tipo (CLIENT ou LEAD)
- `Notes`: Observações
- `UserID`: Referência ao usuário proprietário

### Interaction (Interação)
Registra comunicações com contatos.

**Campos principais:**
- `ID`: Identificador único
- `Type`: Tipo (EMAIL, CALL, MEETING, OTHER)
- `Date`: Data da interação
- `Subject`: Assunto
- `Description`: Descrição detalhada
- `ContactID`: Referência ao contato

### Task (Tarefa)
Gerencia atividades e tarefas.

**Campos principais:**
- `ID`: Identificador único
- `Title`: Título da tarefa
- `Description`: Descrição detalhada
- `DueDate`: Data de vencimento (opcional)
- `Priority`: Prioridade (LOW, MEDIUM, HIGH)
- `Status`: Status (PENDING, COMPLETED)
- `UserID`: Referência ao usuário
- `ContactID`: Referência ao contato (opcional)
- `ProjectID`: Referência ao projeto (opcional)

### Project (Projeto)
Organiza trabalhos maiores.

**Campos principais:**
- `ID`: Identificador único
- `Name`: Nome do projeto
- `Description`: Descrição
- `Status`: Status (IN_PROGRESS, COMPLETED, CANCELLED)
- `UserID`: Referência ao usuário
- `ClientID`: Referência ao cliente (contato)

## Configuração e Instalação

### Pré-requisitos
- Go 1.21 ou superior
- PostgreSQL 12 ou superior
- Git

### Instalação

1. **Clone o repositório:**
```bash
git clone <repository-url>
cd crm-backend
```

2. **Instale as dependências:**
```bash
go mod download
```

3. **Configure as variáveis de ambiente:**
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

4. **Configure o banco de dados:**
```bash
# Crie um banco de dados PostgreSQL
createdb crm_db
```

5. **Execute a aplicação:**
```bash
go run cmd/main.go
```

### Variáveis de Ambiente

| Variável | Descrição | Valor Padrão |
|----------|-----------|--------------|
| `DATABASE_URL` | URL de conexão com PostgreSQL | `postgres://localhost:5432/crm_db?sslmode=disable` |
| `JWT_SECRET` | Chave secreta para JWT | `default-secret-key` |
| `PORT` | Porta do servidor | `8080` |
| `ENVIRONMENT` | Ambiente de execução | `development` |
| `LOG_LEVEL` | Nível de log | `info` |

## Estrutura da API

### Autenticação
- `POST /api/auth/register` - Registro de usuário
- `POST /api/auth/login` - Login de usuário

### Usuários (Protegido)
- `GET /api/users/profile` - Obter perfil do usuário
- `PUT /api/users/profile` - Atualizar perfil

### Contatos (Protegido)
- `POST /api/contacts` - Criar contato
- `GET /api/contacts` - Listar contatos
- `GET /api/contacts/:id` - Obter contato específico
- `PUT /api/contacts/:id` - Atualizar contato
- `DELETE /api/contacts/:id` - Excluir contato

### Interações (Protegido)
- `POST /api/contacts/:id/interactions` - Criar interação
- `GET /api/contacts/:id/interactions` - Listar interações do contato

### Tarefas (Protegido)
- `POST /api/tasks` - Criar tarefa
- `GET /api/tasks` - Listar tarefas
- `GET /api/tasks/:id` - Obter tarefa específica
- `PUT /api/tasks/:id` - Atualizar tarefa
- `DELETE /api/tasks/:id` - Excluir tarefa

### Projetos (Protegido)
- `POST /api/projects` - Criar projeto
- `GET /api/projects` - Listar projetos
- `GET /api/projects/:id` - Obter projeto específico
- `PUT /api/projects/:id` - Atualizar projeto
- `DELETE /api/projects/:id` - Excluir projeto

## Middleware

### CORS
Configurado para permitir requisições de qualquer origem durante desenvolvimento.

### Autenticação JWT
Protege rotas que requerem autenticação, validando tokens JWT no header `Authorization: Bearer <token>`.

### Logging
Registra todas as requisições HTTP com detalhes de método, path, status, latência e IP do cliente.

### Tratamento de Erros
Captura e formata erros de forma consistente, retornando respostas JSON padronizadas.

## Segurança

### Autenticação
- Senhas são hasheadas usando bcrypt
- Autenticação baseada em JWT com expiração
- Tokens devem ser enviados no header Authorization

### Autorização
- Cada usuário acessa apenas seus próprios dados
- Validação de propriedade em todas as operações

### Validação
- Validação de entrada em todos os endpoints
- Sanitização de dados para prevenir injeções
- Validação de tipos e formatos de dados

## Desenvolvimento

### Estrutura de Camadas

1. **Handlers**: Recebem requisições HTTP e delegam para services
2. **Services**: Contêm lógica de negócio e orquestram repositories
3. **Repositories**: Fazem acesso direto ao banco de dados
4. **Models**: Definem estruturas de dados e DTOs

### Padrões Utilizados

- **Repository Pattern**: Abstração da camada de dados
- **Service Layer**: Encapsulamento da lógica de negócio
- **DTO Pattern**: Separação entre modelos de domínio e API
- **Middleware Pattern**: Funcionalidades transversais

### Próximos Passos

Para completar a implementação, os seguintes componentes precisam ser criados:

1. **Repositories**: Implementar acesso a dados para cada entidade
2. **Services**: Implementar lógica de negócio
3. **Handlers**: Implementar controladores HTTP
4. **Auth Service**: Implementar autenticação JWT
5. **Validators**: Implementar validações customizadas
6. **Testes**: Criar testes unitários e de integração
7. **Documentação API**: Gerar documentação Swagger/OpenAPI

## Contribuição

1. Faça fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## Licença

Este projeto está licenciado sob a MIT License - veja o arquivo LICENSE para detalhes.

#   c r m  
 # crm
