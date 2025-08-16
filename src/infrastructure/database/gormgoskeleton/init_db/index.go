package initdb

import (
	"gorm.io/gorm"

	"gormgoskeleton/src/application/contracts"
	"gormgoskeleton/src/domain/defaults"
	"gormgoskeleton/src/infrastructure/database/gormgoskeleton/init_db/setups"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
)

func InitMigrate(db *gorm.DB, logger contracts.ILoggerProvider) {
	logger.Info("Auto migrating models")

	logger.Info("Auto migrating Role model")
	setups.NewSetUpRole().Setup(db, db_models.Role{}, defaults.DefaultRoles, logger)
	logger.Info("Role model migrated")

	logger.Info("Auto migrating User model")
	setups.NewSetupUser().Setup(db, db_models.User{}, defaults.DefaultUsers, logger)
	logger.Info("User model migrated")

	logger.Info("Auto migrating Password model")
	setups.NewSetupPassword().Setup(db, db_models.Password{}, defaults.DefaultPasswords, logger)
	logger.Info("Password model migrated")

	logger.Info("Auto migrating ended")
}
