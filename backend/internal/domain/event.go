package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventVisibility string

const (
	VisibilityPublic   EventVisibility = "PUBLIC"
	VisibilityInternal EventVisibility = "INTERNAL"
)

type Event struct {
	ID          uuid.UUID       `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Visibility  EventVisibility `json:"visibility"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     *time.Time      `json:"end_time,omitempty"`
	Location    string          `json:"location"`
	BannerURL   *string         `json:"banner_url,omitempty"`
	CreatedBy   uuid.UUID       `json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
