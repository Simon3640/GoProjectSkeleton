package repositories

import (
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	RepositoryBase[models.UserCreate, models.UserUpdate, models.User, db_models.User]
}

// var _ contracts_repositories.IUserRepository = (*UserRepository)(nil)

type UserConverter struct{}

var _ ModelConverter[models.UserCreate, models.User, db_models.User] = (*UserConverter)(nil)

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

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		RepositoryBase: RepositoryBase[
			models.UserCreate,
			models.UserUpdate,
			models.User,
			db_models.User,
		]{DB: db, modelConverter: &UserConverter{}},
	}
}
