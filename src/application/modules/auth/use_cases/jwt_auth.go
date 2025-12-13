// Package authusecases contains the use cases for the authentication module.
package authusecases

import (
	"fmt"
	"regexp"
	"time"

	contractproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authservices "github.com/simon3640/goprojectskeleton/src/application/modules/auth/services"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	services "github.com/simon3640/goprojectskeleton/src/application/shared/services"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// AuthenticateUseCase is the use case for the authentication of a user
type AuthenticateUseCase struct {
	usecase.BaseUseCaseValidation[dtos.UserCredentials, dtos.Token]
	log contractproviders.ILoggerProvider

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
	return result
}

func (uc *AuthenticateUseCase) checkRateLimitAndSetError(result *usecase.UseCaseResult[dtos.Token], email string) {
	if uc.cacheProvider == nil {
		return
	}

	exceeded, err := uc.checkRateLimit(email)
	if err != nil {
		uc.log.Error("Error checking rate limit, continuing with authentication", err.ToError())
		return
	}

	if exceeded {
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
	sendOTPService := authservices.NewSendOTPEmailBackgroundService(
		uc.log,
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
		uc.log.Error("Error submitting OTP email service to background executor", err)
	}
}

func (uc *AuthenticateUseCase) validate(input dtos.UserCredentials) (bool, []string) {
	// Validate the input data
	var validationErrors []string
	if input.Email == "" {
		validationErrors = append(validationErrors, uc.AppMessages.Get(uc.Locale, messages.MessageKeysInstance.SOME_PARAMETERS_ARE_MISSING))
	}
	// regex for email validation
	if !regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(input.Email) {
		validationErrors = append(validationErrors, uc.AppMessages.Get(uc.Locale, messages.MessageKeysInstance.INVALID_EMAIL))
	}

	return len(validationErrors) == 0, validationErrors
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
		uc.log.Error("Error getting failed attempts", err.ToError())
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
		uc.log.Error("Error incrementing failed attempts", err.ToError())
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
		uc.log.Error("Error clearing failed attempts", err.ToError())
	}
}

// getRateLimitKey generate the cache key for the rate limiting
func (uc *AuthenticateUseCase) getRateLimitKey(email string) string {
	return fmt.Sprintf("login_attempts:%s", email)
}

func NewAuthenticateUseCase(
	log contractproviders.ILoggerProvider,
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
		log:           log,
		pass:          pass,
		userRepo:      userRepo,
		otpRepo:       otpRepo,
		jwtProvider:   jwtProvider,
		hashProvider:  hashProvider,
		cacheProvider: cacheProvider,
	}
}
