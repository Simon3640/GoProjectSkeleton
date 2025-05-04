package db_models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	ID       int    `gorm:"primaryKey"`
	Key      string `gorm:"type:varchar(100);unique;not null"`
	IsActive string `gorm:"not null;type:boolean;default:true"`
	Priority int    `gorm:"not null;type:int;default:0"`
}

func (Role) TableName() string {
	return "role"
}
