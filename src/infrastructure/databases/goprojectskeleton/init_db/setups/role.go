package setups

import (
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	repositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories"
)

// SetupRole is the setup struct for the role model
type SetupRole struct {
	SetupBase[models.RoleCreate, models.RoleUpdate, models.Role, dbmodels.Role]
}

var _ SetupModel[models.RoleCreate, models.RoleUpdate, models.Role, dbmodels.Role] = (*SetupRole)(nil)

// NewSetUpRole creates a new setup for the role model
func NewSetUpRole() *SetupRole {
	return &SetupRole{
		SetupBase: SetupBase[models.RoleCreate, models.RoleUpdate, models.Role, dbmodels.Role]{
			modelConverter: &repositories.RoleConverter{},
		},
	}
}
