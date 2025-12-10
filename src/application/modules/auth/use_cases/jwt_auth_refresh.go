package authusecases

import (
	"context"
	"regexp"
	"strings"
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	dtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
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

	uc.validateInput(result, input)
	if result.HasError() {
		return result
	}

	claims := uc.parseAndValidateToken(result, input)
	if result.HasError() {
		return result
	}

	subject := uc.validateClaims(result, claims)
	if result.HasError() {
		return result
	}

	token := uc.generateTokens(ctx, result, subject)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result, token)
	return result
}

func (uc *AuthenticationRefreshUseCase) validateInput(result *usecase.UseCaseResult[dtos.Token], input string) {
	validation, msg := uc.validate(input)
	if !validation {
		uc.log.Error("Invalid input", nil)
		result.SetError(
			status.InvalidInput,
			strings.Join(msg, "\n"),
		)
		return
	}
}

func (uc *AuthenticationRefreshUseCase) parseAndValidateToken(result *usecase.UseCaseResult[dtos.Token], token string) contractsProviders.JWTCLaims {
	claims, err := uc.jwtProvider.ParseTokenAndValidate(token)
	if err != nil {
		uc.log.Error("Error parsing or validating token", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return nil
	}
	return claims
}

func (uc *AuthenticationRefreshUseCase) validateClaims(result *usecase.UseCaseResult[dtos.Token], claims contractsProviders.JWTCLaims) string {
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
		return ""
	}

	if claims["typ"] != "refresh" {
		uc.log.Error("Invalid token type", nil)
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
			),
		)
		return ""
	}

	if exp, ok := claims["exp"].(float64); !ok || exp < float64(time.Now().Unix()) {
		uc.log.Error("Token has expired", nil)
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_TOKEN_EXPIRED,
			),
		)
		return ""
	}

	return sub
}

func (uc *AuthenticationRefreshUseCase) generateTokens(ctx context.Context, result *usecase.UseCaseResult[dtos.Token], subject string) dtos.Token {
	var claimsMap map[string]interface{}
	access, exp, err := uc.jwtProvider.GenerateAccessToken(ctx, subject, claimsMap)
	if err != nil {
		uc.log.Error("Error generating access token", err.ToError())
		result.SetError(
			status.InternalError,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
			),
		)
		return dtos.Token{}
	}

	refresh, expRefresh, err := uc.jwtProvider.GenerateRefreshToken(ctx, subject)
	if err != nil {
		uc.log.Error("Error generating refresh token", err.ToError())
		result.SetError(
			status.InternalError,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
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

func (uc *AuthenticationRefreshUseCase) setSuccessResult(result *usecase.UseCaseResult[dtos.Token], token dtos.Token) {
	result.SetData(
		status.Success,
		token,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_GENERATED,
		),
	)
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
