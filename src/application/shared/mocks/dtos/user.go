package dtomocks

import (
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/domain/models"
)

var userStatusActive = models.UserStatusActive

var UserBase = models.UserBase{
	Name:   "Test User",
	Email:  "testuser@example.com",
	Phone:  "123",
	Status: &userStatusActive,
	RoleID: 2,
}

var UserCreate = dtos.UserCreate{
	UserBase: UserBase,
}

var UserAndPasswordCreate = dtos.UserAndPasswordCreate{
	UserCreate: UserCreate,
	Password:   "$trongPassword123",
}

var UserWithRole = models.UserWithRole{
	UserBase: UserBase,
	ID:       1,
}

func init() {
	UserWithRole.SetRole(UserRole)
}
