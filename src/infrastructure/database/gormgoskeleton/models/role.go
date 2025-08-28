package db_models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Key      string `gorm:"type:varchar(100);unique;not null"`
	IsActive bool   `gorm:"not null;type:boolean;default:true"`
	Priority int    `gorm:"not null;type:int;default:0"`
	Users    []User `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func (Role) TableName() string {
	return "role"
}

var _ DBModel = (*Role)(nil)
