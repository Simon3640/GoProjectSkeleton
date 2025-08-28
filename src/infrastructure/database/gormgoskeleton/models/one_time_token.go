package db_models

import (
	"time"

	"gorm.io/gorm"
)

type OneTimeToken struct {
	gorm.Model
	UserID  uint      `gorm:"not null;index"`
	Purpose string    `gorm:"not null:varchar(255)"`
	Hash    string    `gorm:"not null;varchar(255);uniqueIndex"`
	IsUsed  bool      `gorm:"not null"`
	Expires time.Time `gorm:"not null"`
}

func (OneTimeToken) TableName() string {
	return "one_time_token"
}

var _ DBModel = (*OneTimeToken)(nil)
