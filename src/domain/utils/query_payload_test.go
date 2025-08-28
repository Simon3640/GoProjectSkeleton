package domaim_utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestModel is a mock model for testing
type TestModel struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func TestFilter_Validate(t *testing.T) {
	value := "test"

	tests := []struct {
		name           string
		filter         Filter
		expectedErrors int
	}{
		{
			name: "Valid filter",
			filter: Filter{
				Field:    "name",
				Operator: OperatorEqual,
				Value:    &value,
			},
			expectedErrors: 0,
		},
		{
			name: "Missing field",
			filter: Filter{
				Field:    "",
				Operator: OperatorEqual,
				Value:    &value,
			},
			expectedErrors: 1,
		},
		{
			name: "Missing operator",
			filter: Filter{
				Field:    "name",
				Operator: "",
				Value:    &value,
			},
			expectedErrors: 1,
		},
		{
			name: "Missing value for regular operator",
			filter: Filter{
				Field:    "name",
				Operator: OperatorEqual,
				Value:    nil,
			},
			expectedErrors: 1,
		},
		{
			name: "IsNull operator without value (valid)",
			filter: Filter{
				Field:    "name",
				Operator: OperatorIsNull,
				Value:    nil,
			},
			expectedErrors: 0,
		},
		{
			name: "IsNotNull operator without value (valid)",
			filter: Filter{
				Field:    "name",
				Operator: OperatorIsNotNull,
				Value:    nil,
			},
			expectedErrors: 0,
		},
		{
			name: "Multiple errors",
			filter: Filter{
				Field:    "",
				Operator: "",
				Value:    nil,
			},
			expectedErrors: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.filter.Validate()
			assert.Equal(t, tt.expectedErrors, len(errors))
		})
	}
}

func TestSort_Validate(t *testing.T) {
	tests := []struct {
		name           string
		sort           Sort
		expectedErrors int
	}{
		{
			name: "Valid sort ascending",
			sort: Sort{
				Field: "name",
				Order: SortAsc,
			},
			expectedErrors: 0,
		},
		{
			name: "Valid sort descending",
			sort: Sort{
				Field: "name",
				Order: SortDesc,
			},
			expectedErrors: 0,
		},
		{
			name: "Missing field",
			sort: Sort{
				Field: "",
				Order: SortAsc,
			},
			expectedErrors: 1,
		},
		{
			name: "Invalid order",
			sort: Sort{
				Field: "name",
				Order: "invalid",
			},
			expectedErrors: 1,
		},
		{
			name: "Multiple errors",
			sort: Sort{
				Field: "",
				Order: "invalid",
			},
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.sort.Validate()
			assert.Equal(t, tt.expectedErrors, len(errors))
		})
	}
}

func TestPagination_Validate(t *testing.T) {
	tests := []struct {
		name           string
		pagination     Pagination
		expectedErrors int
	}{
		{
			name: "Valid pagination",
			pagination: Pagination{
				Page:     1,
				PageSize: 10,
			},
			expectedErrors: 0,
		},
		{
			name: "Page less than 1",
			pagination: Pagination{
				Page:     0,
				PageSize: 10,
			},
			expectedErrors: 1,
		},
		{
			name: "PageSize less than 1",
			pagination: Pagination{
				Page:     1,
				PageSize: 0,
			},
			expectedErrors: 1,
		},
		{
			name: "PageSize greater than 1000",
			pagination: Pagination{
				Page:     1,
				PageSize: 1001,
			},
			expectedErrors: 1,
		},
		{
			name: "Multiple errors",
			pagination: Pagination{
				Page:     0,
				PageSize: 0,
			},
			expectedErrors: 2,
		},
		{
			name: "All errors",
			pagination: Pagination{
				Page:     -1,
				PageSize: 1001,
			},
			expectedErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.pagination.Validate()
			assert.Equal(t, tt.expectedErrors, len(errors))
		})
	}
}

func TestPagination_GetOffset(t *testing.T) {
	tests := []struct {
		name       string
		pagination Pagination
		expected   int
	}{
		{
			name: "First page",
			pagination: Pagination{
				Page:     1,
				PageSize: 10,
			},
			expected: 0,
		},
		{
			name: "Second page",
			pagination: Pagination{
				Page:     2,
				PageSize: 10,
			},
			expected: 10,
		},
		{
			name: "Fifth page with page size 20",
			pagination: Pagination{
				Page:     5,
				PageSize: 20,
			},
			expected: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pagination.GetOffset()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPagination_GetLimit(t *testing.T) {
	tests := []struct {
		name       string
		pagination Pagination
		expected   int
	}{
		{
			name: "Page size 10",
			pagination: Pagination{
				Page:     1,
				PageSize: 10,
			},
			expected: 10,
		},
		{
			name: "Page size 50",
			pagination: Pagination{
				Page:     1,
				PageSize: 50,
			},
			expected: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pagination.GetLimit()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQueryPayloadBuilder_Validate(t *testing.T) {
	value := "test"

	tests := []struct {
		name           string
		payload        QueryPayloadBuilder[TestModel]
		expectedErrors int
	}{
		{
			name: "Valid payload",
			payload: QueryPayloadBuilder[TestModel]{
				Filters: []Filter{
					{Field: "name", Operator: OperatorEqual, Value: &value},
				},
				Sorts: []Sort{
					{Field: "name", Order: SortAsc},
				},
				Pagination: Pagination{Page: 1, PageSize: 10},
			},
			expectedErrors: 0,
		},
		{
			name: "Invalid filter",
			payload: QueryPayloadBuilder[TestModel]{
				Filters: []Filter{
					{Field: "", Operator: OperatorEqual, Value: &value},
				},
				Sorts: []Sort{
					{Field: "name", Order: SortAsc},
				},
				Pagination: Pagination{Page: 1, PageSize: 10},
			},
			expectedErrors: 1,
		},
		{
			name: "Invalid sort",
			payload: QueryPayloadBuilder[TestModel]{
				Filters: []Filter{
					{Field: "name", Operator: OperatorEqual, Value: &value},
				},
				Sorts: []Sort{
					{Field: "", Order: SortAsc},
				},
				Pagination: Pagination{Page: 1, PageSize: 10},
			},
			expectedErrors: 1,
		},
		{
			name: "Invalid pagination",
			payload: QueryPayloadBuilder[TestModel]{
				Filters: []Filter{
					{Field: "name", Operator: OperatorEqual, Value: &value},
				},
				Sorts: []Sort{
					{Field: "name", Order: SortAsc},
				},
				Pagination: Pagination{Page: 0, PageSize: 10},
			},
			expectedErrors: 1,
		},
		{
			name: "Multiple errors",
			payload: QueryPayloadBuilder[TestModel]{
				Filters: []Filter{
					{Field: "", Operator: "", Value: nil},
					{Field: "name", Operator: OperatorEqual, Value: nil},
				},
				Sorts: []Sort{
					{Field: "", Order: "invalid"},
				},
				Pagination: Pagination{Page: 0, PageSize: 0},
			},
			expectedErrors: 8, // 5 filter errors + 2 sort errors + 2 pagination errors - 1 (because one filter error is covered by another)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := tt.payload.Validate()
			assert.Equal(t, tt.expectedErrors, len(errors))
		})
	}
}

func TestQueryPayloadBuilder_ParseFilter(t *testing.T) {
	qp := &QueryPayloadBuilder[TestModel]{}

	tests := []struct {
		name           string
		filterString   string
		expectedFilter Filter
	}{
		{
			name:         "Valid filter",
			filterString: "name:eq:John",
			expectedFilter: Filter{
				Field:    "Name",
				Operator: OperatorEqual,
				Value:    stringPtr("John"),
			},
		},
		{
			name:         "Valid filter with different operator",
			filterString: "age:gt:25",
			expectedFilter: Filter{
				Field:    "Age",
				Operator: OperatorGreaterThan,
				Value:    stringPtr("25"),
			},
		},
		{
			name:           "Invalid filter - not enough parts",
			filterString:   "Name:eq",
			expectedFilter: Filter{},
		},
		{
			name:           "Invalid filter - field not found",
			filterString:   "InvalidField:eq:value",
			expectedFilter: Filter{},
		},
		{
			name:           "Empty filter string",
			filterString:   "",
			expectedFilter: Filter{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := qp.ParseFilter(tt.filterString)
			assert.Equal(t, tt.expectedFilter, result)
		})
	}
}

func TestQueryPayloadBuilder_ParseFilters(t *testing.T) {
	qp := &QueryPayloadBuilder[TestModel]{}

	filters := []string{
		"Name:eq:John",
		"Age:gt:25",
		"Email:like:@example.com",
	}

	qp.ParseFilters(filters)

	if len(qp.Filters) != 3 {
		t.Errorf("Expected 3 filters, got %d", len(qp.Filters))
	}

	expectedFilters := []Filter{
		{Field: "Name", Operator: OperatorEqual, Value: stringPtr("John")},
		{Field: "Age", Operator: OperatorGreaterThan, Value: stringPtr("25")},
		{Field: "Email", Operator: OperatorLike, Value: stringPtr("@example.com")},
	}

	for i, expected := range expectedFilters {
		assert.Equal(t, expected, qp.Filters[i])
	}
}

func TestQueryPayloadBuilder_ParseSort(t *testing.T) {
	qp := &QueryPayloadBuilder[TestModel]{}

	tests := []struct {
		name         string
		sortString   string
		expectedSort Sort
	}{
		{
			name:       "Valid sort ascending",
			sortString: "Name:asc",
			expectedSort: Sort{
				Field: "Name",
				Order: SortAsc,
			},
		},
		{
			name:       "Valid sort descending",
			sortString: "Age:desc",
			expectedSort: Sort{
				Field: "Age",
				Order: SortDesc,
			},
		},
		{
			name:         "Invalid sort - not enough parts",
			sortString:   "Name",
			expectedSort: Sort{},
		},
		{
			name:         "Invalid sort - field not found",
			sortString:   "InvalidField:asc",
			expectedSort: Sort{},
		},
		{
			name:         "Empty sort string",
			sortString:   "",
			expectedSort: Sort{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := qp.ParseSort(tt.sortString)
			assert.Equal(t, tt.expectedSort, result)
		})
	}
}

func TestQueryPayloadBuilder_ParseSorts(t *testing.T) {
	qp := &QueryPayloadBuilder[TestModel]{}

	sorts := []string{
		"Name:asc",
		"Age:desc",
		"Email:asc",
	}

	qp.ParseSorts(sorts)

	if len(qp.Sorts) != 3 {
		t.Errorf("Expected 3 sorts, got %d", len(qp.Sorts))
	}

	expectedSorts := []Sort{
		{Field: "Name", Order: SortAsc},
		{Field: "Age", Order: SortDesc},
		{Field: "Email", Order: SortAsc},
	}

	for i, expected := range expectedSorts {
		assert.Equal(t, expected, qp.Sorts[i])
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
