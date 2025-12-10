package authusecases

import (
	"context"
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type AuthenticateOTPUseCase struct {
	appMessages *locales.Locale
	log         contractsProviders.ILoggerProvider
	locale      locales.LocaleTypeEnum

	userRepo contracts_repositories.IUserRepository
	otpRepo  contracts_repositories.IOneTimePasswordRepository

	jwtProvider  contractsProviders.IJWTProvider
	hashProvider contractsProviders.IHashProvider
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

	oneTimePassword := uc.validateAndGetOTP(result, input)
	if result.HasError() {
		return result
	}

	user := uc.getUser(result, oneTimePassword.UserID)
	if result.HasError() {
		return result
	}

	token := uc.generateTokens(ctx, result, user)
	if result.HasError() {
		return result
	}

	uc.markOTPAsUsed(result, oneTimePassword.ID)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result, token)
	return result
}

func (uc *AuthenticateOTPUseCase) validateAndGetOTP(result *usecase.UseCaseResult[dtos.Token], otp string) *models.OneTimePassword {
	hash := uc.hashProvider.HashOneTimeToken(otp)
	oneTimePassword, err := uc.otpRepo.GetByPasswordHash(hash)

	if err != nil {
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return nil
	}

	if oneTimePassword == nil || oneTimePassword.IsUsed || oneTimePassword.Expires.Before(time.Now()) {
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_OTP,
			),
		)
		return nil
	}

	return oneTimePassword
}

func (uc *AuthenticateOTPUseCase) getUser(result *usecase.UseCaseResult[dtos.Token], userID uint) *models.UserWithRole {
	user, err := uc.userRepo.GetUserWithRole(userID)
	if err != nil {
		result.SetError(
			status.NotFound,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return nil
	}
	return user
}

func (uc *AuthenticateOTPUseCase) generateTokens(ctx context.Context, result *usecase.UseCaseResult[dtos.Token], user *models.UserWithRole) dtos.Token {
	claims := contractsProviders.JWTCLaims{
		"role": user.GetRoleKey(),
	}

	access, exp, err := uc.jwtProvider.GenerateAccessToken(ctx, user.GetUserIDString(), claims)
	if err != nil {
		result.SetError(
			status.Conflict,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
		return dtos.Token{}
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
		return dtos.Token{}
	}

	return dtos.Token{
		AccessToken:           access,
		RefreshToken:          refresh,
		TokenType:             "Bearer",
		AccessTokenExpiresAt:  exp,
		RefreshTokenExpiresAt: expRefresh,
	}
}

func (uc *AuthenticateOTPUseCase) markOTPAsUsed(result *usecase.UseCaseResult[dtos.Token], otpID uint) {
	_, err := uc.otpRepo.Update(otpID,
		dtos.OneTimePasswordUpdate{IsUsed: true, ID: otpID})
	if err != nil {
		uc.log.Error("Error updating one time password as used", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return
	}
}

func (uc *AuthenticateOTPUseCase) setSuccessResult(result *usecase.UseCaseResult[dtos.Token], token dtos.Token) {
	result.SetData(
		status.Success,
		token,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
		),
	)
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
	log contractsProviders.ILoggerProvider,
	userRepo contracts_repositories.IUserRepository,
	otpRepo contracts_repositories.IOneTimePasswordRepository,
	hashProvider contractsProviders.IHashProvider,
	jwtProvider contractsProviders.IJWTProvider,
) *AuthenticateOTPUseCase {
	return &AuthenticateOTPUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		jwtProvider:  jwtProvider,
		hashProvider: hashProvider,
	}
}
