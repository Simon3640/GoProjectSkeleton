package repositories

import (
	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "goprojectskeleton/src/application/contracts/repositories"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/domain/models"
	dbModels "goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	RepositoryBase[dtos.UserCreate, dtos.UserUpdate, models.User, dbModels.User]
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
	passwordModel := dbModels.Password{
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
	var userWithRole dbModels.User
	if err := ur.DB.Preload("Role").First(&userWithRole, id).Error; err != nil {
		ur.logger.Debug("Error retrieving user with role", err)
		return nil, MapOrmError(err)
	}
	userStatus := models.UserStatus(userWithRole.Status)
	userWithRoleModel := models.UserWithRole{
		UserBase: models.UserBase{
			Name:     userWithRole.Name,
			Email:    userWithRole.Email,
			Phone:    userWithRole.Phone,
			Status:   &userStatus,
			RoleID:   userWithRole.RoleID,
			OTPLogin: userWithRole.OTPLogin,
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
	var user dbModels.User
	if err := ur.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		ur.logger.Debug("Error retrieving user by email or phone", err)
		return nil, MapOrmError(err)
	}
	userModel := ur.modelConverter.ToDomain(&user)
	return userModel, nil
}

var _ contracts_repositories.IUserRepository = (*UserRepository)(nil)

type UserConverter struct{}

var _ ModelConverter[dtos.UserCreate, dtos.UserUpdate, models.User, dbModels.User] = (*UserConverter)(nil)

func (uc *UserConverter) ToGormCreate(model dtos.UserCreate) *dbModels.User {
	return &dbModels.User{
		Name:     model.Name,
		Email:    model.Email,
		Phone:    model.Phone,
		Status:   string(*model.Status),
		RoleID:   model.RoleID,
		OTPLogin: model.OTPLogin,
	}
}

func (uc *UserConverter) ToDomain(ormModel *dbModels.User) *models.User {
	userStatus := models.UserStatus(ormModel.Status)
	return &models.User{
		DBBaseModel: models.DBBaseModel{
			ID:        ormModel.ID,
			CreatedAt: ormModel.CreatedAt,
			UpdatedAt: ormModel.UpdatedAt,
			DeletedAt: ormModel.DeletedAt.Time,
		},
		UserBase: models.UserBase{
			Name:     ormModel.Name,
			Email:    ormModel.Email,
			Phone:    ormModel.Phone,
			Status:   &userStatus,
			RoleID:   ormModel.RoleID,
			OTPLogin: ormModel.OTPLogin,
		},
	}
}

func (uc *UserConverter) ToGormUpdate(model dtos.UserUpdate) *dbModels.User {
	user := &dbModels.User{}

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
		user.Status = string(*model.Status)
	}
	if model.RoleID != nil {
		user.RoleID = *model.RoleID
	}
	if model.OTPLogin != nil {
		user.OTPLogin = *model.OTPLogin
	}
	user.ID = model.ID
	return user
}

func NewUserRepository(db *gorm.DB, logger contractsProviders.ILoggerProvider) *UserRepository {
	return &UserRepository{
		RepositoryBase: RepositoryBase[
			dtos.UserCreate,
			dtos.UserUpdate,
			models.User,
			dbModels.User,
		]{DB: db, modelConverter: &UserConverter{}, logger: logger},
	}
}
