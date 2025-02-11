package usecase

import (
	"gingoskeleton/src/application/shared/status"
	"strconv"
)

type UseCaseResult[T any] struct {
	Data       *T
	Error      *string
	Headers    map[string]string
	Details    string
	StatusCode status.ApplicationStatusEnum
	Success    bool
}

func NewUseCaseResult[T any]() *UseCaseResult[T] {
	return &UseCaseResult[T]{
		StatusCode: status.InternalError,
		Success:    false,
		Headers:    make(map[string]string),
	}
}

func (ucr *UseCaseResult[T]) GetHeaders() map[string]string {
	return ucr.Headers
}

func (ucr *UseCaseResult[T]) GetStatusCode() status.ApplicationStatusEnum {
	return ucr.StatusCode
}

func (ucr *UseCaseResult[T]) SetStatusCode(statusCode status.ApplicationStatusEnum) {
	ucr.StatusCode = statusCode
}

func (ucr *UseCaseResult[T]) SetSuccess(success bool) {
	ucr.StatusCode = status.Success
	ucr.Success = success
}

func (ucr *UseCaseResult[T]) SetOnlyData(data T) *UseCaseResult[T] {
	ucr.Data = &data
	return ucr
}

func (ucr *UseCaseResult[T]) SetDetails(details string) *UseCaseResult[T] {
	ucr.Details = details
	return ucr
}

func (ucr *UseCaseResult[T]) AddHeader(key string, value string) *UseCaseResult[T] {
	ucr.Headers[key] = value
	return ucr
}

func (ucr *UseCaseResult[T]) SetData(
	statusCode status.ApplicationStatusEnum,
	data T,
	details string,
) *UseCaseResult[T] {
	ucr.StatusCode = statusCode
	ucr.Data = &data
	ucr.Details = details
	ucr.Success = true
	return ucr
}

func (ucr *UseCaseResult[T]) SetError(
	statusCode status.ApplicationStatusEnum,
	err string,
) *UseCaseResult[T] {
	ucr.StatusCode = statusCode
	ucr.Error = &err
	ucr.Success = false
	return ucr
}

func (ucr *UseCaseResult[T]) ToResultDTO() map[string]string {
	result := map[string]string{
		"statusCode": string(ucr.StatusCode),
		"success":    strconv.FormatBool(ucr.Success),
		"details":    string(ucr.Details),
	}

	if ucr.Error != nil {
		result["error"] = *ucr.Error
	}

	if ucr.Data != nil {
		result["data"] = "Data"
	}

	return result

}

func (ucr *UseCaseResult[T]) IsSuccess() bool {
	return ucr.Success
}

func (ucr *UseCaseResult[T]) GetData() *T {
	return ucr.Data
}

func (ucr *UseCaseResult[T]) HasError() bool {
	return ucr.Error != nil
}
