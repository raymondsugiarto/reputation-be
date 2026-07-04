package middleware

import (
	config "github.com/raymondsugiarto/reputation-be/config"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/module/authentication"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

// Protected protect routes
func Protected() fiber.Handler {
	cfg := config.GetConfig()
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(cfg.Server.Rest.SecretKey)},
		ErrorHandler:   jwtError,
		SuccessHandler: SuccessHandler,
		Extractor:      extractors.FromAuthHeader("Bearer"),
	})
}

func jwtError(c fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func SuccessHandler(c fiber.Ctx) error {
	authenticationSvc := fiber.MustGetState[authentication.Service](c.App().State(), authentication.ServiceName)
	userSessionDto, err := authenticationSvc.GetSession(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "error", "message": "Failed to get user session", "data": nil})
	}
	c.Locals(entity.UserSessionKey, userSessionDto)

	// funderService := fiber.MustGetState[funder.Service](c.App().State(), funder.ServiceName)
	// funderDto, err := funderService.IdentifySessionFunder(c, userSessionDto)
	// if err != nil {
	// 	return c.Status(fiber.StatusUnauthorized).
	// 		JSON(fiber.Map{"status": "error", "message": "Failed to identify funder from session", "data": nil})
	// }
	// if funderDto != nil {
	// 	c.Locals(entity.FunderSessionKey, funderDto)
	// }

	return c.Next()
}
