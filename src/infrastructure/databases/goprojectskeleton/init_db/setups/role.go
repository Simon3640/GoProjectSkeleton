package setups

import (
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
)

// SetupRole is the setup struct for the role model
type SetupRole struct {
	SetupBase[usermodels.RoleCreate, usermodels.RoleUpdate, usermodels.Role, dbmodels.Role]
}

var _ SetupModel[usermodels.RoleCreate, usermodels.RoleUpdate, usermodels.Role, dbmodels.Role] = (*SetupRole)(nil)

// NewSetUpRole creates a new setup for the role model
func NewSetUpRole() *SetupRole {
	return &SetupRole{
		SetupBase: SetupBase[usermodels.RoleCreate, usermodels.RoleUpdate, usermodels.Role, dbmodels.Role]{
			modelConverter: &userrepositories.RoleConverter{},
		},
	}
}
