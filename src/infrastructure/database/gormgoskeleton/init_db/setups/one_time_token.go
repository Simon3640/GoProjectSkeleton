package setups

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
	"gormgoskeleton/src/infrastructure/repositories"
)

type SetupOneTimeToken struct {
	SetupBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, db_models.OneTimeToken]
}

var _ SetupModel[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, db_models.OneTimeToken] = (*SetupOneTimeToken)(nil)

func NewSetupOneTimeToken() *SetupOneTimeToken {
	return &SetupOneTimeToken{
		SetupBase: SetupBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, db_models.OneTimeToken]{
			modelConverter: &repositories.OneTimeTokenConverter{},
		},
	}
}
