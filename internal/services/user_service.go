package services

import (
	"crm-backend/internal/models"
	"crm-backend/internal/repositories"
	"crm-backend/pkg/errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService define a interface para operações de usuário
type UserService interface {
	GetProfile(userID uint) (*models.UserResponse, error)
	UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error)
	ChangePassword(userID uint, currentPassword, newPassword string) error
	DeleteAccount(userID uint, password string) error
	GetUserStats(userID uint) (*UserStats, error)
}

// UserStats representa estatísticas do usuário
type UserStats struct {
	TotalContacts       int64 `json:"total_contacts"`
	TotalClients        int64 `json:"total_clients"`
	TotalLeads          int64 `json:"total_leads"`
	TotalTasks          int64 `json:"total_tasks"`
	PendingTasks        int64 `json:"pending_tasks"`
	CompletedTasks      int64 `json:"completed_tasks"`
	TotalProjects       int64 `json:"total_projects"`
	ActiveProjects      int64 `json:"active_projects"`
	CompletedProjects   int64 `json:"completed_projects"`
	TotalInteractions   int64 `json:"total_interactions"`
}

// userService implementa UserService
type userService struct {
	userRepo        repositories.UserRepository
	contactRepo     repositories.ContactRepository
	taskRepo        repositories.TaskRepository
	projectRepo     repositories.ProjectRepository
	interactionRepo repositories.InteractionRepository
}

// NewUserService cria uma nova instância do serviço de usuários
func NewUserService(
	userRepo repositories.UserRepository,
	contactRepo repositories.ContactRepository,
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
	interactionRepo repositories.InteractionRepository,
) UserService {
	return &userService{
		userRepo:        userRepo,
		contactRepo:     contactRepo,
		taskRepo:        taskRepo,
		projectRepo:     projectRepo,
		interactionRepo: interactionRepo,
	}
}

// GetProfile obtém o perfil do usuário
func (s *userService) GetProfile(userID uint) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Usuário")
		}
		return nil, errors.ErrInternalServer
	}

	response := user.ToResponse()
	return &response, nil
}

// UpdateProfile atualiza o perfil do usuário
func (s *userService) UpdateProfile(userID uint, req *models.UserUpdateRequest) (*models.UserResponse, error) {
	// Buscar usuário existente
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Usuário")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se o email está sendo alterado e se já existe
	if req.Email != "" && req.Email != user.Email {
		exists, err := s.userRepo.EmailExists(req.Email)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		if exists {
			return nil, errors.NewConflictError("Email já está em uso")
		}
		user.Email = req.Email
	}

	// Atualizar campos fornecidos
	if req.Name != "" {
		user.Name = req.Name
	}

	// Salvar alterações
	if err := s.userRepo.Update(user); err != nil {
		return nil, errors.ErrInternalServer
	}

	response := user.ToResponse()
	return &response, nil
}

// ChangePassword altera a senha do usuário
func (s *userService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Buscar usuário
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("Usuário")
		}
		return errors.ErrInternalServer
	}

	// Verificar senha atual
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return errors.NewUnauthorizedError("Senha atual incorreta")
	}

	// Hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.ErrInternalServer
	}

	// Atualizar senha
	user.Password = string(hashedPassword)
	if err := s.userRepo.Update(user); err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

// DeleteAccount exclui a conta do usuário
func (s *userService) DeleteAccount(userID uint, password string) error {
	// Buscar usuário
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("Usuário")
		}
		return errors.ErrInternalServer
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return errors.NewUnauthorizedError("Senha incorreta")
	}

	// Excluir usuário (soft delete - GORM cuidará das relações)
	if err := s.userRepo.Delete(userID); err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

// GetUserStats obtém estatísticas do usuário
func (s *userService) GetUserStats(userID uint) (*UserStats, error) {
	stats := &UserStats{}

	// Total de contatos
	if s.contactRepo != nil {
		totalContacts, err := s.contactRepo.CountByUserID(userID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.TotalContacts = totalContacts

		// Contatos por tipo
		clients, err := s.contactRepo.CountByType(userID, models.ContactTypeClient)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.TotalClients = clients

		leads, err := s.contactRepo.CountByType(userID, models.ContactTypeLead)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.TotalLeads = leads
	}

	// Estatísticas de tarefas
	if s.taskRepo != nil {
		totalTasks, err := s.taskRepo.CountByUserID(userID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.TotalTasks = totalTasks

		pendingTasks, err := s.taskRepo.CountPendingByUserID(userID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.PendingTasks = pendingTasks
		stats.CompletedTasks = totalTasks - pendingTasks
	}

	// Estatísticas de projetos
	if s.projectRepo != nil {
		totalProjects, err := s.projectRepo.CountByUserID(userID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.TotalProjects = totalProjects

		activeProjects, err := s.projectRepo.CountByStatus(userID, models.ProjectStatusInProgress)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.ActiveProjects = activeProjects

		completedProjects, err := s.projectRepo.CountByStatus(userID, models.ProjectStatusCompleted)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.CompletedProjects = completedProjects
	}

	// Total de interações (através dos contatos do usuário)
	if s.interactionRepo != nil {
		filter := &models.InteractionListFilter{}
		interactions, err := s.interactionRepo.GetByUserID(userID, filter)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		stats.TotalInteractions = int64(len(interactions))
	}

	return stats, nil
}

