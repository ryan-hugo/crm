package models

import (
	"time"

	"gorm.io/gorm"
)

// Priority representa a prioridade de uma tarefa
type Priority string

const (
	PriorityLow    Priority = "LOW"
	PriorityMedium Priority = "MEDIUM"
	PriorityHigh   Priority = "HIGH"
)

// TaskStatus representa o status de uma tarefa
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
)

// Task representa uma tarefa
type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"not null" validate:"required,min=2,max=255"`
	Description string         `json:"description,omitempty"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
	Priority    Priority       `json:"priority" gorm:"not null" validate:"required,oneof=LOW MEDIUM HIGH"`
	Status      TaskStatus     `json:"status" gorm:"not null" validate:"required,oneof=PENDING COMPLETED"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	ContactID   *uint          `json:"contact_id,omitempty"`
	ProjectID   *uint          `json:"project_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relacionamentos
	User    User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Contact *Contact `json:"contact,omitempty" gorm:"foreignKey:ContactID"`
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// TaskCreateRequest representa os dados para criação de tarefa
type TaskCreateRequest struct {
	Title       string     `json:"title" validate:"required,min=2,max=255"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Priority    Priority   `json:"priority" validate:"required,oneof=LOW MEDIUM HIGH"`
	Status      TaskStatus `json:"status,omitempty" validate:"omitempty,oneof=PENDING COMPLETED"` // Opcional, será ignorado
	ContactID   *uint      `json:"contact_id,omitempty"`
	ProjectID   *uint      `json:"project_id,omitempty"`
}

// TaskUpdateRequest representa os dados para atualização de tarefa
type TaskUpdateRequest struct {
	Title       string     `json:"title,omitempty" validate:"omitempty,min=2,max=255"`
	Description string     `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Priority    Priority   `json:"priority,omitempty" validate:"omitempty,oneof=LOW MEDIUM HIGH"`
	Status      TaskStatus `json:"status,omitempty" validate:"omitempty,oneof=PENDING COMPLETED"`
	ContactID   *uint      `json:"contact_id,omitempty"`
	ProjectID   *uint      `json:"project_id,omitempty"`
}

// TaskListFilter representa os filtros para listagem de tarefas
type TaskListFilter struct {
	Status    TaskStatus `form:"status" validate:"omitempty,oneof=PENDING COMPLETED"`
	Priority  Priority   `form:"priority" validate:"omitempty,oneof=LOW MEDIUM HIGH"`
	ContactID *uint      `form:"contact_id"`
	ProjectID *uint      `form:"project_id"`
	DueBefore *time.Time `form:"due_before"`
	DueAfter  *time.Time `form:"due_after"`
	Limit     int        `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset    int        `form:"offset" validate:"omitempty,min=0"`
}

