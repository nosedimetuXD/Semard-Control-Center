package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserRole string

const (
	RoleDirector      UserRole = "DIRECTOR"
	RoleAdministrador UserRole = "ADMINISTRADOR"
	RoleMiembro       UserRole = "MIEMBRO"
)

type RegistrationStatus string

const (
	RegPending  RegistrationStatus = "PENDING"
	RegApproved RegistrationStatus = "APPROVED"
	RegRejected RegistrationStatus = "REJECTED"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	GoogleID      *string   `json:"google_id,omitempty"`
	Email         string    `json:"email"`
	StudentCode   *string   `json:"student_code,omitempty"`
	FullName      string    `json:"full_name"`
	Role          UserRole  `json:"role"`
	CanOperate3D  bool      `json:"can_operate_3d"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	Bio           *string   `json:"bio,omitempty"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type RegistrationRequest struct {
	ID               uuid.UUID          `json:"id"`
	GoogleID         string             `json:"google_id"`
	Email            string             `json:"email"`
	FullName         string             `json:"full_name"`
	StudentCode      string             `json:"student_code"`
	CareerProgram    *string            `json:"career_program,omitempty"`
	MotivationLetter *string            `json:"motivation_letter,omitempty"`
	Status           RegistrationStatus `json:"status"`
	DirectorFeedback *string            `json:"director_feedback,omitempty"`
	ReviewedBy       *uuid.UUID         `json:"reviewed_by,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	ReviewedAt       *time.Time         `json:"reviewed_at,omitempty"`
}

type JWTClaims struct {
	UserID       uuid.UUID `json:"user_id"`
	Email        string    `json:"email"`
	FullName     string    `json:"full_name"`
	Role         UserRole  `json:"role"`
	CanOperate3D bool      `json:"can_operate_3d"`
	jwt.RegisteredClaims
}
