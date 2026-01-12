// Package passworddtos contains the DTOs for the password module
package passworddtos

import (
	"time"

	passwordmodels "github.com/simon3640/goprojectskeleton/src/domain/password/models"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

// PasswordCreate is the DTO for creating a new password
type PasswordCreate struct {
	passwordmodels.PasswordBase
}

// PasswordCreateNoHash is the DTO for creating a new password without a hash
type PasswordCreateNoHash struct {
	UserID           uint       `json:"user_id"`
	NoHashedPassword string     `json:"noHashedPassword"`
	ExpiresAt        *time.Time `json:"expiresAt,omitempty"`
	IsActive         bool       `json:"isActive"`
}

// GetUserID returns the user ID
func (p PasswordCreateNoHash) GetUserID() uint {
	return p.UserID
}

// ExpiresAt is a pointer to allow it to be optional but if not provided, it defaults to 30 days from now.
func (p *PasswordCreate) SetDefaultExpiresAt() {
	if p.ExpiresAt == nil {
		defaultExpiry := time.Now().Add(30 * 24 * time.Hour) // 30 days
		p.ExpiresAt = &defaultExpiry
	}
}

// Validate validates the password create no hash
func (p PasswordCreateNoHash) Validate() []string {
	var errs []string
	if !sharedmodels.IsValidPassword(p.NoHashedPassword) {
		errs = append(errs, "Invalid password")
	}
	return errs
}

// NewPasswordCreate creates a new password create DTO
func NewPasswordCreate(userID uint, hash string, expiresAt *time.Time, isActive bool) PasswordCreate {
	p := PasswordCreate{
		PasswordBase: passwordmodels.PasswordBase{
			UserID:    userID,
			Hash:      hash,
			ExpiresAt: expiresAt,
			IsActive:  isActive,
		},
	}
	p.SetDefaultExpiresAt()
	return p
}

// PasswordTokenCreate is the DTO for creating a password reset token
type PasswordTokenCreate struct {
	Token            string `json:"token"`
	NoHashedPassword string `json:"noHashedPassword"`
}

// Validate validates the password token create
func (p PasswordTokenCreate) Validate() []string {
	var errs []string
	if p.Token == "" {
		errs = append(errs, "Token is required")
	}
	if !sharedmodels.IsValidPassword(p.NoHashedPassword) {
		errs = append(errs, "Invalid password")
	}
	return errs
}

// PasswordUpdateBase is the base DTO for updating a password
type PasswordUpdateBase struct {
	Hash      *string    `json:"hash"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	IsActive  *bool      `json:"isActive,omitempty"`
}

// PasswordUpdate is the DTO for updating a password
type PasswordUpdate struct {
	PasswordUpdateBase
	ID uint `json:"id"`
}
