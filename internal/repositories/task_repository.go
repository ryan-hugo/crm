package repositories

import (
	"crm-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// TaskRepository define a interface para operações de tarefa no banco de dados
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

// taskRepository implementa TaskRepository
type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository cria uma nova instância do repositório de tarefas
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// Create cria uma nova tarefa no banco de dados
func (r *taskRepository) Create(task *models.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		return err
	}
	return nil
}

// GetByID busca uma tarefa pelo ID
func (r *taskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.Preload("Contact").Preload("Project").First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// GetByUserID busca tarefas por ID do usuário com filtros
func (r *taskRepository) GetByUserID(userID uint, filter *models.TaskListFilter) ([]models.Task, error) {
	var tasks []models.Task
	query := r.db.Where("user_id = ?", userID)

	// Aplicar filtros
	if filter != nil {
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.Priority != "" {
			query = query.Where("priority = ?", filter.Priority)
		}
		if filter.ContactID != nil {
			query = query.Where("contact_id = ?", *filter.ContactID)
		}
		if filter.ProjectID != nil {
			query = query.Where("project_id = ?", *filter.ProjectID)
		}
		if filter.DueBefore != nil {
			query = query.Where("due_date <= ?", filter.DueBefore)
		}
		if filter.DueAfter != nil {
			query = query.Where("due_date >= ?", filter.DueAfter)
		}

		// Paginação
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	// Ordenar por prioridade e data de vencimento
	query = query.Order("CASE WHEN priority = 'HIGH' THEN 1 WHEN priority = 'MEDIUM' THEN 2 ELSE 3 END, due_date ASC")

	if err := query.Preload("Contact").Preload("Project").Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByContactID busca tarefas por ID do contato
func (r *taskRepository) GetByContactID(contactID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Where("contact_id = ?", contactID).
		Preload("Contact").
		Preload("Project").
		Order("due_date ASC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetByProjectID busca tarefas por ID do projeto
func (r *taskRepository) GetByProjectID(projectID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Where("project_id = ?", projectID).
		Preload("Contact").
		Preload("Project").
		Order("due_date ASC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// Update atualiza uma tarefa existente
func (r *taskRepository) Update(task *models.Task) error {
	if err := r.db.Save(task).Error; err != nil {
		return err
	}
	return nil
}

// Delete remove uma tarefa do banco de dados (soft delete)
func (r *taskRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Task{}, id).Error; err != nil {
		return err
	}
	return nil
}

// CountByUserID conta o número total de tarefas de um usuário
func (r *taskRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Task{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountPendingByUserID conta o número de tarefas pendentes de um usuário
func (r *taskRepository) CountPendingByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Task{}).
		Where("user_id = ? AND status = ?", userID, models.TaskStatusPending).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetOverdueTasks busca tarefas em atraso de um usuário
func (r *taskRepository) GetOverdueTasks(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	now := time.Now()
	
	if err := r.db.Where("user_id = ? AND status = ? AND due_date < ?", 
		userID, models.TaskStatusPending, now).
		Preload("Contact").
		Preload("Project").
		Order("due_date ASC").
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	
	return tasks, nil
}

