package setups

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/repositories"
)

type SetupPassword struct {
	SetupBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbModels.Password]
}

var _ SetupModel[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbModels.Password] = (*SetupPassword)(nil)

func NewSetupPassword() *SetupPassword {
	return &SetupPassword{
		SetupBase: SetupBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbModels.Password]{
			modelConverter: &repositories.PasswordConverter{},
		},
	}
}
