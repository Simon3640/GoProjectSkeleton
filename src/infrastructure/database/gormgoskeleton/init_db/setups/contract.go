package setups

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"

	"gorm.io/gorm"
)

type SetupModel[ModelCreate any, ModelUpdate any, Model any, DBModel any] interface {
	Setup(db *gorm.DB, db_model DBModel, defaults *[]ModelCreate, logger contracts_providers.ILoggerProvider)
}
