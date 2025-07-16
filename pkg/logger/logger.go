package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

// Init inicializa os loggers
func Init() {
	// Configuração com timestamp mais detalhado
	flags := log.Ldate | log.Ltime | log.Lmicroseconds

	InfoLogger = log.New(os.Stdout, "INFO: ", flags)
	WarningLogger = log.New(os.Stdout, "WARNING: ", flags)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", flags)
	DebugLogger = log.New(os.Stdout, "DEBUG: ", flags)
}

// Info registra uma mensagem de informação
func Info(v ...interface{}) {
	InfoLogger.Println(v...)
}

// Infof registra uma mensagem de informação formatada
func Infof(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// Warning registra uma mensagem de aviso
func Warning(v ...interface{}) {
	WarningLogger.Println(v...)
}

// Warningf registra uma mensagem de aviso formatada
func Warningf(format string, v ...interface{}) {
	WarningLogger.Printf(format, v...)
}

// Error registra uma mensagem de erro
func Error(v ...interface{}) {
	ErrorLogger.Println(v...)
}

// Errorf registra uma mensagem de erro formatada
func Errorf(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}

// Debug registra uma mensagem de debug (apenas em modo desenvolvimento)
func Debug(v ...interface{}) {
	if isDebugMode() {
		DebugLogger.Println(v...)
	}
}

// Debugf registra uma mensagem de debug formatada
func Debugf(format string, v ...interface{}) {
	if isDebugMode() {
		DebugLogger.Printf(format, v...)
	}
}

// WithFields registra uma mensagem com campos estruturados
func WithFields(level string, message string, fields map[string]interface{}) {
	var logger *log.Logger

	switch strings.ToUpper(level) {
	case "INFO":
		logger = InfoLogger
	case "WARNING":
		logger = WarningLogger
	case "ERROR":
		logger = ErrorLogger
	case "DEBUG":
		if !isDebugMode() {
			return
		}
		logger = DebugLogger
	default:
		logger = InfoLogger
	}

	fieldsStr := ""
	for k, v := range fields {
		fieldsStr += fmt.Sprintf(" %s=%v", k, v)
	}

	logger.Printf("%s%s", message, fieldsStr)
}

// Fatal registra uma mensagem de erro e termina o programa
func Fatal(v ...interface{}) {
	ErrorLogger.Fatal(v...)
}

// Fatalf registra uma mensagem de erro formatada e termina o programa
func Fatalf(format string, v ...interface{}) {
	ErrorLogger.Fatalf(format, v...)
}

// LogRequest registra informações de uma requisição HTTP
func LogRequest(method, path string, statusCode int, duration time.Duration, clientIP, userAgent string) {
	WithFields("INFO", "HTTP Request", map[string]interface{}{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration":    duration,
		"client_ip":   clientIP,
		"user_agent":  userAgent,
	})
}

// LogError registra um erro com contexto adicional
func LogError(err error, context string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	fields["error"] = err.Error()
	fields["context"] = context
	fields["file"] = getCallerInfo()

	WithFields("ERROR", "Application Error", fields)
}

// LogServiceCall registra uma chamada de serviço
func LogServiceCall(service, method string, duration time.Duration, success bool) {
	level := "INFO"
	if !success {
		level = "WARNING"
	}

	WithFields(level, "Service Call", map[string]interface{}{
		"service":  service,
		"method":   method,
		"duration": duration,
		"success":  success,
	})
}

// isDebugMode verifica se está em modo debug
func isDebugMode() bool {
	return os.Getenv("DEBUG") == "true" || os.Getenv("GIN_MODE") != "release"
}

// getCallerInfo obtém informações do arquivo que chamou o log
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown"
	}

	// Pega apenas o nome do arquivo sem o caminho completo
	parts := strings.Split(file, "/")
	if len(parts) > 0 {
		file = parts[len(parts)-1]
	}

	return fmt.Sprintf("%s:%d", file, line)
}
