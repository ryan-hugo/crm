package repositories

import (
	"crm-backend/internal/models"

	"gorm.io/gorm"
)

// ContactRepository define a interface para operações de contato no banco de dados
type ContactRepository interface {
	Create(contact *models.Contact) error
	GetByID(id uint) (*models.Contact, error)
	GetByUserID(userID uint, filter *models.ContactListFilter) ([]models.Contact, error)
	Update(contact *models.Contact) error
	Delete(id uint) error
	GetByEmail(email string) (*models.Contact, error)
	CountByUserID(userID uint) (int64, error)
	CountByType(userID uint, contactType models.ContactType) (int64, error)
	SearchByName(userID uint, name string) ([]models.Contact, error)
	GetWithInteractions(id uint) (*models.Contact, error)
	GetWithTasks(id uint) (*models.Contact, error)
	GetWithProjects(id uint) (*models.Contact, error)
}

// contactRepository implementa ContactRepository
type contactRepository struct {
	db *gorm.DB
}

// NewContactRepository cria uma nova instância do repositório de contatos
func NewContactRepository(db *gorm.DB) ContactRepository {
	return &contactRepository{db: db}
}

// Create cria um novo contato no banco de dados
func (r *contactRepository) Create(contact *models.Contact) error {
	if err := r.db.Create(contact).Error; err != nil {
		return err
	}
	return nil
}

// GetByID busca um contato pelo ID
func (r *contactRepository) GetByID(id uint) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.Preload("User").First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetWithInteractions busca um contato com suas interações
func (r *contactRepository) GetWithInteractions(id uint) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.Preload("User").
		Preload("Interactions", func(db *gorm.DB) *gorm.DB {
			return db.Order("date DESC")
		}).
		First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetWithTasks busca um contato com suas tarefas
func (r *contactRepository) GetWithTasks(id uint) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.Preload("User").
		Preload("Tasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("due_date ASC")
		}).
		First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetWithProjects busca um contato com seus projetos
func (r *contactRepository) GetWithProjects(id uint) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.Preload("User").
		Preload("Projects", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetByUserID busca contatos por ID do usuário com filtros
func (r *contactRepository) GetByUserID(userID uint, filter *models.ContactListFilter) ([]models.Contact, error) {
	var contacts []models.Contact
	query := r.db.Where("user_id = ?", userID)

	// Aplicar filtros
	if filter != nil {
		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}
		if filter.Search != "" {
			searchTerm := "%" + filter.Search + "%"
			query = query.Where("name ILIKE ? OR email ILIKE ? OR company ILIKE ?", 
				searchTerm, searchTerm, searchTerm)
		}

		// Paginação
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	// Ordenar por nome
	query = query.Order("name ASC")

	if err := query.Preload("User").Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}

// GetByEmail busca um contato pelo email
func (r *contactRepository) GetByEmail(email string) (*models.Contact, error) {
	var contact models.Contact
	if err := r.db.Where("email = ?", email).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// Update atualiza um contato existente
func (r *contactRepository) Update(contact *models.Contact) error {
	if err := r.db.Save(contact).Error; err != nil {
		return err
	}
	return nil
}

// Delete remove um contato do banco de dados (soft delete)
func (r *contactRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Contact{}, id).Error; err != nil {
		return err
	}
	return nil
}

// CountByUserID conta o número total de contatos de um usuário
func (r *contactRepository) CountByUserID(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Contact{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByType conta o número de contatos por tipo de um usuário
func (r *contactRepository) CountByType(userID uint, contactType models.ContactType) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Contact{}).
		Where("user_id = ? AND type = ?", userID, contactType).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// SearchByName busca contatos por nome (busca parcial)
func (r *contactRepository) SearchByName(userID uint, name string) ([]models.Contact, error) {
	var contacts []models.Contact
	searchTerm := "%" + name + "%"
	
	if err := r.db.Where("user_id = ? AND name ILIKE ?", userID, searchTerm).
		Order("name ASC").
		Preload("User").
		Find(&contacts).Error; err != nil {
		return nil, err
	}
	
	return contacts, nil
}

