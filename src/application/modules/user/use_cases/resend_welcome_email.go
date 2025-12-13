package userusecases

import (
	"strings"

	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contractsrepositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/services"
	emailservices "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	emailmodels "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/application/shared/templates"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// ResendWelcomeEmailUseCase is a use case that resends the welcome email to the user
type ResendWelcomeEmailUseCase struct {
	usecase.BaseUseCaseValidation[userdtos.ResendWelcomeEmailRequest, bool]
	log contractsproviders.ILoggerProvider

	hashProvider contractsproviders.IHashProvider
	userRepo     usercontracts.IUserRepository
	tokenRepo    contractsrepositories.IOneTimeTokenRepository
}

var _ usecase.BaseUseCase[userdtos.ResendWelcomeEmailRequest, bool] = (*ResendWelcomeEmailUseCase)(nil)

// Execute executes the use case
func (uc *ResendWelcomeEmailUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input userdtos.ResendWelcomeEmailRequest,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	// Validar input
	uc.validate(&input, result)
	if result.HasError() {
		return result
	}

	// Get user by email or phone
	user := uc.getByEmailOrPhone(input, result)
	if result.HasError() {
		return result
	}
	// Create new verification token
	token := uc.createOneTimeToken(user, result)
	if result.HasError() {
		return result
	}

	// Send welcome email
	uc.sendWelcomeEmail(*user, *token, result)
	if result.HasError() {
		return result
	}

	result.SetData(
		status.Success,
		true,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.WelcomeEmailResent,
		),
	)
	return result
}

// Attempts to get the user by email or phone.
// Sets errors in the result if the user is not found or is already verified.
// Returns the user if found, or nil if an error occurs.
func (uc *ResendWelcomeEmailUseCase) getByEmailOrPhone(
	input userdtos.ResendWelcomeEmailRequest,
	result *usecase.UseCaseResult[bool],
) *models.User {
	// Search user by email or phone
	user, err := uc.userRepo.GetByEmailOrPhone(input.Email)
	if err != nil {
		uc.log.Error("Error getting user by email or phone", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return nil
	}
	// Check if user is already verified
	if *user.Status == models.UserStatusActive {
		result.SetError(
			status.Conflict,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.UserAlreadyVerified,
			),
		)
		return nil
	}
	return user
}

// createOneTimeToken creates a one time token for the user
// returns the token if created successfully, otherwise returns nil
func (uc *ResendWelcomeEmailUseCase) createOneTimeToken(
	user *models.User,
	result *usecase.UseCaseResult[bool],
) *string {
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
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return nil
	}
	return &token
}

// sendWelcomeEmail sends a welcome email to the user
func (uc *ResendWelcomeEmailUseCase) sendWelcomeEmail(
	user models.User,
	token string,
	result *usecase.UseCaseResult[bool],
) {
	// Prepare email data
	newUserEmailData := emailmodels.NewUserEmailData{
		Name:              user.Name,
		ActivationLink:    settings.AppSettingsInstance.FrontendActivateAccountURL + "?token=" + token,
		ExpirationMinutes: int(settings.AppSettingsInstance.OneTimeTokenEmailVerifyTTL),
		AppName:           settings.AppSettingsInstance.AppName,
		SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
	}

	// Send welcome email
	if err := emailservices.RegisterUserEmailServiceInstance.SendWithTemplate(
		newUserEmailData,
		user.Email,
		uc.Locale,
		templates.TemplateKeysInstance.WelcomeEmail,
		emailservices.SubjectKeysInstance.WelcomeEmail,
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

func (uc *ResendWelcomeEmailUseCase) validate(
	input *userdtos.ResendWelcomeEmailRequest,
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
	log contractsproviders.ILoggerProvider,
	hashProvider contractsproviders.IHashProvider,
	userRepo usercontracts.IUserRepository,
	tokenRepo contractsrepositories.IOneTimeTokenRepository,
) *ResendWelcomeEmailUseCase {
	return &ResendWelcomeEmailUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[userdtos.ResendWelcomeEmailRequest, bool]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		log:          log,
		hashProvider: hashProvider,
		userRepo:     userRepo,
		tokenRepo:    tokenRepo,
	}
}
