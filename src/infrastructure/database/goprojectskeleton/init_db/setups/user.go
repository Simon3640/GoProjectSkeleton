package setups

import (
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/repositories"
)

type SetupUser struct {
	SetupBase[userdtos.UserCreate, userdtos.UserUpdate, models.User, dbModels.User]
}

var _ SetupModel[userdtos.UserCreate, userdtos.UserUpdate, models.User, dbModels.User] = (*SetupUser)(nil)

func NewSetupUser() *SetupUser {
	return &SetupUser{
		SetupBase: SetupBase[userdtos.UserCreate, userdtos.UserUpdate, models.User, dbModels.User]{
			modelConverter: &repositories.UserConverter{},
		},
	}
}
