package db_models

import "gorm.io/gorm"

type DBModel interface {
	gorm.Model
	TableName() string
}
