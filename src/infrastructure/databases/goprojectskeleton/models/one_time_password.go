package dbmodels

import (
	"time"

	"gorm.io/gorm"
)

type OneTimePassword struct {
	gorm.Model
	UserID  uint      `gorm:"not null;index"`
	Purpose string    `gorm:"not null:varchar(255)"`
	Hash    []byte    `gorm:"not null;varchar(255);uniqueIndex"`
	IsUsed  bool      `gorm:"not null"`
	Expires time.Time `gorm:"not null"`
}

func (OneTimePassword) TableName() string {
	return "one_time_password"
}

var _ DBModel = (*OneTimePassword)(nil)
