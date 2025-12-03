package setups

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/domain/models"
	dbModels "goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"goprojectskeleton/src/infrastructure/repositories"
)

type SetupPassword struct {
	SetupBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbModels.Password]
}

var _ SetupModel[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbModels.Password] = (*SetupPassword)(nil)

func NewSetupPassword() *SetupPassword {
	return &SetupPassword{
		SetupBase: SetupBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbModels.Password]{
			modelConverter: &repositories.PasswordConverter{},
		},
	}
}
