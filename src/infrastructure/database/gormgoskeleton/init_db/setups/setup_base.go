package setups

import (
	"fmt"

	"gormgoskeleton/src/application/contracts"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
	"gormgoskeleton/src/infrastructure/repositories"

	"gorm.io/gorm"
)

type SetupBase[CreateModel any, UpdateModel any, Model any, DBModel db_models.DBModel] struct {
	modelConverter repositories.ModelConverter[CreateModel, UpdateModel, Model, DBModel]
}

func (s *SetupBase[CreateModel, UpdateModel, Model, DBModel]) Setup(db *gorm.DB,
	db_model DBModel,
	defaults []CreateModel,
	logger contracts.ILoggerProvider) {
	logger.Info(fmt.Sprintf("Auto migrating the %s table", db_model.TableName()))
	db_has_table := db.Migrator().HasTable(db_model)
	err := db.AutoMigrate(db_model)
	if err != nil {
		logger.Error(fmt.Sprintf("Error auto-migrating %s model", db_model.TableName()), err)
		panic(err) // TODO: Bad practice
	}

	if db_has_table {
		logger.Info(fmt.Sprintf("Table %s exists", db_model.TableName()))
		return
	} else {
		logger.Info(fmt.Sprintf("Table %s does not exist creating defaults if needed", db_model.TableName()))
	}

	if defaults != nil {
		logger.Info("Creating default data")
		for _, item := range defaults {
			db_model_item := s.modelConverter.ToGormCreate(item)
			if err := db.Table(db_model.TableName()).Create(&db_model_item).Error; err != nil {
				logger.Error(fmt.Sprintf("Error creating default %s data", db_model.TableName()), err)
				panic(err) // TODO: Bad practice
			} else {
				logger.Info(fmt.Sprintf("Default %s data created", db_model.TableName()))
			}
		}
	}

}
