package email_service

import (
	"fmt"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/application/shared/templates"
	"gormgoskeleton/src/domain/models"
)

type EmailServiceBase[D any] struct {
	Renderer contracts_providers.IRendererProvider[D]
	Sender   contracts_providers.IEmailProvider
}

func (svc *EmailServiceBase[D]) SetUp(
	renderer contracts_providers.IRendererProvider[D],
	sender contracts_providers.IEmailProvider,
) {
	svc.Renderer = renderer
	svc.Sender = sender
}

func (svc *EmailServiceBase[D]) SendWithTemplate(
	data D,
	user models.User,
	locale locales.LocaleTypeEnum,
	templateKey templates.TemplateKeysEnum,
	subjectKey messages.MessageKeysEnum,
) *application_errors.ApplicationError {
	appMessages := locales.NewLocale(locale)
	template := (settings.AppSettingsInstance.TemplatesPath + "emails/" +
		templates.GetTemplate(
			locale,
			templateKey,
		))
	rendered, err := svc.Renderer.Render(template, data)
	if err != nil {
		return err
	}
	err = svc.Sender.SendEmail(user.Email,
		fmt.Sprintf(
			appMessages.Get(
				locale,
				subjectKey,
			),
			settings.AppSettingsInstance.AppName,
		), rendered)
	if err != nil {
		return err
	}
	return nil
}
