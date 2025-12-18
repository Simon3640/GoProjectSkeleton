// Package database provides the database connection for the GoProjectSkeleton project
package database

import (
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GoProjectSkeletonDB is the database connection for the GoProjectSkeleton project
type GoProjectSkeletonDB struct {
	DB *gorm.DB
}

// SetUp sets up the database connection for the GoProjectSkeleton project
func (gpsbd *GoProjectSkeletonDB) SetUp(
	host string,
	port string,
	user string,
	password string,
	dbname string,
	ssl *bool,
	logger contractsproviders.ILoggerProvider,
) *applicationerrors.ApplicationError {
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
		return applicationerrors.NewApplicationError(status.DatabaseInitializationError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	gpsbd.DB = db
	logger.Info("Database connection established")
	return nil
}

// GoProjectSkeletondb is the database connection for the GoProjectSkeleton project
var GoProjectSkeletondb *GoProjectSkeletonDB

func init() {
	GoProjectSkeletondb = &GoProjectSkeletonDB{}
}
