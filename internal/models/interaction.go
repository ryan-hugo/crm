package models

import (
	"time"

	"gorm.io/gorm"
)

// InteractionType representa o tipo de interação
type InteractionType string

const (
	InteractionTypeEmail   InteractionType = "EMAIL"
	InteractionTypeCall    InteractionType = "CALL"
	InteractionTypeMeeting InteractionType = "MEETING"
	InteractionTypeOther   InteractionType = "OTHER"
)

// Interaction representa uma interação com um contato
type Interaction struct {
	ID          uint               `json:"id" gorm:"primaryKey"`
	Type        InteractionType    `json:"type" gorm:"not null" validate:"required,oneof=EMAIL CALL MEETING OTHER"`
	Date        time.Time          `json:"date" gorm:"not null" validate:"required"`
	Subject     string             `json:"subject,omitempty" validate:"omitempty,max=255"`
	Description string             `json:"description,omitempty"`
	ContactID   uint               `json:"contact_id" gorm:"not null"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `json:"-" gorm:"index"`

	// Relacionamentos
	Contact Contact `json:"contact,omitempty" gorm:"foreignKey:ContactID"`
}

// InteractionCreateRequest representa os dados para criação de interação
type InteractionCreateRequest struct {
	Type        InteractionType `json:"type" validate:"required,oneof=EMAIL CALL MEETING OTHER"`
	Date        time.Time       `json:"date" validate:"required"`
	Subject     string          `json:"subject,omitempty" validate:"omitempty,max=255"`
	Description string          `json:"description,omitempty"`
}

// InteractionUpdateRequest representa os dados para atualização de interação
type InteractionUpdateRequest struct {
	Type        InteractionType `json:"type,omitempty" validate:"omitempty,oneof=EMAIL CALL MEETING OTHER"`
	Date        *time.Time      `json:"date,omitempty"`
	Subject     string          `json:"subject,omitempty" validate:"omitempty,max=255"`
	Description string          `json:"description,omitempty"`
}

// InteractionListFilter representa os filtros para listagem de interações
type InteractionListFilter struct {
	Type      InteractionType `form:"type" validate:"omitempty,oneof=EMAIL CALL MEETING OTHER"`
	DateFrom  *time.Time      `form:"date_from"`
	DateTo    *time.Time      `form:"date_to"`
	ContactID uint            `form:"contact_id"`
	Limit     int             `form:"limit" validate:"omitempty,min=1,max=100"`
	Offset    int             `form:"offset" validate:"omitempty,min=0"`
}

