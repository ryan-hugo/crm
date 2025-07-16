package services

import (
	"crm-backend/internal/models"
	"crm-backend/internal/repositories"
	"crm-backend/pkg/errors"

	"gorm.io/gorm"
)

// TaskService define a interface para operações de tarefa
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

// taskService implementa TaskService
type taskService struct {
	taskRepo    repositories.TaskRepository
	contactRepo repositories.ContactRepository
	projectRepo repositories.ProjectRepository
}

// NewTaskService cria uma nova instância do serviço de tarefas
func NewTaskService(
	taskRepo repositories.TaskRepository,
	contactRepo repositories.ContactRepository,
	projectRepo repositories.ProjectRepository,
) TaskService {
	return &taskService{
		taskRepo:    taskRepo,
		contactRepo: contactRepo,
		projectRepo: projectRepo,
	}
}

// Create cria uma nova tarefa
func (s *taskService) Create(userID uint, req *models.TaskCreateRequest) (*models.Task, error) {
	// Validar associações se fornecidas
	if req.ContactID != nil {
		contact, err := s.contactRepo.GetByID(*req.ContactID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Contato")
			}
			return nil, errors.ErrInternalServer
		}
		if contact.UserID != userID {
			return nil, errors.ErrForbidden
		}
	}

	if req.ProjectID != nil {
		project, err := s.projectRepo.GetByID(*req.ProjectID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Projeto")
			}
			return nil, errors.ErrInternalServer
		}
		if project.UserID != userID {
			return nil, errors.ErrForbidden
		}
	}

	// Criar tarefa
	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		Priority:    req.Priority,
		Status:      req.Status,
		UserID:      userID,
		ContactID:   req.ContactID,
		ProjectID:   req.ProjectID,
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar tarefa criada com relacionamentos
	createdTask, err := s.taskRepo.GetByID(task.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return createdTask, nil
}

// GetByID obtém uma tarefa específica
func (s *taskService) GetByID(userID, taskID uint) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Tarefa")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se a tarefa pertence ao usuário
	if task.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return task, nil
}

// GetByUserID obtém todas as tarefas do usuário
func (s *taskService) GetByUserID(userID uint, filter *models.TaskListFilter) ([]models.Task, error) {
	// Aplicar valores padrão ao filtro se necessário
	if filter == nil {
		filter = &models.TaskListFilter{}
	}
	if filter.Limit == 0 {
		filter.Limit = 50 // Limite padrão
	}

	tasks, err := s.taskRepo.GetByUserID(userID, filter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return tasks, nil
}

// Update atualiza uma tarefa existente
func (s *taskService) Update(userID, taskID uint, req *models.TaskUpdateRequest) (*models.Task, error) {
	// Buscar tarefa existente
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Tarefa")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se a tarefa pertence ao usuário
	if task.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Validar novas associações se fornecidas
	if req.ContactID != nil {
		contact, err := s.contactRepo.GetByID(*req.ContactID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Contato")
			}
			return nil, errors.ErrInternalServer
		}
		if contact.UserID != userID {
			return nil, errors.ErrForbidden
		}
		task.ContactID = req.ContactID
	}

	if req.ProjectID != nil {
		project, err := s.projectRepo.GetByID(*req.ProjectID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.NewNotFoundError("Projeto")
			}
			return nil, errors.ErrInternalServer
		}
		if project.UserID != userID {
			return nil, errors.ErrForbidden
		}
		task.ProjectID = req.ProjectID
	}

	// Atualizar campos fornecidos
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if req.Status != "" {
		task.Status = req.Status
	}

	// Salvar alterações
	if err := s.taskRepo.Update(task); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar tarefa atualizada com relacionamentos
	updatedTask, err := s.taskRepo.GetByID(task.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return updatedTask, nil
}

// Delete exclui uma tarefa
func (s *taskService) Delete(userID, taskID uint) error {
	// Buscar tarefa existente
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("Tarefa")
		}
		return errors.ErrInternalServer
	}

	// Verificar se a tarefa pertence ao usuário
	if task.UserID != userID {
		return errors.ErrForbidden
	}

	// Excluir tarefa
	if err := s.taskRepo.Delete(taskID); err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

// MarkAsCompleted marca uma tarefa como concluída
func (s *taskService) MarkAsCompleted(userID, taskID uint) (*models.Task, error) {
	req := &models.TaskUpdateRequest{
		Status: models.TaskStatusCompleted,
	}
	return s.Update(userID, taskID, req)
}

// MarkAsPending marca uma tarefa como pendente
func (s *taskService) MarkAsPending(userID, taskID uint) (*models.Task, error) {
	req := &models.TaskUpdateRequest{
		Status: models.TaskStatusPending,
	}
	return s.Update(userID, taskID, req)
}

// GetByContactID obtém tarefas de um contato específico
func (s *taskService) GetByContactID(userID, contactID uint) ([]models.Task, error) {
	// Verificar se o contato existe e pertence ao usuário
	contact, err := s.contactRepo.GetByID(contactID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Contato")
		}
		return nil, errors.ErrInternalServer
	}

	if contact.UserID != userID {
		return nil, errors.ErrForbidden
	}

	tasks, err := s.taskRepo.GetByContactID(contactID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return tasks, nil
}

// GetByProjectID obtém tarefas de um projeto específico
func (s *taskService) GetByProjectID(userID, projectID uint) ([]models.Task, error) {
	// Verificar se o projeto existe e pertence ao usuário
	project, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Projeto")
		}
		return nil, errors.ErrInternalServer
	}

	if project.UserID != userID {
		return nil, errors.ErrForbidden
	}

	tasks, err := s.taskRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return tasks, nil
}

// GetOverdueTasks obtém tarefas em atraso do usuário
func (s *taskService) GetOverdueTasks(userID uint) ([]models.Task, error) {
	tasks, err := s.taskRepo.GetOverdueTasks(userID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return tasks, nil
}

// GetUpcomingTasks obtém tarefas próximas do vencimento
func (s *taskService) GetUpcomingTasks(userID uint, days int) ([]models.Task, error) {
	if days <= 0 {
		days = 7 // Padrão: próximos 7 dias
	}

	// Usar filtro para buscar tarefas com vencimento nos próximos dias
	// Implementação simplificada - pode ser melhorada no repository
	filter := &models.TaskListFilter{
		Status: models.TaskStatusPending,
		Limit:  100, // Limite alto para capturar todas as tarefas relevantes
	}

	tasks, err := s.taskRepo.GetByUserID(userID, filter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// Filtrar tarefas com vencimento nos próximos dias (implementação básica)
	// Em uma implementação mais robusta, isso seria feito no repository
	var upcomingTasks []models.Task
	for _, task := range tasks {
		if task.DueDate != nil {
			// Lógica de filtro por data seria implementada aqui
			upcomingTasks = append(upcomingTasks, task)
		}
	}

	return upcomingTasks, nil
}

