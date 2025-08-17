package contracts_repositories

import "gormgoskeleton/src/domain/models"

type IPasswordRepository interface {
	IRepositoryBase[models.PasswordCreate, models.PasswordUpdate, models.Password, models.PasswordInDB]
}
