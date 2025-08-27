package repositories

import (
	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	application_errors "gormgoskeleton/src/application/shared/errors"

	"gorm.io/gorm"
)

type RepositoryBase[CreateModel any, UpdateModel any, Model any, DBModel any] struct {
	DB             *gorm.DB
	logger         contracts_providers.ILoggerProvider
	modelConverter ModelConverter[CreateModel, UpdateModel, Model, DBModel]
}

var _ contracts_repositories.IRepositoryBase[any, any, any, any] = (*RepositoryBase[any, any, any, any])(nil)

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Create(entity CreateModel) (*Model, *application_errors.ApplicationError) {
	// Convertir a modelo de GORM
	_entity := rb.modelConverter.ToGormCreate(entity)
	if err := rb.DB.Create(_entity).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error creating entity", appErr.ToError())
		return nil, appErr
	}
	rb.logger.Debug("Entity created successfully", _entity)
	// Convertir de nuevo a modelo de dominio
	return rb.modelConverter.ToDomain(_entity), nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetByID(id uint) (*Model, *application_errors.ApplicationError) {
	var entity DBModel
	if err := rb.DB.First(&entity, id).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error retrieving entity", appErr.ToError())
		return nil, appErr
	}
	rb.logger.Debug("Entity retrieved successfully", entity)
	return rb.modelConverter.ToDomain(&entity), nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Update(id uint, entity UpdateModel) (*Model, *application_errors.ApplicationError) {
	updateData := rb.modelConverter.ToGormUpdate(entity)

	if err := rb.DB.Model(new(DBModel)).Where("id = ?", id).Updates(updateData).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error updating entity", appErr.ToError())
		return nil, appErr
	}
	rb.logger.Debug("Entity updated successfully", updateData)
	updatedEntity, _ := rb.GetByID(id)
	return updatedEntity, nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Delete(id uint) *application_errors.ApplicationError {
	if err := rb.DB.Delete(new(DBModel), id).Error; err != nil {
		appErr := MapOrmError(err)
		rb.logger.Debug("Error deleting entity", appErr.ToError())
		return appErr
	}
	return nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetAll(payload *map[string]string, skip *int, limit *int) ([]Model, *application_errors.ApplicationError) {
	var entities []DBModel
	// Apply filters from payload
	if payload != nil {
		for key, value := range *payload {
			// Assuming the key is a column name and value is the value to filter by
			rb.DB = rb.DB.Where(key+" = ?", value)
		}
	}
	// Apply pagination
	if skip != nil && *skip > 0 {
		rb.DB = rb.DB.Offset(*skip)
	}
	if limit != nil && *limit > 0 {
		rb.DB = rb.DB.Limit(*limit)
	}
	// Execute the query
	if err := rb.DB.Find(&entities).Error; err != nil {
		if appErr := MapOrmError(err); appErr != nil {
			rb.logger.Debug("Error retrieving entities", appErr.ToError())
			return nil, appErr
		}
		return nil, DefaultORMError
	}

	result := make([]Model, len(entities))

	for i, entity := range entities {
		result[i] = *rb.modelConverter.ToDomain(&entity)
	}

	return result, nil
}
