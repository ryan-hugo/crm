package services

import (
	"crm-backend/internal/models"
	"crm-backend/internal/repositories"
	"crm-backend/pkg/errors"
	"sort"
	"time"

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
	GetDashboardData(userID uint) (*DashboardData, error)
}

// UserStats representa estatísticas do usuário
type UserStats struct {
	TotalContacts      int64 `json:"total_contacts"`
	TotalClients       int64 `json:"total_clients"`
	TotalLeads         int64 `json:"total_leads"`
	TotalTasks         int64 `json:"total_tasks"`
	PendingTasks       int64 `json:"pending_tasks"`
	CompletedTasks     int64 `json:"completed_tasks"`
	OverdueTasks       int64 `json:"overdue_tasks"`
	TotalProjects      int64 `json:"total_projects"`
	ActiveProjects     int64 `json:"active_projects"`
	CompletedProjects  int64 `json:"completed_projects"`
	TotalInteractions  int64 `json:"total_interactions"`
	RecentInteractions int64 `json:"recent_interactions"`
}

// DashboardProject representa um resumo de projeto para o dashboard
type DashboardProject struct {
	ID         uint                 `json:"id"`
	Name       string               `json:"name"`
	Status     models.ProjectStatus `json:"status"`
	ClientName string               `json:"client_name"`
	CreatedAt  time.Time            `json:"created_at"`
}

// DashboardInteraction representa um resumo de interação para o dashboard
type DashboardInteraction struct {
	ID          uint                   `json:"id"`
	Type        models.InteractionType `json:"type"`
	Subject     string                 `json:"subject"`
	ContactName string                 `json:"contact_name"`
	Date        time.Time              `json:"date"`
}

// DashboardTask representa um resumo de tarefa para o dashboard
type DashboardTask struct {
	ID          uint            `json:"id"`
	Title       string          `json:"title"`
	Priority    models.Priority `json:"priority"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
	ContactName string          `json:"contact_name,omitempty"`
	ProjectName string          `json:"project_name,omitempty"`
}

// DashboardContact representa contatos para o dashboard
type DashboardContact struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	Email     string             `json:"email"`
	Type      models.ContactType `json:"type"`
	Company   string             `json:"company,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
}

// DashboardData representa os dados completos para o dashboard
type DashboardData struct {
	Stats              UserStats              `json:"stats"`
	RecentActivities   []models.UserActivity  `json:"recent_activities"`
	RecentProjects     []DashboardProject     `json:"recent_projects"`
	RecentInteractions []DashboardInteraction `json:"recent_interactions"`
	RecentPendingTasks []DashboardTask        `json:"recent_pending_tasks"`
	RecentContacts     []DashboardContact     `json:"recent_contacts"`
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
	stats := &UserStats{
		RecentInteractions: 0, // Inicializar explicitamente
		OverdueTasks:       0, // Inicializar explicitamente
	}

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

		// Contar tarefas em atraso
		overdueTasks, err := s.taskRepo.CountOverdueByUserID(userID)
		if err != nil {
			// Se houver erro, definir como 0 mas incluir no resultado
			stats.OverdueTasks = 0
		} else {
			stats.OverdueTasks = overdueTasks
		}
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

		// Contar interações recentes dos últimos 7 dias
		recentInteractions, err := s.interactionRepo.GetRecentByUserID(userID, 7, 100) // limite alto para contar todas
		if err != nil {
			// Se houver erro, definir como 0 mas incluir no resultado
			stats.RecentInteractions = 0
		} else {
			stats.RecentInteractions = int64(len(recentInteractions))
		}

		// // Para debug: garantir que sempre tenha pelo menos 0
		// if stats.RecentInteractions < 0 {
		// 	stats.RecentInteractions = 0
		// }
	}

	return stats, nil
}

// GetRecentActivities obtém as atividades recentes do usuário
func (s *userService) GetRecentActivities(userID uint, limit int) (*models.RecentActivityResponse, error) {
	if limit <= 0 {
		limit = 20 // Limite padrão aumentado para capturar mais atividades
	}

	activities := []models.UserActivity{}

	// 1. Buscar interações recentes (ordenadas por created_at/updated_at)
	interactions, err := s.interactionRepo.GetRecentByUserID(userID, 30, limit*2) // Buscar mais para filtrar depois
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Converter interações para atividades
	for _, interaction := range interactions {
		// Atividade de criação
		createActivity := createActivityFromInteraction(interaction)
		activities = append(activities, createActivity)

		// Se foi atualizada depois da criação, adicionar atividade de atualização
		if interaction.UpdatedAt.After(interaction.CreatedAt.Add(time.Minute)) {
			updateActivity := createActivity
			updateActivity.Action = models.ActionUpdated
			updateActivity.CreatedAt = interaction.UpdatedAt
			updateActivity.UpdatedAt = interaction.UpdatedAt
			activities = append(activities, updateActivity)
		}
	}

	// 2. Buscar tarefas recentes
	taskFilter := &models.TaskListFilter{
		Limit: limit * 2,
	}
	tasks, err := s.taskRepo.GetByUserID(userID, taskFilter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Converter tarefas para atividades
	for _, task := range tasks {
		// Atividade de criação
		createActivity := createActivityFromTask(task)
		createActivity.Action = models.ActionCreated
		activities = append(activities, createActivity)

		// Se foi atualizada depois da criação, adicionar atividade de atualização
		if task.UpdatedAt.After(task.CreatedAt.Add(time.Minute)) {
			updateActivity := createActivity
			updateActivity.Action = models.ActionUpdated
			updateActivity.CreatedAt = task.UpdatedAt
			updateActivity.UpdatedAt = task.UpdatedAt
			activities = append(activities, updateActivity)
		}

		// Se foi concluída, adicionar atividade de conclusão
		if task.Status == models.TaskStatusCompleted {
			completeActivity := createActivity
			completeActivity.Action = models.ActionCompleted
			completeActivity.CreatedAt = task.UpdatedAt
			completeActivity.UpdatedAt = task.UpdatedAt
			activities = append(activities, completeActivity)
		}
	}

	// 3. Buscar projetos recentes
	projectFilter := &models.ProjectListFilter{
		Limit: limit * 2,
	}
	projects, err := s.projectRepo.GetByUserID(userID, projectFilter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Converter projetos para atividades
	for _, project := range projects {
		// Atividade de criação
		createActivity := createActivityFromProject(project)
		createActivity.Action = models.ActionCreated
		activities = append(activities, createActivity)

		// Se foi atualizado depois da criação, adicionar atividade de atualização
		if project.UpdatedAt.After(project.CreatedAt.Add(time.Minute)) {
			updateActivity := createActivity

			// Determinar o tipo de atualização baseado no status
			switch project.Status {
			case models.ProjectStatusInProgress:
				updateActivity.Action = models.ActionStarted
			case models.ProjectStatusCompleted:
				updateActivity.Action = models.ActionCompleted
			case models.ProjectStatusCancelled:
				updateActivity.Action = models.ActionCancelled
			default:
				updateActivity.Action = models.ActionUpdated
			}

			updateActivity.CreatedAt = project.UpdatedAt
			updateActivity.UpdatedAt = project.UpdatedAt
			activities = append(activities, updateActivity)
		}
	}

	// 4. Buscar contatos recentes
	contactFilter := &models.ContactListFilter{
		Limit: limit * 2,
	}
	contacts, err := s.contactRepo.GetByUserID(userID, contactFilter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Converter contatos para atividades
	for _, contact := range contacts {
		// Atividade de criação
		createActivity := createActivityFromContact(contact)
		activities = append(activities, createActivity)

		// Se foi atualizado depois da criação, adicionar atividade de atualização
		if contact.UpdatedAt.After(contact.CreatedAt.Add(time.Minute)) {
			updateActivity := createActivity
			updateActivity.Action = models.ActionUpdated
			updateActivity.CreatedAt = contact.UpdatedAt
			updateActivity.UpdatedAt = contact.UpdatedAt
			activities = append(activities, updateActivity)
		}
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
}

// Funções auxiliares para criar UserActivity de forma segura

// createActivityFromInteraction cria uma UserActivity a partir de uma Interaction
func createActivityFromInteraction(interaction models.Interaction) models.UserActivity {
	title := interaction.Subject
	if title == "" {
		title = "Interação sem assunto"
	}

	contactID := interaction.ContactID
	activity := models.UserActivity{
		ID:        interaction.ID,
		Type:      models.ActivityTypeInteraction,
		Action:    models.ActionCreated,
		Title:     title,
		Detail:    truncateString(interaction.Description, 100),
		ItemID:    interaction.ID,
		CreatedAt: interaction.CreatedAt,
		UpdatedAt: interaction.UpdatedAt,
		RelatedID: &contactID,
	}

	if interaction.Contact.Name != "" {
		contactName := interaction.Contact.Name
		activity.RelatedName = &contactName
	}

	return activity
}

// createActivityFromTask cria uma UserActivity a partir de uma Task
func createActivityFromTask(task models.Task) models.UserActivity {
	var action models.ActivityAction
	if task.Status == models.TaskStatusCompleted {
		action = models.ActionCompleted
	} else {
		action = models.ActionCreated
	}

	title := task.Title
	if title == "" {
		title = "Tarefa sem título"
	}

	activity := models.UserActivity{
		ID:        task.ID,
		Type:      models.ActivityTypeTask,
		Action:    action,
		Title:     title,
		Detail:    truncateString(task.Description, 100),
		ItemID:    task.ID,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
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

	return activity
}

// createActivityFromProject cria uma UserActivity a partir de um Project
func createActivityFromProject(project models.Project) models.UserActivity {
	var action models.ActivityAction
	switch project.Status {
	case models.ProjectStatusInProgress:
		action = models.ActionStarted
	case models.ProjectStatusCompleted:
		action = models.ActionCompleted
	case models.ProjectStatusCancelled:
		action = models.ActionCancelled
	default:
		action = models.ActionCreated
	}

	title := project.Name
	if title == "" {
		title = "Projeto sem nome"
	}

	activity := models.UserActivity{
		ID:        project.ID,
		Type:      models.ActivityTypeProject,
		Action:    action,
		Title:     title,
		Detail:    truncateString(project.Description, 100),
		ItemID:    project.ID,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
	}

	if project.ClientID != 0 && project.Client.Name != "" {
		clientID := project.ClientID
		activity.RelatedID = &clientID
		clientName := project.Client.Name
		activity.RelatedName = &clientName
	}

	return activity
}

// createActivityFromContact cria uma UserActivity a partir de um Contact
func createActivityFromContact(contact models.Contact) models.UserActivity {
	title := contact.Name
	if title == "" {
		title = "Contato sem nome"
	}

	activity := models.UserActivity{
		ID:        contact.ID,
		Type:      models.ActivityTypeContact,
		Action:    models.ActionCreated,
		Title:     title,
		Detail:    truncateString(contact.Notes, 100),
		ItemID:    contact.ID,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}

	return activity
}

// Helper para truncar strings longas
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// Helper para ordenar atividades por data (mais recentes primeiro)
func sortActivitiesByDate(activities []models.UserActivity) {
	sort.Slice(activities, func(i, j int) bool {
		return activities[i].CreatedAt.After(activities[j].CreatedAt)
	})
}

// GetDashboardData obtém dados específicos para o dashboard
func (s *userService) GetDashboardData(userID uint) (*DashboardData, error) {
	// 1. Obter estatísticas do usuário
	stats, err := s.GetUserStats(userID)
	if err != nil {
		return nil, err
	}

	// 2. Obter atividades recentes (limitado a 10 para o dashboard)
	recentActivitiesResponse, err := s.GetRecentActivities(userID, 10)
	if err != nil {
		return nil, err
	}

	dashboardData := &DashboardData{
		Stats:              *stats,
		RecentActivities:   recentActivitiesResponse.Activities,
		RecentProjects:     []DashboardProject{},
		RecentInteractions: []DashboardInteraction{},
		RecentPendingTasks: []DashboardTask{},
		RecentContacts:     []DashboardContact{},
	}

	// 3. Buscar 5 interações mais recentes para o dashboard
	if s.interactionRepo != nil {
		recentFilter := &models.InteractionListFilter{
			Limit: 5,
		}
		recentInteractions, err := s.interactionRepo.GetByUserID(userID, recentFilter)
		if err == nil {
			for _, interaction := range recentInteractions {
				dashboardInteraction := DashboardInteraction{
					ID:          interaction.ID,
					Type:        interaction.Type,
					Subject:     interaction.Subject,
					ContactName: interaction.Contact.Name,
					Date:        interaction.Date,
				}
				dashboardData.RecentInteractions = append(dashboardData.RecentInteractions, dashboardInteraction)
			}
		}
	}

	// Buscar projetos ativos recentes para o dashboard
	if s.projectRepo != nil {
		activeFilter := &models.ProjectListFilter{
			Status: "IN_PROGRESS",
			Limit:  5,
		}
		activeProjects, err := s.projectRepo.GetByUserID(userID, activeFilter)
		if err == nil {
			for _, project := range activeProjects {
				dashboardProject := DashboardProject{
					ID:         project.ID,
					Name:       project.Name,
					Status:     project.Status,
					ClientName: project.Client.Name,
					CreatedAt:  project.CreatedAt,
				}
				dashboardData.RecentProjects = append(dashboardData.RecentProjects, dashboardProject)
			}
		}
	}

	// Buscar tarefas pendentes recentes para o dashboard
	if s.taskRepo != nil {
		pendingFilter := &models.TaskListFilter{
			Status: models.TaskStatusPending,
			Limit:  5,
		}
		pendingTasks, err := s.taskRepo.GetByUserID(userID, pendingFilter)
		if err == nil {
			for _, task := range pendingTasks {
				dashboardTask := DashboardTask{
					ID:       task.ID,
					Title:    task.Title,
					Priority: task.Priority,
					DueDate:  task.DueDate,
				}

				if task.Contact != nil {
					dashboardTask.ContactName = task.Contact.Name
				}
				if task.Project != nil {
					dashboardTask.ProjectName = task.Project.Name
				}

				dashboardData.RecentPendingTasks = append(dashboardData.RecentPendingTasks, dashboardTask)
			}
		}
	}

	// 4. Buscar 5 contatos mais recentes para o dashboard
	if s.contactRepo != nil {
		recentContactFilter := &models.ContactListFilter{
			Limit: 5,
		}
		contacts, err := s.contactRepo.GetByUserID(userID, recentContactFilter)
		if err == nil {
			for _, contact := range contacts {
				dashboardContact := DashboardContact{
					ID:        contact.ID,
					Name:      contact.Name,
					Email:     contact.Email,
					Type:      contact.Type,
					Company:   contact.Company,
					CreatedAt: contact.CreatedAt,
				}

				dashboardData.RecentContacts = append(dashboardData.RecentContacts, dashboardContact)
			}
		}
	}

	return dashboardData, nil
}
