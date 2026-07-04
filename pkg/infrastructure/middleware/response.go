package middleware

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/response"
)

func DefaultErrorHandler() func(fiber.Ctx, error) error {
	return func(c fiber.Ctx, err error) error {
		log.WithContext(c).Errorf("Error: %v", err)
		resp := response.NewError(err)
		return c.Status(resp.HTTPCode).JSON(resp)
	}

}

func DefaultResponseHandler() func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		err := c.Next()

		if err != nil {
			return err
		}

		respBody := c.Response().Body()

		var resp = response.NewSuccess(c.Response().StatusCode(), respBody)

		var data any
		if len(respBody) > 0 {
			if err := json.Unmarshal(respBody, &data); err == nil {
				resp = response.NewSuccess(c.Response().StatusCode(), data)
			}
		}

		return c.JSON(resp)
	}
}
