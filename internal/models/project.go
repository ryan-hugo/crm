package models

import (
	"time"

	"gorm.io/gorm"
)

// ProjectStatus representa o status de um projeto
type ProjectStatus string

const (
	ProjectStatusInProgress ProjectStatus = "IN_PROGRESS"
	ProjectStatusCompleted  ProjectStatus = "COMPLETED"
	ProjectStatusCancelled  ProjectStatus = "CANCELLED"
)

// Project representa um projeto
type Project struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null" validate:"required,min=2,max=255"`
	Description string         `json:"description,omitempty"`
	Status      ProjectStatus  `json:"status" gorm:"not null" validate:"required,oneof=IN_PROGRESS COMPLETED CANCELLED"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	ClientID    uint           `json:"client_id" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relacionamentos
	User   User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Client Contact `json:"client,omitempty" gorm:"foreignKey:ClientID"`
	Tasks  []Task  `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
}

// ProjectCreateRequest representa os dados para criação de projeto
type ProjectCreateRequest struct {
	Name        string        `json:"name" validate:"required,min=2,max=255"`
	Description string        `json:"description,omitempty"`
	Status      ProjectStatus `json:"status" validate:"required,oneof=IN_PROGRESS COMPLETED CANCELLED"`
	ClientID    uint          `json:"client_id" validate:"required"`
}

// ProjectUpdateRequest representa os dados para atualização de projeto
type ProjectUpdateRequest struct {
	Name        string        `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Description string        `json:"description,omitempty"`
	Status      ProjectStatus `json:"status,omitempty" validate:"omitempty,oneof=IN_PROGRESS COMPLETED CANCELLED"`
	ClientID    uint          `json:"client_id,omitempty"`
}

// ProjectListFilter representa os filtros para listagem de projetos
type ProjectListFilter struct {
	Status   ProjectStatus `form:"status" validate:"omitempty,oneof=IN_PROGRESS COMPLETED CANCELLED"`
	ClientID *uint         `form:"client_id"`
	Limit    int           `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset   int           `form:"offset" validate:"omitempty,min=0"`
}

