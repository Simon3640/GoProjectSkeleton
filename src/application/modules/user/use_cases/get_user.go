// Package userusecases provides use cases for user management
package userusecases

import (
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

type GetUserUseCase struct {
	usecase.BaseUseCaseValidation[uint, usermodels.User]
	repo usercontracts.IUserRepository
}

func (uc *GetUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *GetUserUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input uint,
) *usecase.UseCaseResult[usermodels.User] {
	result := usecase.NewUseCaseResult[usermodels.User]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.Validate(input, result)
	if result.HasError() {
		return result
	}
	res := uc.getUser(result, input)
	if result.HasError() {
		return result
	}
	result.SetData(status.Success, *res, "")
	observability.GetObservabilityComponents().Logger.InfoWithContext("User retrieved successfully", uc.AppContext)
	return result
}

func (uc *GetUserUseCase) getUser(result *usecase.UseCaseResult[usermodels.User], id uint) *usermodels.User {
	res, err := uc.repo.GetByID(id)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error getting user by ID", err.ToError(), uc.AppContext)
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
	}
	return res
}

func NewGetUserUseCase(
	repo usercontracts.IUserRepository,
) *GetUserUseCase {
	return &GetUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[uint, usermodels.User]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin", "user"), guards.UserGetItSelf),
		},
		repo: repo,
	}
}
