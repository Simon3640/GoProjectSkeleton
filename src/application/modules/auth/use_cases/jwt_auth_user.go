package authusecases

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	authcontracts "github.com/simon3640/goprojectskeleton/src/application/modules/auth/contracts"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

// AuthUserUseCase is the use case for authenticating a user with a JWT token
type AuthUserUseCase struct {
	appMessages *locales.Locale
	log         contractsproviders.ILoggerProvider
	locale      locales.LocaleTypeEnum

	userRepository authcontracts.IUserRepository

	jwtProvider authcontracts.IJWTProvider
}

var _ usecase.BaseUseCase[string, models.UserWithRole] = (*AuthUserUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *AuthUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *AuthUserUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input string,
) *usecase.UseCaseResult[models.UserWithRole] {
	result := usecase.NewUseCaseResult[models.UserWithRole]()
	uc.SetLocale(locale)

	uc.validate(input, result)
	if result.HasError() {
		return result
	}

	sub := uc.parseTokenAndValidate(input, result)
	if result.HasError() {
		return result
	}

	userID := uc.convertSubjectToID(result, sub)
	if result.HasError() {
		return result
	}

	user := uc.getUser(result, userID)
	if result.HasError() {
		return result
	}

	uc.setSuccessResult(result, user)
	return result
}

func (uc *AuthUserUseCase) convertSubjectToID(result *usecase.UseCaseResult[models.UserWithRole], sub *string) uint {
	subInt, err := strconv.Atoi(*sub)
	if err != nil {
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
			),
		)
		return 0
	}
	return uint(subInt)
}

func (uc *AuthUserUseCase) getUser(result *usecase.UseCaseResult[models.UserWithRole], userID uint) *models.UserWithRole {
	user, appError := uc.userRepository.GetUserWithRole(userID)
	if appError != nil {
		uc.log.Error("Error getting user with role", appError.ToError())
		result.SetError(
			appError.Code,
			uc.appMessages.Get(
				uc.locale,
				appError.Context,
			),
		)
		return nil
	}
	return user
}

func (uc *AuthUserUseCase) setSuccessResult(result *usecase.UseCaseResult[models.UserWithRole], user *models.UserWithRole) {
	result.SetData(
		status.Success,
		*user,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.PASSWORD_CREATED,
		),
	)
}

func (uc *AuthUserUseCase) validate(input string, result *usecase.UseCaseResult[models.UserWithRole]) {
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
	if len(validationErrors) > 0 {
		result.SetError(
			status.Unauthorized,
			strings.Join(validationErrors, "\n"),
		)
	}
}

func (uc *AuthUserUseCase) parseTokenAndValidate(tokenString string, result *usecase.UseCaseResult[models.UserWithRole]) *string {
	claims, err := uc.jwtProvider.ParseTokenAndValidate(tokenString)
	if err != nil {
		uc.log.Error("Failed to parse and validate token", err.ToError())
		result.SetError(
			err.Code,
			uc.appMessages.Get(
				uc.locale,
				err.Context,
			),
		)
		return nil
	}

	// Validate that is not refresh token
	if claims["typ"] != "access" {
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
			),
		)
		return nil
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
		return nil
	}

	// Extract subject from claims
	sub, ok := claims["sub"].(string)
	if !ok {
		uc.log.Error("Invalid subject in token claims", nil)
		result.SetError(
			status.Unauthorized,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.AUTHORIZATION_HEADER_INVALID,
			),
		)
		return nil
	}

	return &sub
}

func NewAuthUserUseCase(
	log contractsproviders.ILoggerProvider,
	userRepository authcontracts.IUserRepository,
	jwtProvider authcontracts.IJWTProvider,
) *AuthUserUseCase {
	return &AuthUserUseCase{
		appMessages:    locales.NewLocale(locales.EN_US),
		log:            log,
		userRepository: userRepository,
		jwtProvider:    jwtProvider,
	}
}
