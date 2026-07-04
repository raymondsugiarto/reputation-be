package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/response/status"
)

// AdminOnly restricts a route to platform-internal employees only. It must
// run AFTER middleware.Protected() so that UserSessionKey is already
// populated on c.Locals(). Regular org-owner ADMIN users (user_type=ADMIN)
// are rejected — only INTERNAL_ADMIN accounts may proceed.
//
// Unlike the organization middleware, this guard does NOT touch
// x-origin / x-origin-type. Internal admins are not bound to any tenant.
func AdminOnly() fiber.Handler {
	return func(c fiber.Ctx) error {
		sessionRaw := c.Locals(entity.UserSessionKey)
		if sessionRaw == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "invalid session",
				"data":    nil,
			})
		}

		session, ok := sessionRaw.(*entity.UserSessionDto)
		if !ok || session == nil {
			return status.New(status.InvalidSession)
		}

		// UserCredential is populated by authentication.SuccessHandler via
		// authenticationSvc.GetSession. If it's missing, the session is
		// incomplete and we reject.
		if session.UserCredential == nil || session.UserCredential.User == nil {
			return status.New(status.InvalidSession)
		}

		if session.UserCredential.User.UserType != model.INTERNAL_ADMIN {
			log.WithContext(c).Warnf(
				"non-internal-admin user attempted admin route: user_id=%s user_type=%s",
				session.UserCredential.User.ID,
				session.UserCredential.User.UserType,
			)
			return status.New(status.Forbidden)
		}

		return c.Next()
	}
}
