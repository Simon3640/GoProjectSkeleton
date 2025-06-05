package setups

import (
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"
	"gormgoskeleton/src/infrastructure/repositories"
)

type SetupRole struct {
	SetupBase[models.RoleCreate, models.RoleUpdate, models.Role, db_models.Role]
}

var _ SetupModel[models.RoleCreate, models.RoleUpdate, models.Role, db_models.Role] = (*SetupRole)(nil)

func NewSetUpRole() *SetupRole {
	return &SetupRole{
		SetupBase: SetupBase[models.RoleCreate, models.RoleUpdate, models.Role, db_models.Role]{
			modelConverter: &repositories.RoleConverter{},
		},
	}
}
