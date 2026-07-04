package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/module/customer"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/response/status"
)

// CustomerSignUp is the authenticated endpoint for onboarding a new customer
// under the caller's organization. It creates a UserCredential, a User
// (type CUSTOMER), and the Customer profile in one transaction.
func CustomerSignUp(service customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		log.WithContext(c).Infof("CustomerSignUp handler called")
		request := new(entity.CustomerSignUpRequestDto)
		if err := c.Bind().Body(request); err != nil {
			log.WithContext(c).Errorf("error body parser: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "errorInvalidRequest")
		}

		response, err := service.SignUp(c, request)
		if err != nil {
			log.WithContext(c).Errorf("error customer sign up: %v", err)
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func CreateCustomer(service customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		dto := new(entity.CustomerDto)
		if err := c.Bind().Body(dto); err != nil {
			log.WithContext(c).Errorf("error body parser: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "errorInvalidRequest")
		}
		response, err := service.Create(c, dto)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func FindCustomerByID(service customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		response, err := service.FindByID(c, id)
		if err != nil {
			return err
		}
		return c.JSON(response)
	}
}

func FindAllCustomer(service customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		query := new(entity.CustomerFilterDto)
		if err := c.Bind().Query(query); err != nil {
			log.WithContext(c).Errorf("error query parser: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "errorInvalidQuery")
		}
		query.GenerateFilter()
		response, err := service.FindAll(c, query)
		if err != nil {
			// An empty search result is a 404 (not a 500). The FE catches
			// the status code and renders the dedicated "result-not-found"
			// UI from /public/result-not-found.png. We surface this as a
			// typed AppStatus (case code = EntityNotFound) so the wire
			// envelope reads `4041101 entity not found` rather than the
			// generic route-not-found mapping that fiber.NewError produces.
			if errors.Is(err, customer.ErrCustomerNotFound) {
				return status.New(status.EntityNotFound, errors.New("customerNotFound"))
			}
			return err
		}
		return c.JSON(response)
	}
}

func UpdateCustomerByID(service customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		dto := new(entity.CustomerDto)
		if err := c.Bind().Body(dto); err != nil {
			log.WithContext(c).Errorf("error body parser: %v", err)
			return fiber.NewError(fiber.StatusBadRequest, "errorInvalidRequest")
		}
		dto.ID = id
		response, err := service.Update(c, dto)
		if err != nil {
			return err
		}
		return c.JSON(response)
	}
}

func DeleteCustomerByID(service customer.Service) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		if err := service.Delete(c, id); err != nil {
			return err
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
