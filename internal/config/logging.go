package config

import (
	"os"
	"strconv"
	"strings"
)

// LoggingConfig contém as configurações de logging
type LoggingConfig struct {
	Level      string
	Format     string
	Output     string
	MaxSize    int // Em MB
	MaxBackups int
	MaxAge     int // Em dias
	Compress   bool
}

// GetLoggingConfig retorna as configurações de logging
func GetLoggingConfig() *LoggingConfig {
	return &LoggingConfig{
		Level:      getEnvOrDefault("LOG_LEVEL", "INFO"),
		Format:     getEnvOrDefault("LOG_FORMAT", "text"),   // text ou json
		Output:     getEnvOrDefault("LOG_OUTPUT", "stdout"), // stdout, file, ou both
		MaxSize:    getIntEnvOrDefault("LOG_MAX_SIZE", 10),
		MaxBackups: getIntEnvOrDefault("LOG_MAX_BACKUPS", 5),
		MaxAge:     getIntEnvOrDefault("LOG_MAX_AGE", 30),
		Compress:   getBoolEnvOrDefault("LOG_COMPRESS", true),
	}
}

// IsDebugEnabled verifica se o debug está habilitado
func (c *LoggingConfig) IsDebugEnabled() bool {
	return strings.ToUpper(c.Level) == "DEBUG"
}

// IsLevelEnabled verifica se um nível específico está habilitado
func (c *LoggingConfig) IsLevelEnabled(level string) bool {
	levels := map[string]int{
		"DEBUG":   0,
		"INFO":    1,
		"WARNING": 2,
		"ERROR":   3,
	}

	currentLevel, exists := levels[strings.ToUpper(c.Level)]
	if !exists {
		currentLevel = 1 // Default para INFO
	}

	requestedLevel, exists := levels[strings.ToUpper(level)]
	if !exists {
		return false
	}

	return requestedLevel >= currentLevel
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
