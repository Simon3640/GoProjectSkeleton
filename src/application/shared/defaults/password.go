package defaults

import (
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	"github.com/simon3640/goprojectskeleton/src/domain/password/models"
)

var AdminPassword = dtos.PasswordCreate{
	PasswordBase: models.PasswordBase{
		UserID:   1,
		Hash:     "hashed_password_for_admin",
		IsActive: false,
	},
}
var DefaultPasswords = []dtos.PasswordCreate{}
