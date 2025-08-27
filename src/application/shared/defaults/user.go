package defaults

import (
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/domain/models"
)

var AdminUser = dtos.UserCreate{
	UserBase: models.UserBase{
		Name:   "Admin",
		Email:  "admin@gormgoskeleton.com",
		Phone:  "1234567890",
		Status: "active",
		RoleID: 1,
	},
}

var DefaultUsers = []dtos.UserCreate{
	AdminUser,
}
