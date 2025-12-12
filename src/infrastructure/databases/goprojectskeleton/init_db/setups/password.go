package setups

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	passwordrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/password"
)

// SetupPassword is the setup struct for the password model
type SetupPassword struct {
	SetupBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbmodels.Password]
}

var _ SetupModel[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbmodels.Password] = (*SetupPassword)(nil)

// NewSetupPassword creates a new setup for the password model
func NewSetupPassword() *SetupPassword {
	return &SetupPassword{
		SetupBase: SetupBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbmodels.Password]{
			modelConverter: &passwordrepositories.PasswordConverter{},
		},
	}
}
