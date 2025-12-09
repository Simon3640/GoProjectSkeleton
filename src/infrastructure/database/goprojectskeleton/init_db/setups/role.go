package setups

import (
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/repositories"
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
