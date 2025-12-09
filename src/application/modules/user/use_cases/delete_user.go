package userusecases

import (
	"context"
	"go/types"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	contracts_repositories "github.com/simon3640/goprojectskeleton/src/application/contracts/repositories"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
)

// DeleteUserUseCase is a use case that deletes a user
type DeleteUserUseCase struct {
	usecase.BaseUseCaseValidation[uint, types.Nil]
	log  contractsProviders.ILoggerProvider
	repo contracts_repositories.IUserRepository
}

var _ usecase.BaseUseCase[uint, types.Nil] = (*DeleteUserUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *DeleteUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

// Execute executes the use case
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

// NewDeleteUserUseCase creates a new delete user use case
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
