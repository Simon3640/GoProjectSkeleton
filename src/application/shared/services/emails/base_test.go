package email_service

import (
	"testing"

	"gormgoskeleton/src/application/shared/mocks"
	"gormgoskeleton/src/domain/models"

	"github.com/stretchr/testify/assert"
)

type structTestData struct {
	Something string
}

func TestEmailServiceBase(t *testing.T) {
	assert := assert.New(t)

	mockEmailProvider := new(mocks.MockEmailProvider)
	mockRenderProvider := new(mocks.MockRenderProvider[structTestData])

	test_data := structTestData{Something: "<Data>"}

	mockEmailProvider.On("SendEmail", "<Email>", "<Subject>", "<Data>").Return(nil)
	mockRenderProvider.On("Render", "<Template>", test_data).Return("<Data>", nil)

	emailServiceBase := &EmailServiceBase[structTestData]{
		Renderer: mockRenderProvider,
		Sender:   mockEmailProvider,
		template: "<Template>",
		subject:  "<Subject>",
	}

	assert.Nil(
		emailServiceBase.SendWithTemplate(
			test_data,
			models.User{
				UserBase: models.UserBase{
					Email: "<Email>",
				},
			},
		),
	)

}
