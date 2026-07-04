package pagination

import "github.com/raymondsugiarto/reputation-be/pkg/shared/database/view/dto"

type ResultPagination[T any] struct {
	Data        []T   `json:"contents"`
	Count       int64 `json:"totalElements"`
	Page        int   `json:"pageNumber"`
	RowsPerPage int   `json:"size"`
	TotalPages  int   `json:"totalPages"`
}

type Pagination struct {
	Page        int   `json:"page"`
	Count       int64 `json:"totalElements"`
	RowsPerPage int   `json:"size"`
}

type PaginationRequestDto interface {
	GetView() dto.View
	GetPage() int
	GetSize() int
	GetSortBy() string
	GetSortDir() string
	GetQuery() string
	GetFilter() []FilterItem
	AddFilter(FilterItem)
	GenerateFilter()
}

type FilterItem struct {
	Field string
	Op    string
	Val   interface{}
}
