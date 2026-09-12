package domain

import (
	"time"

	"github.com/google/uuid"
)

type LoanStatus string

const (
	LoanPending          LoanStatus = "PENDING"
	LoanApproved         LoanStatus = "APPROVED"
	LoanApprovedModified LoanStatus = "APPROVED_MODIFIED"
	LoanRejected         LoanStatus = "REJECTED"
	LoanReturned         LoanStatus = "RETURNED"
	LoanOverdue          LoanStatus = "OVERDUE"
)

type InventoryItem struct {
	ID             uuid.UUID `json:"id"`
	Code           string    `json:"code"`
	Name           string    `json:"name"`
	Category       string    `json:"category"`
	Description    *string   `json:"description,omitempty"`
	TotalStock     int       `json:"total_stock"`
	AvailableStock int       `json:"available_stock"`
	Location       *string   `json:"location,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type LoanRequest struct {
	ID                 uuid.UUID  `json:"id"`
	ItemID             uuid.UUID  `json:"item_id"`
	RequestedBy        uuid.UUID  `json:"requested_by"`
	Reason             string     `json:"reason"`
	RequestedStartDate time.Time  `json:"requested_start_date"`
	RequestedEndDate   time.Time  `json:"requested_end_date"`
	ApprovedEndDate    *time.Time `json:"approved_end_date,omitempty"`
	Status             LoanStatus `json:"status"`
	ReviewerFeedback   *string    `json:"reviewer_feedback,omitempty"`
	ReviewedBy         *uuid.UUID `json:"reviewed_by,omitempty"`
	ReturnedAt         *time.Time `json:"returned_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Relaciones opcionales cargadas por join
	Item *InventoryItem `json:"item,omitempty"`
	User *User          `json:"user,omitempty"`
}
