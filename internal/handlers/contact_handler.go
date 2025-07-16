package handlers

import (
	"crm-backend/internal/models"
	"crm-backend/internal/services"
	"crm-backend/pkg/errors"
	"crm-backend/pkg/logger"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ContactHandler gerencia as rotas de contatos
type ContactHandler struct {
	contactService services.ContactService
}

// NewContactHandler cria uma nova instância do handler de contatos
func NewContactHandler(contactService services.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

// Create cria um novo contato
// @Summary Criar novo contato
// @Description Cria um novo contato (cliente ou lead)
// @Tags contacts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.ContactCreateRequest true "Dados do contato"
// @Success 201 {object} models.Contact
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 409 {object} map[string]interface{} "Email já existe"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts [post]
func (h *ContactHandler) Create(c *gin.Context) {
	start := time.Now()
	userID := c.GetUint("user_id")
	var req models.ContactCreateRequest

	logger.Debugf("Criando novo contato para usuário %d", userID)

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(errors.NewBadRequestError("Dados de entrada inválidos: "+err.Error()), "Contact Creation", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para criar contato
	contact, err := h.contactService.Create(userID, &req)
	if err != nil {
		logger.LogError(err, "Contact Creation Service", map[string]interface{}{
			"user_id": userID,
			"request": req,
		})
		c.Error(err)
		return
	}

	duration := time.Since(start)
	logger.LogServiceCall("ContactHandler", "Create", duration, true)
	logger.WithFields("INFO", "Contact Created", map[string]interface{}{
		"user_id":    userID,
		"contact_id": contact.ID,
		"email":      contact.Email,
		"duration":   duration,
	})

	c.JSON(http.StatusCreated, contact)
}

// List lista todos os contatos do usuário
// @Summary Listar contatos
// @Description Lista todos os contatos do usuário com filtros opcionais
// @Tags contacts
// @Security BearerAuth
// @Produce json
// @Param type query string false "Tipo de contato (CLIENT ou LEAD)"
// @Param search query string false "Busca por nome, email ou empresa"
// @Param limit query int false "Limite de resultados (padrão: 50)"
// @Param offset query int false "Offset para paginação (padrão: 0)"
// @Success 200 {array} models.Contact
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts [get]
func (h *ContactHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var filter models.ContactListFilter

	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(errors.NewBadRequestError("Parâmetros de consulta inválidos: " + err.Error()))
		return
	}

	// Chamar service para listar contatos
	contacts, err := h.contactService.GetByUserID(userID, &filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// GetByID obtém um contato específico
// @Summary Obter contato por ID
// @Description Obtém os detalhes de um contato específico
// @Tags contacts
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do contato"
// @Success 200 {object} models.Contact
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{id} [get]
func (h *ContactHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do contato da URL
	contactIDStr := c.Param("id")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Chamar service para obter contato
	contact, err := h.contactService.GetByID(userID, uint(contactID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, contact)
}

// GetDetails obtém detalhes completos de um contato
// @Summary Obter detalhes completos do contato
// @Description Obtém um contato com todas as informações relacionadas (interações, tarefas, projetos)
// @Tags contacts
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do contato"
// @Success 200 {object} services.ContactDetails
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{id}/details [get]
func (h *ContactHandler) GetDetails(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do contato da URL
	contactIDStr := c.Param("id")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Chamar service para obter detalhes do contato
	details, err := h.contactService.GetWithDetails(userID, uint(contactID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, details)
}

// Update atualiza um contato existente
// @Summary Atualizar contato
// @Description Atualiza os dados de um contato existente
// @Tags contacts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID do contato"
// @Param request body models.ContactUpdateRequest true "Dados para atualização"
// @Success 200 {object} models.Contact
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 409 {object} map[string]interface{} "Email já existe"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{id} [put]
func (h *ContactHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.ContactUpdateRequest

	// Obter ID do contato da URL
	contactIDStr := c.Param("id")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para atualizar contato
	updatedContact, err := h.contactService.Update(userID, uint(contactID), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedContact)
}

// Delete exclui um contato
// @Summary Excluir contato
// @Description Exclui um contato e todos os dados relacionados
// @Tags contacts
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do contato"
// @Success 204 "Contato excluído com sucesso"
// @Failure 400 {object} map[string]interface{} "ID inválido ou contato tem projetos associados"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{id} [delete]
func (h *ContactHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do contato da URL
	contactIDStr := c.Param("id")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Chamar service para excluir contato
	err = h.contactService.Delete(userID, uint(contactID))
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Search busca contatos por nome
// @Summary Buscar contatos por nome
// @Description Busca contatos do usuário por nome (busca parcial)
// @Tags contacts
// @Security BearerAuth
// @Produce json
// @Param q query string true "Termo de busca (nome)"
// @Success 200 {array} models.Contact
// @Failure 400 {object} map[string]interface{} "Termo de busca obrigatório"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/search [get]
func (h *ContactHandler) Search(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter termo de busca
	searchTerm := c.Query("q")
	if searchTerm == "" {
		c.Error(errors.NewBadRequestError("Termo de busca é obrigatório"))
		return
	}

	// Chamar service para buscar contatos
	contacts, err := h.contactService.SearchByName(userID, searchTerm)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, contacts)
}

// GetSummary obtém resumo de um contato
// @Summary Obter resumo do contato
// @Description Obtém estatísticas e resumo de um contato específico
// @Tags contacts
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do contato"
// @Success 200 {object} services.ContactSummary
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{id}/summary [get]
func (h *ContactHandler) GetSummary(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do contato da URL
	contactIDStr := c.Param("id")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Chamar service para obter resumo do contato
	summary, err := h.contactService.GetContactSummary(userID, uint(contactID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ConvertToClient converte um lead em cliente
// @Summary Converter lead em cliente
// @Description Converte um lead em cliente
// @Tags contacts
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do contato (lead)"
// @Success 200 {object} models.Contact
// @Failure 400 {object} map[string]interface{} "ID inválido ou contato não é lead"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{id}/convert-to-client [put]
func (h *ContactHandler) ConvertToClient(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do contato da URL
	contactIDStr := c.Param("id")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Chamar service para converter lead em cliente
	contact, err := h.contactService.ConvertLeadToClient(userID, uint(contactID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lead convertido em cliente com sucesso",
		"contact": contact,
	})
}
