package dtomocks

import "github.com/simon3640/goprojectskeleton/src/domain/models"

// TODO: Evaluate if Role can be created with DTOs

var UserRoleBase = models.RoleBase{
	Key:      "user",
	IsActive: true,
	Priority: 5,
}

var AdminRoleBase = models.RoleBase{
	Key:      "admin",
	IsActive: true,
	Priority: 1,
}

var UserRoleCreate = models.RoleCreate{
	RoleBase: UserRoleBase,
}

var AdminRoleCreate = models.RoleCreate{
	RoleBase: AdminRoleBase,
}

var UserRole = models.Role{
	RoleBase: UserRoleBase,
	ID:       2,
}

var AdminRole = models.Role{
	RoleBase: AdminRoleBase,
	ID:       1,
}
