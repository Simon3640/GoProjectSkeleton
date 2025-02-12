package database

import (
	"gormgoskeleton/src/application/contracts"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormGoSkeletonDB struct{}

var db *gorm.DB

func (ggsbd GormGoSkeletonDB) SetUp(host string, port string, user string, password string, dbname string, logger contracts.ILoggerProvider) {
	dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + dbname + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Panic("Error connecting to database", err)
	}
	db = db
	logger.Info("Database connection established")
}

var Gormgoskeletondb *GormGoSkeletonDB

func init() {
	Gormgoskeletondb = &GormGoSkeletonDB{}
}
