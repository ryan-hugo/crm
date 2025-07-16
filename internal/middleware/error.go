package middleware

import (
	"crm-backend/pkg/errors"
	"crm-backend/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware para tratamento global de erros
func ErrorHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()

		// Verificar se houve algum erro
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			
			// Verificar se é um erro da aplicação
			if appErr, ok := err.Err.(*errors.AppError); ok {
				logger.Warning("Application error:", appErr.Message, "Details:", appErr.Details)
				c.JSON(appErr.Code, gin.H{
					"error":   appErr.Message,
					"details": appErr.Details,
				})
				return
			}

			// Erro genérico
			logger.Error("Unexpected error:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Erro interno do servidor",
			})
		}
	})
}

