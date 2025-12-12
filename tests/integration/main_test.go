package integrationtest

import (
	"os"
	"testing"

	infrastructure "github.com/simon3640/goprojectskeleton/src/infrastructure"
	database "github.com/simon3640/goprojectskeleton/src/infrastructure/databases/goprojectskeleton"
	"github.com/simon3640/goprojectskeleton/src/infrastructure/providers"
)

func TestMain(m *testing.M) {
	infrastructure.Initialize()
	providers.Logger.Info("Running tests setup...")
	providers.Logger.Info("Migrating DummyEntity for RepositoryBase tests...")
	if err := database.GoProjectSkeletondb.DB.AutoMigrate(&DummyEntity{}); err != nil {
		providers.Logger.Error("Error migrating DummyEntity", err)
		os.Exit(1)
	}
	providers.Logger.Info("Tests setup completed.")
	code := m.Run()
	os.Exit(code)
}
