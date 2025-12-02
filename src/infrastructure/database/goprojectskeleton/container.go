// Package database provides the database connection for the GoProjectSkeleton project
package database

import (
	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	initdb "goprojectskeleton/src/infrastructure/database/goprojectskeleton/init_db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GoProjectSkeletonDB is the database connection for the GoProjectSkeleton project
type GoProjectSkeletonDB struct {
	DB *gorm.DB
}

// SetUp sets up the database connection for the GoProjectSkeleton project
func (gpsbd *GoProjectSkeletonDB) SetUp(host string, port string, user string, password string, dbname string, ssl *bool, logger contractsProviders.ILoggerProvider) {
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
		logger.Panic("Error connecting to database", err)
	}
	gpsbd.DB = db
	logger.Info("Database connection established")
	initdb.InitMigrate(db, logger)
}

// GoProjectSkeletondb is the database connection for the GoProjectSkeleton project
var GoProjectSkeletondb *GoProjectSkeletonDB

func init() {
	GoProjectSkeletondb = &GoProjectSkeletonDB{}
}
