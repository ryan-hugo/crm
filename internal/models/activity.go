package models

import "time"

// ActivityType define o tipo de atividade
type ActivityType string

const (
	// Tipos de atividade
	ActivityTypeTask        ActivityType = "TASK"        // Nova tarefa, tarefa concluída, editada, excluída
	ActivityTypeProject     ActivityType = "PROJECT"     // Novo projeto, atualização de status, editado, excluído
	ActivityTypeContact     ActivityType = "CONTACT"     // Novo contato, atualização de tipo, editado, excluído
	ActivityTypeInteraction ActivityType = "INTERACTION" // Nova interação, editada, excluída
)

// ActivityAction define o tipo de ação realizada
type ActivityAction string

const (
	// Ações de atividade
	ActionCreated   ActivityAction = "CREATED"   // Item criado
	ActionUpdated   ActivityAction = "UPDATED"   // Item atualizado
	ActionCompleted ActivityAction = "COMPLETED" // Item concluído (tarefas)
	ActionDeleted   ActivityAction = "DELETED"   // Item excluído
	ActionStarted   ActivityAction = "STARTED"   // Projeto iniciado
	ActionCancelled ActivityAction = "CANCELLED" // Projeto cancelado
)

// UserActivity representa uma atividade recente do usuário
type UserActivity struct {
	ID          uint           `json:"id"`
	Type        ActivityType   `json:"type"`
	Action      ActivityAction `json:"action"`
	Title       string         `json:"title"`
	Detail      string         `json:"detail,omitempty"`
	ItemID      uint           `json:"item_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	RelatedID   *uint          `json:"related_id,omitempty"`
	RelatedName *string        `json:"related_name,omitempty"`
}

// RecentActivityResponse representa uma resposta de atividades recentes
type RecentActivityResponse struct {
	Activities []UserActivity `json:"activities"`
	Count      int            `json:"count"`
}
