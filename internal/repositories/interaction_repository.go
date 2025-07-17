package repositories

import (
	"crm-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// InteractionRepository define a interface para operações de interação no banco de dados
type InteractionRepository interface {
	Create(interaction *models.Interaction) error
	GetByID(id uint) (*models.Interaction, error)
	GetByContactID(contactID uint, filter *models.InteractionListFilter) ([]models.Interaction, error)
	Update(interaction *models.Interaction) error
	Delete(id uint) error
	GetByUserID(userID uint, filter *models.InteractionListFilter) ([]models.Interaction, error)
	CountByContactID(contactID uint) (int64, error)
	GetRecentByUserID(userID uint, days int, limit int) ([]models.Interaction, error)
}

// interactionRepository implementa InteractionRepository
type interactionRepository struct {
	db *gorm.DB
}

// NewInteractionRepository cria uma nova instância do repositório de interações
func NewInteractionRepository(db *gorm.DB) InteractionRepository {
	return &interactionRepository{db: db}
}

// Create cria uma nova interação no banco de dados
func (r *interactionRepository) Create(interaction *models.Interaction) error {
	if err := r.db.Create(interaction).Error; err != nil {
		return err
	}
	return nil
}

// GetByID busca uma interação pelo ID
func (r *interactionRepository) GetByID(id uint) (*models.Interaction, error) {
	var interaction models.Interaction
	if err := r.db.Preload("Contact").First(&interaction, id).Error; err != nil {
		return nil, err
	}
	return &interaction, nil
}

// GetByContactID busca interações por ID do contato com filtros
func (r *interactionRepository) GetByContactID(contactID uint, filter *models.InteractionListFilter) ([]models.Interaction, error) {
	var interactions []models.Interaction
	query := r.db.Where("contact_id = ?", contactID)

	// Aplicar filtros
	if filter != nil {
		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}
		if filter.DateFrom != nil {
			query = query.Where("date >= ?", filter.DateFrom)
		}
		if filter.DateTo != nil {
			query = query.Where("date <= ?", filter.DateTo)
		}

		// Paginação
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	// Ordenar por data (mais recente primeiro)
	query = query.Order("date DESC")

	if err := query.Preload("Contact").Find(&interactions).Error; err != nil {
		return nil, err
	}

	return interactions, nil
}

// GetByUserID busca interações por ID do usuário (através dos contatos)
func (r *interactionRepository) GetByUserID(userID uint, filter *models.InteractionListFilter) ([]models.Interaction, error) {
	var interactions []models.Interaction
	query := r.db.Joins("JOIN contacts ON interactions.contact_id = contacts.id").
		Where("contacts.user_id = ?", userID)

	// Aplicar filtros
	if filter != nil {
		if filter.Type != "" {
			query = query.Where("interactions.type = ?", filter.Type)
		}
		if filter.DateFrom != nil {
			query = query.Where("interactions.date >= ?", filter.DateFrom)
		}
		if filter.DateTo != nil {
			query = query.Where("interactions.date <= ?", filter.DateTo)
		}
		if filter.ContactID > 0 {
			query = query.Where("interactions.contact_id = ?", filter.ContactID)
		}

		// Paginação
		if filter.Limit > 0 {
			query = query.Limit(filter.Limit)
		}
		if filter.Offset > 0 {
			query = query.Offset(filter.Offset)
		}
	}

	// Ordenar por data (mais recente primeiro)
	query = query.Order("interactions.date DESC")

	if err := query.Preload("Contact").Find(&interactions).Error; err != nil {
		return nil, err
	}

	return interactions, nil
}

// Update atualiza uma interação existente
func (r *interactionRepository) Update(interaction *models.Interaction) error {
	if err := r.db.Save(interaction).Error; err != nil {
		return err
	}
	return nil
}

// Delete remove uma interação do banco de dados (soft delete)
func (r *interactionRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Interaction{}, id).Error; err != nil {
		return err
	}
	return nil
}

// CountByContactID conta o número de interações de um contato
func (r *interactionRepository) CountByContactID(contactID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Interaction{}).Where("contact_id = ?", contactID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}



// GetRecentByUserID busca interações recentes do usuário nos últimos X dias
func (r *interactionRepository) GetRecentByUserID(userID uint, days int, limit int) ([]models.Interaction, error) {
	var interactions []models.Interaction

	// Calcular data de início (X dias atrás)
	startDate := time.Now().AddDate(0, 0, -days)

	query := r.db.Joins("JOIN contacts ON interactions.contact_id = contacts.id").
		Where("contacts.user_id = ? AND interactions.date >= ?", userID, startDate).
		Order("interactions.date DESC").
		Preload("Contact")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&interactions).Error; err != nil {
		return nil, err
	}

	return interactions, nil
}
