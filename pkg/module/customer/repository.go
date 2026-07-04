package customer

import (
	"context"
	"time"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error)
	FindByID(ctx context.Context, id string) (*entity.CustomerDto, error)
	FindByUserID(ctx context.Context, userID string) (*entity.CustomerDto, error)
	FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error)
	// FindApprovalHistory returns paginated rows whose status is
	// APPROVED or REJECTED. `action` further filters by exact status;
	// when empty, both approved and rejected rows are returned.
	FindApprovalHistory(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error)
	Update(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error)
	Delete(ctx context.Context, id string) error
	// CountByStatus fills `*out` with the number of customers whose
	// status column equals the given value. Used for the admin
	// approval stats endpoint.
	CountByStatus(ctx context.Context, status model.CustomerStatus, out *int64) error
	// CountByStatusSince fills `*out` with the number of customers
	// whose `column` (e.g. "approved_at" or "rejected_at") is >= since.
	CountByStatusSince(ctx context.Context, status model.CustomerStatus, column string, since time.Time) int64
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error) {
	m := dto.ToModel()
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	dto.ID = m.ID
	return dto, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (*entity.CustomerDto, error) {
	var m model.Customer
	if err := r.db.WithContext(ctx).Preload("User").First(&m, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return entity.NewCustomerDtoFromModel(&m), nil
}

func (r *repository) FindByUserID(ctx context.Context, userID string) (*entity.CustomerDto, error) {
	var m model.Customer
	if err := r.db.WithContext(ctx).Preload("User").First(&m, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return entity.NewCustomerDtoFromModel(&m), nil
}

func (r *repository) FindAll(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	info, paginationResult, err := pagination.NewTable[*entity.CustomerFilterDto, *model.Customer, entity.CustomerDto]().
		Paginate(ctx, func(req *entity.CustomerFilterDto) *gorm.DB {
			query := r.db.WithContext(ctx).Model(&model.Customer{}).Preload("User")
			if req.OrganizationID != "" {
				query = query.Where("organization_id = ?", req.OrganizationID)
			}
			// Approval status filter. When empty, no WHERE clause is
			// applied — preserves the original admin "all customers"
			// behaviour.
			if req.Status != "" {
				query = query.Where("status = ?", req.Status)
			}
			return query
		}, &pagination.TableRequest[*entity.CustomerFilterDto]{
			Request:       req.(*entity.CustomerFilterDto),
			QueryField:    []string{"nama_lengkap", "nomor_ktp"},
			AllowedFields: []string{"organization_id", "status"},
		})
	if err != nil {
		return nil, err
	}
	result := make([]entity.CustomerDto, len(paginationResult))
	for i, m := range paginationResult {
		result[i] = *entity.NewCustomerDtoFromModel(m)
	}
	info.Data = result
	return info, nil
}

func (r *repository) Update(ctx context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error) {
	if err := r.db.WithContext(ctx).Save(dto.ToModel()).Error; err != nil {
		return nil, err
	}
	return dto, nil
}

func (r *repository) FindApprovalHistory(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	filter, ok := req.(*entity.ApprovalHistoryFilterDto)
	if !ok {
		filter = &entity.ApprovalHistoryFilterDto{GetListRequest: pagination.GetListRequest{}}
	}
	info, paginationResult, err := pagination.NewTable[*entity.ApprovalHistoryFilterDto, *model.Customer, entity.CustomerDto]().
		Paginate(ctx, func(req *entity.ApprovalHistoryFilterDto) *gorm.DB {
			query := r.db.WithContext(ctx).Model(&model.Customer{}).Preload("User")
			// Default: exclude the PENDING_APPROVAL backlog — this
			// endpoint is for the audit/history view, not the queue.
			if req.Action != "" {
				query = query.Where("status = ?", req.Action)
			} else {
				query = query.Where("status IN (?, ?)",
					model.CustomerStatusApproved, model.CustomerStatusRejected)
			}
			if req.CustomerType != "" {
				query = query.Where("customer_type = ?", req.CustomerType)
			}
			return query
		}, &pagination.TableRequest[*entity.ApprovalHistoryFilterDto]{
			Request:       filter,
			QueryField:    []string{"nama_lengkap", "nomor_ktp"},
			AllowedFields: []string{"status", "customer_type"},
		})
	if err != nil {
		return nil, err
	}
	result := make([]entity.CustomerDto, len(paginationResult))
	for i, m := range paginationResult {
		result[i] = *entity.NewCustomerDtoFromModel(m)
	}
	info.Data = result
	return info, nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Customer{}).Error
}

// CountByStatus fills `*out` with the number of customers whose status
// matches. Errors leave *out untouched so callers can fail fast.
func (r *repository) CountByStatus(ctx context.Context, status model.CustomerStatus, out *int64) error {
	return r.db.WithContext(ctx).
		Model(&model.Customer{}).
		Where("status = ?", status).
		Count(out).Error
}

// CountByStatusSince returns the number of customers in the given
// status whose `column` timestamp is >= `since`. Used for the
// "approved today" / "rejected today" counters.
//
// NOTE: `column` is interpolated directly into the SQL — it MUST be a
// hardcoded constant from the caller (we never accept user input).
func (r *repository) CountByStatusSince(ctx context.Context, status model.CustomerStatus, column string, since time.Time) int64 {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&model.Customer{}).
		Where("status = ? AND "+column+" >= ?", status, since).
		Count(&count).Error; err != nil {
		return 0
	}
	return count
}
