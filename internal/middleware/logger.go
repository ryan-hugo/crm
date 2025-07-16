package middleware

import (
	"crm-backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger middleware para registrar requisições HTTP
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Usar a nova função LogRequest do logger
		logger.LogRequest(
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Request.UserAgent(),
		)
		return ""
	})
}

// CustomLogger middleware mais detalhado
func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Processar requisição
		c.Next()

		// Calcular tempo de resposta
		latency := time.Since(start)

		// Obter informações da requisição
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		userAgent := c.Request.UserAgent()

		if raw != "" {
			path = path + "?" + raw
		}

		// Determinar nível de log baseado no status code
		fields := map[string]interface{}{
			"method":      method,
			"path":        path,
			"status_code": statusCode,
			"latency":     latency,
			"client_ip":   clientIP,
			"body_size":   bodySize,
			"user_agent":  userAgent,
		}

		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.String()
		}

		// Log baseado no status code
		if statusCode >= 500 {
			logger.WithFields("ERROR", "HTTP Server Error", fields)
		} else if statusCode >= 400 {
			logger.WithFields("WARNING", "HTTP Client Error", fields)
		} else {
			logger.WithFields("INFO", "HTTP Request", fields)
		}
	}
}
