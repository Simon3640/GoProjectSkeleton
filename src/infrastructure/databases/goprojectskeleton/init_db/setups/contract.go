package setups

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"

	"gorm.io/gorm"
)

// SetupModel is the contract for the setup models
type SetupModel[ModelCreate any, ModelUpdate any, Model any, DBModel any] interface {
	Setup(db *gorm.DB, dbModel DBModel, defaults *[]ModelCreate, logger contractsProviders.ILoggerProvider) error
}
