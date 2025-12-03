package setups

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/domain/models"
	dbModels "goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"goprojectskeleton/src/infrastructure/repositories"
)

type SetupOneTimeToken struct {
	SetupBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbModels.OneTimeToken]
}

var _ SetupModel[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbModels.OneTimeToken] = (*SetupOneTimeToken)(nil)

func NewSetupOneTimeToken() *SetupOneTimeToken {
	return &SetupOneTimeToken{
		SetupBase: SetupBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbModels.OneTimeToken]{
			modelConverter: &repositories.OneTimeTokenConverter{},
		},
	}
}
