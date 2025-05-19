package repositories

import (
	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	RepositoryBase[models.UserCreate, models.UserUpdate, models.User, db_models.User]
}

var _ contracts_repositories.IUserRepository = (*UserRepository)(nil)

type UserConverter struct{}

var _ ModelConverter[models.UserCreate, models.UserUpdate, models.User, db_models.User] = (*UserConverter)(nil)

func (uc *UserConverter) toGormCreate(model models.UserCreate) *db_models.User {
	return &db_models.User{
		Name:   model.Name,
		Email:  model.Email,
		Phone:  model.Phone,
		Status: model.Status,
	}
}

func (uc *UserConverter) toDomain(ormModel *db_models.User) *models.User {
	return &models.User{
		ID: ormModel.ID,
		UserBase: models.UserBase{
			Name:   ormModel.Name,
			Email:  ormModel.Email,
			Phone:  ormModel.Phone,
			Status: ormModel.Status,
		},
	}
}

func (uc *UserConverter) toGormUpdate(model models.UserUpdate) *db_models.User {
	user := &db_models.User{}

	if model.Name != nil {
		user.Name = *model.Name
	}

	if model.Email != nil {
		user.Email = *model.Email
	}

	if model.Phone != nil {
		user.Phone = *model.Phone
	}

	if model.Status != nil {
		user.Status = *model.Status
	}
	user.ID = model.ID
	return user
}

func NewUserRepository(db *gorm.DB, logger contracts.ILoggerProvider) *UserRepository {
	return &UserRepository{
		RepositoryBase: RepositoryBase[
			models.UserCreate,
			models.UserUpdate,
			models.User,
			db_models.User,
		]{DB: db, modelConverter: &UserConverter{}, logger: logger},
	}
}
