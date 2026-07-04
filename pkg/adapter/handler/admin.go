package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/module/admin"
	"github.com/raymondsugiarto/reputation-be/pkg/module/customer"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/response/status"
)

// GetAdminProfile handles GET /api/admin/me. The session is already
// populated by middleware.Protected() → SuccessHandler; we only need to
// shape it for the FE.
func GetAdminProfile(service admin.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionRaw := c.Locals(entity.UserSessionKey)
		session, ok := sessionRaw.(*entity.UserSessionDto)
		if !ok || session == nil {
			return status.New(status.InvalidSession)
		}

		profile := service.GetProfile(c, session)
		if profile == nil {
			return status.New(status.InvalidSession)
		}
		return c.JSON(profile)
	}
}

// GetAdminStats handles GET /api/admin/stats. Returns the dashboard KPI
// payload assembled by admin.Service.GetStats.
func GetAdminStats(service admin.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("GetAdminStats handler called")
		stats, err := service.GetStats(c)
		if err != nil {
			log.WithContext(c).Errorf("error fetching admin stats: %v", err)
			return status.New(status.InternalServerError, err)
		}
		return c.JSON(stats)
	}
}

// FindAllCustomersForAdmin handles GET /api/admin/customers. Internally it
// reuses customer.Service.FindAll so that all filtering / pagination rules
// stay in one place. The admin route simply opts out of the
// organization-scoped middleware above it — no business logic duplication.
func FindAllCustomersForAdmin(customerService customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("FindAllCustomersForAdmin handler called")
		query := new(entity.CustomerFilterDto)
		if err := c.Bind().Query(query); err != nil {
			log.WithContext(c).Errorf("error query parser: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "errorInvalidQuery")
		}
		query.GenerateFilter()
		result, err := customerService.FindAll(c, query)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// FindPendingApprovals handles GET /api/admin/approvals. Returns the
// queue of PENDING_APPROVAL customers, paginated. Optional filters:
//
//	?query=...               — iLIKE on nama_lengkap / nomor_ktp
//	?customerType=INDIVIDUAL|COMPANY
func FindPendingApprovals(customerService customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("FindPendingApprovals handler called")
		query := new(entity.CustomerFilterDto)
		if err := c.Bind().Query(query); err != nil {
			log.WithContext(c).Errorf("error query parser: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "errorInvalidQuery")
		}
		query.GenerateFilter()
		result, err := customerService.FindPendingApprovals(c, query)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// FindApprovalHistory handles GET /api/admin/approvals/history.
// Returns APPROVED and REJECTED customers, paginated. Optional filters:
//
//	?action=APPROVED|REJECTED
//	?customerType=INDIVIDUAL|COMPANY
func FindApprovalHistory(customerService customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("FindApprovalHistory handler called")
		query := new(entity.ApprovalHistoryFilterDto)
		if err := c.Bind().Query(query); err != nil {
			log.WithContext(c).Errorf("error query parser: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "errorInvalidQuery")
		}
		query.GenerateFilter()
		result, err := customerService.FindApprovalHistory(c, query)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}
}

// GetApprovalStats handles GET /api/admin/approvals/stats.
func GetApprovalStats(customerService customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("GetApprovalStats handler called")
		stats, err := customerService.GetApprovalStats(c)
		if err != nil {
			log.WithContext(c).Errorf("error fetching approval stats: %v", err)
			return status.New(status.InternalServerError, err)
		}
		return c.JSON(stats)
	}
}

// ApproveCustomer handles POST /api/admin/approvals/:customerId/approve.
// Body: { "remark": "optional remark" }. The authenticated admin's
// user id (from the session) is recorded as the approver.
func ApproveCustomer(customerService customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("ApproveCustomer handler called")
		customerID := c.Params("customerId")
		if customerID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "errorMissingCustomerId")
		}

		// Body is optional — approve accepts an empty remark. We
		// attempt the parse for forward-compat (so FE can send a
		// remark when desired) but treat parse errors as "no body".
		req := new(entity.ApprovalActionRequestDto)
		if err := c.Bind().Body(req); err != nil {
			log.WithContext(c).Debugf("approve: body parse skipped (%v)", err)
			req = &entity.ApprovalActionRequestDto{}
		}

		adminID, _ := adminUserIDFromSession(c)
		result, err := customerService.Approve(c, customerID, adminID, req.Remark)
		if err != nil {
			return mapApprovalError(err)
		}
		return c.JSON(result)
	}
}

// RejectCustomer handles POST /api/admin/approvals/:customerId/reject.
// Body: { "remark": "required remark" }.
func RejectCustomer(customerService customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("RejectCustomer handler called")
		customerID := c.Params("customerId")
		if customerID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "errorMissingCustomerId")
		}

		req := new(entity.ApprovalActionRequestDto)
		if err := c.Bind().Body(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "errorInvalidRequest")
		}

		adminID, _ := adminUserIDFromSession(c)
		result, err := customerService.Reject(c, customerID, adminID, req.Remark)
		if err != nil {
			return mapApprovalError(err)
		}
		return c.JSON(result)
	}
}

// adminUserIDFromSession extracts the internal-admin's user id from the
// session populated by middleware.Protected. Returns "" when no
// session is present (which AdminOnly would have rejected anyway).
func adminUserIDFromSession(c fiber.Ctx) (string, error) {
	sessionRaw := c.Locals(entity.UserSessionKey)
	session, ok := sessionRaw.(*entity.UserSessionDto)
	if !ok || session == nil || session.UserCredential == nil {
		return "", errors.New("no session")
	}
	return session.UserCredential.UserID, nil
}

// mapApprovalError translates customer.ErrCustomerNotPending /
// ErrApprovalRemarkRequired into typed AppStatus responses. Any other
// error is treated as an internal server error.
func mapApprovalError(err error) error {
	switch {
	case errors.Is(err, customer.ErrCustomerNotPending):
		return status.New(status.EntityConflict, err)
	case errors.Is(err, customer.ErrApprovalRemarkRequired):
		return status.New(status.MandatoryFieldMissing, err)
	default:
		return status.New(status.InternalServerError, err)
	}
}
