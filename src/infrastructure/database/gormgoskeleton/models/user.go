package db_models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name   string `gorm:"type:varchar(100);not null"`
	Email  string `gorm:"type:varchar(100);not null;unique"`
	Phone  string `gorm:"type:varchar(20);not null;unique"`
	Status string `gorm:"type:varchar(20);not null"`
	ID     int    `gorm:"primaryKey"`
}

func (User) TableName() string {
	return "user"
}

var _ DBModel = (*User)(nil)
