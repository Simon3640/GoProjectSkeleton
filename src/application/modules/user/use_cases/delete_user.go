package usecases_user

import (
	"context"
	"go/types"
	"gormgoskeleton/src/application/contracts"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"strings"
)

type DeleteUserUseCase struct {
	appMessages *locales.Locale
	log         contracts.ILoggerProvider
	repo        contracts_repositories.IUserRepository
	locale      locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[uint, types.Nil] = (*DeleteUserUseCase)(nil)

func (uc *DeleteUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input uint,
) *usecase.UseCaseResult[types.Nil] {
	validation, val_msgs := uc.validate(input)
	result := usecase.NewUseCaseResult[types.Nil]()
	uc.SetLocale(locale)
	if !validation {
		result.SetError(
			status.InvalidInput,
			strings.Join(val_msgs, "\n"),
		)
		return result
	}
	err := uc.repo.Delete(input)
	if err != nil {
		result.SetError(
			status.Conflict,
			uc.appMessages.Get(
				uc.locale,
				messages.MessageKeysInstance.SOMETHING_WENT_WRONG,
			),
		)
	}
	result.SetData(
		status.Success,
		types.Nil{},
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.USER_DELETE_SUCCESS,
		),
	)
	return result
}

func (uc *DeleteUserUseCase) validate(input uint) (bool, []string) {
	var val_msgs []string
	if input <= 0 {
		val_msgs = append(val_msgs, uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.INVALID_USER_ID,
		))
	}
	return len(val_msgs) == 0, val_msgs
}

func NewDeleteUserUseCase(
	log contracts.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		appMessages: locales.NewLocale(locales.EN_US),
		log:         log,
		repo:        repo,
	}
}
