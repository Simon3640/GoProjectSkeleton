package setups

import (
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
	"gormgoskeleton/src/infrastructure/repositories"
)

type SetupUser struct {
	SetupBase[models.UserCreate, models.UserUpdate, models.User, db_models.User]
}

var _ SetupModel[models.UserCreate, models.UserUpdate, models.User, db_models.User] = (*SetupUser)(nil)

func NewSetupUser() *SetupUser {
	return &SetupUser{
		SetupBase: SetupBase[models.UserCreate, models.UserUpdate, models.User, db_models.User]{
			modelConverter: &repositories.UserConverter{},
		},
	}
}
