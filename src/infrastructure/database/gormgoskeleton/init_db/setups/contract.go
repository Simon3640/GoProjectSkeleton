package setups

import (
	contractsProviders "gormgoskeleton/src/application/contracts/providers"

	"gorm.io/gorm"
)

type SetupModel[ModelCreate any, ModelUpdate any, Model any, DBModel any] interface {
	Setup(db *gorm.DB, dbModel DBModel, defaults *[]ModelCreate, logger contractsProviders.ILoggerProvider)
}
