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

type EventRegistration struct {
	ID           uuid.UUID `json:"id"`
	EventID      uuid.UUID `json:"event_id"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        *string   `json:"phone,omitempty"`
	Institution  *string   `json:"institution,omitempty"`
	RegisteredAt time.Time `json:"registered_at"`
}

type CreateEventRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Visibility  EventVisibility `json:"visibility"`
	StartTime   time.Time       `json:"start_time"`
	EndTime     *time.Time      `json:"end_time,omitempty"`
	Location    string          `json:"location"`
	BannerURL   *string         `json:"banner_url,omitempty"`
}

type UpdateEventRequest struct {
	Title       *string          `json:"title,omitempty"`
	Description *string          `json:"description,omitempty"`
	Visibility  *EventVisibility `json:"visibility,omitempty"`
	StartTime   *time.Time       `json:"start_time,omitempty"`
	EndTime     *time.Time       `json:"end_time,omitempty"`
	Location    *string          `json:"location,omitempty"`
	BannerURL   *string          `json:"banner_url,omitempty"`
}

type RegisterAttendeeRequest struct {
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	Phone       *string `json:"phone,omitempty"`
	Institution *string `json:"institution,omitempty"`
}

