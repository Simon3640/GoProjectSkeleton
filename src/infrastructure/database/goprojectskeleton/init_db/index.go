package initdb

import (
	"gorm.io/gorm"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	"github.com/simon3640/goprojectskeleton/src/application/shared/defaults"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/init_db/setups"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
)

// InitMigrate initializes the database and returns an application error if it fails
func InitMigrate(db *gorm.DB, logger contractsProviders.ILoggerProvider) error {
	logger.Info("Auto migrating models")

	logger.Info("Auto migrating Role model")
	if err := setups.NewSetUpRole().Setup(db, dbModels.Role{}, &defaults.DefaultRoles, logger); err != nil {
		logger.Error("Error migrating Role model", err)
		return err
	}
	logger.Info("Role model migrated")

	logger.Info("Auto migrating User model")
	if err := setups.NewSetupUser().Setup(db, dbModels.User{}, &defaults.DefaultUsers, logger); err != nil {
		logger.Error("Error migrating User model", err)
		return err
	}
	logger.Info("User model migrated")

	logger.Info("Auto migrating Password model")
	if err := setups.NewSetupPassword().Setup(db, dbModels.Password{}, &defaults.DefaultPasswords, logger); err != nil {
		logger.Error("Error migrating Password model", err)
		return err
	}
	logger.Info("Password model migrated")

	logger.Info("Auto migrating ended")

	logger.Info("Auto migrating OneTimeToken model")
	if err := setups.NewSetupOneTimeToken().Setup(db, dbModels.OneTimeToken{}, nil, logger); err != nil {
		logger.Error("Error migrating OneTimeToken model", err)
		return err
	}
	logger.Info("OneTimeToken model migrated")

	logger.Info("Auto migrating OneTimePassword model")
	if err := setups.NewSetupOneTimePassword().Setup(db, dbModels.OneTimePassword{}, nil, logger); err != nil {
		logger.Error("Error migrating OneTimePassword model", err)
		return err
	}
	logger.Info("OneTimePassword model migrated")

	return nil
}
