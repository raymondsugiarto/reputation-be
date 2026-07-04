package organization

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
)

// Config defines the config for middleware.
type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// Header is the header key where to get/set the unique request ID
	//
	// Optional. Default: "X-Request-ID"
	HeaderOriginKey     string
	HeaderOriginTypeKey string
}

// ConfigDefault is the default config
// It uses a fast UUID generator which will expose the number of
// requests made to the server. To conceal this value for better
// privacy, use the "utils.UUIDv4" generator.
var ConfigDefault = Config{
	Next:                nil,
	HeaderOriginKey:     entity.OriginKey,
	HeaderOriginTypeKey: entity.OriginTypeKey,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	fmt.Printf("cfg: %v\n", config)

	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.HeaderOriginTypeKey == "" {
		cfg.HeaderOriginTypeKey = ConfigDefault.HeaderOriginTypeKey
	}
	if cfg.HeaderOriginKey == "" {
		cfg.HeaderOriginKey = ConfigDefault.HeaderOriginKey
	}
	return cfg
}
