package userusecases

import (
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	"github.com/simon3640/goprojectskeleton/src/domain/models"
	domainutils "github.com/simon3640/goprojectskeleton/src/domain/utils"
)

// GetAllUserUseCase is a use case that gets all users
type GetAllUserUseCase struct {
	usecase.BaseUseCaseValidation[domainutils.QueryPayloadBuilder[models.User], userdtos.UserMultiResponse]
	log   contractsProviders.ILoggerProvider
	repo  usercontracts.IUserRepository
	cache contractsProviders.ICacheProvider
}

var _ usecase.BaseUseCase[domainutils.QueryPayloadBuilder[models.User], userdtos.UserMultiResponse] = (*GetAllUserUseCase)(nil)

// Execute executes the use case
func (uc *GetAllUserUseCase) Execute(
	ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input domainutils.QueryPayloadBuilder[models.User],
) *usecase.UseCaseResult[userdtos.UserMultiResponse] {
	result := usecase.NewUseCaseResult[userdtos.UserMultiResponse]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.Validate(input, result)
	if result.HasError() {
		return result
	}

	uc.getUsersFromCache(input, result)
	if result.Data != nil {
		return result
	}

	data, total, err := uc.getUsersFromRepository(input, result)
	if err != nil {
		uc.log.Error("Error getting all users from repository", err.ToError())
		return result
	}

	uc.setCache(input, data, total)

	// Build MultiResponse
	result.SetData(
		status.Success,
		uc.buildMultiResponse(data, total, input, false),
		uc.AppMessages.Get(
			uc.Locale,
			messages.MessageKeysInstance.USER_LIST_SUCCESS,
		),
	)
	return result
}

// buildMultiResponse builds the multi response for the users
// it builds the response with the data, total, meta and links
func (uc *GetAllUserUseCase) buildMultiResponse(
	data []models.User, total int64,
	input domainutils.QueryPayloadBuilder[models.User],
	cached bool,
) userdtos.UserMultiResponse {
	var response userdtos.UserMultiResponse
	response.Records = data
	hasNext, hasPrev := input.HasNextPrev(total)
	response.Meta = shareddtos.NewMetaMultiResponse(len(data), total, hasNext, hasPrev, cached)
	response.Meta.BuildLinks(
		"/user",
		input.Pagination.Page,
		input.Pagination.PageSize, input.BuildQueryParamsURL(),
	)
	return response
}

// getUsersFromCache gets the users from the cache
// it checks if the cache is hit and if it is, it sets the result with a complete UserMultiResponse object (including records and meta information such as total)
func (uc *GetAllUserUseCase) getUsersFromCache(
	input domainutils.QueryPayloadBuilder[models.User],
	result *usecase.UseCaseResult[userdtos.UserMultiResponse],
) {
	// Check Cache
	cacheKey := uc.cacheKey(input)
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
				uc.buildMultiResponse(data, total, input, true),
				uc.AppMessages.Get(
					uc.Locale,
					messages.MessageKeysInstance.USER_LIST_SUCCESS,
				),
			)
		}
	}
}

// cacheKey builds the cache key for the users
// it builds the key with the input query key
func (uc *GetAllUserUseCase) cacheKey(input domainutils.QueryPayloadBuilder[models.User]) string {
	return "users:" + input.GetQueryKey()
}

// getUsersFromRepository gets the users from the repository
// it returns the users and total, or sets an error in the result if the repository call fails
func (uc *GetAllUserUseCase) getUsersFromRepository(
	input domainutils.QueryPayloadBuilder[models.User],
	result *usecase.UseCaseResult[userdtos.UserMultiResponse],
) ([]models.User, int64, *applicationerrors.ApplicationError) {
	data, total, err := uc.repo.GetAll(&input, input.Pagination.GetOffset(), input.Pagination.GetLimit())
	if err != nil {
		uc.log.Error("Error getting all users", err.ToError())
		result.SetError(
			err.Code,
			uc.AppMessages.Get(uc.Locale, err.Context),
		)
		return nil, 0, err
	}
	return data, total, nil
}

// setCache sets the cache for the users
// it sets the cache for the users with the data and total
func (uc *GetAllUserUseCase) setCache(input domainutils.QueryPayloadBuilder[models.User], data []models.User, total int64) {
	if err := uc.cache.Set(uc.cacheKey(input), data, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second); err != nil {
		uc.log.Error("Error setting cache for users", err.ToError())
	}
	if err := uc.cache.Set(uc.cacheKey(input)+":total", total, time.Duration(settings.AppSettingsInstance.RedisTTL)*time.Second); err != nil {
		uc.log.Error("Error setting cache for users total", err.ToError())
	}
}

// NewGetAllUserUseCase creates a new get all user use case
func NewGetAllUserUseCase(
	log contractsProviders.ILoggerProvider,
	repo usercontracts.IUserRepository,
	cache contractsProviders.ICacheProvider,
) *GetAllUserUseCase {
	return &GetAllUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[domainutils.QueryPayloadBuilder[models.User], userdtos.UserMultiResponse]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin")),
		},
		log:   log,
		repo:  repo,
		cache: cache,
	}
}
