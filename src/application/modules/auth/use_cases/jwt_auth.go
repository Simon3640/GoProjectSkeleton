// Package authusecases contains the use cases for the authentication module.
package authusecases

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	application_errors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/services"
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
	log         contractsProviders.ILoggerProvider
	locale      locales.LocaleTypeEnum

	pass     contracts_repositories.IPasswordRepository
	userRepo contracts_repositories.IUserRepository
	otpRepo  contracts_repositories.IOneTimePasswordRepository

	jwtProvider   contractsProviders.IJWTProvider
	hashProvider  contractsProviders.IHashProvider
	cacheProvider contractsProviders.ICacheProvider
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
	validation, msg := uc.validate(input)

	if !validation {
		result.SetError(
			status.InvalidInput,
			strings.Join(msg, "\n"),
		)
		return result
	}

	// Check rate limiting before attempting authentication
	if uc.cacheProvider != nil {
		exceeded, err := uc.checkRateLimit(input.Email)
		if err != nil {
			uc.log.Error("Error checking rate limit, continuing with authentication", err.ToError())
		} else if exceeded {
			result.SetError(
				status.TooManyRequests,
				uc.appMessages.Get(
					uc.locale,
					messages.MessageKeysInstance.LoginMaxAttemptsExceeded,
				),
			)
			return result
		}
	}

	// Get password from repository
	password, err := uc.pass.GetActivePassword(input.Email)

	if err != nil {
		// Increment failed attempts counter
		if uc.cacheProvider != nil {
			uc.incrementFailedAttempts(input.Email)
		}
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
		// Increment failed attempts counter
		if uc.cacheProvider != nil {
			uc.incrementFailedAttempts(input.Email)
		}
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
		// Increment failed attempts counter
		if uc.cacheProvider != nil {
			uc.incrementFailedAttempts(input.Email)
		}
		result.SetError(
			status.NotFound,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.INVALID_USER_OR_PASSWORD,
			),
		)
		return result
	}

	// Clear failed attempts counter on successful authentication
	if uc.cacheProvider != nil {
		uc.clearFailedAttempts(input.Email)
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
	log contractsProviders.ILoggerProvider,
	pass contracts_repositories.IPasswordRepository,
	userRepo contracts_repositories.IUserRepository,
	otpRepo contracts_repositories.IOneTimePasswordRepository,
	hashProvider contractsProviders.IHashProvider,
	jwtProvider contractsProviders.IJWTProvider,
	cacheProvider contractsProviders.ICacheProvider,
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
