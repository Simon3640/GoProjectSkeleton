package contracts_repositories

import "github.com/simon3640/goprojectskeleton/src/domain/models"

type IRoleRepository interface {
	IRepositoryBase[models.RoleCreate, models.RoleUpdate, models.Role, models.RoleInDB]
}
