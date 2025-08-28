package setups

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
	"gormgoskeleton/src/infrastructure/repositories"
)

type SetupPassword struct {
	SetupBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, db_models.Password]
}

var _ SetupModel[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, db_models.Password] = (*SetupPassword)(nil)

func NewSetupPassword() *SetupPassword {
	return &SetupPassword{
		SetupBase: SetupBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, db_models.Password]{
			modelConverter: &repositories.PasswordConverter{},
		},
	}
}
