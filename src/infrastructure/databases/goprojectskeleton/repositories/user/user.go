// Package userrepositories contains the repository for the user model
package userrepositories

import (
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"

	// Decouple modules dependencies
	passworddtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	passwordmodels "github.com/simon3640/goprojectskeleton/src/domain/password/models"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
	dbmodels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	reposhared "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/shared"

	"gorm.io/gorm"
)

// UserRepository is the repository for the user model
type UserRepository struct {
	reposhared.RepositoryBase[userdtos.UserCreate, userdtos.UserUpdate, usermodels.User, dbmodels.User]
}

// CreateWithPassword creates a new user with a password
func (ur *UserRepository) CreateWithPassword(input userdtos.UserAndPasswordCreate) (*usermodels.User, *applicationerrors.ApplicationError) {
	// Convert input to 2 models UserCreate and UserInDB
	userCreate := ur.ModelConverter.ToGormCreate(input.UserCreate)
	tx := ur.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Create(userCreate).Error; err != nil {
		ur.Logger.Debug("Error creating user", err)
		tx.Rollback()
		return nil, reposhared.MapOrmError(err)
	}
	userInDB := ur.ModelConverter.ToDomain(userCreate)
	// Create password
	passwordCreate := passworddtos.PasswordCreate{
		PasswordBase: passwordmodels.PasswordBase{
			UserID:   userInDB.ID,
			Hash:     input.Password,
			IsActive: true,
		},
	}
	passwordCreate.SetDefaultExpiresAt()
	passwordModel := dbmodels.Password{
		Hash:      passwordCreate.Hash,
		ExpiresAt: passwordCreate.ExpiresAt,
		IsActive:  passwordCreate.IsActive,
		UserID:    passwordCreate.UserID,
	}
	if err := tx.Create(&passwordModel).Error; err != nil {
		ur.Logger.Debug("Error creating password for user", err)
		tx.Rollback()
		return nil, reposhared.MapOrmError(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, reposhared.DefaultORMError
	}

	return userInDB, nil
}

// GetUserWithRole gets a user with their role
func (ur *UserRepository) GetUserWithRole(id uint) (*usermodels.UserWithRole, *applicationerrors.ApplicationError) {
	var userWithRole dbmodels.User
	if err := ur.DB.Preload("Role").First(&userWithRole, id).Error; err != nil {
		ur.Logger.Debug("Error retrieving user with role", err)
		return nil, reposhared.MapOrmError(err)
	}
	userStatus := usermodels.UserStatus(userWithRole.Status)
	userWithRoleModel := usermodels.UserWithRole{
		UserBase: usermodels.UserBase{
			Name:     userWithRole.Name,
			Email:    userWithRole.Email,
			Phone:    userWithRole.Phone,
			Status:   &userStatus,
			RoleID:   userWithRole.RoleID,
			OTPLogin: userWithRole.OTPLogin,
		},
		ID: userWithRole.ID,
	}
	userWithRoleModel.SetRole(usermodels.Role{
		ID: userWithRole.Role.ID,
		RoleBase: usermodels.RoleBase{
			Key:      userWithRole.Role.Key,
			IsActive: userWithRole.Role.IsActive,
			Priority: userWithRole.Role.Priority,
		},
	})
	return &userWithRoleModel, nil
}

// GetByEmailOrPhone gets a user by email or phone
func (ur *UserRepository) GetByEmailOrPhone(emailOrPhone string) (*usermodels.User, *applicationerrors.ApplicationError) {
	var user dbmodels.User
	if err := ur.DB.Where("email = ? OR phone = ?", emailOrPhone, emailOrPhone).First(&user).Error; err != nil {
		ur.Logger.Debug("Error retrieving user by email or phone", err)
		return nil, reposhared.MapOrmError(err)
	}
	userModel := ur.ModelConverter.ToDomain(&user)
	return userModel, nil
}

var _ usercontracts.IUserRepository = (*UserRepository)(nil)

// UserConverter is the converter for the user model
type UserConverter struct{}

var _ reposhared.ModelConverter[userdtos.UserCreate, userdtos.UserUpdate, usermodels.User, dbmodels.User] = (*UserConverter)(nil)

// ToGormCreate converts a user create model to a user gorm model
func (uc *UserConverter) ToGormCreate(model userdtos.UserCreate) *dbmodels.User {
	return &dbmodels.User{
		Name:     model.Name,
		Email:    model.Email,
		Phone:    model.Phone,
		Status:   string(*model.Status),
		RoleID:   model.RoleID,
		OTPLogin: model.OTPLogin,
	}
}

// ToDomain converts a user gorm model to a user domain model
func (uc *UserConverter) ToDomain(ormModel *dbmodels.User) *usermodels.User {
	userStatus := usermodels.UserStatus(ormModel.Status)
	return &usermodels.User{
		DBBaseModel: sharedmodels.DBBaseModel{
			ID:        ormModel.ID,
			CreatedAt: ormModel.CreatedAt,
			UpdatedAt: ormModel.UpdatedAt,
			DeletedAt: ormModel.DeletedAt.Time,
		},
		UserBase: usermodels.UserBase{
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
func (uc *UserConverter) ToGormUpdate(model userdtos.UserUpdate) *dbmodels.User {
	user := &dbmodels.User{}

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
func NewUserRepository(db *gorm.DB, logger contractsproviders.ILoggerProvider) *UserRepository {
	return &UserRepository{
		RepositoryBase: reposhared.RepositoryBase[
			userdtos.UserCreate,
			userdtos.UserUpdate,
			usermodels.User,
			dbmodels.User,
		]{
			DB:             db,
			ModelConverter: &UserConverter{},
			Logger:         logger,
		},
	}
}
