package setups

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
	dbModels "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
	"gormgoskeleton/src/infrastructure/repositories"
)

type SetupOneTimePassword struct {
	SetupBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, dbModels.OneTimePassword]
}

var _ SetupModel[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, dbModels.OneTimePassword] = (*SetupOneTimePassword)(nil)

func NewSetupOneTimePassword() *SetupOneTimePassword {
	return &SetupOneTimePassword{
		SetupBase: SetupBase[dtos.OneTimePasswordCreate, dtos.OneTimePasswordUpdate, models.OneTimePassword, dbModels.OneTimePassword]{
			modelConverter: &repositories.OneTimePasswordConverter{},
		},
	}
}
