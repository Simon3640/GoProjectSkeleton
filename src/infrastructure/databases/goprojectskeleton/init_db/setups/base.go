// Package setups contains the setup functions for the database models
package setups

import (
	"fmt"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	repositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/shared"

	"gorm.io/gorm"
)

// SetupBase is the abstract base setup struct for the database models
type SetupBase[CreateModel any, UpdateModel any, Model any, DBModel dbmodels.DBModel] struct {
	modelConverter repositories.ModelConverter[CreateModel, UpdateModel, Model, DBModel]
}

// Setup sets up the database models
func (s *SetupBase[CreateModel, UpdateModel, Model, DBModel]) Setup(db *gorm.DB,
	dbModel DBModel,
	defaults *[]CreateModel,
	logger contractsProviders.ILoggerProvider) error {
	logger.Info(fmt.Sprintf("Auto migrating the %s table", dbModel.TableName()))
	dbHasTable := db.Migrator().HasTable(dbModel)
	err := db.AutoMigrate(dbModel)
	if err != nil {
		logger.Error(fmt.Sprintf("Error auto-migrating %s model", dbModel.TableName()), err)
		return err
	}

	if dbHasTable {
		logger.Info(fmt.Sprintf("Table %s exists", dbModel.TableName()))
		return nil
	} else {
		logger.Info(fmt.Sprintf("Table %s does not exist creating defaults if needed", dbModel.TableName()))
	}

	if defaults != nil {
		logger.Info("Creating default data")
		for _, item := range *defaults {
			dbModelItem := s.modelConverter.ToGormCreate(item)
			if err := db.Table(dbModel.TableName()).Create(&dbModelItem).Error; err != nil {
				logger.Error(fmt.Sprintf("Error creating default %s data", dbModel.TableName()), err)
				return err
			} else {
				logger.Info(fmt.Sprintf("Default %s data created", dbModel.TableName()))
			}
		}
	}
	return nil
}
