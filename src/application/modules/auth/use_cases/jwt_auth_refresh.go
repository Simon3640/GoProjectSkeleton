package authusecases

import (
	"context"
	"regexp"
	"strings"
	"time"

	contractsProviders "goprojectskeleton/src/application/contracts/providers"
	dtos "goprojectskeleton/src/application/shared/DTOs"
	"goprojectskeleton/src/application/shared/locales"
	"goprojectskeleton/src/application/shared/locales/messages"
	"goprojectskeleton/src/application/shared/status"
	usecase "goprojectskeleton/src/application/shared/use_case"
)

type AuthenticationRefreshUseCase struct {
	appMessages *locales.Locale
	log         contractsProviders.ILoggerProvider
	locale      locales.LocaleTypeEnum

	jwtProvider contractsProviders.IJWTProvider
}

var _ usecase.BaseUseCase[string, dtos.Token] = (*AuthenticationRefreshUseCase)(nil)

func (uc *AuthenticationRefreshUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *AuthenticationRefreshUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input string,
) *usecase.UseCaseResult[dtos.Token] {
	result := usecase.NewUseCaseResult[dtos.Token]()
	uc.SetLocale(locale)
	validation, msg := uc.validate(input)

	if !validation {
		uc.log.Error("Invalid input", nil)
		result.SetError(
			status.InvalidInput,
			strings.Join(msg, "\n"),
		)
		return result
	}

	// Validate the access token
	claims, err := uc.jwtProvider.ParseTokenAndValidate(input)

	if err != nil {
		uc.log.Error("Error parsing or validating token", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return result
	}

	// Generate new access and refresh tokens
	sub, ok := claims["sub"].(string)
	if !ok {
		uc.log.Error("Invalid subject in claims", nil)
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
			),
		)
		return result
	}

	// Validate expiration Refresh token
	if claims["typ"] != "refresh" {
		uc.log.Error("Invalid token type", nil)
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
			),
		)
		return result
	}

	// Validate expiration date
	// define if exp is int or float

	if exp, ok := claims["exp"].(float64); !ok || exp < float64(time.Now().Unix()) {
		uc.log.Error("Token has expired", nil)
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_TOKEN_EXPIRED,
			),
		)
		return result
	}

	var claimsMap map[string]interface{}
	access, exp, err := uc.jwtProvider.GenerateAccessToken(ctx, sub, claimsMap)

	if err != nil {
		uc.log.Error("Error generating access token", err.ToError())
		result.SetError(
			status.InternalError,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
			),
		)
		return result
	}

	refresh, expRefresh, err := uc.jwtProvider.GenerateRefreshToken(ctx, sub)
	if err != nil {
		uc.log.Error("Error generating refresh token", err.ToError())
		result.SetError(
			status.InternalError,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
			),
		)
		return result
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

func (uc *AuthenticationRefreshUseCase) validate(input string) (bool, []string) {
	var validationErrors []string

	if input == "" {
		validationErrors = append(validationErrors, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.SOME_PARAMETERS_ARE_MISSING))
	}
	// regex for JWT token validation
	jwtRegex := `^[A-Za-z0-9-_=]+\.([A-Za-z0-9-_=]+\.?)*$`
	if !regexp.MustCompile(jwtRegex).MatchString(input) {
		validationErrors = append(validationErrors, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.INVALID_JWT_TOKEN))
	}
	return len(validationErrors) == 0, validationErrors
}

func NewAuthenticationRefreshUseCase(
	log contractsProviders.ILoggerProvider,
	jwtProvider contractsProviders.IJWTProvider,
) *AuthenticationRefreshUseCase {
	return &AuthenticationRefreshUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
		jwtProvider: jwtProvider,
	}
}
