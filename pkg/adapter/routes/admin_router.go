package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/adapter/handler"
	"github.com/raymondsugiarto/reputation-be/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/reputation-be/pkg/module/admin"
	"github.com/raymondsugiarto/reputation-be/pkg/module/customer"
)

// AdminRouter mounts internal-admin-only endpoints.
//
// IMPORTANT: this router must be attached BEFORE the organization
// middleware is added to the parent group. Internal admins are not
// bound to any tenant — the org-context lookup would otherwise fail
// for requests whose X-Origin header does not match a seeded org.
//
// All routes in this group:
//   - require a valid JWT (handled by middleware.Protected on the parent)
//   - require user_type = INTERNAL_ADMIN (handled by AdminOnly)
func AdminRouter(app *fiber.App, router fiber.Router) {
	adminSvc := fiber.MustGetState[admin.Service](app.State(), admin.ServiceName)
	customerSvc := fiber.MustGetState[customer.Service](app.State(), customer.ServiceName)

	// AdminOnly must run AFTER middleware.Protected (parent) so the
	// session is populated, and BEFORE any handler.
	router.Use(middleware.AdminOnly())

	router.Get("/me", handler.GetAdminProfile(adminSvc))
	router.Get("/stats", handler.GetAdminStats(adminSvc))
	router.Get("/customers", handler.FindAllCustomersForAdmin(customerSvc))

	// Customer approval workflow. All routes are platform-wide — they
	// intentionally bypass organization scoping because internal admins
	// triage customers across tenants.
	approvals := router.Group("/approvals")
	approvals.Get("/", handler.FindPendingApprovals(customerSvc))
	approvals.Get("/history", handler.FindApprovalHistory(customerSvc))
	approvals.Get("/stats", handler.GetApprovalStats(customerSvc))
	approvals.Post("/:customerId/approve", handler.ApproveCustomer(customerSvc))
	approvals.Post("/:customerId/reject", handler.RejectCustomer(customerSvc))
}
