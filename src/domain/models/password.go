package models

import (
	"fmt"
	"time"
	"unicode"
)

type PasswordBase struct {
	UserID    uint       `json:"user_id"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsActive  bool       `json:"is_active"`
	Hash      string     `json:"hash"`
}

func (p PasswordBase) UserIDString() string {
	return fmt.Sprintf("%d", p.UserID)
}

type PasswordCreate struct {
	PasswordBase
}

type PasswordCreateNoHash struct {
	UserID           uint       `json:"user_id"`
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

func IsValidPassword(p string) bool {
	var hasMinLen, hasUpper, hasLower, hasNumber, hasSpecial bool
	if len(p) >= 8 {
		hasMinLen = true
	}

	for _, char := range p {
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

func (p PasswordCreateNoHash) IsValidPassword() bool {
	return IsValidPassword(p.NoHashedPassword)
}

func NewPasswordCreate(userID uint, hash string, expiresAt *time.Time, isActive bool) PasswordCreate {
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
	ID uint `json:"id"`
}

type Password struct {
	PasswordBase
	ID uint `json:"id"`
}

type PasswordInDB struct {
	PasswordBase
	DBBaseModel
}
