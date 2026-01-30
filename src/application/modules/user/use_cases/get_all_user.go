package userusecases

import (
	"encoding/json"
	"time"

	contractsProviders "github.com/simon3640/goprojectskeleton/src/application/contracts/providers"
	usercontracts "github.com/simon3640/goprojectskeleton/src/application/modules/user/contracts"
	userdtos "github.com/simon3640/goprojectskeleton/src/application/modules/user/dtos"
	shareddtos "github.com/simon3640/goprojectskeleton/src/application/shared/DTOs"
	"github.com/simon3640/goprojectskeleton/src/application/shared/cache"
	app_context "github.com/simon3640/goprojectskeleton/src/application/shared/context"
	applicationerrors "github.com/simon3640/goprojectskeleton/src/application/shared/errors"
	"github.com/simon3640/goprojectskeleton/src/application/shared/guards"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"
	"github.com/simon3640/goprojectskeleton/src/application/shared/locales/messages"
	"github.com/simon3640/goprojectskeleton/src/application/shared/observability"
	"github.com/simon3640/goprojectskeleton/src/application/shared/settings"
	"github.com/simon3640/goprojectskeleton/src/application/shared/status"
	usecase "github.com/simon3640/goprojectskeleton/src/application/shared/use_case"
	domainutils "github.com/simon3640/goprojectskeleton/src/domain/shared/utils"
	usermodels "github.com/simon3640/goprojectskeleton/src/domain/user/models"
)

type inputPayload = domainutils.QueryPayloadBuilder[usermodels.User]

// GetAllUserUseCase is a use case that gets all users
type GetAllUserUseCase struct {
	usecase.BaseUseCaseValidation[inputPayload, userdtos.UserMultiResponse]
	repo          usercontracts.IUserRepository
	cacheExecutor *cache.Executor[userdtos.UserMultiResponse, inputPayload]
}

var _ usecase.BaseUseCase[inputPayload, userdtos.UserMultiResponse] = (*GetAllUserUseCase)(nil)

// Execute ejecuta el caso de uso:
// - Validate input and guards
// - Execute the use case
// - Set the success result
// - Return the result
// - If there is an error, set the error result
// - Return the result
func (uc *GetAllUserUseCase) Execute(
	ctx *app_context.AppContext,
	locale locales.LocaleTypeEnum,
	input inputPayload,
) *usecase.UseCaseResult[userdtos.UserMultiResponse] {
	result := usecase.NewUseCaseResult[userdtos.UserMultiResponse]()
	uc.SetLocale(locale)
	uc.SetAppContext(ctx)
	uc.Validate(input, result)
	if result.HasError() {
		return result
	}

	cachePolicy := uc.createCachePolicy()
	multipleResponse, err := uc.cacheExecutor.Execute(cachePolicy, input, ctx)
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error getting users", err.ToError(), uc.AppContext)
		result.SetError(err.Code, uc.AppMessages.Get(uc.Locale, err.Context))
		return result
	}
	uc.setSuccessResult(result, *multipleResponse)

	return result
}

// createCachePolicy creates the cache policy for the use case
func (uc *GetAllUserUseCase) createCachePolicy() cache.ICachePolicy[userdtos.UserMultiResponse, inputPayload] {
	return &getAllUserCachePolicy{uc: uc}
}

// getAllUserCachePolicy implements cache.ICachePolicy for GetAllUserUseCase
type getAllUserCachePolicy struct {
	uc *GetAllUserUseCase
}

func (p *getAllUserCachePolicy) BuildKey(input inputPayload, _ *app_context.AppContext) string {
	return "users:" + input.GetQueryKey()
}

func (p *getAllUserCachePolicy) FetchData(input inputPayload, _ *app_context.AppContext) (userdtos.UserMultiResponse, *applicationerrors.ApplicationError) {
	data, total, err := p.uc.repo.GetAll(&input, input.Pagination.GetOffset(), input.Pagination.GetLimit())
	if err != nil {
		observability.GetObservabilityComponents().Logger.ErrorWithContext("Error getting all users", err.ToError(), p.uc.AppContext)
		return userdtos.UserMultiResponse{}, err
	}
	return p.uc.buildMultiResponse(data, total, input, false), nil
}

func (p *getAllUserCachePolicy) Serialize(data userdtos.UserMultiResponse) (any, *applicationerrors.ApplicationError) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, applicationerrors.NewApplicationError(status.InternalError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
	}
	return bytes, nil
}

func (p *getAllUserCachePolicy) Deserialize(data any) (userdtos.UserMultiResponse, *applicationerrors.ApplicationError) {
	switch v := data.(type) {
	case userdtos.UserMultiResponse:
		return v, nil
	case []byte:
		var resp userdtos.UserMultiResponse
		if err := json.Unmarshal(v, &resp); err != nil {
			return userdtos.UserMultiResponse{}, applicationerrors.NewApplicationError(status.InternalError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
		}
		return resp, nil
	default:
		bytes, err := json.Marshal(data)
		if err != nil {
			return userdtos.UserMultiResponse{}, applicationerrors.NewApplicationError(status.InternalError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
		}
		var resp userdtos.UserMultiResponse
		if err := json.Unmarshal(bytes, &resp); err != nil {
			return userdtos.UserMultiResponse{}, applicationerrors.NewApplicationError(status.InternalError, messages.MessageKeysInstance.SOMETHING_WENT_WRONG, err.Error())
		}
		return resp, nil
	}
}

func (p *getAllUserCachePolicy) GetTTL() time.Duration {
	return time.Duration(settings.AppSettingsInstance.RedisTTL) * time.Second
}

func (p *getAllUserCachePolicy) OnCacheHit(data userdtos.UserMultiResponse) *userdtos.UserMultiResponse {
	data.Meta.Cached = true
	return &data
}

// buildMultiResponse builds the multi response for the users
func (uc *GetAllUserUseCase) buildMultiResponse(
	data []usermodels.User, total int64,
	input inputPayload,
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

func (uc *GetAllUserUseCase) setSuccessResult(
	result *usecase.UseCaseResult[userdtos.UserMultiResponse],
	multipleResponse userdtos.UserMultiResponse,
) {
	result.SetData(
		status.Success,
		multipleResponse,
		uc.AppMessages.Get(uc.Locale, messages.MessageKeysInstance.USER_LIST_SUCCESS),
	)
}

// NewGetAllUserUseCase creates a new get all user use case
func NewGetAllUserUseCase(
	repo usercontracts.IUserRepository,
	cacheProvider contractsProviders.ICacheProvider,
) *GetAllUserUseCase {
	return &GetAllUserUseCase{
		BaseUseCaseValidation: usecase.BaseUseCaseValidation[inputPayload, userdtos.UserMultiResponse]{
			AppMessages: locales.NewLocale(locales.EN_US),
			Guards:      usecase.NewGuards(guards.RoleGuard("admin")),
		},
		repo:          repo,
		cacheExecutor: cache.NewExecutor[userdtos.UserMultiResponse, inputPayload](cacheProvider),
	}
}
