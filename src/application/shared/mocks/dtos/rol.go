package dtomocks

import usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"

// TODO: Evaluate if Role can be created with DTOs

// UserRoleBase is a mock user role base for testing
var UserRoleBase = usermodels.RoleBase{
	Key:      "user",
	IsActive: true,
	Priority: 5,
}

// AdminRoleBase is a mock admin role base for testing
var AdminRoleBase = usermodels.RoleBase{
	Key:      "admin",
	IsActive: true,
	Priority: 1,
}

// UserRoleCreate is a mock user role create for testing
var UserRoleCreate = usermodels.RoleCreate{
	RoleBase: UserRoleBase,
}

// AdminRoleCreate is a mock admin role create for testing
var AdminRoleCreate = usermodels.RoleCreate{
	RoleBase: AdminRoleBase,
}

// UserRole is a mock user role for testing
var UserRole = usermodels.Role{
	RoleBase: UserRoleBase,
	ID:       2,
}

// AdminRole is a mock admin role for testing
var AdminRole = usermodels.Role{
	RoleBase: AdminRoleBase,
	ID:       1,
}
