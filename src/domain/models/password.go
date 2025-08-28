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

type Password struct {
	PasswordBase
	ID uint `json:"id"`
}

type PasswordInDB struct {
	PasswordBase
	DBBaseModel
}
