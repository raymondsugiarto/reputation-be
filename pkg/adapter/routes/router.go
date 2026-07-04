package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/reputation-be/pkg/infrastructure/middleware/organization"
	// "github.com/raymondsugiarto/reputation-be/pkg/module/dashboard"
)

func InitRouter(app *fiber.App) {
	auth := app.Group("/api/auth")
	AuthRouter(app, auth)

	// Public customer self-service endpoints (no JWT required). Mounted on
	// the same /api prefix as the protected group so the path surface stays
	// stable for the FE.
	customer := app.Group("/api/customer")
	CustomerRouter(app, customer)

	// Protected routes - requires JWT auth
	api := app.Group("/api", middleware.Protected())

	// Internal-admin endpoints. Mounted BEFORE organization.New() so
	// internal admins (who are not bound to any tenant) don't fail the
	// x-origin lookup. The AdminOnly middleware inside AdminRouter still
	// enforces that the caller has user_type = INTERNAL_ADMIN.
	adminGroup := api.Group("/admin")
	AdminRouter(app, adminGroup)

	// // Organization middleware - requires x-origin header for protected routes
	api.Use(organization.New())

	// userCredentialSvc := fiber.MustGetState[usercredential.Service](app.State(), usercredential.ServiceName)
	// api.Put("user-credential/password", handler.ChangePassword(userCredentialSvc))

	// FunderRouter(app, api)
	// ContractRouter(app, api)
	// ContractPaymentRouter(app, api)

	// Customer CRUD (find / update / delete) — sign-up is exposed publicly
	// above; everything else still needs JWT.
	CustomerAdminRouter(app, api)
	// MerchantRouter(app, api)
	// CategoryRouter(app, api)
	// AccountRouter(app, api)
	// TransactionRouter(app, api)

	// Dashboard
	// dashboardSvc := fiber.MustGetState[dashboard.Service](app.State(), dashboard.ServiceName)
	// api.Get("/dashboard/summary", handler.GetDashboardSummary(dashboardSvc))

}
