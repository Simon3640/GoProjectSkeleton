package setups

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/repositories"
)

type SetupUser struct {
	SetupBase[dtos.UserCreate, dtos.UserUpdate, models.User, dbModels.User]
}

var _ SetupModel[dtos.UserCreate, dtos.UserUpdate, models.User, dbModels.User] = (*SetupUser)(nil)

func NewSetupUser() *SetupUser {
	return &SetupUser{
		SetupBase: SetupBase[dtos.UserCreate, dtos.UserUpdate, models.User, dbModels.User]{
			modelConverter: &repositories.UserConverter{},
		},
	}
}
