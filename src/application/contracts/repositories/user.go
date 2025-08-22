package contracts_repositories

import "gormgoskeleton/src/domain/models"

type IUserRepository interface {
	IRepositoryBase[models.UserCreate, models.UserUpdate, models.User, models.UserInDB]
	CreateWithPassword(input models.UserAndPasswordCreate) (*models.User, error)
	GetUserWithRole(id uint) (*models.UserWithRole, error)
}
