package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectStatus string

const (
	ProjectDraft     ProjectStatus = "DRAFT"
	ProjectActive    ProjectStatus = "ACTIVE"
	ProjectPaused    ProjectStatus = "PAUSED"
	ProjectCompleted ProjectStatus = "COMPLETED"
)

type UpdateStatus string

const (
	UpdatePending          UpdateStatus = "PENDING"
	UpdateApproved         UpdateStatus = "APPROVED"
	UpdateChangesRequested UpdateStatus = "CHANGES_REQUESTED"
)

type ResourceType string

const (
	ResourceDigital   ResourceType = "DIGITAL"
	ResourceEconomic  ResourceType = "ECONOMIC"
	ResourceKnowledge ResourceType = "KNOWLEDGE"
	ResourceHardware  ResourceType = "HARDWARE"
)

type ResourceStatus string

const (
	ResourcePending                 ResourceStatus = "PENDING"
	ResourceApproved                ResourceStatus = "APPROVED"
	ResourceRejected                ResourceStatus = "REJECTED"
	ResourceReturnedForModification ResourceStatus = "RETURNED_FOR_MODIFICATION"
)

type Project struct {
	ID            uuid.UUID     `json:"id"`
	Title         string        `json:"title"`
	Description   string        `json:"description"`
	ResearchLine  string        `json:"research_line"`
	Status        ProjectStatus `json:"status"`
	CreatedBy     uuid.UUID     `json:"created_by"`
	StartDate     time.Time     `json:"start_date"`
	TargetEndDate *time.Time    `json:"target_end_date,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`

	// Relaciones opcionales para vistas enriquecidas
	Members []ProjectMember `json:"members,omitempty"`
}

type ProjectMember struct {
	ProjectID  uuid.UUID `json:"project_id"`
	UserID     uuid.UUID `json:"user_id"`
	IsLead     bool      `json:"is_lead"`
	AssignedAt time.Time `json:"assigned_at"`

	// Datos del usuario cargados en join
	User *User `json:"user,omitempty"`
}

type ProjectUpdate struct {
	ID               uuid.UUID       `json:"id"`
	ProjectID        uuid.UUID       `json:"project_id"`
	SubmittedBy      uuid.UUID       `json:"submitted_by"`
	Title            string          `json:"title"`
	Content          string          `json:"content"`
	AttachmentsURL   json.RawMessage `json:"attachments_url"`
	Status           UpdateStatus    `json:"status"`
	DirectorFeedback *string         `json:"director_feedback,omitempty"`
	ReviewedBy       *uuid.UUID      `json:"reviewed_by,omitempty"`
	ReviewedAt       *time.Time      `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}

type ResourceRequest struct {
	ID               uuid.UUID      `json:"id"`
	ProjectID        uuid.UUID      `json:"project_id"`
	RequestedBy      uuid.UUID      `json:"requested_by"`
	ResourceType     ResourceType   `json:"resource_type"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	EstimatedCost    float64        `json:"estimated_cost"`
	Status           ResourceStatus `json:"status"`
	DirectorFeedback *string        `json:"director_feedback,omitempty"`
	ReviewedBy       *uuid.UUID     `json:"reviewed_by,omitempty"`
	ReviewedAt       *time.Time     `json:"reviewed_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
}
