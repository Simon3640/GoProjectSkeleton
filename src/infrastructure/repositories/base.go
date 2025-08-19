package repositories

import (
	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"

	"gorm.io/gorm"
)

type RepositoryBase[CreateModel any, UpdateModel any, Model any, DBModel any] struct {
	DB             *gorm.DB
	logger         contracts.ILoggerProvider
	modelConverter ModelConverter[CreateModel, UpdateModel, Model, DBModel]
}

var _ contracts_repositories.IRepositoryBase[any, any, any, any] = (*RepositoryBase[any, any, any, any])(nil)

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Create(entity CreateModel) (*Model, error) {
	// Convertir a modelo de GORM
	_entity := rb.modelConverter.ToGormCreate(entity)
	err := rb.DB.Create(_entity).Error
	if err != nil {
		return nil, err
	}
	rb.logger.Debug("Entity created successfully", _entity)
	// Convertir de nuevo a modelo de dominio
	return rb.modelConverter.ToDomain(_entity), nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetByID(id uint) (*Model, error) {
	var entity DBModel
	err := rb.DB.First(&entity, id).Error
	rb.logger.Debug("Entity retrieved successfully", entity)
	return rb.modelConverter.ToDomain(&entity), err
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Update(id uint, entity UpdateModel) (*Model, error) {
	updateData := rb.modelConverter.ToGormUpdate(entity)
	err := rb.DB.Model(new(DBModel)).Where("id = ?", id).Updates(updateData).Error

	if err != nil {
		return nil, err
	}
	rb.logger.Debug("Entity updated successfully", updateData)
	updatedEntity, _ := rb.GetByID(id)
	return updatedEntity, nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Delete(id uint) error {
	err := rb.DB.Delete(new(DBModel), id).Error
	rb.logger.Debug("Entity deleted", id)
	return err
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetAll(payload *map[string]string, skip *int, limit *int) ([]Model, error) {
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
	err := rb.DB.Find(&entities).Error
	if err != nil {
		return nil, err
	}

	result := make([]Model, len(entities))

	for i, entity := range entities {
		result[i] = *rb.modelConverter.ToDomain(&entity)
	}

	return result, nil
}
