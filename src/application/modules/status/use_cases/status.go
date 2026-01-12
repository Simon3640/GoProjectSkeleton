package usecases

import (
	"time"

	statuscontracts "github.com/simon3640/goprojectskeleton/src/application/modules/status/contracts"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	locales "github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	messages "github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	status "github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	statusmodels "github.com/simon3640/goprojectskeleton/src/domain/status/models"
)

type GetStatusUseCase struct {
	usecase.BaseUseCaseValidation[time.Time, statusmodels.Status]
	apiStatusProvider statuscontracts.IApiStatusProvider
}

func (uc *GetStatusUseCase) Execute(ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input time.Time,
) *usecase.UseCaseResult[statusmodels.Status] {
	result := usecase.NewUseCaseResult[statusmodels.Status]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)

	result.SetData(status.Success,
		uc.apiStatusProvider.Get(input),
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.APPLICATION_STATUS_OK,
		))
	observability.GetObservabilityComponents().Logger.InfoWithContext("status_retrieved", uc.AppContext)
	return result
}

func NewGetStatusUseCase(
	apiStatusProvider statuscontracts.IApiStatusProvider,
) *GetStatusUseCase {
	return &GetStatusUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[time.Time, statusmodels.Status]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(),
		},
		apiStatusProvider: apiStatusProvider,
	}
}
