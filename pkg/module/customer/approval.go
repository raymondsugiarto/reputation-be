package customer

import (
	"context"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
)

// Sentinel errors used by the admin approval flow. The handler maps
// these to typed AppStatus responses (404 / 409 / 422) so the FE can
// react appropriately without inspecting raw error strings.
var (
	// ErrCustomerNotPending — admin tried to approve/reject a customer
	// whose status is not PENDING_APPROVAL. Surfaced as 409 to signal
	// "state conflict" rather than "not found".
	ErrCustomerNotPending = errors.New("customerNotPending")

	// ErrApprovalRemarkRequired — admin tried to reject a customer
	// without providing a remark. Surfaced as 422.
	ErrApprovalRemarkRequired = errors.New("approvalRemarkRequired")
)

// ApprovalResultDto is the response payload returned by both Approve
// and Reject. It wraps the updated CustomerDto so the FE can render
// the new status without an extra GET.
type ApprovalResultDto struct {
	Customer *entity.CustomerDto `json:"customer"`
}

// Approve transitions a customer from PENDING_APPROVAL to APPROVED.
// `adminUserID` is recorded as `approved_by` so the audit trail is
// preserved. `remark` is optional.
//
// Returns ErrCustomerNotPending if the customer is already APPROVED or
// REJECTED.
func (s *service) Approve(
	ctx context.Context,
	customerID string,
	adminUserID string,
	remark string,
) (*ApprovalResultDto, error) {
	existing, err := s.repo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if existing.Status != model.CustomerStatusPendingApproval {
		return nil, ErrCustomerNotPending
	}

	now := time.Now().UTC()
	dto := &entity.CustomerDto{
		ID:             existing.ID,
		OrganizationID: existing.OrganizationID,
		UserID:         existing.UserID,
		CustomerType:   existing.CustomerType,
		Status:         model.CustomerStatusApproved,
		ApprovedBy:     adminUserID,
		ApprovedAt:     &now,
		RejectedBy:     "",
		RejectedAt:     nil,
		Remark:         remark,
		// Carry through all the demographic fields. ToModel only writes
		// the columns it knows about, so we need to populate them here.
		NamaLengkap:        existing.NamaLengkap,
		NomorKtp:           existing.NomorKtp,
		NomorNpwp:          existing.NomorNpwp,
		TanggalLahir:       existing.TanggalLahir,
		KotaLahir:          existing.KotaLahir,
		StatusPernikahan:   existing.StatusPernikahan,
		PendidikanTerakhir: existing.PendidikanTerakhir,
		LamaTinggal:        existing.LamaTinggal,
		AlamatJalan:        existing.AlamatJalan,
		Kecamatan:          existing.Kecamatan,
		KotaKabupaten:      existing.KotaKabupaten,
		Provinsi:           existing.Provinsi,
		KodePos:            existing.KodePos,
		SamaDenganDomisili: existing.SamaDenganDomisili,
		AlamatKtpJalan:     existing.AlamatKtpJalan,
		KecamatanKtp:       existing.KecamatanKtp,
		KotaKabupatenKtp:   existing.KotaKabupatenKtp,
		NamaPt:             existing.NamaPt,
		TanggalPendirian:   existing.TanggalPendirian,
		SektorUsaha:        existing.SektorUsaha,
		KodeKbli:           existing.KodeKbli,
		NpwpPerusahaan:     existing.NpwpPerusahaan,
		PemegangSaham:      existing.PemegangSaham,
	}

	var updated *entity.CustomerDto
	err = s.txManager.Execute(ctx, func(txCtx context.Context) error {
		res, err := s.repo.Update(txCtx, dto)
		if err != nil {
			log.WithContext(ctx).Errorf("approve customer: %v", err)
			return err
		}
		updated = res
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ApprovalResultDto{Customer: updated}, nil
}

// Reject transitions a customer from PENDING_APPROVAL to REJECTED.
// `remark` is mandatory — admins must explain why the customer was
// rejected. Returns ErrApprovalRemarkRequired when empty, and
// ErrCustomerNotPending when the customer is not in PENDING_APPROVAL.
func (s *service) Reject(
	ctx context.Context,
	customerID string,
	adminUserID string,
	remark string,
) (*ApprovalResultDto, error) {
	if remark == "" {
		return nil, ErrApprovalRemarkRequired
	}

	existing, err := s.repo.FindByID(ctx, customerID)
	if err != nil {
		return nil, err
	}
	if existing.Status != model.CustomerStatusPendingApproval {
		return nil, ErrCustomerNotPending
	}

	now := time.Now().UTC()
	dto := &entity.CustomerDto{
		ID:             existing.ID,
		OrganizationID: existing.OrganizationID,
		UserID:         existing.UserID,
		CustomerType:   existing.CustomerType,
		Status:         model.CustomerStatusRejected,
		ApprovedBy:     "",
		ApprovedAt:     nil,
		RejectedBy:     adminUserID,
		RejectedAt:     &now,
		Remark:         remark,
		// Carry through demographic fields (same reasoning as Approve).
		NamaLengkap:        existing.NamaLengkap,
		NomorKtp:           existing.NomorKtp,
		NomorNpwp:          existing.NomorNpwp,
		TanggalLahir:       existing.TanggalLahir,
		KotaLahir:          existing.KotaLahir,
		StatusPernikahan:   existing.StatusPernikahan,
		PendidikanTerakhir: existing.PendidikanTerakhir,
		LamaTinggal:        existing.LamaTinggal,
		AlamatJalan:        existing.AlamatJalan,
		Kecamatan:          existing.Kecamatan,
		KotaKabupaten:      existing.KotaKabupaten,
		Provinsi:           existing.Provinsi,
		KodePos:            existing.KodePos,
		SamaDenganDomisili: existing.SamaDenganDomisili,
		AlamatKtpJalan:     existing.AlamatKtpJalan,
		KecamatanKtp:       existing.KecamatanKtp,
		KotaKabupatenKtp:   existing.KotaKabupatenKtp,
		NamaPt:             existing.NamaPt,
		TanggalPendirian:   existing.TanggalPendirian,
		SektorUsaha:        existing.SektorUsaha,
		KodeKbli:           existing.KodeKbli,
		NpwpPerusahaan:     existing.NpwpPerusahaan,
		PemegangSaham:      existing.PemegangSaham,
	}

	var updated *entity.CustomerDto
	err = s.txManager.Execute(ctx, func(txCtx context.Context) error {
		res, err := s.repo.Update(txCtx, dto)
		if err != nil {
			log.WithContext(ctx).Errorf("reject customer: %v", err)
			return err
		}
		updated = res
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ApprovalResultDto{Customer: updated}, nil
}

// FindPendingApprovals returns the queue of customers awaiting admin
// review. Backed by customer.Service.FindAll with status filter.
//
// We re-use FindAll instead of writing a custom SQL because the
// pagination library already handles filtering, sorting, and the
// organisation-scoping rule that FindAll enforces.
func (s *service) FindPendingApprovals(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	filter, ok := req.(*entity.CustomerFilterDto)
	if !ok {
		filter = &entity.CustomerFilterDto{GetListRequest: pagination.GetListRequest{}}
	}
	filter.Status = string(model.CustomerStatusPendingApproval)
	return s.repo.FindAll(ctx, filter)
}

// FindApprovalHistory returns paginated customers that have already
// been approved or rejected. Uses the dedicated repo method that
// always excludes PENDING_APPROVAL rows. `req` should be a
// *entity.ApprovalHistoryFilterDto so the action / customerType
// filters are honoured.
func (s *service) FindApprovalHistory(ctx context.Context, req pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	return s.repo.FindApprovalHistory(ctx, req)
}

// GetApprovalStats returns the counters used by the admin approval
// page's top-of-page strip.
func (s *service) GetApprovalStats(ctx context.Context) (*entity.CustomerApprovalStatsDto, error) {
	stats := &entity.CustomerApprovalStatsDto{}

	if err := s.repo.CountByStatus(ctx, model.CustomerStatusPendingApproval, &stats.PendingApprovals); err != nil {
		return nil, err
	}
	if err := s.repo.CountByStatus(ctx, model.CustomerStatusApproved, &stats.TotalApproved); err != nil {
		return nil, err
	}
	if err := s.repo.CountByStatus(ctx, model.CustomerStatusRejected, &stats.TotalRejected); err != nil {
		return nil, err
	}

	// Today counters. We reuse today-start UTC for consistency with the
	// rest of the platform (timestamps are stored in UTC).
	now := time.Now().UTC()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	stats.ApprovedToday = s.repo.CountByStatusSince(ctx, model.CustomerStatusApproved, "approved_at", todayStart)
	stats.RejectedToday = s.repo.CountByStatusSince(ctx, model.CustomerStatusRejected, "rejected_at", todayStart)

	return stats, nil
}
