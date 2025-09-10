package usecases_user

import (
	"context"
	"time"

	contractsProviders "gormgoskeleton/src/application/contracts/providers"
	contracts_repositories "gormgoskeleton/src/application/contracts/repositories"
	dtos "gormgoskeleton/src/application/shared/DTOs"
	"gormgoskeleton/src/application/shared/guards"
	"gormgoskeleton/src/application/shared/locales"
	"gormgoskeleton/src/application/shared/locales/messages"
	"gormgoskeleton/src/application/shared/settings"
	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
	"gormgoskeleton/src/domain/models"
	domain_utils "gormgoskeleton/src/domain/utils"
)

type GetAllUserUseCase struct {
	usecase.BaseUseCaseValidation[domain_utils.QueryPayloadBuilder[models.User], dtos.UserMultiResponse]
	log   contractsProviders.ILoggerProvider
	repo  contracts_repositories.IUserRepository
	cache contractsProviders.ICacheProvider
}

var _ usecase.BaseUseCase[domain_utils.QueryPayloadBuilder[models.User], dtos.UserMultiResponse] = (*GetAllUserUseCase)(nil)

func (uc *GetAllUserUseCase) SetLocale(locale locales.LocaleTypeEnum) {
	if locale != "" {
		uc.Locale = locale
	}
}

func (uc *GetAllUserUseCase) Execute(
	ctx context.Context,
	locale locales.LocaleTypeEnum,
	input domain_utils.QueryPayloadBuilder[models.User],
) *usecase.UseCaseResult[dtos.UserMultiResponse] {
	result := usecase.NewUseCaseResult[dtos.UserMultiResponse]()
	uc.SetLocale(locale)
	uc.Validate(ctx, input, result)
	if result.HasError() {
		return result
	}

	// Check Cache
	cacheKey := "users:" + input.GetQueryKey()
	var data []models.User
	found, err := uc.cache.Get(cacheKey, &data)

	if err != nil {
		uc.log.Error("Error getting cache for users", err.ToError())
	}

	if found {
		var total int64
		found, err = uc.cache.Get(cacheKey+":total", &total)
		if err != nil {
			uc.log.Error("Error getting cache for users total", err.ToError())
		}
		uc.log.Debug("Cache hit for users", map[string]any{"cacheKey": cacheKey, "total": total})
		if found {
			result.SetData(
				status.Success,
				uc.buildMultiRespose(data, total, input, true),
				uc.AppMessages.Get(
					uc.Locale,
					messages.MessageKeysInstance.USER_LIST_SUCCESS,
				),
			)
			return result
		}
	}

	data, total, err := uc.repo.GetAll(&input, input.Pagination.GetOffset(), input.Pagination.GetLimit())
	if err != nil {
		uc.log.Error("Error getting all users", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(
				uc.Locale,
				err.Context,
			),
		)
		return result
	}

	// Set Cache
	if err := uc.cache.Set(cacheKey, data, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second); err != nil {
		uc.log.Error("Error setting cache for users", err.ToError())
	}
	if err := uc.cache.Set(cacheKey+":total", total, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second); err != nil {
		uc.log.Error("Error setting cache for users total", err.ToError())
	}

	// Build MultiResponse
	result.SetData(
		status.Success,
		uc.buildMultiRespose(data, total, input, false),
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_LIST_SUCCESS,
		),
	)
	return result
}

func (uc *GetAllUserUseCase) buildMultiRespose(data []models.User, total int64, input domain_utils.QueryPayloadBuilder[models.User], cached bool) dtos.UserMultiResponse {
	var response dtos.UserMultiResponse
	response.Records = data
	hasNext, hasPrev := input.HasNextPrev(total)
	response.Meta = dtos.NewMetaMultiResponse(len(data), total, hasNext, hasPrev, cached)
	response.Meta.BuildLinks(
		"/user",
		input.Pagination.Page,
		input.Pagination.PageSize, input.BuildQueryParamsURL(),
	)
	return response
}

func NewGetAllUserUseCase(
	log contractsProviders.ILoggerProvider,
	repo contracts_repositories.IUserRepository,
	cache contractsProviders.ICacheProvider,
) *GetAllUserUseCase {
	return &GetAllUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[domain_utils.QueryPayloadBuilder[models.User], dtos.UserMultiResponse]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin")),
		},
		log:   log,
		repo:  repo,
		cache: cache,
	}
}
