package contracts_repositories

import "goprojectskeleton/src/domain/models"

type IRoleRepository interface {
	IRepositoryBase[models.RoleCreate, models.RoleUpdate, models.Role, models.RoleInDB]
}
