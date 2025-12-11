package repositories

// ModelConverter is the interface for the model converter
type ModelConverter[ModelCreate any, ModelUpdate any, Model any, DBModel any] interface {
	ToGormCreate(model ModelCreate) *DBModel
	ToDomain(ormModel *DBModel) *Model
	ToGormUpdate(model ModelUpdate) *DBModel
}
