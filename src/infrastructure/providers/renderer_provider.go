package providers

import (
	"bytes"
	html_template "html/template"

	"gormgoskeleton/src/application/contracts"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales/messages"
	email_models "gormgoskeleton/src/application/shared/services/emails/models"
	"gormgoskeleton/src/application/shared/status"
)

type RendererBase[T any] struct {
	Data T
}

var _ contracts.IRendererProvider[any] = (*RendererBase[any])(nil)

func (r RendererBase[T]) Render(templatePath string, data T) (string, *application_errors.ApplicationError) {
	tmpl, err := html_template.ParseFiles(templatePath)
	if err != nil {
		return "", application_errors.NewApplicationError(
			status.InternalError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", application_errors.NewApplicationError(
			status.InternalError,
			messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			err.Error(),
		)
	}

	return rendered.String(), nil
}

type RenderNewUserEmail struct {
	RendererBase[email_models.NewUserEmailData]
}

var RenderNewUserEmailInstance *RenderNewUserEmail

func init() {
	RenderNewUserEmailInstance = &RenderNewUserEmail{}
}
