package database

import (
	"gormgoskeleton/src/application/contracts"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"

	"gorm.io/gorm"
)

func InitMigrate(db *gorm.DB, logger contracts.ILoggerProvider) {
	logger.Info("Auto migrating models")

	logger.Info("Auto migrating User model")
	db.AutoMigrate(&db_models.User{})
	logger.Info("User model migrated")

	logger.Info("Auto migrating ended")
}
