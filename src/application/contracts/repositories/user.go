package contracts_repositories

import "gormgoskeleton/src/domain/models"

type IUserRepository interface {
	IRepositoryBase[models.UserCreate, models.UserUpdate, models.User, models.UserInDB]
}
