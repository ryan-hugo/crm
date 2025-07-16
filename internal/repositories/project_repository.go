package repositories

import (
	"crm-backend/internal/models"

	"gorm.io/gorm"
)

// ProjectRepository define a interface para operações de projeto no banco de dados
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

// projectRepository implementa ProjectRepository
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository cria uma nova instância do repositório de projetos
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// Create cria um novo projeto no banco de dados
func (r *projectRepository) Create(project *models.Project) error {
	if err := r.db.Create(project).Error; err != nil {
		return err
	}
	return nil
}

// GetByID busca um projeto pelo ID
func (r *projectRepository) GetByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.Preload("Client").Preload("User").First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// GetWithTasks busca um projeto pelo ID incluindo suas tarefas
func (r *projectRepository) GetWithTasks(id uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.Preload("Client").
		Preload("User").
		Preload("Tasks").
		First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

// GetByUserID busca projetos por ID do usuário com filtros
func (r *projectRepository) GetByUserID(userID uint, filter *models.ProjectListFilter) ([]models.Project, error) {
	var projects []models.Project
	query := r.db.Where("user_id = ?", userID)

	// Aplicar filtros
	if filter != nil {
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
		if filter.ClientID != nil {
			query = query.Where("client_id = ?", *filter.ClientID)
		}

		// Paginação
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	// Ordenar por data de criação (mais recente primeiro)
	query = query.Order("created_at DESC")

	if err := query.Preload("Client").Preload("User").Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

// GetByClientID busca projetos por ID do cliente
func (r *projectRepository) GetByClientID(clientID uint) ([]models.Project, error) {
	var projects []models.Project
	if err := r.db.Where("client_id = ?", clientID).
		Preload("Client").
		Preload("User").
		Order("created_at DESC").
		Find(&projects).Error; err != nil {
		return nil, err
	}
	return projects, nil
}

// Update atualiza um projeto existente
func (r *projectRepository) Update(project *models.Project) error {
	if err := r.db.Save(project).Error; err != nil {
		return err
	}
	return nil
}

// Delete remove um projeto do banco de dados (soft delete)
func (r *projectRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Project{}, id).Error; err != nil {
		return err
	}
	return nil
}

// CountByUserID conta o número total de projetos de um usuário
func (r *projectRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Project{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByStatus conta o número de projetos por status de um usuário
func (r *projectRepository) CountByStatus(userID uint, status models.ProjectStatus) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Project{}).
		Where("user_id = ? AND status = ?", userID, status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

