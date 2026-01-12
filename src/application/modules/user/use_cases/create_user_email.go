package userusecases

import (
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/services"
	emailservice "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	emailmodels "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/application/shared/templates"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	sharedmodels "github.com/simon3640/goprojectskeleton/src/domain/shared/models"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

// CreateUserSendEmailUseCase is a use case that sends an email to a user
type CreateUserSendEmailUseCase struct {
	usecase.BaseUseCaseValidation[usermodels.User, usermodels.User]

	hashProvider contractsproviders.IHashProvider

	tokenRepo contractsrepositories.IOneTimeTokenRepository
}

var _ usecase.BaseUseCase[usermodels.User, usermodels.User] = (*CreateUserSendEmailUseCase)(nil)

// Execute executes the use case
func (uc *CreateUserSendEmailUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input usermodels.User,
) *usecase.UseCaseResult[usermodels.User] {
	result := usecase.NewUseCaseResult[usermodels.User]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	token := uc.createOneTimeToken(input, result)
	if result.HasError() {
		return result
	}

	uc.sendWelcomeEmail(input, *token, result)
	if result.HasError() {
		return result
	}

	result.SetData(
		status.Created,
		input,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_WAS_CREATED,
		),
	)
	observability.GetObservabilityComponents().Logger.InfoWithContext("user_created_and_email_sent", uc.AppContext)
	return result
}

// createOneTimeToken creates a one time token for the user
// returns the token if created successfully, otherwise returns nil
func (uc *CreateUserSendEmailUseCase) createOneTimeToken(input usermodels.User, result *usecase.UseCaseResult[usermodels.User]) *string {
	token, err := services.CreateOneTimeTokenService(
		input.ID,
		sharedmodels.OneTimeTokenPurposeEmailVerify,
		uc.hashProvider,
		uc.tokenRepo,
	)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error creating one time token", err.ToError(), uc.AppContext)
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return nil
	}
	return &token
}

// sendWelcomeEmail sends a welcome email to the user.
// If sending fails, sets an error in the result.
func (uc *CreateUserSendEmailUseCase) sendWelcomeEmail(input usermodels.User, token string, result *usecase.UseCaseResult[usermodels.User]) {
	newUserEmailData := emailmodels.NewUserEmailData{
		Name:              input.Name,
		ActivationLink:    settings.AppSettingsInstance.FrontendActivateAccountURL + "?token=" + token,
		ExpirationMinutes: int(settings.AppSettingsInstance.OneTimeTokenEmailVerifyTTL),
		AppName:           settings.AppSettingsInstance.AppName,
		SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
	}
	if err := emailservice.RegisterUserEmailServiceInstance.SendWithTemplate(
		newUserEmailData,
		input.Email,
		uc.Locale,
		templates.TemplateKeysInstance.WelcomeEmail,
		emailservice.SubjectKeysInstance.WelcomeEmail,
	); err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error sending email", err.ToError(), uc.AppContext)
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}
}

// NewCreateUserSendEmailUseCase creates a new create user send email use case
func NewCreateUserSendEmailUseCase(
	hashProvider contractsproviders.IHashProvider,
	tokenRepo contractsrepositories.IOneTimeTokenRepository,
) *CreateUserSendEmailUseCase {
	return &CreateUserSendEmailUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[usermodels.User, usermodels.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		hashProvider: hashProvider,
		tokenRepo:    tokenRepo,
	}
}
