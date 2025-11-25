package usecases_user

import (
	"context"
	"go/types"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	"gormgoskeleton/src/application/shared/guards"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
)

type DeleteUserUseCase struct {
	usecase.BaseUseCaseValidation[uint, types.Nil]
	log  contractsProviders.ILoggerProvider
	repo contracts_repositories.IUserRepository
}

var _ usecase.BaseUseCase[uint, types.Nil] = (*DeleteUserUseCase)(nil)

func (uc *DeleteUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input uint,
) *usecase.UseCaseResult[types.Nil] {
	result := usecase.NewUseCaseResult[types.Nil]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	err := uc.repo.SoftDelete(input)
	if err != nil {
		uc.log.Error("Error deleting user", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return result
	}
	result.SetData(
		status.Success,
		types.Nil{},
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_DELETE_SUCCESS,
		),
	)
	return result
}

func NewDeleteUserUseCase(
	log contractsProviders.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[uint, types.Nil]{
			Guards:      usecase.NewGuards(guards.RoleGuard("admin", "user"), guards.UserGetItSelf),
			AppMessages: locales.NewLocale(locales.EN_US),
		},
		log:  log,
		repo: repo,
	}
}
