package authusecases

import (
	"context"
	"time"

	contractproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// AuthenticateOTPUseCase is the use case for authenticating a user with an OTP
type AuthenticateOTPUseCase struct {
	usecase.BaseUseCaseValidation[string, dtos.Token]

	userRepo authcontracts.IUserRepository
	otpRepo  authcontracts.IOneTimePasswordRepository

	jwtProvider  authcontracts.IJWTProvider
	hashProvider contractproviders.IHashProvider
}

var _ usecase.BaseUseCase[string, dtos.Token] = (*AuthenticateOTPUseCase)(nil)

// Execute authenticates a user with an OTP
func (uc *AuthenticateOTPUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input string,
) *usecase.UseCaseResult[dtos.Token] {
	result := usecase.NewUseCaseResult[dtos.Token]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
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
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return nil
	}

	if oneTimePassword == nil || oneTimePassword.IsUsed || oneTimePassword.Expires.Before(time.Now()) {
		result.SetError(
			status.Unauthorized,
			uc.AppMessages.Get(
				uc.Locale,
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
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return nil
	}
	return user
}

func (uc *AuthenticateOTPUseCase) generateTokens(ctx context.Context, result *usecase.UseCaseResult[dtos.Token], user *models.UserWithRole) dtos.Token {
	claims := authcontracts.JWTCLaims{
		"role": user.GetRoleKey(),
	}

	access, exp, err := uc.jwtProvider.GenerateAccessToken(ctx, user.GetUserIDString(), claims)
	if err != nil {
		result.SetError(
			status.Conflict,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
		return dtos.Token{}
	}

	refresh, expRefresh, err := uc.jwtProvider.GenerateRefreshToken(ctx, user.GetUserIDString())
	if err != nil {
		result.SetError(
			status.Conflict,
			uc.AppMessages.Get(
				uc.Locale,
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
		observability.GetObservabilityComponents().Logger.Error("Error updating one time password as used", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
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
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
		),
	)
}

func (uc *AuthenticateOTPUseCase) Validate(input string, result *usecase.UseCaseResult[dtos.Token]) {
	if input == "" {
		result.SetError(
			status.InvalidInput,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.INVALID_DATA,
			),
		)
	}

}

func NewAuthenticateOTPUseCase(
	userRepo authcontracts.IUserRepository,
	otpRepo authcontracts.IOneTimePasswordRepository,
	hashProvider contractproviders.IHashProvider,
	jwtProvider authcontracts.IJWTProvider,
) *AuthenticateOTPUseCase {
	return &AuthenticateOTPUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[string, dtos.Token]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		userRepo:     userRepo,
		otpRepo:      otpRepo,
		jwtProvider:  jwtProvider,
		hashProvider: hashProvider,
	}
}
