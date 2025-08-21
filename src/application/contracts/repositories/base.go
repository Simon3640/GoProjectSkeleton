package contracts_repositories

type IRepositoryBase[CreateDomainModel any, UpdateDomainModel any, DomainModel any, DBModel any] interface {
	Create(entity CreateDomainModel) (*DomainModel, error)
	GetByID(id uint) (*DomainModel, error)
	Update(id uint, entity UpdateDomainModel) (*DomainModel, error)
	Delete(id uint) error
	GetAll(payload *map[string]string, skip *int, limit *int) ([]DomainModel, error)
}
