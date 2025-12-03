package usecases_user

import (
	"context"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "goprojectskeleton/src/application/contracts/repositories"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/locales/messages"
	"goprojectskeleton/src/application/shared/services"
	email_service "goprojectskeleton/src/application/shared/services/emails"
	email_models "goprojectskeleton/src/application/shared/services/emails/models"
	"goprojectskeleton/src/application/shared/settings"
	"goprojectskeleton/src/application/shared/status"
	"goprojectskeleton/src/application/shared/templates"
	usecase "goprojectskeleton/src/application/shared/use_case"
	"goprojectskeleton/src/domain/models"
)

type CreateUserSendEmailUseCase struct {
	appMessages *locales.Locale
	log         contractsProviders.ILoggerProvider
	locale      locales.LocaleTypeEnum

	hashProvider contractsProviders.IHashProvider

	tokenRepo contracts_repositories.IOneTimeTokenRepository
}

var _ usecase.BaseUseCase[models.User, models.User] = (*CreateUserSendEmailUseCase)(nil)

func (uc *CreateUserSendEmailUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *CreateUserSendEmailUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input models.User,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)

	token, err := services.CreateOneTimeTokenService(
		input.ID,
		models.OneTimeTokenPurposeEmailVerify,
		uc.hashProvider,
		uc.tokenRepo,
	)
	if err != nil {
		uc.log.Error("Error creating one time token", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return result
	}

	newUserEmailData := email_models.NewUserEmailData{
		Name:              input.Name,
		ActivationLink:    settings.AppSettingsInstance.FrontendActivateAccountURL + "?token=" + token,
		ExpirationMinutes: int(settings.AppSettingsInstance.OneTimeTokenEmailVerifyTTL),
		AppName:           settings.AppSettingsInstance.AppName,
		SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
	}

	if err := email_service.RegisterUserEmailServiceInstance.SendWithTemplate(
		newUserEmailData,
		input.Email,
		locale,
		templates.TemplateKeysInstance.WelcomeEmail,
		email_service.SubjectKeysInstance.WelcomeEmail,
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
		input,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	return result
}

func NewCreateUserSendEmailUseCase(
	log contractsProviders.ILoggerProvider,
	hashProvider contractsProviders.IHashProvider,
	tokenRepo contracts_repositories.IOneTimeTokenRepository,
) *CreateUserSendEmailUseCase {
	return &CreateUserSendEmailUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		hashProvider: hashProvider,
		tokenRepo:    tokenRepo,
	}
}
