package auth

import (
	"context"
	"regexp"
	"strings"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/services"
	email_service "gormgoskeleton/src/application/shared/services/emails"
	email_models "gormgoskeleton/src/application/shared/services/emails/models"
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/application/shared/status"
	"gormgoskeleton/src/application/shared/templates"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type AuthenticateUseCase struct {
	appMessages *locales.Locale
	log         contractsProviders.ILoggerProvider
	locale      locales.LocaleTypeEnum

	pass     contracts_repositories.IPasswordRepository
	userRepo contracts_repositories.IUserRepository
	otpRepo  contracts_repositories.IOneTimePasswordRepository

	jwtProvider  contractsProviders.IJWTProvider
	hashProvider contractsProviders.IHashProvider
}

var _ usecase.BaseUseCase[dtos.UserCredentials, dtos.Token] = (*AuthenticateUseCase)(nil)

func (uc *AuthenticateUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *AuthenticateUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.UserCredentials,
) *usecase.UseCaseResult[dtos.Token] {
	result := usecase.NewUseCaseResult[dtos.Token]()
	uc.SetLocale(locale)
	validation, msg := uc.validate(input)

	if !validation {
		result.SetError(
			status.InvalidInput,
			strings.Join(msg, "\n"),
		)
		return result
	}

	// Get password from repository
	password, err := uc.pass.GetActivePassword(input.Email)

	if err != nil {
		result.SetError(
			status.NotFound,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return result
	}

	// Get user with role from repository
	user, err := uc.userRepo.GetUserWithRole(password.UserID)

	if err != nil {
		result.SetError(
			status.NotFound,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return result
	}

	// Validate password
	valid, verifyErr := uc.hashProvider.VerifyPassword(password.Hash, input.Password)
	if !valid || verifyErr != nil {
		result.SetError(
			status.NotFound,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return result
	}

	// OTP Login
	if user.OTPLogin {
		otp, err := services.CreateOneTimePasswordService(user.ID, models.OneTimePasswordLogin, uc.hashProvider, uc.otpRepo)
		otpEmailData := email_models.OneTimePasswordEmailData{
			Name:              user.Name,
			OTPCode:           otp,
			ExpirationMinutes: int(settings.AppSettingsInstance.OneTimeTokenPasswordTTL),
			AppName:           settings.AppSettingsInstance.AppName,
			SupportEmail:      settings.AppSettingsInstance.AppSupportEmail,
		}

		if err := email_service.OneTimePasswordEmailServiceInstance.SendWithTemplate(
			otpEmailData,
			user.Email,
			uc.locale,
			templates.TemplateKeysInstance.OTPEmail,
			email_service.SubjectKeysInstance.OTPEmail,
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

		if err != nil {
			uc.log.Error("Error creating OTP", err.ToError())
			result.SetError(
				status.Conflict,
				uc.appMessages.Get(
					uc.locale,
					messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
				),
			)
			return result
		}

		result.SetSuccess(true)
		result.SetDetails(uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.OTP_LOGIN_ENABLED,
		))
		return result
	}

	claims := contractsProviders.JWTCLaims{
		"role": user.GetRoleKey(),
	}

	// Generate JWT tokens
	access, exp, err := uc.jwtProvider.GenerateAccessToken(ctx, password.UserIDString(), claims)

	if err != nil {
		result.SetError(
			status.Conflict,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
	}
	refresh, expRefresh, err := uc.jwtProvider.GenerateRefreshToken(ctx, password.UserIDString())
	if err != nil {
		result.SetError(
			status.Conflict,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
	}

	// Response

	token := dtos.Token{
		AccessToken:           access,
		RefreshToken:          refresh,
		TokenType:             "Bearer",
		AccessTokenExpiresAt:  exp,
		RefreshTokenExpiresAt: expRefresh,
	}

	result.SetData(
		status.Success,
		token,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
		),
	)
	return result
}

func (uc *AuthenticateUseCase) validate(input dtos.UserCredentials) (bool, []string) {
	// Validate the input data
	var validationErrors []string
	if input.Email == "" {
		validationErrors = append(validationErrors, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.SOME_PARAMETERS_ARE_MISSING))
	}
	// regex for email validation
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(input.Email) {
		validationErrors = append(validationErrors, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.INVALID_EMAIL))
	}

	return len(validationErrors) == 0, validationErrors
}

func NewAuthenticateUseCase(
	log contractsProviders.ILoggerProvider,
	pass contracts_repositories.IPasswordRepository,
	userRepo contracts_repositories.IUserRepository,
	otpRepo contracts_repositories.IOneTimePasswordRepository,
	hashProvider contractsProviders.IHashProvider,
	jwtProvider contractsProviders.IJWTProvider,
) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		pass:         pass,
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		jwtProvider:  jwtProvider,
		hashProvider: hashProvider,
	}
}
