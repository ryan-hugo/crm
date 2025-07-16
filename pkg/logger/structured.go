package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

// StructuredLogger provides structured logging capabilities
type StructuredLogger struct {
	logger *log.Logger
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	Source    string                 `json:"source,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// NewStructuredLogger creates a new structured logger
func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{
		logger: log.New(os.Stdout, "", 0),
	}
}

// Log logs a structured entry
func (sl *StructuredLogger) Log(level, message string, fields map[string]interface{}, err error) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   message,
		Fields:    fields,
		Source:    getCallerInfo(),
	}

	if err != nil {
		entry.Error = err.Error()
	}

	jsonData, jsonErr := json.Marshal(entry)
	if jsonErr != nil {
		// Fallback to regular logging if JSON marshaling fails
		sl.logger.Printf("ERROR: Failed to marshal log entry: %v", jsonErr)
		return
	}

	sl.logger.Println(string(jsonData))
}

// Convenience methods for structured logging
func (sl *StructuredLogger) Info(message string, fields map[string]interface{}) {
	sl.Log("INFO", message, fields, nil)
}

func (sl *StructuredLogger) Warning(message string, fields map[string]interface{}) {
	sl.Log("WARNING", message, fields, nil)
}

func (sl *StructuredLogger) Error(message string, fields map[string]interface{}, err error) {
	sl.Log("ERROR", message, fields, err)
}

func (sl *StructuredLogger) Debug(message string, fields map[string]interface{}) {
	if isDebugMode() {
		sl.Log("DEBUG", message, fields, nil)
	}
}

// Global structured logger instance
var StructuredLog *StructuredLogger

// InitStructuredLogger initializes the structured logger
func InitStructuredLogger() {
	StructuredLog = NewStructuredLogger()
}

// Example usage functions
func LogDatabaseOperation(operation string, table string, duration time.Duration, success bool, err error) {
	fields := map[string]interface{}{
		"operation": operation,
		"table":     table,
		"duration":  duration.String(),
		"success":   success,
	}

	if success {
		StructuredLog.Info("Database operation completed", fields)
	} else {
		StructuredLog.Error("Database operation failed", fields, err)
	}
}

func LogAPICall(method, endpoint string, statusCode int, duration time.Duration, userID uint) {
	fields := map[string]interface{}{
		"method":      method,
		"endpoint":    endpoint,
		"status_code": statusCode,
		"duration":    duration.String(),
		"user_id":     userID,
	}

	level := "INFO"
	if statusCode >= 400 {
		level = "WARNING"
	}
	if statusCode >= 500 {
		level = "ERROR"
	}

	StructuredLog.Log(level, "API call", fields, nil)
}

func LogBusinessEvent(event string, entityType string, entityID uint, userID uint, details map[string]interface{}) {
	fields := map[string]interface{}{
		"event":       event,
		"entity_type": entityType,
		"entity_id":   entityID,
		"user_id":     userID,
	}

	// Merge additional details
	for k, v := range details {
		fields[k] = v
	}

	StructuredLog.Info("Business event", fields)
}

// LogPerformanceMetrics logs performance metrics
func LogPerformanceMetrics(component string, metrics map[string]interface{}) {
	fields := map[string]interface{}{
		"component": component,
	}

	// Merge metrics
	for k, v := range metrics {
		fields[k] = v
	}

	StructuredLog.Info("Performance metrics", fields)
}

// Example of how to use in your handlers
func ExampleUsage() {
	// Initialize structured logger
	InitStructuredLogger()

	// Log a database operation
	start := time.Now()
	// ... database operation ...
	LogDatabaseOperation("SELECT", "contacts", time.Since(start), true, nil)

	// Log an API call
	LogAPICall("POST", "/api/contacts", 201, time.Since(start), 123)

	// Log a business event
	LogBusinessEvent("contact_created", "contact", 456, 123, map[string]interface{}{
		"email": "user@example.com",
		"type":  "CLIENT",
	})

	// Log performance metrics
	LogPerformanceMetrics("contact_service", map[string]interface{}{
		"total_contacts":    1000,
		"avg_response_time": "50ms",
		"memory_usage":      "25MB",
	})
}
