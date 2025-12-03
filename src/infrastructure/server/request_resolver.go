package api

import (
	"goprojectskeleton/src/application/shared/status"
	usecase "goprojectskeleton/src/application/shared/use_case"

	gin "github.com/gin-gonic/gin"
)

type RequestResolver[D any] struct {
	statusMapping map[status.ApplicationStatusEnum]int
}

func NewRequestResolver[D any]() *RequestResolver[D] {
	return &RequestResolver[D]{
		statusMapping: map[status.ApplicationStatusEnum]int{
			status.Success:                   200,
			status.Updated:                   200,
			status.Created:                   201,
			status.PartialContent:            206,
			status.InvalidInput:              400,
			status.Unauthorized:              401,
			status.NotFound:                  404,
			status.Conflict:                  409,
			status.TooManyRequests:            429,
			status.InternalError:             500,
			status.NotImplemented:            501,
			status.ProviderError:             502,
			status.ChatProviderError:         502,
			status.ProviderEmptyResponse:     502,
			status.ProviderEmptyCacheContext: 502,
		},
	}
}

func (rr *RequestResolver[D]) ResolveDTO(
	ctx *gin.Context,
	result *usecase.UseCaseResult[D],
	headersToAdd map[HTTPHeaderTypeEnum]string,
) (*gin.H, int) {
	content := gin.H{}

	if result.HasError() {
		content["details"] = result.Error
		rr.getHeaders(ctx, headersToAdd)

		return &content, rr.statusMapping[result.StatusCode]
	}

	if result.Data == nil && !result.HasError() {
		rr.getHeaders(ctx, headersToAdd)

		return nil, 204
	}

	content["data"] = result.Data
	content["details"] = result.Details
	rr.getHeaders(ctx, headersToAdd)

	return &content, rr.statusMapping[result.StatusCode]
}

func (rr *RequestResolver[D]) getHeaders(
	ctx *gin.Context, headersToAdd map[HTTPHeaderTypeEnum]string,
) {
	for key, value := range headersToAdd {
		ctx.Header(key.String(), value)
	}
}
