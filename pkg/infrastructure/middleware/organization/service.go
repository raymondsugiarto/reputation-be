package organization

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/infrastructure/database"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
)

func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	return func(c fiber.Ctx) error {

		// Get id from request, else we generate one
		origin := c.Get(cfg.HeaderOriginKey)
		if origin == "" {
			log.WithContext(c).Errorf("Missing origin in request")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing or malformed origin", "data": nil})
		}

		originType := c.Get(cfg.HeaderOriginTypeKey)
		if originType == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Missing or malformed origin type", "data": nil})
		}

		log.WithContext(c).Infof("Origin: %v", origin)
		org, err := getOrganizationByOrigin(c, origin)
		if err != nil {
			log.WithContext(c).Errorf("error org not found")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Organization not found", "data": nil})

		}
		c.Locals(entity.OriginKey, origin)
		c.Locals(entity.OriginTypeKey, originType)
		c.Locals(entity.OrganizationKey, org)

		return c.Next()
	}
}

func getOrganizationByOrigin(ctx context.Context, origin string) (*entity.OrganizationDto, error) {
	db := database.DBConn
	log.WithContext(ctx).Infof("find org with origin %s", origin)

	var organization *model.Organization
	err := db.WithContext(ctx).Where("origin = ?", origin).Find(&organization).Error
	if err != nil {
		log.WithContext(ctx).Errorf("Error: %v", err)
		return nil, errors.New("organization not found")
	}
	if organization.ID == "" {
		return nil, errors.New("organization not found")
	}
	log.WithContext(ctx).Infof("Organization found: %v", organization.ID)

	return entity.NewOrganizationDtoFromModel(organization), nil
}
