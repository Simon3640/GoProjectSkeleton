package usecases_password

import (
	"context"
	"strings"

	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
)

type CreatePasswordUseCase struct {
	appMessages  *locales.Locale
	log          contracts.ILoggerProvider
	repo         contracts_repositories.IPasswordRepository
	locale       locales.LocaleTypeEnum
	hashProvider contracts.IHashProvider
}

var _ usecase.BaseUseCase[models.PasswordCreateNoHash, bool] = (*CreatePasswordUseCase)(nil)

func (uc *CreatePasswordUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *CreatePasswordUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input models.PasswordCreateNoHash,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	validation, msg := uc.validate(input)

	if !validation {
		result.SetError(
			status.InvalidInput,
			strings.Join(msg, "\n"),
		)
		return result
	}

	hashedPassword, err := uc.hashProvider.HashPassword(input.NoHashedPassword)
	if err != nil {
		result.SetError(
			status.InternalError,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
		return result
	}

	passwordCreate := models.NewPasswordCreate(
		input.UserID,
		hashedPassword,
		input.ExpiresAt,
		input.IsActive,
	)

	res, err := uc.repo.Create(passwordCreate)

	if err != nil {
		result.SetError(
			status.Conflict,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
	}

	var success bool
	if res != nil {
		success = true
	} else {
		success = false
	}

	result.SetData(
		status.Success,
		success,
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.PASSWORD_CREATED,
		),
	)
	return result
}

func (uc *CreatePasswordUseCase) validate(input models.PasswordCreateNoHash) (bool, []string) {
	// Validate the input data
	var validationErrors []string
	if input.UserID <= 0 {
		validationErrors = append(validationErrors, uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.INVALID_USER_ID,
		))
	}
	if input.NoHashedPassword == "" {
		validationErrors = append(validationErrors, uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.PASSWORD_REQUIRED,
		))
	}
	// Must have at least one uppercase letter, one lowercase letter, one number, and one special character
	if input.NoHashedPassword != "" && !input.IsValidPassword() {
		validationErrors = append(validationErrors, uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.PASSWORD_UNDERMINED_STRENGTH,
		))
	}
	return len(validationErrors) == 0, validationErrors
}

func NewCreatePasswordUseCase(
	log contracts.ILoggerProvider,
	repo contracts_repositories.IPasswordRepository,
	hashProvider contracts.IHashProvider,
) *CreatePasswordUseCase {
	return &CreatePasswordUseCase{
		appMessages:  locales.NewLocale(locales.EN_US),
		log:          log,
		repo:         repo,
		hashProvider: hashProvider,
	}
}
