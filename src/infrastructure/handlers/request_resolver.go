package handlers

import (
	"encoding/json"
	"net/http"

	"gormgoskeleton/src/application/shared/status"
	usecase "gormgoskeleton/src/application/shared/use_case"
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
	w http.ResponseWriter,
	result *usecase.UseCaseResult[D],
	headersToAdd map[HTTPHeaderTypeEnum]string,
) {
	rr.setHeaders(w, headersToAdd)

	if result.HasError() {
		w.WriteHeader(rr.statusMapping[result.StatusCode])
		resp := map[string]any{
			"details": result.Error,
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	if result.Data == nil && !result.HasError() {
		w.WriteHeader(204)
		return
	}

	w.WriteHeader(rr.statusMapping[result.StatusCode])
	resp := map[string]any{
		"data":    result.Data,
		"details": result.Details,
	}
	json.NewEncoder(w).Encode(resp)
}

func (rr *RequestResolver[D]) setHeaders(
	w http.ResponseWriter, headersToAdd map[HTTPHeaderTypeEnum]string,
) {
	for key, value := range headersToAdd {
		w.Header().Set(key.String(), value)
	}
}
