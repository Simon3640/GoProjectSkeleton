package repositories

import (
	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/domain/models"
	db_models "gormgoskeleton/src/infrastructure/database/gormgoskeleton/models"

	"gorm.io/gorm"
)

type PasswordRepository struct {
	RepositoryBase[models.PasswordCreate, models.PasswordUpdate, models.Password, db_models.Password]
}

var _ contracts_repositories.IPasswordRepository = (*PasswordRepository)(nil)

type PasswordConverter struct{}

var _ ModelConverter[models.PasswordCreate, models.PasswordUpdate, models.Password, db_models.Password] = (*PasswordConverter)(nil)

func (uc *PasswordConverter) ToGormCreate(model models.PasswordCreate) *db_models.Password {
	return &db_models.Password{
		Hash:      model.Hash,
		ExpiresAt: model.ExpiresAt,
		IsActive:  model.IsActive,
		UserID:    model.UserID,
	}
}

func (uc *PasswordConverter) ToDomain(ormModel *db_models.Password) *models.Password {
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

func (uc *PasswordConverter) ToGormUpdate(model models.PasswordUpdate) *db_models.Password {
	password := &db_models.Password{}

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

func NewPasswordRepository(db *gorm.DB, logger contracts.ILoggerProvider) *PasswordRepository {
	return &PasswordRepository{
		RepositoryBase: RepositoryBase[
			models.PasswordCreate,
			models.PasswordUpdate,
			models.Password,
			db_models.Password,
		]{DB: db, modelConverter: &PasswordConverter{}, logger: logger},
	}
}

func (r *PasswordRepository) Create(model models.PasswordCreate) (*models.Password, error) {
	// start a transaction thay clean all previous passwords for the user setting is_active to false
	// and then create the new password
	tx := r.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := tx.Model(&db_models.Password{}).Where(
		"user_id = ? AND is_active = ?", model.UserID, true,
	).Updates(map[string]interface{}{"is_active": false}).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_entity := r.modelConverter.ToGormCreate(model)

	r.logger.Debug("Creating new password", _entity)

	if err := tx.Create(_entity).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return r.modelConverter.ToDomain(_entity), tx.Commit().Error
}
