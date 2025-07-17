package models

import (
	"time"

	"gorm.io/gorm"
)

// ContactType representa o tipo de contato
type ContactType string

const (
	ContactTypeClient ContactType = "CLIENT"
	ContactTypeLead   ContactType = "LEAD"
)

// Contact representa um contato (cliente ou lead)
type Contact struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=255"`
	Email     string         `json:"email" gorm:"not null" validate:"required,email"`
	Phone     string         `json:"phone,omitempty" validate:"omitempty,max=50"`
	Company   string         `json:"company,omitempty" validate:"omitempty,max=255"`
	Position  string         `json:"position,omitempty" validate:"omitempty,max=255"`
	Type      ContactType    `json:"type" gorm:"not null" validate:"required,oneof=CLIENT LEAD"`
	Notes     string         `json:"notes,omitempty"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relacionamentos
	User         User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Interactions []Interaction `json:"interactions,omitempty" gorm:"foreignKey:ContactID"`
	Tasks        []Task        `json:"tasks,omitempty" gorm:"foreignKey:ContactID"`
	Projects     []Project     `json:"projects,omitempty" gorm:"foreignKey:ClientID"`
}

// ContactCreateRequest representa os dados para criação de contato
type ContactCreateRequest struct {
	Name     string      `json:"name" validate:"required,min=2,max=255"`
	Email    string      `json:"email" validate:"required,email"`
	Phone    string      `json:"phone,omitempty" validate:"omitempty,max=50"`
	Company  string      `json:"company,omitempty" validate:"omitempty,max=255"`
	Position string      `json:"position,omitempty" validate:"omitempty,max=255"`
	Type     ContactType `json:"type" validate:"required,oneof=CLIENT LEAD"`
	Notes    string      `json:"notes,omitempty"`
}

// ContactUpdateRequest representa os dados para atualização de contato
type ContactUpdateRequest struct {
	Name     string      `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Email    string      `json:"email,omitempty" validate:"omitempty,email"`
	Phone    string      `json:"phone,omitempty" validate:"omitempty,max=50"`
	Company  string      `json:"company,omitempty" validate:"omitempty,max=255"`
	Position string      `json:"position,omitempty" validate:"omitempty,max=255"`
	Type     ContactType `json:"type,omitempty" validate:"omitempty,oneof=CLIENT LEAD"`
	Notes    string      `json:"notes,omitempty"`
}

// ContactListFilter representa os filtros para listagem de contatos
type ContactListFilter struct {
	Type   ContactType `form:"type" validate:"omitempty,oneof=CLIENT LEAD"`
	Search string      `form:"search"`
	Limit  int         `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset int         `form:"offset" validate:"omitempty,min=0"`
}
