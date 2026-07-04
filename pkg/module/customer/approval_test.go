package customer_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/module/customer"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
)

// fakeApprovalService implements customer.Service for unit-testing the
// approval workflow without a live DB. Every method that's called from
// Approve / Reject / FindPendingApprovals / FindApprovalHistory /
// GetApprovalStats is stubbed here.
//
// Methods we don't exercise (SignUp, Create, Update, etc.) return
// their zero value so the test compiles cleanly.
type fakeApprovalService struct {
	// State — what the in-memory "DB" thinks each customer looks like.
	byID map[string]*entity.CustomerDto

	// Spies — record the last action taken.
	lastApproveID string
	lastRejectID  string
	lastRemark    string
	lastAdminID   string
}

func newFakeApprovalService() *fakeApprovalService {
	return &fakeApprovalService{
		byID: map[string]*entity.CustomerDto{
			"pending-1": {
				ID:           "pending-1",
				CustomerType: model.CustomerTypeIndividual,
				NamaLengkap:  "Budi Pending",
				Status:       model.CustomerStatusPendingApproval,
			},
			"approved-1": {
				ID:           "approved-1",
				CustomerType: model.CustomerTypeCompany,
				NamaPt:       "PT Sudah Disetujui",
				Status:       model.CustomerStatusApproved,
			},
			"rejected-1": {
				ID:           "rejected-1",
				CustomerType: model.CustomerTypeIndividual,
				NamaLengkap:  "Andi Ditolak",
				Status:       model.CustomerStatusRejected,
				Remark:       "Dokumen tidak lengkap",
			},
		},
	}
}

// Methods exercised by the approval tests.

func (f *fakeApprovalService) FindByID(_ context.Context, id string) (*entity.CustomerDto, error) {
	c, ok := f.byID[id]
	if !ok {
		return nil, errors.New("customerNotFound")
	}
	// Return a copy so callers can't mutate the spy state through the
	// returned pointer.
	cp := *c
	return &cp, nil
}

func (f *fakeApprovalService) Update(_ context.Context, dto *entity.CustomerDto) (*entity.CustomerDto, error) {
	if _, ok := f.byID[dto.ID]; !ok {
		return nil, errors.New("customerNotFound")
	}
	cp := *dto
	f.byID[dto.ID] = &cp
	return dto, nil
}

func (f *fakeApprovalService) Approve(_ context.Context, customerID, adminID, remark string) (*customer.ApprovalResultDto, error) {
	existing, err := f.FindByID(context.Background(), customerID)
	if err != nil {
		return nil, err
	}
	if existing.Status != model.CustomerStatusPendingApproval {
		return nil, customer.ErrCustomerNotPending
	}
	now := time.Now().UTC()
	updated := *existing
	updated.Status = model.CustomerStatusApproved
	updated.ApprovedBy = adminID
	updated.ApprovedAt = &now
	updated.Remark = remark
	f.byID[customerID] = &updated

	f.lastApproveID = customerID
	f.lastAdminID = adminID
	f.lastRemark = remark

	return &customer.ApprovalResultDto{Customer: &updated}, nil
}

func (f *fakeApprovalService) Reject(_ context.Context, customerID, adminID, remark string) (*customer.ApprovalResultDto, error) {
	if remark == "" {
		return nil, customer.ErrApprovalRemarkRequired
	}
	existing, err := f.FindByID(context.Background(), customerID)
	if err != nil {
		return nil, err
	}
	if existing.Status != model.CustomerStatusPendingApproval {
		return nil, customer.ErrCustomerNotPending
	}
	now := time.Now().UTC()
	updated := *existing
	updated.Status = model.CustomerStatusRejected
	updated.RejectedBy = adminID
	updated.RejectedAt = &now
	updated.Remark = remark
	f.byID[customerID] = &updated

	f.lastRejectID = customerID
	f.lastAdminID = adminID
	f.lastRemark = remark

	return &customer.ApprovalResultDto{Customer: &updated}, nil
}

func (f *fakeApprovalService) FindPendingApprovals(_ context.Context, _ pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	rows := make([]entity.CustomerDto, 0)
	for _, c := range f.byID {
		if c.Status == model.CustomerStatusPendingApproval {
			rows = append(rows, *c)
		}
	}
	return &pagination.ResultPagination[entity.CustomerDto]{
		Data:        rows,
		Count:       int64(len(rows)),
		Page:        0,
		RowsPerPage: 10,
		TotalPages:  1,
	}, nil
}

func (f *fakeApprovalService) FindApprovalHistory(_ context.Context, _ pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	rows := make([]entity.CustomerDto, 0)
	for _, c := range f.byID {
		if c.Status == model.CustomerStatusApproved ||
			c.Status == model.CustomerStatusRejected {
			rows = append(rows, *c)
		}
	}
	return &pagination.ResultPagination[entity.CustomerDto]{
		Data:        rows,
		Count:       int64(len(rows)),
		Page:        0,
		RowsPerPage: 10,
		TotalPages:  1,
	}, nil
}

func (f *fakeApprovalService) GetApprovalStats(_ context.Context) (*entity.CustomerApprovalStatsDto, error) {
	stats := &entity.CustomerApprovalStatsDto{}
	for _, c := range f.byID {
		switch c.Status {
		case model.CustomerStatusPendingApproval:
			stats.PendingApprovals++
		case model.CustomerStatusApproved:
			stats.TotalApproved++
		case model.CustomerStatusRejected:
			stats.TotalRejected++
		}
	}
	return stats, nil
}

// Stubs for methods we don't exercise.

func (f *fakeApprovalService) SignUp(_ context.Context, _ *entity.CustomerSignUpRequestDto) (*entity.CustomerSignUpResponseDto, error) {
	return nil, nil
}
func (f *fakeApprovalService) Create(_ context.Context, _ *entity.CustomerDto) (*entity.CustomerDto, error) {
	return nil, nil
}
func (f *fakeApprovalService) FindByUserID(_ context.Context, _ string) (*entity.CustomerDto, error) {
	return nil, nil
}
func (f *fakeApprovalService) FindAll(_ context.Context, _ pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	return nil, nil
}
func (f *fakeApprovalService) Delete(_ context.Context, _ string) error { return nil }

// Compile-time interface guard.
var _ customer.Service = (*fakeApprovalService)(nil)

// ---------------------------------------------------------------------------
// Approve
// ---------------------------------------------------------------------------

func TestApprove_PendingCustomerTransitionsToApproved(t *testing.T) {
	svc := newFakeApprovalService()

	result, err := svc.Approve(context.Background(), "pending-1", "admin-1", "lolos KYC")
	if err != nil {
		t.Fatalf("approve: %v", err)
	}
	if result.Customer.Status != model.CustomerStatusApproved {
		t.Fatalf("expected APPROVED, got %s", result.Customer.Status)
	}
	if result.Customer.ApprovedBy != "admin-1" {
		t.Fatalf("expected approvedBy=admin-1, got %q", result.Customer.ApprovedBy)
	}
	if result.Customer.ApprovedAt == nil {
		t.Fatalf("expected approvedAt to be set")
	}
	if result.Customer.Remark != "lolos KYC" {
		t.Fatalf("expected remark persisted, got %q", result.Customer.Remark)
	}
}

func TestApprove_NonPendingReturnsConflict(t *testing.T) {
	svc := newFakeApprovalService()

	// Already APPROVED.
	if _, err := svc.Approve(context.Background(), "approved-1", "admin-1", ""); !errors.Is(err, customer.ErrCustomerNotPending) {
		t.Fatalf("expected ErrCustomerNotPending for approved customer, got %v", err)
	}
	// Already REJECTED.
	if _, err := svc.Approve(context.Background(), "rejected-1", "admin-1", ""); !errors.Is(err, customer.ErrCustomerNotPending) {
		t.Fatalf("expected ErrCustomerNotPending for rejected customer, got %v", err)
	}
}

func TestApprove_NotFoundReturnsError(t *testing.T) {
	svc := newFakeApprovalService()
	if _, err := svc.Approve(context.Background(), "missing", "admin-1", ""); err == nil {
		t.Fatalf("expected error for missing customer")
	}
}

// ---------------------------------------------------------------------------
// Reject
// ---------------------------------------------------------------------------

func TestReject_PendingCustomerTransitionsToRejected(t *testing.T) {
	svc := newFakeApprovalService()

	result, err := svc.Reject(context.Background(), "pending-1", "admin-1", "KTP tidak terbaca")
	if err != nil {
		t.Fatalf("reject: %v", err)
	}
	if result.Customer.Status != model.CustomerStatusRejected {
		t.Fatalf("expected REJECTED, got %s", result.Customer.Status)
	}
	if result.Customer.RejectedBy != "admin-1" {
		t.Fatalf("expected rejectedBy=admin-1, got %q", result.Customer.RejectedBy)
	}
	if result.Customer.RejectedAt == nil {
		t.Fatalf("expected rejectedAt to be set")
	}
	if result.Customer.Remark != "KTP tidak terbaca" {
		t.Fatalf("expected remark persisted, got %q", result.Customer.Remark)
	}
}

func TestReject_WithoutRemarkReturnsError(t *testing.T) {
	svc := newFakeApprovalService()
	_, err := svc.Reject(context.Background(), "pending-1", "admin-1", "")
	if !errors.Is(err, customer.ErrApprovalRemarkRequired) {
		t.Fatalf("expected ErrApprovalRemarkRequired, got %v", err)
	}
}

func TestReject_NonPendingReturnsConflict(t *testing.T) {
	svc := newFakeApprovalService()
	_, err := svc.Reject(context.Background(), "approved-1", "admin-1", "late remark")
	if !errors.Is(err, customer.ErrCustomerNotPending) {
		t.Fatalf("expected ErrCustomerNotPending, got %v", err)
	}
}

// ---------------------------------------------------------------------------
// FindPendingApprovals / FindApprovalHistory
// ---------------------------------------------------------------------------

func TestFindPendingApprovals_ExcludesNonPending(t *testing.T) {
	svc := newFakeApprovalService()
	result, err := svc.FindPendingApprovals(context.Background(), nil)
	if err != nil {
		t.Fatalf("find pending: %v", err)
	}
	if result.Count != 1 {
		t.Fatalf("expected 1 pending, got %d", result.Count)
	}
	for _, r := range result.Data {
		if r.Status != model.CustomerStatusPendingApproval {
			t.Fatalf("non-pending row leaked: %s", r.Status)
		}
	}
}

func TestFindApprovalHistory_ExcludesPending(t *testing.T) {
	svc := newFakeApprovalService()
	result, err := svc.FindApprovalHistory(context.Background(), nil)
	if err != nil {
		t.Fatalf("find history: %v", err)
	}
	if result.Count != 2 {
		t.Fatalf("expected 2 historical rows, got %d", result.Count)
	}
	for _, r := range result.Data {
		if r.Status != model.CustomerStatusApproved && r.Status != model.CustomerStatusRejected {
			t.Fatalf("pending row leaked into history: %s", r.Status)
		}
	}
}

// ---------------------------------------------------------------------------
// GetApprovalStats
// ---------------------------------------------------------------------------

func TestGetApprovalStats_CountsAllBuckets(t *testing.T) {
	svc := newFakeApprovalService()
	stats, err := svc.GetApprovalStats(context.Background())
	if err != nil {
		t.Fatalf("stats: %v", err)
	}
	if stats.PendingApprovals != 1 {
		t.Errorf("pending: got %d want 1", stats.PendingApprovals)
	}
	if stats.TotalApproved != 1 {
		t.Errorf("approved: got %d want 1", stats.TotalApproved)
	}
	if stats.TotalRejected != 1 {
		t.Errorf("rejected: got %d want 1", stats.TotalRejected)
	}
}
