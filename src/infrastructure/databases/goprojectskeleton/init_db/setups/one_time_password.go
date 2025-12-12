package setups

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	repositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories"
)

// SetupOneTimePassword is the setup struct for the one time password model
type SetupOneTimePassword struct {
	SetupBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, dbmodels.OneTimePassword]
}

var _ SetupModel[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, dbmodels.OneTimePassword] = (*SetupOneTimePassword)(nil)

// NewSetupOneTimePassword creates a new setup for the one time password model
func NewSetupOneTimePassword() *SetupOneTimePassword {
	return &SetupOneTimePassword{
		SetupBase: SetupBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, dbmodels.OneTimePassword]{
			modelConverter: &repositories.OneTimePasswordConverter{},
		},
	}
}
