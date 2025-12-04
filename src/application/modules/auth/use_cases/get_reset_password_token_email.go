package authusecases

import (
	"context"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/locales/messages"
	email_service "goprojectskeleton/src/application/shared/services/emails"
	email_models "goprojectskeleton/src/application/shared/services/emails/models"
	"goprojectskeleton/src/application/shared/settings"
	"goprojectskeleton/src/application/shared/status"
	"goprojectskeleton/src/application/shared/templates"
	usecase "goprojectskeleton/src/application/shared/use_case"
)

type GetResetPasswordSendEmailUseCase struct {
	appMessages *locales.Locale
	log         contractsProviders.ILoggerProvider
	locale      locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[dtos.OneTimeTokenUser, bool] = (*GetResetPasswordSendEmailUseCase)(nil)

func (uc *GetResetPasswordSendEmailUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *GetResetPasswordSendEmailUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.OneTimeTokenUser,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)

	newUserEmailData := email_models.ResetPasswordEmailData{
		Name:              input.User.Name,
		ResetLink:         input.BuildURL(settings.AppSettingsInstance.FrontendResetPasswordURL),
		ExpirationMinutes: settings.AppSettingsInstance.OneTimeTokenPasswordTTL,
		AppName:           settings.AppSettingsInstance.AppName,
		SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
	}

	if err := email_service.ResetPasswordEmailServiceInstance.SendWithTemplate(
		newUserEmailData,
		input.User.Email,
		locale,
		templates.TemplateKeysInstance.PasswordResetEmail,
		email_service.SubjectKeysInstance.PasswordResetEmail,
	); err != nil {
		uc.log.Error("Error sending email", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return result
	}
	result.SetData(
		status.Success,
		true,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.RESET_PASSWORD_EMAIL_SENT_SUCCESSFULLY,
		),
	)
	return result
}

func NewGetResetPasswordSendEmailUseCase(
	log contractsProviders.ILoggerProvider,
) *GetResetPasswordSendEmailUseCase {
	return &GetResetPasswordSendEmailUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
	}
}
