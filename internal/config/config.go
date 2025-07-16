package config

import (
	"os"
)

// Config representa as configurações da aplicação
type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	Environment string
	LogLevel    string
}

// Load carrega as configurações das variáveis de ambiente
func Load() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://ryan:secure123@localhost:5433/crm-tcc?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key"),
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv obtém uma variável de ambiente ou retorna um valor padrão
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
