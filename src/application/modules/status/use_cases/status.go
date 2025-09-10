package usecases

import (
	"context"
	"time"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	locales "gormgoskeleton/src/application/shared/locales"
	messages "gormgoskeleton/src/application/shared/locales/messages"
	status "gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	models "gormgoskeleton/src/domain/models"
)

type GetStatusUseCase struct {
	appMessages       *locales.Locale
	log               contractsProviders.ILoggerProvider
	apiStatusProvider contractsProviders.IApiStatusProvider
	locale            locales.LocaleTypeEnum
}

var _ usecase.BaseUseCase[time.Time, models.Status] = (*GetStatusUseCase)(nil)

func (uc *GetStatusUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.locale = locale
	}
}

func (uc *GetStatusUseCase) Execute(ctx context.Context,
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
	apiStatusProvider contractsProviders.IApiStatusProvider,
) *GetStatusUseCase {
	return &GetStatusUseCase{
		appMessages:       locales.NewLocale(locales.EN_US),
		log:               log,
		apiStatusProvider: apiStatusProvider,
	}
}
