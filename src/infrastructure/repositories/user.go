package repositories

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	RepositoryBase[dtos.UserCreate, dtos.UserUpdate, models.User, db_models.User]
}

func (ur *UserRepository) CreateWithPassword(input dtos.UserAndPasswordCreate) (*models.User, *application_errors.ApplicationError) {
	// Convert input to 2 models UserCreate and UserInDB
	userCreate := ur.modelConverter.ToGormCreate(input.UserCreate)
	tx := ur.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(userCreate).Error; err != nil {
		ur.logger.Debug("Error creating user", err)
		tx.Rollback()
		return nil, MapOrmError(err)
	}
	userInDB := ur.modelConverter.ToDomain(userCreate)
	// Create password
	passwordCreate := dtos.PasswordCreate{
		PasswordBase: models.PasswordBase{
			UserID:   userInDB.ID,
			Hash:     input.Password,
			IsActive: true,
		},
	}
	passwordCreate.SetDefaultExpiresAt()
	passwordModel := db_models.Password{
		Hash:      passwordCreate.Hash,
		ExpiresAt: passwordCreate.ExpiresAt,
		IsActive:  passwordCreate.IsActive,
		UserID:    passwordCreate.UserID,
	}
	if err := tx.Create(&passwordModel).Error; err != nil {
		ur.logger.Debug("Error creating password for user", err)
		tx.Rollback()
		return nil, MapOrmError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, DefaultORMError
	}

	return userInDB, nil
}

func (ur *UserRepository) GetUserWithRole(id uint) (*models.UserWithRole, *application_errors.ApplicationError) {
	var userWithRole db_models.User
	if err := ur.DB.Preload("Role").First(&userWithRole, id).Error; err != nil {
		ur.logger.Debug("Error retrieving user with role", err)
		return nil, MapOrmError(err)
	}
	userWithRoleModel := models.UserWithRole{
		UserBase: models.UserBase{
			Name:   userWithRole.Name,
			Email:  userWithRole.Email,
			Phone:  userWithRole.Phone,
			Status: userWithRole.Status,
			RoleID: userWithRole.RoleID,
		},
		ID: userWithRole.ID,
	}
	userWithRoleModel.SetRole(models.Role{
		ID: userWithRole.Role.ID,
		RoleBase: models.RoleBase{
			Key:      userWithRole.Role.Key,
			IsActive: userWithRole.Role.IsActive,
			Priority: userWithRole.Role.Priority,
		},
	})
	return &userWithRoleModel, nil
}

func (ur *UserRepository) GetByEmailOrPhone(emailOrPhone string) (*models.User, *application_errors.ApplicationError) {
	var user db_models.User
	if err := ur.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		ur.logger.Debug("Error retrieving user by email or phone", err)
		return nil, MapOrmError(err)
	}
	userModel := ur.modelConverter.ToDomain(&user)
	return userModel, nil
}

var _ contracts_repositories.IUserRepository = (*UserRepository)(nil)

type UserConverter struct{}

var _ ModelConverter[dtos.UserCreate, dtos.UserUpdate, models.User, db_models.User] = (*UserConverter)(nil)

func (uc *UserConverter) ToGormCreate(model dtos.UserCreate) *db_models.User {
	return &db_models.User{
		Name:   model.Name,
		Email:  model.Email,
		Phone:  model.Phone,
		Status: model.Status,
		RoleID: model.RoleID,
	}
}

func (uc *UserConverter) ToDomain(ormModel *db_models.User) *models.User {
	return &models.User{
		DBBaseModel: models.DBBaseModel{
			ID:        ormModel.ID,
			CreatedAt: ormModel.CreatedAt,
			UpdatedAt: ormModel.UpdatedAt,
			DeletedAt: ormModel.DeletedAt.Time,
		},
		UserBase: models.UserBase{
			Name:   ormModel.Name,
			Email:  ormModel.Email,
			Phone:  ormModel.Phone,
			Status: ormModel.Status,
			RoleID: ormModel.RoleID,
		},
	}
}

func (uc *UserConverter) ToGormUpdate(model dtos.UserUpdate) *db_models.User {
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
	if model.RoleID != nil {
		user.RoleID = *model.RoleID
	}
	user.ID = model.ID
	return user
}

func NewUserRepository(db *gorm.DB, logger contracts_providers.ILoggerProvider) *UserRepository {
	return &UserRepository{
		RepositoryBase: RepositoryBase[
			dtos.UserCreate,
			dtos.UserUpdate,
			models.User,
			db_models.User,
		]{DB: db, modelConverter: &UserConverter{}, logger: logger},
	}
}
