// Package authusecases contains the use cases for the authentication module.
package authusecases

import (
	"errors"
	"fmt"
	"time"

	contractproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authservices "github.com/simon3640/goprojectskeleton/src/application/modules/auth/services"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	services "github.com/simon3640/goprojectskeleton/src/application/shared/services"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// AuthenticateUseCase is the use case for the authentication of a user
type AuthenticateUseCase struct {
	usecase.BaseUseCaseValidation[dtos.UserCredentials, dtos.Token]

	pass     authcontracts.IPasswordRepository
	userRepo authcontracts.IUserRepository
	otpRepo  authcontracts.IOneTimePasswordRepository

	jwtProvider   authcontracts.IJWTProvider
	hashProvider  contractproviders.IHashProvider
	cacheProvider contractproviders.ICacheProvider
}

var _ usecase.BaseUseCase[dtos.UserCredentials, dtos.Token] = (*AuthenticateUseCase)(nil)

// Execute execute the use case
// - Check the rate limit: if the user has exceeded the login failed attempts limit, set the error and return the result
// - Get the password: get the password from the database
// - Get the user: get the user from the database
// - Validate the password: validate the password
// - Generate the tokens: generate the tokens
// - Set the success result: set the success result
// - Send the OTP email in background: send the OTP email in background
// - Return the result: return the result
func (uc *AuthenticateUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input dtos.UserCredentials,
) *usecase.UseCaseResult[dtos.Token] {
	result := usecase.NewUseCaseResult[dtos.Token]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.Validate(input, result)
	if result.HasError() {
		return result
	}

	uc.checkRateLimitAndSetError(result, input.Email)
	if result.HasError() {
		return result
	}

	password := uc.getPassword(result, input.Email)
	if result.HasError() {
		return result
	}

	user := uc.getUser(result, password.UserID, input.Email)
	if result.HasError() {
		return result
	}

	uc.validatePassword(result, password.Hash, input.Password, input.Email)
	if result.HasError() {
		return result
	}

	uc.clearFailedAttempts(input.Email)

	if user.OTPLogin {
		// OTP login: send OTP email in background
		uc.sendOTPEmailInBackground(ctx, user, locale)
		result.SetSuccess(true)
		result.SetDetails(uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.OTP_LOGIN_ENABLED,
		))
		return result
	}

	token := uc.generateTokens(ctx, result, password.UserIDString(), user)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result, token)
	observability.GetObservabilityComponents().Logger.InfoWithContext("Authentication successful", uc.AppContext)
	return result
}

func (uc *AuthenticateUseCase) checkRateLimitAndSetError(result *usecase.UseCaseResult[dtos.Token], email string) {
	if uc.cacheProvider == nil {
		return
	}

	exceeded, err := uc.checkRateLimit(email)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error checking rate limit, continuing with authentication", err.ToError(), uc.AppContext)
		return
	}

	if exceeded {
		observability.GetObservabilityComponents().Logger.WarningWithContext("Rate limit exceeded, continuing with authentication", uc.AppContext)
		result.SetError(
			status.TooManyRequests,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.LoginMaxAttemptsExceeded,
			),
		)
		return
	}
}

func (uc *AuthenticateUseCase) getPassword(result *usecase.UseCaseResult[dtos.Token], email string) *models.Password {
	password, err := uc.pass.GetActivePassword(email)
	if err != nil {
		if uc.cacheProvider != nil {
			uc.incrementFailedAttempts(email)
		}
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error getting password", err.ToError(), uc.AppContext)
		result.SetError(
			status.NotFound,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return nil
	}
	return password
}

func (uc *AuthenticateUseCase) getUser(result *usecase.UseCaseResult[dtos.Token], userID uint, email string) *models.UserWithRole {
	user, err := uc.userRepo.GetUserWithRole(userID)
	if err != nil {
		if uc.cacheProvider != nil {
			uc.incrementFailedAttempts(email)
		}
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error getting user", err.ToError(), uc.AppContext)
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

func (uc *AuthenticateUseCase) validatePassword(result *usecase.UseCaseResult[dtos.Token], passwordHash string, inputPassword string, email string) {
	valid, verifyErr := uc.hashProvider.VerifyPassword(passwordHash, inputPassword)
	if !valid || verifyErr != nil {
		if uc.cacheProvider != nil {
			uc.incrementFailedAttempts(email)
		}
		var err error
		if verifyErr != nil {
			err = verifyErr.ToError()
		} else {
			err = errors.New("password validation failed: invalid password")
		}
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error validating password", err, uc.AppContext)
		result.SetError(
			status.NotFound,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return
	}
}

func (uc *AuthenticateUseCase) generateTokens(ctx *app_context.AppContext, result *usecase.UseCaseResult[dtos.Token], userIDString string, user *models.UserWithRole) dtos.Token {
	claims := authcontracts.JWTCLaims{
		"role": user.GetRoleKey(),
	}

	access, exp, err := uc.jwtProvider.GenerateAccessToken(ctx, userIDString, claims)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error generating access token", err.ToError(), uc.AppContext)
		result.SetError(
			status.Conflict,
			uc.AppMessages.Get(
				uc.Locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
		return dtos.Token{}
	}

	refresh, expRefresh, err := uc.jwtProvider.GenerateRefreshToken(ctx, userIDString)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error generating refresh token", err.ToError(), uc.AppContext)
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

func (uc *AuthenticateUseCase) setSuccessResult(result *usecase.UseCaseResult[dtos.Token], token dtos.Token) {
	result.SetData(
		status.Success,
		token,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
		),
	)
}

// sendOTPEmailInBackground sends an OTP email to the user in the background
func (uc *AuthenticateUseCase) sendOTPEmailInBackground(
	ctx *app_context.AppContext,
	user *models.UserWithRole,
	locale locales.LocaleTypeEnum,
) {
	// Create the background service
	observability.GetObservabilityComponents().Logger.InfoWithContext("Creating OTP email background service", ctx)
	sendOTPService := authservices.NewSendOTPEmailBackgroundService(
		observability.GetObservabilityComponents(),
		uc.otpRepo,
		uc.hashProvider,
	)

	// Prepare the input
	input := authservices.SendOTPEmailInput{
		UserID:   user.ID,
		Email:    user.Email,
		UserName: user.Name,
	}

	// Execute the service in background (fire-and-forget)
	if err := services.ExecuteBackgroundService(sendOTPService, ctx, locale, input); err != nil {
		// Log error but don't fail the authentication
		// The service will log its own errors internally
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error submitting OTP email service to background executor", err, ctx)
	}
}

// checkRateLimit verify if the user has exceeded the login failed attempts limit
func (uc *AuthenticateUseCase) checkRateLimit(email string) (bool, *applicationerrors.ApplicationError) {
	if uc.cacheProvider == nil {
		return false, nil
	}

	maxAttempts := settings.AppSettingsInstance.LoginMaxAttempts
	if maxAttempts <= 0 {
		return false, nil
	}

	key := uc.getRateLimitKey(email)
	attempts, err := uc.cacheProvider.GetInt64(key)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error getting failed attempts", err.ToError(), uc.AppContext)
		return false, err
	}

	return attempts >= int64(maxAttempts), nil
}

// incrementFailedAttempts increment the failed attempts counter for an email
func (uc *AuthenticateUseCase) incrementFailedAttempts(email string) {
	if uc.cacheProvider == nil {
		return
	}

	key := uc.getRateLimitKey(email)
	_, err := uc.cacheProvider.Increment(key, time.Duration(settings.AppSettingsInstance.LoginAttemptsWindowMinutes)*time.Minute)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error incrementing failed attempts", err.ToError(), uc.AppContext)
		return
	}
}

// clearFailedAttempts clear the failed attempts counter for an email
func (uc *AuthenticateUseCase) clearFailedAttempts(email string) {
	if uc.cacheProvider == nil {
		return
	}

	key := uc.getRateLimitKey(email)
	if err := uc.cacheProvider.Delete(key); err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error clearing failed attempts", err.ToError(), uc.AppContext)
	}
}

// getRateLimitKey generate the cache key for the rate limiting
func (uc *AuthenticateUseCase) getRateLimitKey(email string) string {
	return fmt.Sprintf("login_attempts:%s", email)
}

func NewAuthenticateUseCase(
	pass authcontracts.IPasswordRepository,
	userRepo authcontracts.IUserRepository,
	otpRepo authcontracts.IOneTimePasswordRepository,
	hashProvider contractproviders.IHashProvider,
	jwtProvider authcontracts.IJWTProvider,
	cacheProvider contractproviders.ICacheProvider,
) *AuthenticateUseCase {
	return &AuthenticateUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[dtos.UserCredentials, dtos.Token]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		pass:          pass,
		userRepo:      userRepo,
		otpRepo:       otpRepo,
		jwtProvider:   jwtProvider,
		hashProvider:  hashProvider,
		cacheProvider: cacheProvider,
	}
}
