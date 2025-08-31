package auth

import (
	"context"
	"time"

	contracts_providers "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
)

type AuthenticateOTPUseCase struct {
	appMessages *locales.Locale
	log         contracts_providers.ILoggerProvider
	locale      locales.LocaleTypeEnum

	pass     contracts_repositories.IPasswordRepository
	userRepo contracts_repositories.IUserRepository
	otpRepo  contracts_repositories.IOneTimePasswordRepository

	jwtProvider  contracts_providers.IJWTProvider
	hashProvider contracts_providers.IHashProvider
}

var _ usecase.BaseUseCase[string, dtos.Token] = (*AuthenticateOTPUseCase)(nil)

func (uc *AuthenticateOTPUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *AuthenticateOTPUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input string,
) *usecase.UseCaseResult[dtos.Token] {
	result := usecase.NewUseCaseResult[dtos.Token]()
	uc.SetLocale(locale)
	uc.Validate(input, result)
	if result.HasError() {
		return result
	}

	// Get hash from input
	hash := uc.hashProvider.HashOneTimeToken(input)
	oneTimePassword, err := uc.otpRepo.GetByPasswordHash(hash)

	if err != nil {
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return result
	}
	if oneTimePassword == nil || oneTimePassword.IsUsed || oneTimePassword.Expires.Before(time.Now()) {
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_OTP,
			),
		)
		return result
	}

	// Get user with role from repository
	user, err := uc.userRepo.GetUserWithRole(oneTimePassword.UserID)

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

	claims := contracts_providers.JWTCLaims{
		"role": user.GetRoleKey(),
	}

	// Generate JWT tokens
	access, exp, err := uc.jwtProvider.GenerateAccessToken(ctx, user.GetUserIDString(), claims)

	if err != nil {
		result.SetError(
			status.Conflict,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
	}
	refresh, expRefresh, err := uc.jwtProvider.GenerateRefreshToken(ctx, user.GetUserIDString())
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

	// Mark OTP as used
	_, err = uc.otpRepo.Update(oneTimePassword.ID,
		dtos.OneTimePasswordUpdate{IsUsed: true, ID: oneTimePassword.ID})
	if err != nil {
		uc.log.Error("Error updating one time password as used", err.ToError())
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
		token,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
		),
	)

	return result
}

func (uc *AuthenticateOTPUseCase) Validate(input string, result *usecase.UseCaseResult[dtos.Token]) {
	if input == "" {
		result.SetError(
			status.InvalidInput,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_DATA,
			),
		)
	}

}

func NewAuthenticateOTPUseCase(
	log contracts_providers.ILoggerProvider,
	pass contracts_repositories.IPasswordRepository,
	userRepo contracts_repositories.IUserRepository,
	otpRepo contracts_repositories.IOneTimePasswordRepository,
	hashProvider contracts_providers.IHashProvider,
	jwtProvider contracts_providers.IJWTProvider,
) *AuthenticateOTPUseCase {
	return &AuthenticateOTPUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		pass:         pass,
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		jwtProvider:  jwtProvider,
		hashProvider: hashProvider,
	}
}
