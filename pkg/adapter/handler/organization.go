package handler

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/module/organization"
)

func SignUp(service organization.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("SignUp handler called")
		request := new(entity.SignUpRequestDto)
		if err := c.Bind().Body(request); err != nil {
			log.WithContext(c).Errorf("error body parser")
			return fiber.NewError(fiber.StatusBadRequest, "errorSignUp")
		}

		response, err := service.SignUp(c, request)
		if err != nil {
			log.WithContext(c).Errorf("error sign up: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(response)
	}
}
