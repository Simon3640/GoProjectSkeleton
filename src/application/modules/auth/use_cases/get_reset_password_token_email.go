package authusecases

import (
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	emailservice "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	emailmodels "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/application/shared/templates"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
)

// GetResetPasswordSendEmailUseCase is the use case for sending a reset password email
type GetResetPasswordSendEmailUseCase struct {
	usecase.BaseUseCaseValidation[bool, bool]
	log contractsproviders.ILoggerProvider
}

// Execute sends a reset password email to the user
func (uc *GetResetPasswordSendEmailUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input bool,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.Validate(input, result)
	if result.HasError() {
		return result
	}

	emailData := uc.buildEmailData(*ctx.OneTimeToken)

	uc.sendResetPasswordEmail(result, emailData, ctx.OneTimeToken.User.Email, locale)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result)
	return result
}

func (uc *GetResetPasswordSendEmailUseCase) buildEmailData(input shareddtos.OneTimeTokenUser) emailmodels.ResetPasswordEmailData {
	return emailmodels.ResetPasswordEmailData{
		Name:              input.User.Name,
		ResetLink:         input.BuildURL(settings.AppSettingsInstance.FrontendResetPasswordURL),
		ExpirationMinutes: settings.AppSettingsInstance.OneTimeTokenPasswordTTL,
		AppName:           settings.AppSettingsInstance.AppName,
		SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
	}
}

func (uc *GetResetPasswordSendEmailUseCase) sendResetPasswordEmail(
	result *usecase.UseCaseResult[bool],
	emailData emailmodels.ResetPasswordEmailData,
	userEmail string,
	locale locales.LocaleTypeEnum,
) {
	if err := emailservice.ResetPasswordEmailServiceInstance.SendWithTemplate(
		emailData,
		userEmail,
		locale,
		templates.TemplateKeysInstance.PasswordResetEmail,
		emailservice.SubjectKeysInstance.PasswordResetEmail,
	); err != nil {
		uc.log.Error("Error sending email", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}
}

func (uc *GetResetPasswordSendEmailUseCase) setSuccessResult(result *usecase.UseCaseResult[bool]) {
	result.SetData(
		status.Success,
		true,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.RESET_PASSWORD_EMAIL_SENT_SUCCESSFULLY,
		),
	)
}

func (uc *GetResetPasswordSendEmailUseCase) Validate(
	input bool,
	result *usecase.UseCaseResult[bool],
) {
	if uc.AppContext.OneTimeToken == nil || !input || uc.AppContext.OneTimeToken.User.Email == "" {
		result.SetError(
			status.InvalidInput,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.INVALID_DATA,
			),
		)
	}
}

func NewGetResetPasswordSendEmailUseCase(
	log contractsproviders.ILoggerProvider,
) *GetResetPasswordSendEmailUseCase {
	return &GetResetPasswordSendEmailUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[bool, bool]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log: log,
	}
}
