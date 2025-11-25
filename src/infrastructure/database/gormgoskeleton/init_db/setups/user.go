package setups

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
	dbModels "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
	"gormgoskeleton/src/infrastructure/repositories"
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
