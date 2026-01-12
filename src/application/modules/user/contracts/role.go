package usercontracts

import (
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"

	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

// IRoleRepository is the interface for the role repository
type IRoleRepository interface {
	contractsrepositories.IRepositoryBase[usermodels.RoleCreate, usermodels.RoleUpdate, usermodels.Role, usermodels.RoleInDB]
}
