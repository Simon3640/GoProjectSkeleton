package defaults

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
)

var AdminPassword = dtos.PasswordCreate{
	PasswordBase: models.PasswordBase{
		UserID:   1,
		Hash:     "hashed_password_for_admin",
		IsActive: false,
	},
}
var DefaultPasswords = []dtos.PasswordCreate{}
