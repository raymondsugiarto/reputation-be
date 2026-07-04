package handler

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/module/account"
)

func CreateAccount(service account.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(entity.AccountRequest)

		if err := c.Bind().Body(request); err != nil {
			log.WithContext(c).Errorf("error body parser")
			return fiber.NewError(fiber.StatusBadRequest, "error creating account")
		}

		response, err := service.Create(c, request.ToDto())
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}

func UpdateAccountByID(service account.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		request := new(entity.AccountRequest)
		id := c.Params("id")

		if err := c.Bind().Body(request); err != nil {
			log.WithContext(c).Errorf("error body parser")
			return fiber.NewError(fiber.StatusBadRequest, "error updating account")
		}

		dto := request.ToDto()
		dto.ID = id

		response, err := service.Update(c, dto)
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}

func FindAccountByID(service account.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")

		response, err := service.FindByID(c, id)
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}

func DeleteAccountByID(service account.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")

		err := service.Delete(c, id)
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func FindAllAccount(service account.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		query := new(entity.AccountFilterDto)
		if err := c.Bind().Query(query); err != nil {
			log.WithContext(c).Errorf("error query parser", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		query.GenerateFilter()

		response, err := service.FindAll(c, query)
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}

func FindAccountByCustomerID(service account.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		customerID := c.Params("customerId")

		response, err := service.FindByCustomerID(c, customerID)
		if err != nil {
			return err
		}

		return c.JSON(response)
	}
}
