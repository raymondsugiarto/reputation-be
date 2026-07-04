package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/adapter/handler"
	"github.com/raymondsugiarto/reputation-be/pkg/module/customer"
)

// CustomerRouter mounts the public customer self-service endpoints. Only
// sign-up is public today; the rest of the lifecycle stays admin-gated.
func CustomerRouter(app *fiber.App, router fiber.Router) {
	customerSvc := fiber.MustGetState[customer.Service](app.State(), customer.ServiceName)

	// Public: customer self-service sign-up. OrganizationID is taken from
	// the request body, not the JWT session.
	router.Post("/sign-up", handler.CustomerSignUp(customerSvc))
}

// CustomerAdminRouter mounts the admin-facing customer CRUD endpoints.
// All routes require a valid JWT.
func CustomerAdminRouter(app *fiber.App, router fiber.Router) {
	customerSvc := fiber.MustGetState[customer.Service](app.State(), customer.ServiceName)

	router.Post("/", handler.CreateCustomer(customerSvc))
	router.Get("/", handler.FindAllCustomer(customerSvc))
	router.Get("/:id", handler.FindCustomerByID(customerSvc))
	router.Put("/:id", handler.UpdateCustomerByID(customerSvc))
	router.Delete("/:id", handler.DeleteCustomerByID(customerSvc))
}
