package dbmodels

import (
	"time"

	"gorm.io/gorm"
)

type Password struct {
	gorm.Model
	UserID    uint       `gorm:"not null;index"`
	Hash      string     `gorm:"type:varchar(255);not null"`
	ExpiresAt *time.Time `gorm:"type:timestamp"`
	IsActive  bool       `gorm:"not null;type:boolean;default:true"`
}

func (Password) TableName() string {
	return "password"
}

var _ DBModel = (*Password)(nil)
