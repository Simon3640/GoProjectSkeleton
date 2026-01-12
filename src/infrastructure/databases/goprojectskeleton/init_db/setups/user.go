package setups

import (
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	userrepositories "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/user"
)

// SetupUser is the setup struct for the user model
type SetupUser struct {
	SetupBase[userdtos.UserCreate, userdtos.UserUpdate, usermodels.User, dbmodels.User]
}

var _ SetupModel[userdtos.UserCreate, userdtos.UserUpdate, usermodels.User, dbmodels.User] = (*SetupUser)(nil)

// NewSetupUser creates a new setup for the user model
func NewSetupUser() *SetupUser {
	return &SetupUser{
		SetupBase: SetupBase[userdtos.UserCreate, userdtos.UserUpdate, usermodels.User, dbmodels.User]{
			modelConverter: &userrepositories.UserConverter{},
		},
	}
}
