package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/adapter/handler"
	"github.com/raymondsugiarto/reputation-be/pkg/module/authentication"
	"github.com/raymondsugiarto/reputation-be/pkg/module/organization"
)

func AuthRouter(app *fiber.App, router fiber.Router) {
	authSvc := fiber.MustGetState[authentication.Service](app.State(), authentication.ServiceName)
	router.Post("/sign-in", handler.SignIn(authSvc))

	organizationSvc := fiber.MustGetState[organization.Service](app.State(), organization.ServiceName)
	router.Post("/sign-up", handler.SignUp(organizationSvc))
}
