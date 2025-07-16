package services

import (
	"crm-backend/internal/models"
	"crm-backend/internal/repositories"
	"crm-backend/pkg/errors"

	"gorm.io/gorm"
)

// ProjectService define a interface para operações de projeto
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

// ProjectSummary representa um resumo do projeto
type ProjectSummary struct {
	Project        *models.Project `json:"project"`
	TotalTasks     int64           `json:"total_tasks"`
	CompletedTasks int64           `json:"completed_tasks"`
	PendingTasks   int64           `json:"pending_tasks"`
	OverdueTasks   int64           `json:"overdue_tasks"`
	TasksProgress  float64         `json:"tasks_progress"`
}

// projectService implementa ProjectService
type projectService struct {
	projectRepo repositories.ProjectRepository
	contactRepo repositories.ContactRepository
	taskRepo    repositories.TaskRepository
}

// NewProjectService cria uma nova instância do serviço de projetos
func NewProjectService(
	projectRepo repositories.ProjectRepository,
	contactRepo repositories.ContactRepository,
	taskRepo repositories.TaskRepository,
) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
		contactRepo: contactRepo,
		taskRepo:    taskRepo,
	}
}

// Create cria um novo projeto
func (s *projectService) Create(userID uint, req *models.ProjectCreateRequest) (*models.Project, error) {
	// Verificar se o cliente existe e pertence ao usuário
	client, err := s.contactRepo.GetByID(req.ClientID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Cliente")
		}
		return nil, errors.ErrInternalServer
	}

	if client.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Verificar se o cliente é do tipo CLIENT
	if client.Type != models.ContactTypeClient {
		return nil, errors.NewBadRequestError("O contato deve ser do tipo CLIENT para ser associado a um projeto")
	}

	// Criar projeto
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		UserID:      userID,
		ClientID:    req.ClientID,
	}

	if err := s.projectRepo.Create(project); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar projeto criado com relacionamentos
	createdProject, err := s.projectRepo.GetByID(project.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return createdProject, nil
}

// GetByID obtém um projeto específico
func (s *projectService) GetByID(userID, projectID uint) (*models.Project, error) {
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Projeto")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se o projeto pertence ao usuário
	if project.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return project, nil
}

// GetWithTasks obtém um projeto com suas tarefas
func (s *projectService) GetWithTasks(userID, projectID uint) (*models.Project, error) {
	// Verificar se o projeto pertence ao usuário
	_, err := s.GetByID(userID, projectID)
	if err != nil {
		return nil, err
	}

	// Buscar projeto com tarefas
	projectWithTasks, err := s.projectRepo.GetWithTasks(projectID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return projectWithTasks, nil
}

// GetByUserID obtém todos os projetos do usuário
func (s *projectService) GetByUserID(userID uint, filter *models.ProjectListFilter) ([]models.Project, error) {
	// Aplicar valores padrão ao filtro se necessário
	if filter == nil {
		filter = &models.ProjectListFilter{}
	}
	if filter.Limit == 0 {
		filter.Limit = 50 // Limite padrão
	}

	projects, err := s.projectRepo.GetByUserID(userID, filter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return projects, nil
}

// Update atualiza um projeto existente
func (s *projectService) Update(userID, projectID uint, req *models.ProjectUpdateRequest) (*models.Project, error) {
	// Buscar projeto existente
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Projeto")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se o projeto pertence ao usuário
	if project.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Validar novo cliente se fornecido
	if req.ClientID != 0 {
		client, err := s.contactRepo.GetByID(req.ClientID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Cliente")
			}
			return nil, errors.ErrInternalServer
		}
		if client.UserID != userID {
			return nil, errors.ErrForbidden
		}
		if client.Type != models.ContactTypeClient {
			return nil, errors.NewBadRequestError("O contato deve ser do tipo CLIENT")
		}
		project.ClientID = req.ClientID
	}

	// Atualizar campos fornecidos
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.Status != "" {
		project.Status = req.Status
	}

	// Salvar alterações
	if err := s.projectRepo.Update(project); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar projeto atualizado com relacionamentos
	updatedProject, err := s.projectRepo.GetByID(project.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return updatedProject, nil
}

// Delete exclui um projeto
func (s *projectService) Delete(userID, projectID uint) error {
	// Buscar projeto existente
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("Projeto")
		}
		return errors.ErrInternalServer
	}

	// Verificar se o projeto pertence ao usuário
	if project.UserID != userID {
		return errors.ErrForbidden
	}

	// Verificar se há tarefas associadas
	tasks, err := s.taskRepo.GetByProjectID(projectID)
	if err != nil {
		return errors.ErrInternalServer
	}

	if len(tasks) > 0 {
		return errors.NewBadRequestError("Não é possível excluir projeto com tarefas associadas. Exclua as tarefas primeiro.")
	}

	// Excluir projeto
	if err := s.projectRepo.Delete(projectID); err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

// GetByClientID obtém projetos de um cliente específico
func (s *projectService) GetByClientID(userID, clientID uint) ([]models.Project, error) {
	// Verificar se o cliente existe e pertence ao usuário
	client, err := s.contactRepo.GetByID(clientID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Cliente")
		}
		return nil, errors.ErrInternalServer
	}

	if client.UserID != userID {
		return nil, errors.ErrForbidden
	}

	projects, err := s.projectRepo.GetByClientID(clientID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return projects, nil
}

// ChangeStatus altera o status de um projeto
func (s *projectService) ChangeStatus(userID, projectID uint, status models.ProjectStatus) (*models.Project, error) {
	req := &models.ProjectUpdateRequest{
		Status: status,
	}
	return s.Update(userID, projectID, req)
}

// GetProjectSummary obtém um resumo detalhado do projeto
func (s *projectService) GetProjectSummary(userID, projectID uint) (*ProjectSummary, error) {
	// Buscar projeto
	project, err := s.GetByID(userID, projectID)
	if err != nil {
		return nil, err
	}

	// Buscar tarefas do projeto
	tasks, err := s.taskRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Calcular estatísticas
	summary := &ProjectSummary{
		Project:    project,
		TotalTasks: int64(len(tasks)),
	}

	var completedTasks, pendingTasks, overdueTasks int64
	for _, task := range tasks {
		if task.Status == models.TaskStatusCompleted {
			completedTasks++
		} else {
			pendingTasks++
			// Verificar se está em atraso (implementação básica)
			// Em uma implementação mais robusta, isso seria feito no repository
			if task.DueDate != nil {
				// Lógica para verificar se está em atraso
				// overdueTasks++
			}
		}
	}

	summary.CompletedTasks = completedTasks
	summary.PendingTasks = pendingTasks
	summary.OverdueTasks = overdueTasks

	// Calcular progresso
	if summary.TotalTasks > 0 {
		summary.TasksProgress = float64(completedTasks) / float64(summary.TotalTasks) * 100
	}

	return summary, nil
}
