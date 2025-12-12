package userusecases

import (
	contractsproviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
)

// DeleteUserUseCase is a use case that deletes a user
type DeleteUserUseCase struct {
	usecase.BaseUseCaseValidation[uint, bool]
	log  contractsproviders.ILoggerProvider
	repo usercontracts.IUserRepository
}

var _ usecase.BaseUseCase[uint, bool] = (*DeleteUserUseCase)(nil)

// SetLocale sets the locale for the use case
func (uc *DeleteUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

// Execute executes the use case
func (uc *DeleteUserUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input uint,
) *usecase.UseCaseResult[bool] {
	result := usecase.NewUseCaseResult[bool]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	uc.deleteUser(input, result)
	if result.HasError() {
		return result
	}

	result.SetData(
		status.Success,
		true,
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_DELETE_SUCCESS,
		),
	)
	return result
}

func (uc *DeleteUserUseCase) deleteUser(id uint, result *usecase.UseCaseResult[bool]) {
	err := uc.repo.SoftDelete(id)
	if err != nil {
		uc.log.Error("Error deleting user", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(uc.Locale, err.Context),
		)
	}
}

// NewDeleteUserUseCase creates a new delete user use case
func NewDeleteUserUseCase(
	log contractsproviders.ILoggerProvider,
	repo usercontracts.IUserRepository,
) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[uint, bool]{
			Guards:      usecase.NewGuards(guards.RoleGuard("admin", "user"), guards.UserGetItSelf),
			AppMessages: locales.NewLocale(locales.EN_US),
		},
		log:  log,
		repo: repo,
	}
}
