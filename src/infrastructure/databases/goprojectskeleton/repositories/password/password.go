// Package passwordrepositories contains the repository for the password model
package passwordrepositories

import (
	contractproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	passwordcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/password/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/password/dtos"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	passwordmodels "github.com/simon3640/goprojectskeleton/src/domain/password/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/models"
	reposhared "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton/repositories/shared"

	"gorm.io/gorm"
)

// PasswordRepository is the repository for the password model
type PasswordRepository struct {
	reposhared.RepositoryBase[dtos.PasswordCreate, dtos.PasswordUpdate, passwordmodels.Password, dbModels.Password]
}

var _ passwordcontracts.IPasswordRepository = (*PasswordRepository)(nil)

// PasswordConverter is the converter for the password model
type PasswordConverter struct{}

var _ reposhared.ModelConverter[dtos.PasswordCreate, dtos.PasswordUpdate, passwordmodels.Password, dbModels.Password] = (*PasswordConverter)(nil)

// Create creates a new password in transaction that cleans all previous passwords for the user setting is_active to false
func (r *PasswordRepository) Create(model dtos.PasswordCreate) (*passwordmodels.Password, *applicationerrors.ApplicationError) {
	// start a transaction thay clean all previous passwords for the user setting is_active to false
	// and then create the new password
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Model(&dbModels.Password{}).Where(
		"user_id = ? AND is_active = ?", model.UserID, true,
	).Updates(map[string]interface{}{"is_active": false}).Error

	if err != nil {
		tx.Rollback()
		r.Logger.Debug("Error deactivating previous passwords", err)
		return nil, reposhared.MapOrmError(err)
	}

	_entity := r.ModelConverter.ToGormCreate(model)

	r.Logger.Debug("Creating new password", _entity)

	if err := tx.Create(_entity).Error; err != nil {
		tx.Rollback()
		return nil, reposhared.MapOrmError(err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, reposhared.DefaultORMError
	}

	return r.ModelConverter.ToDomain(_entity), nil
}

// GetActivePassword retrieves the active password for a user by email
func (r *PasswordRepository) GetActivePassword(userEmail string) (*passwordmodels.Password, *applicationerrors.ApplicationError) {
	var password dbModels.Password
	// Select the user by email, then take the first active password
	if err := r.DB.Joins(`JOIN "user" u ON u.id = password.user_id`).Where("u.email = ? AND password.is_active = ?", userEmail, true).First(&password).Error; err != nil {
		r.Logger.Debug("Error retrieving active password", err)
		return nil, reposhared.MapOrmError(err)
	}
	return r.ModelConverter.ToDomain(&password), nil
}

// ToGormCreate converts a password create model to a password gorm model
func (uc *PasswordConverter) ToGormCreate(model dtos.PasswordCreate) *dbModels.Password {
	return &dbModels.Password{
		Hash:      model.Hash,
		ExpiresAt: model.ExpiresAt,
		IsActive:  model.IsActive,
		UserID:    model.UserID,
	}
}

// ToDomain converts a password gorm model to a password domain model
func (uc *PasswordConverter) ToDomain(ormModel *dbModels.Password) *passwordmodels.Password {
	return &passwordmodels.Password{
		ID: ormModel.ID,
		PasswordBase: passwordmodels.PasswordBase{
			UserID:    ormModel.UserID,
			Hash:      ormModel.Hash,
			ExpiresAt: ormModel.ExpiresAt,
			IsActive:  ormModel.IsActive,
		},
	}
}

// ToGormUpdate converts a password update model to a password gorm model
func (uc *PasswordConverter) ToGormUpdate(model dtos.PasswordUpdate) *dbModels.Password {
	password := &dbModels.Password{}

	if model.Hash != nil {
		password.Hash = *model.Hash
	}
	if model.ExpiresAt != nil {
		password.ExpiresAt = model.ExpiresAt
	}
	if model.IsActive != nil {
		password.IsActive = *model.IsActive
	}

	password.ID = model.ID
	return password
}

// NewPasswordRepository creates a new password repository
func NewPasswordRepository(db *gorm.DB, logger contractproviders.ILoggerProvider) *PasswordRepository {
	return &PasswordRepository{
		RepositoryBase: reposhared.RepositoryBase[
			dtos.PasswordCreate,
			dtos.PasswordUpdate,
			passwordmodels.Password,
			dbModels.Password,
		]{
			DB:             db,
			ModelConverter: &PasswordConverter{},
			Logger:         logger,
		},
	}
}
