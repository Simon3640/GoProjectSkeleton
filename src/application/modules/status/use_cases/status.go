package usecases

import (
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	statuscontracts "github.com/simon3640/goprojectskeleton/src/application/modules/status/contracts"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	locales "github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	messages "github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	status "github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	models "github.com/simon3640/goprojectskeleton/src/domain/models"
)

type GetStatusUseCase struct {
	appMessages       *locales.Locale
	log               contractsProviders.ILoggerProvider
	apiStatusProvider statuscontracts.IApiStatusProvider
	locale            locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[time.Time, models.Status] = (*GetStatusUseCase)(nil)

func (uc *GetStatusUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *GetStatusUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input time.Time,
) *usecase.UseCaseResult[models.Status] {
	result := usecase.NewUseCaseResult[models.Status]()
	uc.SetLocale(locale)

	result.SetData(status.Success,
		uc.apiStatusProvider.Get(input),
		uc.appMessages.Get(
			uc.locale,
			messages.MessageKeysInstance.APPLICATION_STATUS_OK,
		))
	return result
}

func NewGetStatusUseCase(
	log contractsProviders.ILoggerProvider,
	apiStatusProvider statuscontracts.IApiStatusProvider,
) *GetStatusUseCase {
	return &GetStatusUseCase{
		appMessages:       locales.NewLocale(locales.EN_US),
		log:               log,
		apiStatusProvider: apiStatusProvider,
	}
}
