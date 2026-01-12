// Package dtomocks contains mock data for DTOs for testing
package dtomocks

import (
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

var userStatusActive = usermodels.UserStatusActive

// UserBase is the base model for a user
var UserBase = usermodels.UserBase{
	Name:   "Test User",
	Email:  "testuser@example.com",
	Phone:  "123",
	Status: &userStatusActive,
	RoleID: 2,
}

// UserCreate is the create model for a user
var UserCreate = userdtos.UserCreate{
	UserBase: UserBase,
}

// UserAndPasswordCreate is the create model for a user and password
var UserAndPasswordCreate = userdtos.UserAndPasswordCreate{
	UserCreate: UserCreate,
	Password:   "$trongPassword123",
}

// UserWithRole is the model for a user with a role
var UserWithRole = usermodels.UserWithRole{
	UserBase: UserBase,
	ID:       1,
}

func init() {
	UserWithRole.SetRole(UserRole)
}
