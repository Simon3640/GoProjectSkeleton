// Package userusecases provides use cases for user management
package userusecases

import (
	"context"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
)

type GetUserUseCase struct {
	usecase.BaseUseCaseValidation[uint, models.User]
	log  contractsProviders.ILoggerProvider
	repo usercontracts.IUserRepository
}

func (uc *GetUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *GetUserUseCase) Execute(ctx context.Context,
	locale locales.LocaleTypeEnum,
	input uint,
) *usecase.UseCaseResult[models.User] {
	result := usecase.NewUseCaseResult[models.User]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}
	uc.GetUser(ctx, result, input)
	return result
}

func (uc *GetUserUseCase) GetUser(ctx context.Context, result *usecase.UseCaseResult[models.User], id uint) {
	res, err := uc.repo.GetByID(id)
	if err != nil {
		uc.log.Error("Error getting user by ID", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return
	}
	result.SetData(
		status.Success,
		*res,
		"",
	)
}

func NewGetUserUseCase(
	log contractsProviders.ILoggerProvider,
	repo usercontracts.IUserRepository,
) *GetUserUseCase {
	return &GetUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[uint, models.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin", "user"), guards.UserGetItSelf),
		},
		log:  log,
		repo: repo,
	}
}
