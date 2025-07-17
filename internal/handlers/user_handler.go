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

// UserHandler gerencia as rotas de usuários
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler cria uma nova instância do handler de usuários
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile obtém o perfil do usuário autenticado
// @Summary Obter perfil do usuário
// @Description Retorna os dados do perfil do usuário autenticado
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 404 {object} map[string]interface{} "Usuário não encontrado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	start := time.Now()

	// Debug: verificar se user_id está no contexto
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		logger.Error("user_id não encontrado no contexto")
		c.Error(errors.NewUnauthorizedError("Usuário não autenticado"))
		return
	}

	logger.Debugf("user_id do contexto: %v (tipo: %T)", userIDInterface, userIDInterface)

	userID := c.GetUint("user_id")
	if userID == 0 {
		logger.Errorf("user_id é zero após conversão: %d", userID)
		c.Error(errors.NewUnauthorizedError("ID do usuário inválido"))
		return
	}

	logger.Debugf("Buscando perfil do usuário ID: %d", userID)

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		logger.LogError(err, "Erro ao buscar perfil do usuário", map[string]interface{}{
			"user_id": userID,
		})
		c.Error(err)
		return
	}

	duration := time.Since(start)
	logger.WithFields("INFO", "User Profile Retrieved", map[string]interface{}{
		"user_id":  userID,
		"duration": duration,
	})

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile atualiza o perfil do usuário
// @Summary Atualizar perfil do usuário
// @Description Atualiza os dados do perfil do usuário autenticado
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.UserUpdateRequest true "Dados para atualização"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 409 {object} map[string]interface{} "Email já existe"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req models.UserUpdateRequest

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Chamar service para atualizar perfil
	updatedProfile, err := h.userService.UpdateProfile(userID, &req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Perfil atualizado com sucesso",
		"user":    updatedProfile,
	})
}

// ChangePassword altera a senha do usuário
// @Summary Alterar senha do usuário
// @Description Altera a senha do usuário autenticado
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body ChangePasswordRequest true "Dados para alteração de senha"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Senha atual incorreta"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/users/change-password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req ChangePasswordRequest

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Validar campos obrigatórios
	if req.CurrentPassword == "" || req.NewPassword == "" {
		c.Error(errors.NewBadRequestError("Senha atual e nova senha são obrigatórias"))
		return
	}

	// Validar confirmação de senha
	if req.NewPassword != req.ConfirmPassword {
		c.Error(errors.NewBadRequestError("Nova senha e confirmação não conferem"))
		return
	}

	// Chamar service para alterar senha
	err := h.userService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Senha alterada com sucesso",
	})
}

// DeleteAccount exclui a conta do usuário
// @Summary Excluir conta do usuário
// @Description Exclui permanentemente a conta do usuário autenticado
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body DeleteAccountRequest true "Confirmação de senha"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{} "Dados inválidos"
// @Failure 401 {object} map[string]interface{} "Senha incorreta"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/users/delete-account [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req DeleteAccountRequest

	// Validar entrada JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewBadRequestError("Dados de entrada inválidos: " + err.Error()))
		return
	}

	// Validar campo obrigatório
	if req.Password == "" {
		c.Error(errors.NewBadRequestError("Senha é obrigatória para confirmar exclusão"))
		return
	}

	// Chamar service para excluir conta
	err := h.userService.DeleteAccount(userID, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Conta excluída com sucesso",
	})
}

// GetStats obtém estatísticas do usuário
// @Summary Obter estatísticas do usuário
// @Description Retorna estatísticas consolidadas do usuário (contatos, tarefas, projetos)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} services.UserStats
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/users/stats [get]
func (h *UserHandler) GetStats(c *gin.Context) {
	userID := c.GetUint("user_id")

	stats, err := h.userService.GetUserStats(userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetRecentActivities obtém as atividades recentes do usuário
// @Summary Obter atividades recentes do usuário
// @Description Retorna as atividades recentes do usuário autenticado (tarefas, projetos, contatos e interações)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param limit query int false "Limite de resultados (padrão: 10)"
// @Success 200 {object} models.RecentActivityResponse
// @Failure 401 {object} map[string]interface{} "Não autorizado"
// @Failure 500 {object} map[string]interface{} "Erro interno"
// @Router /api/users/activities [get]
func (h *UserHandler) GetRecentActivities(c *gin.Context) {
	start := time.Now()
	userID := c.GetUint("user_id")

	// Obter limite da query string
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	activities, err := h.userService.GetRecentActivities(userID, limit)
	if err != nil {
		logger.LogError(err, "Erro ao buscar atividades recentes", map[string]interface{}{
			"user_id": userID,
			"limit":   limit,
		})
		c.Error(err)
		return
	}

	duration := time.Since(start)
	logger.WithFields("INFO", "User Recent Activities Retrieved", map[string]interface{}{
		"user_id":        userID,
		"limit":          limit,
		"activity_count": activities.Count,
		"duration":       duration,
	})

	c.JSON(http.StatusOK, activities)
}

// ChangePasswordRequest representa os dados para alteração de senha
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required" example:"senhaAtual123"`
	NewPassword     string `json:"new_password" binding:"required,min=6" example:"novaSenha456"`
	ConfirmPassword string `json:"confirm_password" binding:"required" example:"novaSenha456"`
}

// DeleteAccountRequest representa os dados para exclusão de conta
type DeleteAccountRequest struct {
	Password string `json:"password" binding:"required" example:"minhaSenh123"`
}
