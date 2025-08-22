package auth

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type AuthUserUseCase struct {
	appMessages *locales.Locale
	log         contracts.ILoggerProvider
	locale      locales.LocaleTypeEnum

	userRepository contracts_repositories.IUserRepository

	jwtProvider contracts.IJWTProvider
}

var _ usecase.BaseUseCase[string, models.UserWithRole] = (*AuthUserUseCase)(nil)

func (uc *AuthUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *AuthUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input string,
) *usecase.UseCaseResult[models.UserWithRole] {
	result := usecase.NewUseCaseResult[models.UserWithRole]()
	uc.SetLocale(locale)
	validation, msg := uc.validate(input)

	if !validation {
		result.SetError(
			status.Unauthorized,
			strings.Join(msg, "\n"),
		)
		return result
	}

	sub, err := uc.parseTokenAndValidate(input)
	if err != nil {
		result.SetError(
			status.Unauthorized,
			err.Error(),
		)
		return result
	}

	// convert subject to uint

	subInt, err := strconv.Atoi(sub)
	if err != nil {
		result.SetError(
			status.Unauthorized,
			"Invalid subject in token",
		)
		return result
	}
	subID := uint(subInt)

	user, err := uc.userRepository.GetUserWithRole(subID)

	if err != nil {
		result.SetError(
			status.NotFound,
			err.Error(),
		)
		return result
	}

	result.SetData(
		status.Success,
		*user,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.PASSWORD_CREATED,
		),
	)
	return result
}

func (uc *AuthUserUseCase) validate(input string) (bool, []string) {
	// Validate the input data
	var validationErrors []string

	if input == "" {
		validationErrors = append(validationErrors, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.AUTHORIZATION_REQUIRED))
	}
	// regex for JWT token validation
	jwtRegex := `^[A-Za-z0-9-_=]+\.([A-Za-z0-9-_=]+\.?)*$`
	if !regexp.MustCompile(jwtRegex).MatchString(input) {
		validationErrors = append(validationErrors, uc.appMessages.Get(uc.locale, messages.MessageKeysInstance.INVALID_JWT_TOKEN))
	}
	return len(validationErrors) == 0, validationErrors
}

func (uc *AuthUserUseCase) parseTokenAndValidate(tokenString string) (string, error) {
	claims, err := uc.jwtProvider.ParseTokenAndValidate(tokenString)
	if err != nil {
		uc.log.Error("Error parsing or validating token", err)
		return "", errors.New(uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
		))
	}

	// Validate that is not refresh token
	if claims["typ"] != "access" {
		uc.log.Error("Invalid token type", nil)
		return "", errors.New(uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
		))
	}

	if exp, ok := claims["exp"].(float64); !ok || exp < float64(time.Now().Unix()) {
		uc.log.Error("Token has expired", nil)
		return "", errors.New(uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
		))
	}

	// Extract subject from claims
	sub, ok := claims["sub"].(string)
	if !ok {
		uc.log.Error("Invalid subject in token claims", nil)
		return "", errors.New(uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
		))
	}

	return sub, nil
}

func NewAuthUserUseCase(
	log contracts.ILoggerProvider,
	userRepository contracts_repositories.IUserRepository,
	jwtProvider contracts.IJWTProvider,
) *AuthUserUseCase {
	return &AuthUserUseCase{
		appMessages:    locales.NewLocale(locales.EN_US),
		log:            log,
		userRepository: userRepository,
		jwtProvider:    jwtProvider,
	}
}
