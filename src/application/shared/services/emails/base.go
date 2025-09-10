package email_service

import (
	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	application_errors "gormgoskeleton/src/application/shared/errors"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/application/shared/templates"
)

type EmailServiceBase[D any] struct {
	Renderer contractsProviders.IRendererProvider[D]
	Sender   contractsProviders.IEmailProvider
}

func (svc *EmailServiceBase[D]) SetUp(
	renderer contractsProviders.IRendererProvider[D],
	sender contractsProviders.IEmailProvider,
) {
	svc.Renderer = renderer
	svc.Sender = sender
}

func (svc *EmailServiceBase[D]) SendWithTemplate(
	data D,
	userEmail string,
	locale locales.LocaleTypeEnum,
	templateKey templates.TemplateKeysEnum,
	subjectKey SubjectKeysEnum,
) *application_errors.ApplicationError {
	template := (settings.AppSettingsInstance.TemplatesPath + "emails/" +
		templates.GetTemplate(
			locale,
			templateKey,
		))
	rendered, err := svc.Renderer.Render(template, data)
	if err != nil {
		return err
	}
	err = svc.Sender.SendEmail(userEmail, GetSubject(locale, subjectKey), rendered)
	if err != nil {
		return err
	}
	return nil
}
