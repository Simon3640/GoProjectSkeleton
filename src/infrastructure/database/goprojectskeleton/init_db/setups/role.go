package setups

import (
	"goprojectskeleton/src/domain/models"
	dbModels "goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"goprojectskeleton/src/infrastructure/repositories"
)

type SetupRole struct {
	SetupBase[models.RoleCreate, models.RoleUpdate, models.Role, dbModels.Role]
}

var _ SetupModel[models.RoleCreate, models.RoleUpdate, models.Role, dbModels.Role] = (*SetupRole)(nil)

func NewSetUpRole() *SetupRole {
	return &SetupRole{
		SetupBase: SetupBase[models.RoleCreate, models.RoleUpdate, models.Role, dbModels.Role]{
			modelConverter: &repositories.RoleConverter{},
		},
	}
}
