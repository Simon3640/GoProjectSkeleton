package authusecases

import (
	"context"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	email_service "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	email_models "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/application/shared/templates"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
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

// Execute sends a reset password email to the user
func (uc *GetResetPasswordSendEmailUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.OneTimeTokenUser,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)

	emailData := uc.buildEmailData(input)

	uc.sendResetPasswordEmail(result, emailData, input.User.Email, locale)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result)
	return result
}

func (uc *GetResetPasswordSendEmailUseCase) buildEmailData(input dtos.OneTimeTokenUser) email_models.ResetPasswordEmailData {
	return email_models.ResetPasswordEmailData{
		Name:              input.User.Name,
		ResetLink:         input.BuildURL(settings.AppSettingsInstance.FrontendResetPasswordURL),
		ExpirationMinutes: settings.AppSettingsInstance.OneTimeTokenPasswordTTL,
		AppName:           settings.AppSettingsInstance.AppName,
		SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
	}
}

func (uc *GetResetPasswordSendEmailUseCase) sendResetPasswordEmail(
	result *usecase.UseCaseResult[bool],
	emailData email_models.ResetPasswordEmailData,
	userEmail string,
	locale locales.LocaleTypeEnum,
) {
	if err := email_service.ResetPasswordEmailServiceInstance.SendWithTemplate(
		emailData,
		userEmail,
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
	}
}

func (uc *GetResetPasswordSendEmailUseCase) setSuccessResult(result *usecase.UseCaseResult[bool]) {
	result.SetData(
		status.Success,
		true,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.RESET_PASSWORD_EMAIL_SENT_SUCCESSFULLY,
		),
	)
}

func NewGetResetPasswordSendEmailUseCase(
	log contractsProviders.ILoggerProvider,
) *GetResetPasswordSendEmailUseCase {
	return &GetResetPasswordSendEmailUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
	}
}
