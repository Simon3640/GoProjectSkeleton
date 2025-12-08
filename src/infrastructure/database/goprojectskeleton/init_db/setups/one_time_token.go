package setups

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/repositories"
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
