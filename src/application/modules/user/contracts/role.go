package usercontracts

import (
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"

	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// IRoleRepository is the interface for the role repository
type IRoleRepository interface {
	contractsrepositories.IRepositoryBase[models.RoleCreate, models.RoleUpdate, models.Role, models.RoleInDB]
}
