package defaults

import "gormgoskeleton/src/domain/models"

var AdminPassword = models.PasswordCreate{
	PasswordBase: models.PasswordBase{
		UserID:   1,
		Hash:     "hashed_password_for_admin",
		IsActive: false,
	},
}
var DefaultPasswords = []models.PasswordCreate{}
