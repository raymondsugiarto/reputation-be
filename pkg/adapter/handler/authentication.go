package handler

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/module/authentication"
)

func SignIn(service authentication.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("SignIn handler called")
		request := new(entity.LoginRequestDto)
		if err := c.Bind().Body(request); err != nil {
			log.WithContext(c).Errorf("error body parser")
			return fiber.NewError(fiber.StatusBadRequest, "errorSignIn")
		}

		response, err := service.SignIn(c, request)
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}
