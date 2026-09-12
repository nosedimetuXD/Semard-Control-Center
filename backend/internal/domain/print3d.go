package domain

import (
	"time"

	"github.com/google/uuid"
)

type Print3DStatus string

const (
	Print3DPending    Print3DStatus = "PENDING"
	Print3DApproved   Print3DStatus = "APPROVED"
	Print3DInProgress Print3DStatus = "IN_PROGRESS"
	Print3DCompleted  Print3DStatus = "COMPLETED"
	Print3DDelivered  Print3DStatus = "DELIVERED"
	Print3DRejected   Print3DStatus = "REJECTED"
)

type Print3DRequest struct {
	ID               uuid.UUID     `json:"id"`
	RequestedBy      uuid.UUID     `json:"requested_by"`
	ProjectID        *uuid.UUID    `json:"project_id,omitempty"`
	FilePath         string        `json:"file_path"`
	FileName         string        `json:"file_name"`
	Material         string        `json:"material"`
	Color            *string       `json:"color,omitempty"`
	InfillPercentage int           `json:"infill_percentage"`
	Notes            *string       `json:"notes,omitempty"`
	Status           Print3DStatus `json:"status"`
	ReviewerFeedback *string       `json:"reviewer_feedback,omitempty"`
	ApprovedBy       *uuid.UUID    `json:"approved_by,omitempty"`
	OperatorID       *uuid.UUID    `json:"operator_id,omitempty"`
	ApprovedAt       *time.Time    `json:"approved_at,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`

	// Relaciones opcionales cargadas por join
	User     *User    `json:"user,omitempty"`
	Project  *Project `json:"project,omitempty"`
	Operator *User    `json:"operator,omitempty"`
}
