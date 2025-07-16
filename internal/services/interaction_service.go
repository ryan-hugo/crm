package services

import (
	"crm-backend/internal/models"
	"crm-backend/internal/repositories"
	"crm-backend/pkg/errors"

	"gorm.io/gorm"
)

// InteractionService define a interface para operações de interação
type InteractionService interface {
	Create(userID, contactID uint, req *models.InteractionCreateRequest) (*models.Interaction, error)
	GetByID(userID, interactionID uint) (*models.Interaction, error)
	GetByContactID(userID, contactID uint, filter *models.InteractionListFilter) ([]models.Interaction, error)
	GetByUserID(userID uint, filter *models.InteractionListFilter) ([]models.Interaction, error)
	Update(userID, interactionID uint, req *models.InteractionUpdateRequest) (*models.Interaction, error)
	Delete(userID, interactionID uint) error
	GetRecentInteractions(userID uint, limit int) ([]models.Interaction, error)
}

// interactionService implementa InteractionService
type interactionService struct {
	interactionRepo repositories.InteractionRepository
	contactRepo     repositories.ContactRepository
}

// NewInteractionService cria uma nova instância do serviço de interações
func NewInteractionService(
	interactionRepo repositories.InteractionRepository,
	contactRepo repositories.ContactRepository,
) InteractionService {
	return &interactionService{
		interactionRepo: interactionRepo,
		contactRepo:     contactRepo,
	}
}

// Create cria uma nova interação
func (s *interactionService) Create(userID, contactID uint, req *models.InteractionCreateRequest) (*models.Interaction, error) {
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

	// Criar interação
	interaction := &models.Interaction{
		Type:        req.Type,
		Date:        req.Date,
		Subject:     req.Subject,
		Description: req.Description,
		ContactID:   contactID,
	}

	if err := s.interactionRepo.Create(interaction); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar interação criada com relacionamentos
	createdInteraction, err := s.interactionRepo.GetByID(interaction.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return createdInteraction, nil
}

// GetByID obtém uma interação específica
func (s *interactionService) GetByID(userID, interactionID uint) (*models.Interaction, error) {
	interaction, err := s.interactionRepo.GetByID(interactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Interação")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se a interação pertence a um contato do usuário
	if interaction.Contact.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return interaction, nil
}

// GetByContactID obtém interações de um contato específico
func (s *interactionService) GetByContactID(userID, contactID uint, filter *models.InteractionListFilter) ([]models.Interaction, error) {
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

	// Aplicar valores padrão ao filtro se necessário
	if filter == nil {
		filter = &models.InteractionListFilter{}
	}
	if filter.Limit == 0 {
		filter.Limit = 50 // Limite padrão
	}

	interactions, err := s.interactionRepo.GetByContactID(contactID, filter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return interactions, nil
}

// GetByUserID obtém todas as interações do usuário
func (s *interactionService) GetByUserID(userID uint, filter *models.InteractionListFilter) ([]models.Interaction, error) {
	// Aplicar valores padrão ao filtro se necessário
	if filter == nil {
		filter = &models.InteractionListFilter{}
	}
	if filter.Limit == 0 {
		filter.Limit = 50 // Limite padrão
	}

	interactions, err := s.interactionRepo.GetByUserID(userID, filter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return interactions, nil
}

// Update atualiza uma interação existente
func (s *interactionService) Update(userID, interactionID uint, req *models.InteractionUpdateRequest) (*models.Interaction, error) {
	// Buscar interação existente
	interaction, err := s.interactionRepo.GetByID(interactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Interação")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se a interação pertence a um contato do usuário
	if interaction.Contact.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Atualizar campos fornecidos
	if req.Type != "" {
		interaction.Type = req.Type
	}
	if req.Date != nil {
		interaction.Date = *req.Date
	}
	if req.Subject != "" {
		interaction.Subject = req.Subject
	}
	if req.Description != "" {
		interaction.Description = req.Description
	}

	// Salvar alterações
	if err := s.interactionRepo.Update(interaction); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar interação atualizada com relacionamentos
	updatedInteraction, err := s.interactionRepo.GetByID(interaction.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return updatedInteraction, nil
}

// Delete exclui uma interação
func (s *interactionService) Delete(userID, interactionID uint) error {
	// Buscar interação existente
	interaction, err := s.interactionRepo.GetByID(interactionID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("Interação")
		}
		return errors.ErrInternalServer
	}

	// Verificar se a interação pertence a um contato do usuário
	if interaction.Contact.UserID != userID {
		return errors.ErrForbidden
	}

	// Excluir interação
	if err := s.interactionRepo.Delete(interactionID); err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

// GetRecentInteractions obtém as interações mais recentes do usuário
func (s *interactionService) GetRecentInteractions(userID uint, limit int) ([]models.Interaction, error) {
	if limit <= 0 {
		limit = 10 // Limite padrão
	}

	filter := &models.InteractionListFilter{
		Limit: limit,
	}

	interactions, err := s.interactionRepo.GetByUserID(userID, filter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return interactions, nil
}

