package errors

import (
	"fmt"
	"net/http"
)

// AppError representa um erro da aplicação
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implementa a interface error
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError cria um novo erro da aplicação
func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Erros comuns
var (
	ErrInternalServer = NewAppError(http.StatusInternalServerError, "Erro interno do servidor", "")
	ErrBadRequest     = NewAppError(http.StatusBadRequest, "Requisição inválida", "")
	ErrUnauthorized   = NewAppError(http.StatusUnauthorized, "Não autorizado", "")
	ErrForbidden      = NewAppError(http.StatusForbidden, "Acesso negado", "")
	ErrNotFound       = NewAppError(http.StatusNotFound, "Recurso não encontrado", "")
	ErrConflict       = NewAppError(http.StatusConflict, "Conflito de dados", "")
)

// NewBadRequestError cria um erro de requisição inválida
func NewBadRequestError(details string) *AppError {
	return NewAppError(http.StatusBadRequest, "Requisição inválida", details)
}

// NewNotFoundError cria um erro de recurso não encontrado
func NewNotFoundError(resource string) *AppError {
	return NewAppError(http.StatusNotFound, fmt.Sprintf("%s não encontrado", resource), "")
}

// NewConflictError cria um erro de conflito
func NewConflictError(details string) *AppError {
	return NewAppError(http.StatusConflict, "Conflito de dados", details)
}

// NewUnauthorizedError cria um erro de não autorizado
func NewUnauthorizedError(details string) *AppError {
	return NewAppError(http.StatusUnauthorized, "Não autorizado", details)
}

