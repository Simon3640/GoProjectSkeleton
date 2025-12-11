// Package authusecases contains the use cases for the authentication module.
package authusecases

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	contractproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	dtos "github.com/simon3640/goprojectskeleton/src/application/modules/auth/dtos"
	authservices "github.com/simon3640/goprojectskeleton/src/application/modules/auth/services"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	email_service "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails"
	email_models "github.com/simon3640/goprojectskeleton/src/application/shared/services/emails/models"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	"github.com/simon3640/goprojectskeleton/src/application/shared/templates"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// AuthenticateUseCase is the use case for the authentication of a user
type AuthenticateUseCase struct {
	appMessages *locales.Locale
	log         contractproviders.ILoggerProvider
	locale      locales.LocaleTypeEnum

	pass     authcontracts.IPasswordRepository
	userRepo authcontracts.IUserRepository
	otpRepo  authcontracts.IOneTimePasswordRepository

	jwtProvider   authcontracts.IJWTProvider
	hashProvider  contractproviders.IHashProvider
	cacheProvider contractproviders.ICacheProvider
}

var _ usecase.BaseUseCase[dtos.UserCredentials, dtos.Token] = (*AuthenticateUseCase)(nil)

// SetLocale set the locale for the use case
func (uc *AuthenticateUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

// Execute execute the use case
func (uc *AuthenticateUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input dtos.UserCredentials,
) *usecase.UseCaseResult[dtos.Token] {
	result := usecase.NewUseCaseResult[dtos.Token]()
	uc.SetLocale(locale)

	uc.validateInput(result, input)
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
		uc.handleOTPLogin(result, user)
		return result
	}

	token := uc.generateTokens(ctx, result, password.UserIDString(), user)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result, token)
	return result
}

func (uc *AuthenticateUseCase) validateInput(result *usecase.UseCaseResult[dtos.Token], input dtos.UserCredentials) {
	validation, msg := uc.validate(input)
	if !validation {
		result.SetError(
			status.InvalidInput,
			strings.Join(msg, "\n"),
		)
		return
	}
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
			uc.appMessages.Get(
				uc.locale,
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
			uc.appMessages.Get(
				uc.locale,
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
			uc.appMessages.Get(
				uc.locale,
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
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return
	}
}

func (uc *AuthenticateUseCase) handleOTPLogin(result *usecase.UseCaseResult[dtos.Token], user *models.UserWithRole) {
	otp, err := authservices.CreateOneTimePasswordService(user.ID, models.OneTimePasswordLogin, uc.hashProvider, uc.otpRepo)
	if err != nil {
		uc.log.Error("Error creating OTP", err.ToError())
		result.SetError(
			status.Conflict,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
		return
	}

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
		return
	}

	result.SetSuccess(true)
	result.SetDetails(uc.appMessages.Get(
		uc.locale,
		messages.MessageKeysInstance.OTP_LOGIN_ENABLED,
	))
}

func (uc *AuthenticateUseCase) generateTokens(ctx context.Context, result *usecase.UseCaseResult[dtos.Token], userIDString string, user *models.UserWithRole) dtos.Token {
	claims := authcontracts.JWTCLaims{
		"role": user.GetRoleKey(),
	}

	access, exp, err := uc.jwtProvider.GenerateAccessToken(ctx, userIDString, claims)
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

	refresh, expRefresh, err := uc.jwtProvider.GenerateRefreshToken(ctx, userIDString)
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

func (uc *AuthenticateUseCase) setSuccessResult(result *usecase.UseCaseResult[dtos.Token], token dtos.Token) {
	result.SetData(
		status.Success,
		token,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
		),
	)
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

// checkRateLimit verify if the user has exceeded the login failed attempts limit
func (uc *AuthenticateUseCase) checkRateLimit(email string) (bool, *application_errors.ApplicationError) {
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
		appMessages:   locales.NewLocale(locales.EN_US),
		log:           log,
		pass:          pass,
		userRepo:      userRepo,
		otpRepo:       otpRepo,
		jwtProvider:   jwtProvider,
		hashProvider:  hashProvider,
		cacheProvider: cacheProvider,
	}
}
