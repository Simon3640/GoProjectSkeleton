package integrationtest

import (
	"os"
	"testing"

	infrastructure "gormgoskeleton/src/infrastructure"
	database "gormgoskeleton/src/infrastructure/database/gormgoskeleton"
	"gormgoskeleton/src/infrastructure/providers"
)

func TestMain(m *testing.M) {
	infrastructure.Initialize()
	providers.Logger.Info("Running tests setup...")
	providers.Logger.Info("Migrating DummyEntity for RepositoryBase tests...")
	database.DB.AutoMigrate(&DummyEntity{})
	providers.Logger.Info("Tests setup completed.")
	code := m.Run()
	os.Exit(code)
}
