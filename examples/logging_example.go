package main

import (
	"crm-backend/pkg/logger"
	"errors"
	"time"
)

func main() {
	// Inicializar logger
	logger.Init()
	logger.InitStructuredLogger()

	// Exemplo de logs básicos
	logger.Info("Aplicação iniciada")
	logger.Warning("Esta é uma mensagem de aviso")
	logger.Debug("Esta mensagem só aparece em modo debug")

	// Exemplo de logs formatados
	logger.Infof("Usuário %s logado às %s", "admin", time.Now().Format("15:04:05"))

	// Exemplo de logs estruturados
	logger.WithFields("INFO", "User Login", map[string]interface{}{
		"user_id":   123,
		"email":     "admin@example.com",
		"ip":        "192.168.1.100",
		"timestamp": time.Now(),
	})

	// Exemplo de log de erro
	err := errors.New("erro de exemplo")
	logger.LogError(err, "Database Connection", map[string]interface{}{
		"host":     "localhost",
		"port":     5432,
		"database": "crm_db",
	})

	// Exemplo de log de chamada de serviço
	start := time.Now()
	time.Sleep(50 * time.Millisecond) // Simular operação
	logger.LogServiceCall("UserService", "GetProfile", time.Since(start), true)

	// Exemplo de logging estruturado em JSON
	logger.StructuredLog.Info("Structured log example", map[string]interface{}{
		"component": "test",
		"action":    "demo",
		"count":     42,
	})

	// Exemplo de logs de diferentes tipos
	logger.LogDatabaseOperation("SELECT", "users", 25*time.Millisecond, true, nil)
	logger.LogAPICall("GET", "/api/users", 200, 30*time.Millisecond, 123)
	logger.LogBusinessEvent("user_profile_viewed", "user", 123, 123, map[string]interface{}{
		"profile_id": 456,
		"view_type":  "full",
	})

	logger.Info("Exemplo de logging concluído")
}
