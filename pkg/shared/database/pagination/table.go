package pagination

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

const DEFAULT_PAGINATION_SIZE = 10

type TableRequest[T PaginationRequestDto] struct {
	QueryField    []string
	Request       T
	AllowedFields []string
	MapFields     map[string]string
}

type Table[T PaginationRequestDto, U, V any] struct{}

func NewTable[T PaginationRequestDto, U, V any]() *Table[T, U, V] {
	return &Table[T, U, V]{}
}

func (t *Table[T, U, V]) Paginate(ctx context.Context, query func(T) *gorm.DB, req *TableRequest[T]) (*ResultPagination[V], []U, error) {
	var data []U = make([]U, 0)
	req.Request.GenerateFilter()
	req.AllowedFields = append(req.AllowedFields, "id", "organization_id", "created_at")

	reqSize := req.Request.GetSize()
	if reqSize == 0 {
		reqSize = DEFAULT_PAGINATION_SIZE
	}

	var (
		wg       sync.WaitGroup
		count    = make(chan int64, 1)
		results  = make(chan []U, 1)
		errQuery = make(chan error, 2)
		offset   = req.Request.GetPage() * reqSize
	)
	// err := sortValidation(req)
	// if err != nil {
	// 	return nil, err
	// }
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		var cnt int64
		q := query(req.Request)
		q = whereConditions(q, req)
		q, err = whereFilterConditions(q, req)
		if err != nil {
			errQuery <- err
			return
		}

		err := q.Count(&cnt).Error
		if err != nil {
			errQuery <- err
			return
		}
		count <- cnt

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		q := query(req.Request)
		q = whereConditions(q, req)
		q, err = whereFilterConditions(q, req)
		if err != nil {
			errQuery <- err
			return
		}
		size := req.Request.GetSize()
		if size == 0 {
			size = -1
		}
		if size > 0 {
			q = q.Offset(int(offset)).Limit(int(size))
		}
		if req.Request.GetSortBy() != "" && req.Request.GetSortDir() != "" {
			q = q.Order(getMappingField(req, req.Request.GetSortBy()) + " " + req.Request.GetSortDir())
		}
		err := q.Find(&data).Error
		if err != nil {
			errQuery <- err
			return
		}

		results <- data
	}()

	go func() {
		wg.Wait()
		close(count)
		close(results)
		close(errQuery)
	}()

	for err := range errQuery {
		log.Errorf("Pagination err %v", err)
		return nil, nil, errors.New("dbErrorPagination")
	}

	totalData := <-count
	return &ResultPagination[V]{
		Page:        req.Request.GetPage(),
		RowsPerPage: req.Request.GetSize(),
		Count:       totalData,
		TotalPages:  calculateTotalPages(int(totalData), reqSize),
	}, <-results, nil
}

func getMappingField[T PaginationRequestDto](req *TableRequest[T], field string) string {
	value, exists := req.MapFields[field]

	if !exists {
		return field
	}
	return value
}

// Function untuk menghitung total halaman
func calculateTotalPages(count, rowsPerPage int) int {
	if count == 0 || rowsPerPage == 0 {
		return 0
	}
	totalPages := count / rowsPerPage
	if count%rowsPerPage != 0 {
		totalPages++
	}
	return totalPages
}

func sortValidation[T PaginationRequestDto](req *TableRequest[T]) error {
	if req.Request.GetSortBy() == "" {
		return errors.New("dbErrorSortBy")
	}
	if req.Request.GetSortDir() == "" {
		return errors.New("dbErrorSortDir")
	}
	return nil
}

func whereConditions[T PaginationRequestDto](db *gorm.DB, req *TableRequest[T]) *gorm.DB {
	if len(req.Request.GetQuery()) == 0 {
		return db
	}
	condStr := []string{}
	values := make([]interface{}, 0)
	for _, v := range req.QueryField {
		condStr = append(condStr, getMappingField(req, v)+" iLIKE ?")
		values = append(values, "%"+req.Request.GetQuery()+"%")
	}
	return db.Where(strings.Join(condStr, " OR "), values...)
}

func whereFilterConditions[T PaginationRequestDto](db *gorm.DB, req *TableRequest[T]) (*gorm.DB, error) {
	filters := req.Request.GetFilter()
	if filters == nil {
		return db, nil
	}
	for _, v := range filters {
		operator, err := getOperator(v.Op)
		if err != nil {
			return nil, err
		}
		err = validateAllowedFields(v.Field, req)
		if err != nil {
			return nil, err
		}
		db = db.Where(getMappingField(req, v.Field)+" "+operator+" ?", v.Val)
	}
	return db, nil
}

func validateAllowedFields[T PaginationRequestDto](field string, req *TableRequest[T]) error {
	isExists := false
	for _, v := range req.AllowedFields {
		if v == field {
			isExists = true
			break
		}
	}
	if !isExists {
		log.Errorf("Field %s not allowed", field)
		return errors.New("dbErrorFieldNotAllowed")
	}
	return nil
}

func getOperator(op string) (string, error) {
	switch op {
	case "eq":
		return "=", nil
	case "lte":
		return "<=", nil
	case "lt":
		return "<", nil
	case "gte":
		return ">=", nil
	case "gt":
		return ">", nil
	}
	return "", errors.New("dbErrorOp")
}
