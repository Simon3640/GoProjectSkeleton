package defaults

import "gormgoskeleton/src/domain/models"

// list of default roles

var AdminRole = models.RoleCreate{
	RoleBase: models.RoleBase{
		Key:      "admin",
		IsActive: true,
		Priority: 0,
	},
}

var UserRole = models.RoleCreate{
	RoleBase: models.RoleBase{
		Key:      "user",
		IsActive: true,
		Priority: 5,
	},
}

var DefaultRoles = []models.RoleCreate{
	AdminRole,
	UserRole,
}
