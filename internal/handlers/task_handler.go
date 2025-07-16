package handlers

import (
	"crm-backend/internal/models"
	"crm-backend/internal/services"
	"crm-backend/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// TaskHandler gerencia as rotas de tarefas
type TaskHandler struct {
	taskService services.TaskService
}

// NewTaskHandler cria uma nova instância do handler de tarefas
func NewTaskHandler(taskService services.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// Create cria uma nova tarefa
// @Summary Criar nova tarefa
// @Description Cria uma nova tarefa para o usuário
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.TaskCreateRequest true "Dados da tarefa"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato ou projeto não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.TaskCreateRequest

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para criar tarefa
	task, err := h.taskService.Create(userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, task)
}

// List lista todas as tarefas do usuário
// @Summary Listar tarefas
// @Description Lista todas as tarefas do usuário com filtros opcionais
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param status query string false "Status da tarefa (PENDING, COMPLETED)"
// @Param priority query string false "Prioridade (LOW, MEDIUM, HIGH)"
// @Param contact_id query int false "ID do contato específico"
// @Param project_id query int false "ID do projeto específico"
// @Param due_before query string false "Vencimento antes de (formato: 2006-01-02T15:04:05Z)"
// @Param due_after query string false "Vencimento depois de (formato: 2006-01-02T15:04:05Z)"
// @Param limit query int false "Limite de resultados (padrão: 50)"
// @Param offset query int false "Offset para paginação (padrão: 0)"
// @Success 200 {array} models.Task
// @Failure 400 {object} map[string]interface{} "Parâmetros inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	var filter models.TaskListFilter

	// Bind query parameters
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(errors.NewBadRequestError("Parâmetros de consulta inválidos: " + err.Error()))
		return
	}

	// Chamar service para listar tarefas
	tasks, err := h.taskService.GetByUserID(userID, &filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetByID obtém uma tarefa específica
// @Summary Obter tarefa por ID
// @Description Obtém os detalhes de uma tarefa específica
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID da tarefa"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Tarefa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Obter ID da tarefa da URL
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID da tarefa inválido"))
		return
	}

	// Chamar service para obter tarefa
	task, err := h.taskService.GetByID(userID, uint(taskID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, task)
}

// Update atualiza uma tarefa existente
// @Summary Atualizar tarefa
// @Description Atualiza os dados de uma tarefa existente
// @Tags tasks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID da tarefa"
// @Param request body models.TaskUpdateRequest true "Dados para atualização"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Tarefa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.TaskUpdateRequest

	// Obter ID da tarefa da URL
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID da tarefa inválido"))
		return
	}

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para atualizar tarefa
	updatedTask, err := h.taskService.Update(userID, uint(taskID), &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// Delete exclui uma tarefa
// @Summary Excluir tarefa
// @Description Exclui uma tarefa específica
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID da tarefa"
// @Success 204 "Tarefa excluída com sucesso"
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Tarefa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Obter ID da tarefa da URL
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID da tarefa inválido"))
		return
	}

	// Chamar service para excluir tarefa
	err = h.taskService.Delete(userID, uint(taskID))
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// MarkAsCompleted marca uma tarefa como concluída
// @Summary Marcar tarefa como concluída
// @Description Marca uma tarefa específica como concluída
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID da tarefa"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Tarefa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks/{id}/complete [put]
func (h *TaskHandler) MarkAsCompleted(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Obter ID da tarefa da URL
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID da tarefa inválido"))
		return
	}

	// Chamar service para marcar como concluída
	task, err := h.taskService.MarkAsCompleted(userID, uint(taskID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tarefa marcada como concluída",
		"task":    task,
	})
}

// MarkAsPending marca uma tarefa como pendente
// @Summary Marcar tarefa como pendente
// @Description Marca uma tarefa específica como pendente
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID da tarefa"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Tarefa não encontrada"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks/{id}/pending [put]
func (h *TaskHandler) MarkAsPending(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Obter ID da tarefa da URL
	taskIDStr := c.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID da tarefa inválido"))
		return
	}

	// Chamar service para marcar como pendente
	task, err := h.taskService.MarkAsPending(userID, uint(taskID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tarefa marcada como pendente",
		"task":    task,
	})
}

// GetByContact lista tarefas de um contato específico
// @Summary Listar tarefas de um contato
// @Description Lista todas as tarefas associadas a um contato específico
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param contactId path int true "ID do contato"
// @Success 200 {array} models.Task
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Contato não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/contacts/{contactId}/tasks [get]
func (h *TaskHandler) GetByContact(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Obter ID do contato da URL
	contactIDStr := c.Param("contactId")
	contactID, err := strconv.ParseUint(contactIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do contato inválido"))
		return
	}

	// Chamar service para obter tarefas do contato
	tasks, err := h.taskService.GetByContactID(userID, uint(contactID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetByProject lista tarefas de um projeto específico
// @Summary Listar tarefas de um projeto
// @Description Lista todas as tarefas associadas a um projeto específico
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param projectId path int true "ID do projeto"
// @Success 200 {array} models.Task
// @Failure 400 {object} map[string]interface{} "ID inválido"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 403 {object} map[string]interface{} "Acesso negado"
// @Failure 404 {object} map[string]interface{} "Projeto não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/projects/{projectId}/tasks [get]
func (h *TaskHandler) GetByProject(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Obter ID do projeto da URL
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		c.Error(errors.NewBadRequestError("ID do projeto inválido"))
		return
	}

	// Chamar service para obter tarefas do projeto
	tasks, err := h.taskService.GetByProjectID(userID, uint(projectID))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetOverdue obtém tarefas em atraso do usuário
// @Summary Obter tarefas em atraso
// @Description Obtém todas as tarefas em atraso do usuário
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks/overdue [get]
func (h *TaskHandler) GetOverdue(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Chamar service para obter tarefas em atraso
	tasks, err := h.taskService.GetOverdueTasks(userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetUpcoming obtém tarefas próximas do vencimento
// @Summary Obter tarefas próximas do vencimento
// @Description Obtém tarefas que vencem nos próximos dias
// @Tags tasks
// @Security BearerAuth
// @Produce json
// @Param days query int false "Número de dias para buscar (padrão: 7)"
// @Success 200 {array} models.Task
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/tasks/upcoming [get]
func (h *TaskHandler) GetUpcoming(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	// Obter número de dias da query string
	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7
	}

	// Chamar service para obter tarefas próximas do vencimento
	tasks, err := h.taskService.GetUpcomingTasks(userID, days)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

