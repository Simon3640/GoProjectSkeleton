package database

import (
	"fmt"
	"reflect"

	"gorm.io/gorm"

	"gormgoskeleton/src/application/contracts"
	"gormgoskeleton/src/domain/defaults"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
)

func InitMigrate(db *gorm.DB, logger contracts.ILoggerProvider) {
	logger.Info("Auto migrating models")

	logger.Info("Auto migrating User model")
	SetUpModel(db, &db_models.User{}, defaults.DefaultUsers, logger, "User")
	logger.Info("User model migrated")

	logger.Info("Auto migrating Role model")
	SetUpModel(db, &db_models.Role{}, defaults.DefaultRoles, logger, "Role")
	logger.Info("Role model migrated")

	logger.Info("Auto migrating ended")
}

func SetUpModel(db *gorm.DB, model any, defaults any, logger contracts.ILoggerProvider, name string) {
	err := db.AutoMigrate(model)
	if err != nil {
		logger.Error(fmt.Sprintf("Error auto-migrating %s model", name), err)
		return
	}

	if db.Migrator().HasTable(model) {
		logger.Info(fmt.Sprintf("Table %s exists", name))
	} else {
		logger.Error(fmt.Sprintf("Table %s does not exist even after AutoMigrate", name), nil)
		return
	}

	if defaults != nil {
		logger.Info("Creating default data")
		val := reflect.ValueOf(defaults)
		if val.Kind() == reflect.Slice {
			for i := 0; i < val.Len(); i++ {
				item := val.Index(i).Interface()
				if err := db.Create(item).Error; err != nil {
					logger.Error(fmt.Sprintf("Error creating default %s data", name), err)
				} else {
					logger.Info(fmt.Sprintf("Default %s data created", name))
				}
			}
		}
	}

	logger.Info("Model setup completed")
}
