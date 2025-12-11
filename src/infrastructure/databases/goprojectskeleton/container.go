// Package database provides the database connection for the GoProjectSkeleton project
package database

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	initdb "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/init_db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GoProjectSkeletonDB is the database connection for the GoProjectSkeleton project
type GoProjectSkeletonDB struct {
	DB *gorm.DB
}

// SetUp sets up the database connection for the GoProjectSkeleton project
func (gpsbd *GoProjectSkeletonDB) SetUp(host string, port string, user string, password string, dbname string, ssl *bool, logger contractsProviders.ILoggerProvider) *application_errors.ApplicationError {
	var sslmode string
	if ssl != nil && *ssl {
		logger.Info("SSL is enabled")
		sslmode = "require"
	} else {
		logger.Info("SSL is disabled")
		sslmode = "disable"
	}
	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=" + sslmode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return application_errors.NewApplicationError(status.DatabaseInitializationError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	gpsbd.DB = db
	logger.Info("Database connection established")
	if err := initdb.InitMigrate(db, logger); err != nil {
		return application_errors.NewApplicationError(status.DatabaseInitializationError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	return nil
}

// GoProjectSkeletondb is the database connection for the GoProjectSkeleton project
var GoProjectSkeletondb *GoProjectSkeletonDB

func init() {
	GoProjectSkeletondb = &GoProjectSkeletonDB{}
}
