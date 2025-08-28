package database

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	initdb "gormgoskeleton/src/infrastructure/database/gormgoskeleton/init_db"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormGoSkeletonDB struct{}

var DB *gorm.DB

func (ggsbd GormGoSkeletonDB) SetUp(host string, port string, user string, password string, dbname string, ssl *bool, logger contracts_providers.ILoggerProvider) {
	sslmode := "disable"
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
	DB = db
	logger.Info("Database connection established")
	initdb.InitMigrate(db, logger)
}

var Gormgoskeletondb *GormGoSkeletonDB

func init() {
	Gormgoskeletondb = &GormGoSkeletonDB{}
}
