package defaults

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/domain/models"
)

var userStatusActive = models.UserStatusActive

var AdminUser = dtos.UserCreate{
	UserBase: models.UserBase{
		Name:   "Admin",
		Email:  "admin@goprojectskeleton.com",
		Phone:  "1234567890",
		Status: &userStatusActive,
		RoleID: 1,
	},
}

var DefaultUsers = []dtos.UserCreate{
	AdminUser,
}
