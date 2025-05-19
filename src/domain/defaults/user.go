package defaults

import "gormgoskeleton/src/domain/models"

var AdminUser = models.UserCreate{
	UserBase: models.UserBase{
		Name:   "Admin",
		Email:  "admin@gormgoskeleton.com",
		Phone:  "1234567890",
		Status: "active",
	},
}

var DefaultUsers = []models.UserCreate{
	AdminUser,
}
