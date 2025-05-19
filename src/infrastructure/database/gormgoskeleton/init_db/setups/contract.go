package setups

import (
	"gormgoskeleton/src/application/contracts"

	"gorm.io/gorm"
)

type SetupModel[ModelCreate any, ModelUpdate any, Model any, DBModel any] interface {
	Setup(db *gorm.DB, db_model DBModel, defaults []ModelCreate, logger contracts.ILoggerProvider)
}
