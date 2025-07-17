package services

import (
	"crm-backend/internal/models"
	"crm-backend/internal/repositories"
	"crm-backend/pkg/errors"

	"gorm.io/gorm"
)

// ContactService define a interface para operações de contato
type ContactService interface {
	Create(userID uint, req *models.ContactCreateRequest) (*models.Contact, error)
	GetByID(userID, contactID uint) (*models.Contact, error)
	GetWithDetails(userID, contactID uint) (*ContactDetails, error)
	GetByUserID(userID uint, filter *models.ContactListFilter) ([]models.Contact, error)
	Update(userID, contactID uint, req *models.ContactUpdateRequest) (*models.Contact, error)
	Delete(userID, contactID uint) error
	SearchByName(userID uint, name string) ([]models.Contact, error)
	GetContactSummary(userID, contactID uint) (*ContactSummary, error)
	ConvertLeadToClient(userID, contactID uint) (*models.Contact, error)
}

// ContactDetails representa detalhes completos de um contato
type ContactDetails struct {
	Contact      *models.Contact      `json:"contact"`
	Interactions []models.Interaction `json:"interactions"`
	Tasks        []models.Task        `json:"tasks"`
	Projects     []models.Project     `json:"projects"`
}

// ContactSummary representa um resumo do contato
type ContactSummary struct {
	Contact             *models.Contact `json:"contact"`
	TotalInteractions   int64           `json:"total_interactions"`
	TotalTasks          int64           `json:"total_tasks"`
	CompletedTasks      int64           `json:"completed_tasks"`
	PendingTasks        int64           `json:"pending_tasks"`
	TotalProjects       int64           `json:"total_projects"`
	ActiveProjects      int64           `json:"active_projects"`
	CompletedProjects   int64           `json:"completed_projects"`
	LastInteractionDate *string         `json:"last_interaction_date"`
}

// contactService implementa ContactService
type contactService struct {
	contactRepo     repositories.ContactRepository
	interactionRepo repositories.InteractionRepository
	taskRepo        repositories.TaskRepository
	projectRepo     repositories.ProjectRepository
}

// NewContactService cria uma nova instância do serviço de contatos
func NewContactService(
	contactRepo repositories.ContactRepository,
	interactionRepo repositories.InteractionRepository,
	taskRepo repositories.TaskRepository,
	projectRepo repositories.ProjectRepository,
) ContactService {
	return &contactService{
		contactRepo:     contactRepo,
		interactionRepo: interactionRepo,
		taskRepo:        taskRepo,
		projectRepo:     projectRepo,
	}
}

// Create cria um novo contato
func (s *contactService) Create(userID uint, req *models.ContactCreateRequest) (*models.Contact, error) {
	// Verificar se já existe um contato com o mesmo email para este usuário
	existingContact, err := s.contactRepo.GetByEmail(req.Email)
	if err == nil && existingContact.UserID == userID {
		return nil, errors.NewConflictError("Já existe um contato com este email")
	}

	// Criar contato
	contact := &models.Contact{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Company:  req.Company,
		Position: req.Position,
		Type:     req.Type,
		Notes:    req.Notes,
		UserID:   userID,
	}

	if err := s.contactRepo.Create(contact); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar contato criado com relacionamentos
	createdContact, err := s.contactRepo.GetByID(contact.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return createdContact, nil
}

// GetByID obtém um contato específico
func (s *contactService) GetByID(userID, contactID uint) (*models.Contact, error) {
	contact, err := s.contactRepo.GetByID(contactID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Contato")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se o contato pertence ao usuário
	if contact.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return contact, nil
}

// GetWithDetails obtém um contato com todos os detalhes relacionados
func (s *contactService) GetWithDetails(userID, contactID uint) (*ContactDetails, error) {
	// Verificar se o contato pertence ao usuário
	contact, err := s.GetByID(userID, contactID)
	if err != nil {
		return nil, err
	}

	details := &ContactDetails{
		Contact: contact,
	}

	// Buscar interações
	if s.interactionRepo != nil {
		interactions, err := s.interactionRepo.GetByContactID(contactID, &models.InteractionListFilter{
			Limit: 50, // Últimas 50 interações
		})
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		details.Interactions = interactions
	}

	// Buscar tarefas
	if s.taskRepo != nil {
		tasks, err := s.taskRepo.GetByContactID(contactID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		details.Tasks = tasks
	}

	// Buscar projetos
	if s.projectRepo != nil {
		projects, err := s.projectRepo.GetByClientID(contactID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		details.Projects = projects
	}

	return details, nil
}

// GetByUserID obtém todos os contatos do usuário
func (s *contactService) GetByUserID(userID uint, filter *models.ContactListFilter) ([]models.Contact, error) {
	// Aplicar valores padrão ao filtro se necessário
	if filter == nil {
		filter = &models.ContactListFilter{}
	}
	if filter.Limit == 0 {
		filter.Limit = 50 // Limite padrão
	}

	contacts, err := s.contactRepo.GetByUserID(userID, filter)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return contacts, nil
}

// Update atualiza um contato existente
func (s *contactService) Update(userID, contactID uint, req *models.ContactUpdateRequest) (*models.Contact, error) {
	// Buscar contato existente
	contact, err := s.contactRepo.GetByID(contactID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Contato")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se o contato pertence ao usuário
	if contact.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Verificar se o email está sendo alterado e se já existe
	if req.Email != "" && req.Email != contact.Email {
		existingContact, err := s.contactRepo.GetByEmail(req.Email)
		if err == nil && existingContact.UserID == userID && existingContact.ID != contactID {
			return nil, errors.NewConflictError("Já existe um contato com este email")
		}
	}

	// Atualizar campos fornecidos
	if req.Name != "" {
		contact.Name = req.Name
	}
	if req.Email != "" {
		contact.Email = req.Email
	}
	if req.Phone != "" {
		contact.Phone = req.Phone
	}
	if req.Company != "" {
		contact.Company = req.Company
	}
	if req.Position != "" {
		contact.Position = req.Position
	}
	if req.Type != "" {
		contact.Type = req.Type
	}
	if req.Notes != "" {
		contact.Notes = req.Notes
	}

	// Salvar alterações
	if err := s.contactRepo.Update(contact); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar contato atualizado com relacionamentos
	updatedContact, err := s.contactRepo.GetByID(contact.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return updatedContact, nil
}

// Delete exclui um contato
func (s *contactService) Delete(userID, contactID uint) error {
	// Buscar contato existente
	contact, err := s.contactRepo.GetByID(contactID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("Contato")
		}
		return errors.ErrInternalServer
	}

	// Verificar se o contato pertence ao usuário
	if contact.UserID != userID {
		return errors.ErrForbidden
	}

	// Verificar se há projetos associados (apenas para clientes)
	if contact.Type == models.ContactTypeClient && s.projectRepo != nil {
		projects, err := s.projectRepo.GetByClientID(contactID)
		if err != nil {
			return errors.ErrInternalServer
		}
		if len(projects) > 0 {
			return errors.NewBadRequestError("Não é possível excluir cliente com projetos associados. Exclua os projetos primeiro.")
		}
	}

	// Excluir contato (soft delete - GORM cuidará das relações)
	if err := s.contactRepo.Delete(contactID); err != nil {
		return errors.ErrInternalServer
	}

	return nil
}

// SearchByName busca contatos por nome
func (s *contactService) SearchByName(userID uint, name string) ([]models.Contact, error) {
	if name == "" {
		return []models.Contact{}, nil
	}

	contacts, err := s.contactRepo.SearchByName(userID, name)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return contacts, nil
}

// GetContactSummary obtém um resumo detalhado do contato
func (s *contactService) GetContactSummary(userID, contactID uint) (*ContactSummary, error) {
	// Buscar contato
	contact, err := s.GetByID(userID, contactID)
	if err != nil {
		return nil, err
	}

	summary := &ContactSummary{
		Contact: contact,
	}

	// Estatísticas de interações
	if s.interactionRepo != nil {
		interactionCount, err := s.interactionRepo.CountByContactID(contactID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}
		summary.TotalInteractions = interactionCount

		// Buscar última interação para obter a data
		interactions, err := s.interactionRepo.GetByContactID(contactID, &models.InteractionListFilter{
			Limit: 1,
		})
		if err == nil && len(interactions) > 0 {
			lastDate := interactions[0].Date.Format("2006-01-02 15:04:05")
			summary.LastInteractionDate = &lastDate
		}
	}

	// Estatísticas de tarefas
	if s.taskRepo != nil {
		tasks, err := s.taskRepo.GetByContactID(contactID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}

		summary.TotalTasks = int64(len(tasks))
		for _, task := range tasks {
			if task.Status == models.TaskStatusCompleted {
				summary.CompletedTasks++
			} else {
				summary.PendingTasks++
			}
		}
	}

	// Estatísticas de projetos (apenas para clientes)
	if contact.Type == models.ContactTypeClient && s.projectRepo != nil {
		projects, err := s.projectRepo.GetByClientID(contactID)
		if err != nil {
			return nil, errors.ErrInternalServer
		}

		summary.TotalProjects = int64(len(projects))
		for _, project := range projects {
			switch project.Status {
			case models.ProjectStatusInProgress:
				summary.ActiveProjects++
			case models.ProjectStatusCompleted:
				summary.CompletedProjects++
			}
		}
	}

	return summary, nil
}

// ConvertLeadToClient converte um lead em cliente
func (s *contactService) ConvertLeadToClient(userID, contactID uint) (*models.Contact, error) {
	// Buscar contato existente
	contact, err := s.contactRepo.GetByID(contactID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("Contato")
		}
		return nil, errors.ErrInternalServer
	}

	// Verificar se o contato pertence ao usuário
	if contact.UserID != userID {
		return nil, errors.ErrForbidden
	}

	// Verificar se é um lead
	if contact.Type != models.ContactTypeLead {
		return nil, errors.NewBadRequestError("Apenas leads podem ser convertidos em clientes")
	}

	// Converter para cliente
	contact.Type = models.ContactTypeClient

	// Salvar alterações
	if err := s.contactRepo.Update(contact); err != nil {
		return nil, errors.ErrInternalServer
	}

	// Buscar contato atualizado
	updatedContact, err := s.contactRepo.GetByID(contact.ID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return updatedContact, nil
}
