package setups

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	authrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/auth"
)

// SetupOneTimeToken is the setup struct for the one time token model
type SetupOneTimeToken struct {
	SetupBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbmodels.OneTimeToken]
}

var _ SetupModel[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbmodels.OneTimeToken] = (*SetupOneTimeToken)(nil)

// NewSetupOneTimeToken creates a new setup for the one time token model
func NewSetupOneTimeToken() *SetupOneTimeToken {
	return &SetupOneTimeToken{
		SetupBase: SetupBase[dtos.OneTimeTokenCreate, dtos.OneTimeTokenUpdate, models.OneTimeToken, dbmodels.OneTimeToken]{
			modelConverter: &authrepositories.OneTimeTokenConverter{},
		},
	}
}
