package email_service

import (
	"fmt"

	"gormgoskeleton/src/application/contracts"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/application/shared/templates"
	"gormgoskeleton/src/domain/models"
)

type EmailServiceBase[D any] struct {
	Renderer contracts.IRendererProvider[D]
	Sender   contracts.IEmailProvider
}

func (svc *EmailServiceBase[D]) SendWithTemplate(
	data D,
	user models.User,
	locale locales.LocaleTypeEnum,
) *application_errors.ApplicationError {
	appMessages := locales.NewLocale(locale)
	template := (settings.AppSettingsInstance.TemplatesPath + "/emails/" +
		templates.GetTemplate(
			locale,
			templates.TemplateKeysInstance.WelcomeEmail,
		))
	rendered, err := svc.Renderer.Render(template, data)
	if err != nil {
		return err
	}
	err = svc.Sender.SendEmail(user.Email,
		fmt.Sprintf(
			appMessages.Get(
				locale,
				messages.MessageKeysInstance.NEW_USER_WELCOME,
			),
			settings.AppSettingsInstance.AppName,
		), rendered)
	if err != nil {
		return err
	}
	return nil
}
