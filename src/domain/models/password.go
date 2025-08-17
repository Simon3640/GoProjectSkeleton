package models

import (
	"time"
	"unicode"
)

type PasswordBase struct {
	UserID    int        `json:"user_id"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsActive  bool       `json:"is_active"`
	Hash      string     `json:"hash"`
}

type PasswordCreate struct {
	PasswordBase
}

type PasswordCreateNoHash struct {
	UserID           int        `json:"user_id"`
	NoHashedPassword string     `json:"no_hashed_password"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	IsActive         bool       `json:"is_active"`
}

// ExpiresAt is a pointer to allow it to be optional but if not provided, it defaults to 30 days from now.
func (p *PasswordCreate) SetDefaultExpiresAt() {
	if p.ExpiresAt == nil {
		defaultExpiry := time.Now().Add(30 * 24 * time.Hour) // 30 days
		p.ExpiresAt = &defaultExpiry
	}
}

func (p PasswordCreateNoHash) IsValidPassword() bool {
	// regex for password validation
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool
	if len(p.NoHashedPassword) >= 8 {
		hasMinLen = true
	}

	for _, char := range p.NoHashedPassword {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func NewPasswordCreate(userID int, hash string, expiresAt *time.Time, isActive bool) PasswordCreate {
	p := PasswordCreate{
		PasswordBase: PasswordBase{
			UserID:    userID,
			Hash:      hash,
			ExpiresAt: expiresAt,
			IsActive:  isActive,
		},
	}
	p.SetDefaultExpiresAt()
	return p
}

type PasswordUpdateBase struct {
	Hash      *string    `json:"hash"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
}

type PasswordUpdate struct {
	PasswordUpdateBase
	ID int `json:"id"`
}

type Password struct {
	PasswordBase
	ID int `json:"id"`
}

type PasswordInDB struct {
	PasswordBase
	DBBaseModel
}
