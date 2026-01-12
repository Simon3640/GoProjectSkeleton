// Package defaults contains the default values for the user module
package defaults

import (
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	"github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

var userStatusActive = models.UserStatusActive

// AdminUser is the admin user
var AdminUser = userdtos.UserCreate{
	UserBase: models.UserBase{
		Name:   "Admin",
		Email:  "admin@goprojectskeleton.com",
		Phone:  "1234567890",
		Status: &userStatusActive,
		RoleID: 1,
	},
}

// DefaultUsers is the default users
var DefaultUsers = []userdtos.UserCreate{
	AdminUser,
}
