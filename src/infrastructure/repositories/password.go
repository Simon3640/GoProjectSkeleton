package repositories

import (
	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	dbModels "github.com/simon3640/goprojectskeleton/src/infrastructure/database/goprojectskeleton/models"

	"gorm.io/gorm"
)

type PasswordRepository struct {
	RepositoryBase[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbModels.Password]
}

var _ contracts_repositories.IPasswordRepository = (*PasswordRepository)(nil)

type PasswordConverter struct{}

var _ ModelConverter[dtos.PasswordCreate, dtos.PasswordUpdate, models.Password, dbModels.Password] = (*PasswordConverter)(nil)

func (r *PasswordRepository) Create(model dtos.PasswordCreate) (*models.Password, *application_errors.ApplicationError) {
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
		r.logger.Debug("Error deactivating previous passwords", err)
		return nil, MapOrmError(err)
	}

	_entity := r.modelConverter.ToGormCreate(model)

	r.logger.Debug("Creating new password", _entity)

	if err := tx.Create(_entity).Error; err != nil {
		tx.Rollback()
		return nil, MapOrmError(err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, DefaultORMError
	}

	return r.modelConverter.ToDomain(_entity), nil
}

func (r *PasswordRepository) GetActivePassword(userEmail string) (*models.Password, *application_errors.ApplicationError) {
	var password dbModels.Password
	// Select the user by email, then take the first active password
	if err := r.DB.Joins(`JOIN "user" u ON u.id = password.user_id`).Where("u.email = ? AND password.is_active = ?", userEmail, true).First(&password).Error; err != nil {
		r.logger.Debug("Error retrieving active password", err)
		return nil, MapOrmError(err)
	}
	return r.modelConverter.ToDomain(&password), nil
}

func (uc *PasswordConverter) ToGormCreate(model dtos.PasswordCreate) *dbModels.Password {
	return &dbModels.Password{
		Hash:      model.Hash,
		ExpiresAt: model.ExpiresAt,
		IsActive:  model.IsActive,
		UserID:    model.UserID,
	}
}

func (uc *PasswordConverter) ToDomain(ormModel *dbModels.Password) *models.Password {
	return &models.Password{
		ID: ormModel.ID,
		PasswordBase: models.PasswordBase{
			UserID:    ormModel.UserID,
			Hash:      ormModel.Hash,
			ExpiresAt: ormModel.ExpiresAt,
			IsActive:  ormModel.IsActive,
		},
	}
}

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

func NewPasswordRepository(db *gorm.DB, logger contractsProviders.ILoggerProvider) *PasswordRepository {
	return &PasswordRepository{
		RepositoryBase: RepositoryBase[
			dtos.PasswordCreate,
			dtos.PasswordUpdate,
			models.Password,
			dbModels.Password,
		]{DB: db, modelConverter: &PasswordConverter{}, logger: logger},
	}
}
