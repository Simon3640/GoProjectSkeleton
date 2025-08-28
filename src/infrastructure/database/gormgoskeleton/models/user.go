package db_models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(100);not null"`
	Email     string     `gorm:"type:varchar(100);not null;unique"`
	Phone     string     `gorm:"type:varchar(20);not null;unique"`
	Status    string     `gorm:"type:varchar(20);not null"`
	RoleID    uint       `gorm:"not null;index"`
	Role      Role       `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Passwords []Password `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (User) TableName() string {
	return "user"
}

var _ DBModel = (*User)(nil)
