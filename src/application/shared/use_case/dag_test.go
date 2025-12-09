package usecase

import (
	"context"
	"strconv"
	"testing"

	"github.com/simon3640/goprojectskeleton/src/application/shared/status"

	"github.com/simon3640/goprojectskeleton/src/application/shared/locales"

	"github.com/stretchr/testify/assert"
)

type UCStringToInt struct{}

func (uc *UCStringToInt) SetLocale(locale locales.LocaleTypeEnum) {}
func (uc *UCStringToInt) Execute(ctx context.Context, locale locales.LocaleTypeEnum, input string) *UseCaseResult[int] {
	result := NewUseCaseResult[int]()
	result.SetData(status.Success, 42, "Converted string to int")
	return result
}

var _ BaseUseCase[string, int] = (*UCStringToInt)(nil)

type UCIntExponent struct{}

func (uc *UCIntExponent) SetLocale(locale locales.LocaleTypeEnum) {}
func (uc *UCIntExponent) Execute(ctx context.Context, locale locales.LocaleTypeEnum, input int) *UseCaseResult[int] {
	result := NewUseCaseResult[int]()
	result.SetData(status.Success, input*2, "Calculated exponent")
	return result
}

type UCIntToString struct{}

func (uc *UCIntToString) SetLocale(locale locales.LocaleTypeEnum) {}
func (uc *UCIntToString) Execute(ctx context.Context, locale locales.LocaleTypeEnum, input int) *UseCaseResult[string] {
	result := NewUseCaseResult[string]()
	result.SetData(status.Success, strconv.Itoa(input), "Converted int to string")
	return result
}

func TestDagExecution(t *testing.T) {
	assert := assert.New(t)

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	UC3 := &UCIntToString{}

	dag := NewDag(NewStep(UC1), locales.EN_US, context.Background())
	dag2 := Then(dag, NewStep(UC2))
	dag3 := Then(dag2, NewStep(UC3))

	input := "5"
	result := dag3.Execute(input)
	// assert.Nil(err)
	assert.NotNil(result)

}

func TestDagConcurrentExecution(t *testing.T) {
	assert := assert.New(t)

	UC1 := &UCStringToInt{}
	UC2 := &UCIntExponent{}
	UC3 := &UCIntToString{}
	ParallelUC := NewUseCaseParallelDag[string, int]()
	ParallelUC.Usecases = []BaseUseCase[string, int]{UC1, UC1, UC1, UC1, UC1}
	dag := NewDag(NewStep(UC1), locales.EN_US, context.Background())
	dag2 := Then(dag, NewStep(UC2))
	dag3 := Then(dag2, NewStep(UC3))
	dagParallel := Then(dag3, NewStep(ParallelUC))

	input := "5"
	result := dagParallel.Execute(input)
	// assert.Nil(err)
	assert.NotNil(result)
	assert.Equal(5, len(*result.Data))
	for _, val := range *result.Data {
		assert.Equal(42, val)
	}
}
