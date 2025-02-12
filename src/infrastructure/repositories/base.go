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
	_entity := rb.modelConverter.toGormCreate(entity)
	err := rb.DB.Create(_entity).Error
	if err != nil {
		return nil, err
	}

	// Convertir de nuevo a modelo de dominio
	return rb.modelConverter.toDomain(_entity), nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetByID(id int) (*Model, error) {
	var entity DBModel
	err := rb.DB.First(&entity, id).Error
	return rb.modelConverter.toDomain(&entity), err
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Update(id int, entity UpdateModel) (*Model, error) {
	updateData := rb.modelConverter.toGormUpdate(entity)
	err := rb.DB.Model(new(DBModel)).Where("id = ?", id).Updates(updateData).Error

	if err != nil {
		return nil, err
	}

	updatedEntity, _ := rb.GetByID(id)
	return updatedEntity, nil
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) Delete(id int) error {
	return rb.DB.Delete(new(DBModel), id).Error
}

func (rb *RepositoryBase[CreateModel, UpdateModel, Model, DBModel]) GetAll() ([]Model, error) {
	var entities []Model
	err := rb.DB.Find(&entities).Error
	return entities, err
}
