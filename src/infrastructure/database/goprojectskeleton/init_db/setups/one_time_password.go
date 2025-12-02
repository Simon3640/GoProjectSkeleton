package setups

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/domain/models"
	dbModels "goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"goprojectskeleton/src/infrastructure/repositories"
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
