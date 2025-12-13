package usecases

import (
	"testing"
	"time"

	statusmocks "github.com/simon3640/goprojectskeleton/src/application/modules/status/mocks"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	providersmocks "github.com/simon3640/goprojectskeleton/src/application/shared/mocks/providers"
	"github.com/simon3640/goprojectskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestStatusUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := app_context.NewVoidAppContext()

	testLogger := new(providersmocks.MockLoggerProvider)
	testStatusProvider := new(statusmocks.MockStatusProvider)
	testTime := time.Now()
	testStatusProvider.On(
		"Get",
		testTime,
	).Return(models.Status{
		AppName: "Test",
		Version: "1.0.0",
		Status:  "Testing",
		Date:    testTime.Format("2006-01-02 15:04:05"),
	})

	uc := NewGetStatusUseCase(testLogger, testStatusProvider)

	result_en := uc.Execute(ctx, locales.EN_US, testTime)
	result_es := uc.Execute(ctx, locales.ES_ES, testTime)

	assert.NotNil(result_en)
	assert.Equal(result_en.Data.Status == "Testing", true)
	assert.Equal(result_en.Data.AppName == "Test", true)
	assert.Equal(result_en.Data.Date == testTime.Format("2006-01-02 15:04:05"), true)
	assert.Equal(result_en.HasError(), false)

	// En
	assert.Equal(result_en.Details == messages.EnMessages[messages.MessageKeysInstance.APPLICATION_STATUS_OK], true)

	// Es
	assert.Equal(result_es.Details == messages.EsMessages[messages.MessageKeysInstance.APPLICATION_STATUS_OK], true)
}
