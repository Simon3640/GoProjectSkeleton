package providers

import (
	"bytes"
	html_template "html/template"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales/messages"
	email_models "gormgoskeleton/src/application/shared/services/emails/models"
	"gormgoskeleton/src/application/shared/status"
)

type RendererBase[T any] struct {
	Data T
}

var _ contractsProviders.IRendererProvider[any] = (*RendererBase[any])(nil)

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

type RenderResetPasswordEmail struct {
	RendererBase[email_models.ResetPasswordEmailData]
}

type RenderOTPEmail struct {
	RendererBase[email_models.OneTimePasswordEmailData]
}

var RenderNewUserEmailInstance *RenderNewUserEmail
var RenderResetPasswordEmailInstance *RenderResetPasswordEmail
var RenderOTPEmailInstance *RenderOTPEmail

func init() {
	RenderNewUserEmailInstance = &RenderNewUserEmail{}
	RenderResetPasswordEmailInstance = &RenderResetPasswordEmail{}
	RenderOTPEmailInstance = &RenderOTPEmail{}
}
