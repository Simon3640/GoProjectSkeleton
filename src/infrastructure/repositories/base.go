package repositories

import (
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"

	"gorm.io/gorm"
)

type RepositoryBase[Model any] struct {
	DB *gorm.DB
}

var _ contracts_repositories.IRepositoryBase[any] = (*RepositoryBase[any])(nil)

func (rb *RepositoryBase[Model]) Create(entity Model) error {
	return rb.DB.Create(&entity).Error
}

func (rb *RepositoryBase[Model]) GetByID(id int) (*Model, error) {
	var entity Model
	err := rb.DB.First(&entity, id).Error
	return &entity, err
}

func (rb *RepositoryBase[Model]) Update(entity Model) error {
	return rb.DB.Save(&entity).Error
}

func (rb *RepositoryBase[Model]) Delete(id int) error {
	return rb.DB.Delete(new(Model), id).Error
}

func (rb *RepositoryBase[Model]) GetAll() ([]Model, error) {
	var entities []Model
	err := rb.DB.Find(&entities).Error
	return entities, err
}

// func (rb *RepositoryBase[Model]) toORMModel(domainModel Model) DBModel {
// 	return domainModel.(DBModel)
// }

// func (rb *RepositoryBase[Model]) toDomainModel(ormModel DBModel) Model {
// 	return ormModel.(Model)
// }
