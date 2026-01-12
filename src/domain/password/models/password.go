// Package models contains the password models
package models

import (
	"fmt"
	"time"

	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
)

type PasswordBase struct {
	UserID    uint       `json:"user_id"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	IsActive  bool       `json:"isActive"`
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
	sharedmodels.DBBaseModel
}
