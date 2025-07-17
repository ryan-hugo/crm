package handlers

import (
	"crm-backend/internal/models"
	"crm-backend/internal/services"
	"crm-backend/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProjectHandler gerencia as rotas de projetos
type ProjectHandler struct {
	projectService services.ProjectService
}

// NewProjectHandler cria uma nova instância do handler de projetos
func NewProjectHandler(projectService services.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// Create cria um novo projeto
// @Summary Criar novo projeto
// @Description Cria um novo projeto associado a um cliente
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.ProjectCreateRequest true "Dados do projeto"
// @Success 201 {object} models.Project
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Cliente não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects [post]
func (h *ProjectHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.ProjectCreateRequest

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para criar projeto
	project, err := h.projectService.Create(userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, project)
}

// List lista todos os projetos do usuário
// @Summary Listar projetos
// @Description Lista todos os projetos do usuário com filtros opcionais
// @Tags projects
// @Security BearerAuth
// @Produce json
// @Param status query string false "Status do projeto (IN_PROGRESS, COMPLETED, CANCELLED)"
// @Param client_id query int false "ID do cliente específico"
// @Param limit query int false "Limite de resultados (padrão: 50)"
// @Param offset query int false "Offset para paginação (padrão: 0)"
// @Success 200 {array} models.Project
// @Failure 400 {object} map[string]interface{} "Parâmetros inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects [get]
func (h *ProjectHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var filter models.ProjectListFilter

	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(errors.NewBadRequestError("Parâmetros de consulta inválidos: " + err.Error()))
		return
	}

	// Validar status se fornecido
	if filter.Status != "" {
		validStatuses := []string{"IN_PROGRESS", "COMPLETED", "CANCELLED"}
		isValid := false
		for _, status := range validStatuses {
			if filter.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			c.Error(errors.NewBadRequestError("Status inválido. Use: IN_PROGRESS, COMPLETED ou CANCELLED"))
			return
		}
	}

	// Chamar service para listar projetos
	projects, err := h.projectService.GetByUserID(userID, &filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetByID obtém um projeto específico
// @Summary Obter projeto por ID
// @Description Obtém os detalhes de um projeto específico
// @Tags projects
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do projeto"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Projeto não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects/{id} [get]
func (h *ProjectHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do projeto da URL
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do projeto inválido"))
		return
	}

	// Chamar service para obter projeto
	project, err := h.projectService.GetByID(userID, uint(projectID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, project)
}

// GetWithTasks obtém um projeto com suas tarefas
// @Summary Obter projeto com tarefas
// @Description Obtém um projeto específico incluindo todas as suas tarefas
// @Tags projects
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do projeto"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Projeto não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects/{id}/with-tasks [get]
func (h *ProjectHandler) GetWithTasks(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do projeto da URL
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do projeto inválido"))
		return
	}

	// Chamar service para obter projeto com tarefas
	project, err := h.projectService.GetWithTasks(userID, uint(projectID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, project)
}

// Update atualiza um projeto existente
// @Summary Atualizar projeto
// @Description Atualiza os dados de um projeto existente
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID do projeto"
// @Param request body models.ProjectUpdateRequest true "Dados para atualização"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Projeto não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects/{id} [put]
func (h *ProjectHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.ProjectUpdateRequest

	// Obter ID do projeto da URL
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do projeto inválido"))
		return
	}

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para atualizar projeto
	updatedProject, err := h.projectService.Update(userID, uint(projectID), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedProject)
}

// Delete exclui um projeto
// @Summary Excluir projeto
// @Description Exclui um projeto e todos os dados relacionados
// @Tags projects
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do projeto"
// @Success 204 "Projeto excluído com sucesso"
// @Failure 400 {object} map[string]interface{} "ID inválido ou projeto tem tarefas associadas"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Projeto não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects/{id} [delete]
func (h *ProjectHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do projeto da URL
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do projeto inválido"))
		return
	}

	// Chamar service para excluir projeto
	err = h.projectService.Delete(userID, uint(projectID))
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetByClient lista projetos de um cliente específico
// @Summary Listar projetos de um cliente
// @Description Lista todos os projetos associados a um cliente específico
// @Tags projects
// @Security BearerAuth
// @Produce json
// @Param clientId path int true "ID do cliente"
// @Success 200 {array} models.Project
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Cliente não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/clients/{clientId}/projects [get]
func (h *ProjectHandler) GetByClient(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do cliente da URL
	clientIDStr := c.Param("clientId")
	clientID, err := strconv.ParseUint(clientIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do cliente inválido"))
		return
	}

	// Chamar service para obter projetos do cliente
	projects, err := h.projectService.GetByClientID(userID, uint(clientID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, projects)
}

// ChangeStatus altera o status de um projeto
// @Summary Alterar status do projeto
// @Description Altera o status de um projeto específico
// @Tags projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID do projeto"
// @Param request body ChangeStatusRequest true "Novo status"
// @Success 200 {object} models.Project
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Projeto não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects/{id}/status [put]
func (h *ProjectHandler) ChangeStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req ChangeStatusRequest

	// Obter ID do projeto da URL
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do projeto inválido"))
		return
	}

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Validar status
	if req.Status == "" {
		c.Error(errors.NewBadRequestError("Status é obrigatório"))
		return
	}

	// Chamar service para alterar status
	project, err := h.projectService.ChangeStatus(userID, uint(projectID), req.Status)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status do projeto alterado com sucesso",
		"project": project,
	})
}

// GetSummary obtém resumo de um projeto
// @Summary Obter resumo do projeto
// @Description Obtém estatísticas e resumo detalhado de um projeto específico
// @Tags projects
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID do projeto"
// @Success 200 {object} services.ProjectSummary
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Projeto não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects/{id}/summary [get]
func (h *ProjectHandler) GetSummary(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Obter ID do projeto da URL
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do projeto inválido"))
		return
	}

	// Chamar service para obter resumo do projeto
	summary, err := h.projectService.GetProjectSummary(userID, uint(projectID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, summary)
}

// ChangeStatusRequest representa os dados para alteração de status
type ChangeStatusRequest struct {
	Status models.ProjectStatus `json:"status" binding:"required" example:"COMPLETED"`
}
