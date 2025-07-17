package services

import (
	"crm-backend/internal/models"
	"crm-backend/internal/repositories"
	"crm-backend/pkg/errors"
	"fmt"

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
	GetRecentActivities(userID uint, limit int) (*models.RecentActivityResponse, error)
}

// UserStats representa estatísticas do usuário
type UserStats struct {
	TotalContacts     int64 `json:"total_contacts"`
	TotalClients      int64 `json:"total_clients"`
	TotalLeads        int64 `json:"total_leads"`
	TotalTasks        int64 `json:"total_tasks"`
	PendingTasks      int64 `json:"pending_tasks"`
	CompletedTasks    int64 `json:"completed_tasks"`
	TotalProjects     int64 `json:"total_projects"`
	ActiveProjects    int64 `json:"active_projects"`
	CompletedProjects int64 `json:"completed_projects"`
	TotalInteractions int64 `json:"total_interactions"`
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

// GetRecentActivities obtém as atividades recentes do usuário
func (s *userService) GetRecentActivities(userID uint, limit int) (*models.RecentActivityResponse, error) {
	if limit <= 0 {
		limit = 10 // Limite padrão
	}

	activities := []models.UserActivity{}

	// 1. Buscar interações recentes
	interactionFilter := &models.InteractionListFilter{
		Limit: limit,
	}
	interactions, err := s.interactionRepo.GetByUserID(userID, interactionFilter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Converter interações para atividades
	for _, interaction := range interactions {
		contactID := interaction.ContactID
		activity := models.UserActivity{
			ID:        interaction.ID,
			Type:      models.ActivityTypeInteraction,
			Action:    string(interaction.Type),
			Title:     interaction.Subject,
			Detail:    truncateString(interaction.Description, 100),
			ItemID:    interaction.ID,
			CreatedAt: interaction.CreatedAt,
			RelatedID: &contactID,
		}

		if interaction.Contact.Name != "" {
			contactName := interaction.Contact.Name
			activity.RelatedName = &contactName
		}

		activities = append(activities, activity)
	}

	// 2. Buscar tarefas recentes
	taskFilter := &models.TaskListFilter{
		Limit: limit,
	}
	tasks, err := s.taskRepo.GetByUserID(userID, taskFilter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Converter tarefas para atividades
	for _, task := range tasks {
		action := "Adicionada"
		if task.Status == models.TaskStatusCompleted {
			action = "Concluída"
		}

		activity := models.UserActivity{
			ID:        task.ID,
			Type:      models.ActivityTypeTask,
			Action:    action,
			Title:     task.Title,
			Detail:    truncateString(task.Description, 100),
			ItemID:    task.ID,
			CreatedAt: task.CreatedAt,
		}

		if task.ContactID != nil && task.Contact != nil && task.Contact.Name != "" {
			activity.RelatedID = task.ContactID
			contactName := task.Contact.Name
			activity.RelatedName = &contactName
		} else if task.ProjectID != nil && task.Project != nil && task.Project.Name != "" {
			activity.RelatedID = task.ProjectID
			projectName := task.Project.Name
			activity.RelatedName = &projectName
		}

		activities = append(activities, activity)
	}

	// 3. Buscar projetos recentes
	projectFilter := &models.ProjectListFilter{
		Limit: limit,
	}
	projects, err := s.projectRepo.GetByUserID(userID, projectFilter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Converter projetos para atividades
	for _, project := range projects {
		var action string
		switch project.Status {
		case models.ProjectStatusInProgress:
			action = "Em andamento"
		case models.ProjectStatusCompleted:
			action = "Concluído"
		case models.ProjectStatusCancelled:
			action = "Cancelado"
		default:
			action = "Criado"
		}

		activity := models.UserActivity{
			ID:        project.ID,
			Type:      models.ActivityTypeProject,
			Action:    action,
			Title:     project.Name,
			Detail:    truncateString(project.Description, 100),
			ItemID:    project.ID,
			CreatedAt: project.CreatedAt,
		}

		if project.ClientID != 0 && project.Client.Name != "" {
			clientID := project.ClientID
			activity.RelatedID = &clientID
			clientName := project.Client.Name
			activity.RelatedName = &clientName
		}

		activities = append(activities, activity)
	}

	// 4. Buscar contatos recentes
	contactFilter := &models.ContactListFilter{
		Limit: limit,
	}
	contacts, err := s.contactRepo.GetByUserID(userID, contactFilter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Converter contatos para atividades
	for _, contact := range contacts {
		action := fmt.Sprintf("Novo %s", contact.Type)

		activity := models.UserActivity{
			ID:        contact.ID,
			Type:      models.ActivityTypeContact,
			Action:    action,
			Title:     contact.Name,
			Detail:    truncateString(contact.Notes, 100),
			ItemID:    contact.ID,
			CreatedAt: contact.CreatedAt,
		}

		activities = append(activities, activity)
	}

	// Ordenar todas as atividades por data (mais recente primeiro)
	sortActivitiesByDate(activities)

	// Limitar ao número solicitado
	if len(activities) > limit {
		activities = activities[:limit]
	}

	response := &models.RecentActivityResponse{
		Activities: activities,
		Count:      len(activities),
	}

	return response, nil
} // Helper para truncar strings longas
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// Helper para ordenar atividades por data (mais recentes primeiro)
func sortActivitiesByDate(activities []models.UserActivity) {
	// Simple bubble sort (pode ser substituído por sort.Slice para melhor performance)
	for i := 0; i < len(activities)-1; i++ {
		for j := 0; j < len(activities)-i-1; j++ {
			if activities[j].CreatedAt.Before(activities[j+1].CreatedAt) {
				activities[j], activities[j+1] = activities[j+1], activities[j]
			}
		}
	}
}
