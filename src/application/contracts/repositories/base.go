package contracts_repositories

type IRepositoryBase[DomainModel any] interface {
	Create(entity DomainModel) error
	GetByID(id int) (*DomainModel, error)
	Update(entity DomainModel) error
	Delete(id int) error
	GetAll() ([]DomainModel, error)
}
