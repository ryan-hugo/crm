package models

import (
	"time"

	"gorm.io/gorm"
)

// User representa um usuário do sistema
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null" validate:"required,min=2,max=255"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password  string         `json:"-" gorm:"not null" validate:"required,min=6"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relacionamentos
	Contacts     []Contact `json:"contacts,omitempty" gorm:"foreignKey:UserID"`
	Tasks        []Task    `json:"tasks,omitempty" gorm:"foreignKey:UserID"`
	Projects     []Project `json:"projects,omitempty" gorm:"foreignKey:UserID"`
}

// UserCreateRequest representa os dados para criação de usuário
type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserUpdateRequest representa os dados para atualização de usuário
type UserUpdateRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponse representa a resposta de usuário (sem senha)
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converte User para UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

