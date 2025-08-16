package setups

import (
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
	"gormgoskeleton/src/infrastructure/repositories"
)

type SetupPassword struct {
	SetupBase[models.PasswordCreate, models.PasswordUpdate, models.Password, db_models.Password]
}

var _ SetupModel[models.PasswordCreate, models.PasswordUpdate, models.Password, db_models.Password] = (*SetupPassword)(nil)

func NewSetupPassword() *SetupPassword {
	return &SetupPassword{
		SetupBase: SetupBase[models.PasswordCreate, models.PasswordUpdate, models.Password, db_models.Password]{
			modelConverter: &repositories.PasswordConverter{},
		},
	}
}
