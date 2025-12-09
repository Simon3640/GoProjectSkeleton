package defaults

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

var AdminPassword = dtos.PasswordCreate{
	PasswordBase: models.PasswordBase{
		UserID:   1,
		Hash:     "hashed_password_for_admin",
		IsActive: false,
	},
}
var DefaultPasswords = []dtos.PasswordCreate{}
