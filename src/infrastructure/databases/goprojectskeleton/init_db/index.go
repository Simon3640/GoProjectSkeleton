package initdb

import (
	"gorm.io/gorm"

	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	"github.com/simon3640/goprojectskeleton/src/application/shared/defaults"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/init_db/setups"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
)

// InitMigrate initializes the database and returns an error if it fails
func InitMigrate(db *gorm.DB, logger contractsproviders.ILoggerProvider) *applicationerrors.ApplicationError {
	logger.Info("Auto migrating models")

	logger.Info("Auto migrating Role model")
	if err := setups.NewSetUpRole().Setup(db, dbmodels.Role{}, &defaults.DefaultRoles, logger); err != nil {
		return applicationerrors.NewApplicationError(status.DatabaseInitializationError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	logger.Info("Role model migrated")

	logger.Info("Auto migrating User model")
	if err := setups.NewSetupUser().Setup(db, dbmodels.User{}, &defaults.DefaultUsers, logger); err != nil {
		return applicationerrors.NewApplicationError(status.DatabaseInitializationError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	logger.Info("User model migrated")

	logger.Info("Auto migrating Password model")
	if err := setups.NewSetupPassword().Setup(db, dbmodels.Password{}, &defaults.DefaultPasswords, logger); err != nil {
		return applicationerrors.NewApplicationError(status.DatabaseInitializationError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	logger.Info("Password model migrated")

	logger.Info("Auto migrating ended")

	logger.Info("Auto migrating OneTimeToken model")
	if err := setups.NewSetupOneTimeToken().Setup(db, dbmodels.OneTimeToken{}, nil, logger); err != nil {
		return applicationerrors.NewApplicationError(status.DatabaseInitializationError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	logger.Info("OneTimeToken model migrated")

	logger.Info("Auto migrating OneTimePassword model")
	if err := setups.NewSetupOneTimePassword().Setup(db, dbmodels.OneTimePassword{}, nil, logger); err != nil {
		return applicationerrors.NewApplicationError(status.DatabaseInitializationError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	logger.Info("OneTimePassword model migrated")

	return nil
}
