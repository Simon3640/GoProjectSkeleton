package userusecases

import (
	"context"
	"strings"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "goprojectskeleton/src/application/contracts/repositories"
	dtos "goprojectskeleton/src/application/shared/DTOs"
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

// ResendWelcomeEmailUseCase is a use case that resends the welcome email to the user
type ResendWelcomeEmailUseCase struct {
	appMessages *locales.Locale
	log         contractsProviders.ILoggerProvider
	locale      locales.LocaleTypeEnum

	hashProvider contractsProviders.IHashProvider
	userRepo     contracts_repositories.IUserRepository
	tokenRepo    contracts_repositories.IOneTimeTokenRepository
}

var _ usecase.BaseUseCase[dtos.ResendWelcomeEmailRequest, bool] = (*ResendWelcomeEmailUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *ResendWelcomeEmailUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

// Execute executes the use case
func (uc *ResendWelcomeEmailUseCase) Execute(_ context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.ResendWelcomeEmailRequest,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)

	// Validar input
	uc.validate(&input, result)
	if result.HasError() {
		return result
	}

	// Buscar usuario por email
	user, err := uc.userRepo.GetByEmailOrPhone(input.Email)
	if err != nil {
		uc.log.Error("Error getting user by email", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return result
	}

	// Generar nuevo token de verificación
	token, err := services.CreateOneTimeTokenService(
		user.ID,
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

	// Preparar datos del email
	newUserEmailData := email_models.NewUserEmailData{
		Name:              user.Name,
		ActivationLink:    settings.AppSettingsInstance.FrontendActivateAccountURL + "?token=" + token,
		ExpirationMinutes: int(settings.AppSettingsInstance.OneTimeTokenEmailVerifyTTL),
		AppName:           settings.AppSettingsInstance.AppName,
		SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
	}

	// Enviar correo de bienvenida
	if err := email_service.RegisterUserEmailServiceInstance.SendWithTemplate(
		newUserEmailData,
		user.Email,
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
		true,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.WelcomeEmailResent,
		),
	)
	return result
}

func (uc *ResendWelcomeEmailUseCase) validate(
	input *dtos.ResendWelcomeEmailRequest,
	result *usecase.UseCaseResult[bool]) {
	msgs := input.Validate()

	if len(msgs) > 0 {
		result.SetError(
			status.InvalidInput,
			strings.Join(msgs, "\n"),
		)
	}
}

// NewResendWelcomeEmailUseCase creates a new ResendWelcomeEmailUseCase
func NewResendWelcomeEmailUseCase(
	log contractsProviders.ILoggerProvider,
	hashProvider contractsProviders.IHashProvider,
	userRepo contracts_repositories.IUserRepository,
	tokenRepo contracts_repositories.IOneTimeTokenRepository,
) *ResendWelcomeEmailUseCase {
	return &ResendWelcomeEmailUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		hashProvider: hashProvider,
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
	}
}
