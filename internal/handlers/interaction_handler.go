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

// InteractionHandler gerencia as rotas de interações
type InteractionHandler struct {
	interactionService services.InteractionService
}

// NewInteractionHandler cria uma nova instância do handler de interações
func NewInteractionHandler(interactionService services.InteractionService) *InteractionHandler {
	return &InteractionHandler{
		interactionService: interactionService,
	}
}

// Create cria uma nova interação para um contato
// @Summary Criar nova interação
// @Description Cria uma nova interação para um contato específico
// @Tags interactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param contactId path int true "ID do contato"
// @Param request body models.InteractionCreateRequest true "Dados da interação"
// @Success 201 {object} models.Interaction
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{contactId}/interactions [post]
func (h *InteractionHandler) Create(c *gin.Context) {
	start := time.Now()
	userID := c.GetUint("user_id")
	var req models.InteractionCreateRequest

	// Obter ID do contato da URL (parâmetro :id)
	contactIDStr := c.Param("id")
	logger.Debugf("Criando interação para contato ID: %s (usuário: %d)", contactIDStr, userID)

	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		logger.LogError(err, "Erro ao converter ID do contato", map[string]interface{}{
			"contact_id_str": contactIDStr,
			"user_id":        userID,
		})
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.LogError(err, "Erro ao validar dados de entrada", map[string]interface{}{
			"contact_id": contactID,
			"user_id":    userID,
		})
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para criar interação
	interaction, err := h.interactionService.Create(userID, uint(contactID), &req)
	if err != nil {
		logger.LogError(err, "Erro ao criar interação", map[string]interface{}{
			"contact_id": contactID,
			"user_id":    userID,
			"request":    req,
		})
		c.Error(err)
		return
	}

	duration := time.Since(start)
	logger.WithFields("INFO", "Interaction Created", map[string]interface{}{
		"user_id":        userID,
		"contact_id":     contactID,
		"interaction_id": interaction.ID,
		"duration":       duration,
	})

	c.JSON(http.StatusCreated, interaction)
}

// ListByContact lista interações de um contato específico
// @Summary Listar interações de um contato
// @Description Lista todas as interações de um contato específico
// @Tags interactions
// @Security BearerAuth
// @Produce json
// @Param contactId path int true "ID do contato"
// @Param type query string false "Tipo de interação (EMAIL, CALL, MEETING, OTHER)"
// @Param date_from query string false "Data inicial (formato: 2006-01-02T15:04:05Z)"
// @Param date_to query string false "Data final (formato: 2006-01-02T15:04:05Z)"
// @Param limit query int false "Limite de resultados (padrão: 50)"
// @Param offset query int false "Offset para paginação (padrão: 0)"
// @Success 200 {array} models.Interaction
// @Failure 400 {object} map[string]interface{} "Parâmetros inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{contactId}/interactions [get]
func (h *InteractionHandler) ListByContact(c *gin.Context) {
	start := time.Now()
	userID := c.GetUint("user_id")
	var filter models.InteractionListFilter

	// Obter ID do contato da URL (parâmetro :id)
	contactIDStr := c.Param("id")
	logger.Debugf("Listando interações para contato ID: %s (usuário: %d)", contactIDStr, userID)

	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		logger.LogError(err, "Erro ao converter ID do contato", map[string]interface{}{
			"contact_id_str": contactIDStr,
			"user_id":        userID,
		})
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		logger.LogError(err, "Erro ao validar parâmetros de consulta", map[string]interface{}{
			"contact_id": contactID,
			"user_id":    userID,
		})
		c.Error(errors.NewBadRequestError("Parâmetros de consulta inválidos: " + err.Error()))
		return
	}

	// Chamar service para listar interações do contato
	interactions, err := h.interactionService.GetByContactID(userID, uint(contactID), &filter)
	if err != nil {
		logger.LogError(err, "Erro ao listar interações", map[string]interface{}{
			"contact_id": contactID,
			"user_id":    userID,
			"filter":     filter,
		})
		c.Error(err)
		return
	}

	duration := time.Since(start)
	logger.WithFields("INFO", "Interactions Listed", map[string]interface{}{
		"user_id":      userID,
		"contact_id":   contactID,
		"interactions": len(interactions),
		"duration":     duration,
	})

	c.JSON(http.StatusOK, interactions)
}

// List lista todas as interações do usuário
// @Summary Listar todas as interações
// @Description Lista todas as interações do usuário com filtros opcionais
// @Tags interactions
// @Security BearerAuth
// @Produce json
// @Param type query string false "Tipo de interação (EMAIL, CALL, MEETING, OTHER)"
// @Param contact_id query int false "ID do contato específico"
// @Param date_from query string false "Data inicial (formato: 2006-01-02T15:04:05Z)"
// @Param date_to query string false "Data final (formato: 2006-01-02T15:04:05Z)"
// @Param limit query int false "Limite de resultados (padrão: 50)"
// @Param offset query int false "Offset para paginação (padrão: 0)"
// @Success 200 {array} models.Interaction
// @Failure 400 {object} map[string]interface{} "Parâmetros inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/interactions [get]
func (h *InteractionHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var filter models.InteractionListFilter

	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(errors.NewBadRequestError("Parâmetros de consulta inválidos: " + err.Error()))
		return
	}

	// Chamar service para listar interações do usuário
	interactions, err := h.interactionService.GetByUserID(userID, &filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, interactions)
}

// GetByID obtém uma interação específica
// @Summary Obter interação por ID
// @Description Obtém os detalhes de uma interação específica
// @Tags interactions
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID da interação"
// @Success 200 {object} models.Interaction
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Interação não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/interactions/{id} [get]
func (h *InteractionHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID da interação da URL
	interactionIDStr := c.Param("id")
	interactionID, err := strconv.ParseUint(interactionIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID da interação inválido"))
		return
	}

	// Chamar service para obter interação
	interaction, err := h.interactionService.GetByID(userID, uint(interactionID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, interaction)
}

// Update atualiza uma interação existente
// @Summary Atualizar interação
// @Description Atualiza os dados de uma interação existente
// @Tags interactions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID da interação"
// @Param request body models.InteractionUpdateRequest true "Dados para atualização"
// @Success 200 {object} models.Interaction
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Interação não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/interactions/{id} [put]
func (h *InteractionHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.InteractionUpdateRequest

	// Obter ID da interação da URL
	interactionIDStr := c.Param("id")
	interactionID, err := strconv.ParseUint(interactionIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID da interação inválido"))
		return
	}

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para atualizar interação
	updatedInteraction, err := h.interactionService.Update(userID, uint(interactionID), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedInteraction)
}

// Delete exclui uma interação
// @Summary Excluir interação
// @Description Exclui uma interação específica
// @Tags interactions
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID da interação"
// @Success 204 "Interação excluída com sucesso"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Interação não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/interactions/{id} [delete]
func (h *InteractionHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID da interação da URL
	interactionIDStr := c.Param("id")
	interactionID, err := strconv.ParseUint(interactionIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID da interação inválido"))
		return
	}

	// Chamar service para excluir interação
	err = h.interactionService.Delete(userID, uint(interactionID))
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetRecent obtém interações recentes do usuário
// @Summary Obter interações recentes
// @Description Obtém as interações mais recentes do usuário
// @Tags interactions
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limite de resultados (padrão: 10)"
// @Success 200 {array} models.Interaction
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/interactions/recent [get]
func (h *InteractionHandler) GetRecent(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter limite da query string
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Chamar service para obter interações recentes
	interactions, err := h.interactionService.GetRecentInteractions(userID, limit)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, interactions)
}

// GetRecentInteractions obtém interações recentes dos últimos 7 dias
// @Summary Obter interações recentes
// @Description Obtém interações recentes do usuário dos últimos 7 dias
// @Tags interactions
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limite de resultados (padrão: 10)"
// @Success 200 {object} map[string]interface{} "Lista de interações recentes"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/interactions/recent [get]
// GetRecentInteractionsCount retorna apenas o número de interações recentes dos últimos 7 dias
// @Summary Contar interações recentes
// @Description Retorna o número de interações recentes do usuário dos últimos 7 dias
// @Tags interactions
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limite de resultados (padrão: 10)"
// @Success 200 {object} map[string]int "Quantidade de interações recentes"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/interactions/recent/count [get]
func (h *InteractionHandler) GetRecentInteractionsCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter limite da query string (padrão: 10)
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 50 {
		limit = 10
	}

	// Chamar service para obter interações recentes
	interactions, err := h.interactionService.GetRecentInteractions(userID, limit)
	if err != nil {
		logger.LogError(err, "Erro ao buscar interações recentes", map[string]interface{}{
			"user_id": userID,
			"limit":   limit,
		})
		c.Error(err)
		return
	}

	// Retornar apenas o count numérico
	c.JSON(http.StatusOK, map[string]int{
		"count": len(interactions),
	})
}
