package repositories

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"

	// Decouple modules dependencies
	passworddtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"

	"gorm.io/gorm"
)

// UserRepository is the repository for the user model
type UserRepository struct {
	RepositoryBase[userdtos.UserCreate, userdtos.UserUpdate, models.User, dbModels.User]
}

// CreateWithPassword creates a new user with a password
func (ur *UserRepository) CreateWithPassword(input userdtos.UserAndPasswordCreate) (*models.User, *application_errors.ApplicationError) {
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
	passwordCreate := passworddtos.PasswordCreate{
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

// GetUserWithRole gets a user with their role
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

// GetByEmailOrPhone gets a user by email or phone
func (ur *UserRepository) GetByEmailOrPhone(emailOrPhone string) (*models.User, *application_errors.ApplicationError) {
	var user dbModels.User
	if err := ur.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		ur.logger.Debug("Error retrieving user by email or phone", err)
		return nil, MapOrmError(err)
	}
	userModel := ur.modelConverter.ToDomain(&user)
	return userModel, nil
}

var _ usercontracts.IUserRepository = (*UserRepository)(nil)

// UserConverter is the converter for the user model
type UserConverter struct{}

var _ ModelConverter[userdtos.UserCreate, userdtos.UserUpdate, models.User, dbModels.User] = (*UserConverter)(nil)

// ToGormCreate converts a user create model to a user gorm model
func (uc *UserConverter) ToGormCreate(model userdtos.UserCreate) *dbModels.User {
	return &dbModels.User{
		Name:     model.Name,
		Email:    model.Email,
		Phone:    model.Phone,
		Status:   string(*model.Status),
		RoleID:   model.RoleID,
		OTPLogin: model.OTPLogin,
	}
}

// ToDomain converts a user gorm model to a user domain model
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

// ToGormUpdate converts a user update model to a user gorm model
func (uc *UserConverter) ToGormUpdate(model userdtos.UserUpdate) *dbModels.User {
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

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB, logger contractsProviders.ILoggerProvider) *UserRepository {
	return &UserRepository{
		RepositoryBase: RepositoryBase[
			userdtos.UserCreate,
			userdtos.UserUpdate,
			models.User,
			dbModels.User,
		]{DB: db, modelConverter: &UserConverter{}, logger: logger},
	}
}
