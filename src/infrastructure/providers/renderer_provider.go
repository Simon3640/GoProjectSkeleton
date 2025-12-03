package providers

import (
	"bytes"
	html_template "html/template"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	application_errors "goprojectskeleton/src/application/shared/errors"
	"goprojectskeleton/src/application/shared/locales/messages"
	email_models "goprojectskeleton/src/application/shared/services/emails/models"
	"goprojectskeleton/src/application/shared/status"
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
