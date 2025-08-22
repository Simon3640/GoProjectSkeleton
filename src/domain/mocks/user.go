package domain_mocks

import "gormgoskeleton/src/domain/models"

var UserBase = models.UserBase{
	Name:   "Test User",
	Email:  "testuser@example.com",
	Phone:  "123",
	Status: "active",
	RoleID: 2,
}

var UserCreate = models.UserCreate{
	UserBase: UserBase,
}

var UserAndPasswordCreate = models.UserAndPasswordCreate{
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
