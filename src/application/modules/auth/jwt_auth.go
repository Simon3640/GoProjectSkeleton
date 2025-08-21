package auth

import (
	"context"
	"regexp"
	"strings"

	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type AuthenticateUseCase struct {
	appMessages *locales.Locale
	log         contracts.ILoggerProvider
	locale      locales.LocaleTypeEnum

	pass contracts_repositories.IPasswordRepository

	jwtProvider  contracts.IJWTProvider
	hashProvider contracts.IHashProvider
}

var _ usecase.BaseUseCase[models.UserCredentials, models.Token] = (*AuthenticateUseCase)(nil)

func (uc *AuthenticateUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *AuthenticateUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input models.UserCredentials,
) *usecase.UseCaseResult[models.Token] {
	result := usecase.NewUseCaseResult[models.Token]()
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

	// TODO: GO for user info
	claims := contracts.JWTCLaims{
		"role": "user",
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

	token := models.Token{
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

func (uc *AuthenticateUseCase) validate(input models.UserCredentials) (bool, []string) {
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
	log contracts.ILoggerProvider,
	pass contracts_repositories.IPasswordRepository,
	hashProvider contracts.IHashProvider,
	jwtProvider contracts.IJWTProvider,
) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		pass:         pass,
		jwtProvider:  jwtProvider,
		hashProvider: hashProvider,
	}
}
