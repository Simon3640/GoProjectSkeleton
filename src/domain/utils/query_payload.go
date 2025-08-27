package domaim_utils

import (
	"fmt"
	"reflect"
	"strings"
)

type FilterOperator string

const (
	OperatorEqual        FilterOperator = "eq"
	OperatorNotEqual     FilterOperator = "ne"
	OperatorGreaterThan  FilterOperator = "gt"
	OperatorGreaterEqual FilterOperator = "ge"
	OperatorLessThan     FilterOperator = "lt"
	OperatorLessEqual    FilterOperator = "le"
	OperatorLike         FilterOperator = "like"
	OperatorILike        FilterOperator = "ilike"
	OperatorIn           FilterOperator = "in"
	OperatorNotIn        FilterOperator = "not_in"
	OperatorIsNull       FilterOperator = "is_null"
	OperatorIsNotNull    FilterOperator = "is_not_null"
)

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type Filter struct {
	Field    string
	Operator FilterOperator
	Value    *string
}

// Validate the filter
func (f Filter) Validate() []string {
	var errors []string

	if f.Field == "" {
		errors = append(errors, "filter field is required")
	}

	if f.Operator == "" {
		errors = append(errors, "filter operator is required")
	}

	// Para algunos operadores, el valor no es requerido
	if f.Operator != OperatorIsNull && f.Operator != OperatorIsNotNull && f.Value == nil {
		errors = append(errors, "filter value is required for this operator")
	}

	return errors
}

// Sort representa un campo de ordenamiento
type Sort struct {
	Field string
	Order SortOrder
}

// Validate valida que el sort tenga los campos requeridos
func (s Sort) Validate() []string {
	var errors []string

	if s.Field == "" {
		errors = append(errors, "sort field is required")
	}

	if s.Order != SortAsc && s.Order != SortDesc {
		errors = append(errors, "sort order must be 'asc' or 'desc'")
	}

	return errors
}

// Pagination representa los parámetros de paginación
type Pagination struct {
	Page     int
	PageSize int
}

// Validate valida los parámetros de paginación
func (p Pagination) Validate() []string {
	var errors []string

	if p.Page < 1 {
		errors = append(errors, "page must be greater than 0")
	}

	if p.PageSize < 1 {
		errors = append(errors, "page_size must be greater than 0")
	}

	if p.PageSize > 1000 {
		errors = append(errors, "page_size must be less than or equal to 1000")
	}

	return errors
}

// GetOffset calcula el offset para GORM
func (p Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit retorna el límite para GORM
func (p Pagination) GetLimit() int {
	return p.PageSize
}

// QueryPayload es el payload genérico para consultas con filtros, ordenamiento y paginación
type QueryPayloadBuilder[DBModel any] struct {
	Filters    []Filter
	Sorts      []Sort
	Pagination Pagination
}

// Validate valida todo el payload
func (qp QueryPayloadBuilder[DBModel]) Validate() []string {
	var errors []string

	// Validar filtros
	for i, filter := range qp.Filters {
		if filterErrors := filter.Validate(); len(filterErrors) > 0 {
			for _, err := range filterErrors {
				errors = append(errors, fmt.Sprintf("filter[%d]: %s", i, err))
			}
		}
	}

	// Validar sorts
	for i, sort := range qp.Sorts {
		if sortErrors := sort.Validate(); len(sortErrors) > 0 {
			for _, err := range sortErrors {
				errors = append(errors, fmt.Sprintf("sort[%d]: %s", i, err))
			}
		}
	}

	// Validar paginación
	if paginationErrors := qp.Pagination.Validate(); len(paginationErrors) > 0 {
		for _, err := range paginationErrors {
			errors = append(errors, fmt.Sprintf("pagination: %s", err))
		}
	}

	return errors
}

func (qp *QueryPayloadBuilder[DBModel]) ParseFilter(filter string) Filter {
	// <column>:<operator>:<value>
	parts := strings.SplitN(filter, ":", 3)
	if len(parts) < 3 {
		return Filter{}
	}
	operator := FilterOperator(parts[1])
	value := parts[2]
	// field on DBModel?
	var field string
	var found bool
	field, found = getFieldName(reflect.TypeOf(*new(DBModel)), parts[0])
	if !found {
		return Filter{}
	}
	return Filter{
		Field:    field,
		Operator: operator,
		Value:    &value,
	}
}

func getFieldName(t reflect.Type, search string) (string, bool) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// Si es una struct embebida, recursivamente busca en ella
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			if name, found := getFieldName(field.Type, search); found {
				return name, true
			}
		}
		if strings.EqualFold(field.Name, search) {
			return field.Name, true
		}
	}
	return "", false
}

func (qp *QueryPayloadBuilder[DBModel]) ParseFilters(filter []string) {
	for _, f := range filter {
		qp.Filters = append(qp.Filters, qp.ParseFilter(f))
	}
}

func (qp *QueryPayloadBuilder[DBModel]) ParseSort(sort string) Sort {
	// <field>:<order>
	parts := strings.SplitN(sort, ":", 2)
	if len(parts) < 2 {
		return Sort{}
	}
	// field on DBModel?
	var field string
	var found bool
	field, found = getFieldName(reflect.TypeOf(*new(DBModel)), parts[0])
	if !found {
		return Sort{}
	}
	return Sort{
		Field: field,
		Order: SortOrder(parts[1]),
	}
}

func (qp *QueryPayloadBuilder[DBModel]) ParseSorts(sort []string) {
	for _, s := range sort {
		qp.Sorts = append(qp.Sorts, qp.ParseSort(s))
	}
}
