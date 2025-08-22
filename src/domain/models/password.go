package models

import (
	"fmt"
	"time"
)

type PasswordBase struct {
	UserID    uint       `json:"user_id"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	IsActive  bool       `json:"is_active"`
	Hash      string     `json:"hash"`
}

func (p PasswordBase) GetUserID() uint {
	return p.UserID
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

func (p PasswordCreateNoHash) Validate() []string {
	var errs []string
	if !IsValidPassword(p.NoHashedPassword) {
		errs = append(errs, "Invalid password")
	}
	return errs
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
