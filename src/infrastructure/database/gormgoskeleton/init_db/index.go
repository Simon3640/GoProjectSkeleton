package initdb

import (
	"gorm.io/gorm"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	"gormgoskeleton/src/application/shared/defaults"
	"gormgoskeleton/src/infrastructure/database/gormgoskeleton/init_db/setups"
	dbModels "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
)

func InitMigrate(db *gorm.DB, logger contractsProviders.ILoggerProvider) {
	logger.Info("Auto migrating models")

	logger.Info("Auto migrating Role model")
	setups.NewSetUpRole().Setup(db, dbModels.Role{}, &defaults.DefaultRoles, logger)
	logger.Info("Role model migrated")

	logger.Info("Auto migrating User model")
	setups.NewSetupUser().Setup(db, dbModels.User{}, &defaults.DefaultUsers, logger)
	logger.Info("User model migrated")

	logger.Info("Auto migrating Password model")
	setups.NewSetupPassword().Setup(db, dbModels.Password{}, &defaults.DefaultPasswords, logger)
	logger.Info("Password model migrated")

	logger.Info("Auto migrating ended")

	logger.Info("Auto migrating OneTimeToken model")
	setups.NewSetupOneTimeToken().Setup(db, dbModels.OneTimeToken{}, nil, logger)
	logger.Info("OneTimeToken model migrated")

	logger.Info("Auto migrating OneTimePassword model")
	setups.NewSetupOneTimePassword().Setup(db, dbModels.OneTimePassword{}, nil, logger)
	logger.Info("OneTimePassword model migrated")
}
