package integrationtest

import (
	"testing"

	"goprojectskeleton/src/application/shared/status"
	domain_utils "goprojectskeleton/src/domain/utils"
	database "goprojectskeleton/src/infrastructure/database/goprojectskeleton"
	"goprojectskeleton/src/infrastructure/providers"
	"goprojectskeleton/src/infrastructure/repositories"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// DummyEntity para probar RepositoryBase
type DummyEntity struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255;not null;unique"`
}

type DummyCreate struct {
	Name string
}

type DummyUpdate struct {
	Name string
}

type DummyDomain struct {
	ID   uint
	Name string
}

type DummyModelConverter struct{}

func (mc DummyModelConverter) ToGormCreate(model DummyCreate) *DummyEntity {
	return &DummyEntity{Name: model.Name}
}

func (mc DummyModelConverter) ToGormUpdate(model DummyUpdate) *DummyEntity {
	return &DummyEntity{Name: model.Name}
}

func (mc DummyModelConverter) ToDomain(entity *DummyEntity) *DummyDomain {
	return &DummyDomain{ID: entity.ID, Name: entity.Name}
}

type DummyRepository struct {
	repositories.RepositoryBase[DummyCreate, DummyUpdate, DummyDomain, DummyEntity]
}

func NewDummyRepository(db *gorm.DB) *DummyRepository {
	return &DummyRepository{
		RepositoryBase: repositories.SetUpRepositoryBase[DummyCreate, DummyUpdate, DummyDomain, DummyEntity](
			db,
			providers.Logger,
			DummyModelConverter{},
		),
	}
}

// Test Suite
func setupTestRepo() *DummyRepository {
	dummyRepo := NewDummyRepository(database.GoProjectSkeletondb.DB)
	return dummyRepo
}

func TestRepositoryBase_CRUD(t *testing.T) {
	repo := setupTestRepo()
	assert := assert.New(t)

	// Create
	entity, appErr := repo.Create(DummyCreate{Name: "TestName"})
	assert.Nil(appErr)
	assert.NotNil(entity)
	assert.Equal("TestName", entity.Name)

	// GetByID
	got, appErr := repo.GetByID(entity.ID)
	assert.Nil(appErr)
	assert.Equal(entity.ID, got.ID)

	// Update
	updated, appErr := repo.Update(entity.ID, DummyUpdate{Name: "UpdatedName"})
	assert.Nil(appErr)
	assert.Equal("UpdatedName", updated.Name)

	// GetAll
	payload := domain_utils.NewQueryPayloadBuilder[DummyDomain](
		[]string{"Id:asc"},
		[]string{"Name:eq:UpdatedName"},
		nil,
		nil,
	)
	results, total, appErr := repo.GetAll(&payload, 0, 10)
	assert.Nil(appErr)
	assert.Equal(int64(1), total)
	assert.Len(results, 1)

	// SoftDelete
	appErr = repo.SoftDelete(entity.ID)
	assert.Nil(appErr)

	// Should not find it anymore
	_, appErr = repo.GetByID(entity.ID)
	assert.NotNil(appErr)
	assert.Equal(status.NotFound, appErr.Code)

	// Create again to test hard delete
	entity2, _ := repo.Create(DummyCreate{Name: "ToDelete"})
	appErr = repo.Delete(entity2.ID)
	assert.Nil(appErr)
	_, appErr = repo.GetByID(entity2.ID)
	assert.NotNil(appErr)
}
