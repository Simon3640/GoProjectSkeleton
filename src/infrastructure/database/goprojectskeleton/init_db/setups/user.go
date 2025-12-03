package setups

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/domain/models"
	dbModels "goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"goprojectskeleton/src/infrastructure/repositories"
)

type SetupUser struct {
	SetupBase[dtos.UserCreate, dtos.UserUpdate, models.User, dbModels.User]
}

var _ SetupModel[dtos.UserCreate, dtos.UserUpdate, models.User, dbModels.User] = (*SetupUser)(nil)

func NewSetupUser() *SetupUser {
	return &SetupUser{
		SetupBase: SetupBase[dtos.UserCreate, dtos.UserUpdate, models.User, dbModels.User]{
			modelConverter: &repositories.UserConverter{},
		},
	}
}
