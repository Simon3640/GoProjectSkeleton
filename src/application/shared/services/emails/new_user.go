package email_service

import (
	"fmt"
	"time"

	"gormgoskeleton/src/application/contracts"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/settings"
)

type NewUserEmailData struct {
	Name            string
	ActivationToken string
	Expiration      time.Time
}

type RegisterUserEmailService struct {
	EmailServiceBase[NewUserEmailData]
}

func NewRegisterUserEmailService(
	renderer contracts.IRenderProvider[NewUserEmailData],
	sender contracts.IEmailProvider,
	locale locales.LocaleTypeEnum,
	appMessages locales.Locale,
) *RegisterUserEmailService {
	return &RegisterUserEmailService{
		EmailServiceBase: EmailServiceBase[NewUserEmailData]{
			Renderer: renderer,
			Sender:   sender,
			template: TemplatesPath + GetTemplate(locale, TemplateKeysInstance.WelcomeEmail),
			subject: fmt.Sprintf(
				appMessages.Get(
					locale,
					messages.MessageKeysInstance.NEW_USER_WELCOME,
				),
				settings.AppSettingsInstance.AppName,
			),
		},
	}
}
