package contracts_repositories

type IRepositoryBase[CreateDomainModel any, UpdateDomainModel any, DomainModel any, DBModel any] interface {
	Create(entity CreateDomainModel) (*DomainModel, error)
	GetByID(id int) (*DomainModel, error)
	Update(id int, entity UpdateDomainModel) (*DomainModel, error)
	Delete(id int) error
	GetAll(payload *map[string]string, skip *int, limit *int) ([]DomainModel, error)
}
