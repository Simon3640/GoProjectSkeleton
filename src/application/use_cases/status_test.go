package usecases

import (
	"context"
	"testing"
	"time"

	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/use_cases/mocks"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

func TestStatusUseCase(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()

	testLogger := new(mocks.MockLoggerProvider)
	testStatusProvider := new(mocks.MockStatusProvider)
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
