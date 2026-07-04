package middleware_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/infrastructure/middleware"
	"github.com/raymondsugiarto/reputation-be/pkg/model"
)

// putSession writes a synthetic session into the request locals so that
// AdminOnly has something to inspect. AdminOnly only ever reads
// c.Locals(entity.UserSessionKey), so we set it via a tiny stub handler.
func putSession(s *entity.UserSessionDto) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals(entity.UserSessionKey, s)
		return c.Next()
	}
}

func adminSession() *entity.UserSessionDto {
	return &entity.UserSessionDto{
		ID:             "cred-1",
		UserID:         "user-1",
		OrganizationID: "org-1",
		UserCredential: &entity.UserCredentialDto{
			ID:             "cred-1",
			OrganizationID: "org-1",
			UserID:         "user-1",
			Username:       "admin",
			User: &entity.UserDto{
				ID:             "user-1",
				OrganizationID: "org-1",
				UserType:       model.INTERNAL_ADMIN,
			},
		},
	}
}

func orgAdminSession() *entity.UserSessionDto {
	return &entity.UserSessionDto{
		ID:             "cred-2",
		UserID:         "user-2",
		OrganizationID: "org-2",
		UserCredential: &entity.UserCredentialDto{
			ID:             "cred-2",
			OrganizationID: "org-2",
			UserID:         "user-2",
			Username:       "owner",
			User: &entity.UserDto{
				ID:             "user-2",
				OrganizationID: "org-2",
				UserType:       model.ADMIN, // org-owner, NOT internal
			},
		},
	}
}

func customerSession() *entity.UserSessionDto {
	return &entity.UserSessionDto{
		ID:             "cred-3",
		UserID:         "user-3",
		OrganizationID: "",
		UserCredential: &entity.UserCredentialDto{
			ID:             "cred-3",
			OrganizationID: "",
			UserID:         "user-3",
			Username:       "customer",
			User: &entity.UserDto{
				ID:       "user-3",
				UserType: model.CUSTOMER,
			},
		},
	}
}

func newAppWith(session *entity.UserSessionDto) *fiber.App {
	// Mirror production: wire DefaultErrorHandler so typed AppStatus
	// errors are mapped to the right HTTP code + JSON envelope.
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.DefaultErrorHandler(),
	})
	app.Use(putSession(session))
	app.Use(middleware.AdminOnly())
	app.Get("/probe", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})
	return app
}

func TestAdminOnly_AllowsInternalAdmin(t *testing.T) {
	app := newAppWith(adminSession())
	req := httptest.NewRequest("GET", "/probe", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected 200 for INTERNAL_ADMIN, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ok" {
		t.Fatalf("expected body 'ok', got %q", string(body))
	}
}

func TestAdminOnly_RejectsOrgOwner(t *testing.T) {
	app := newAppWith(orgAdminSession())
	req := httptest.NewRequest("GET", "/probe", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403 for org-owner ADMIN, got %d", resp.StatusCode)
	}
}

func TestAdminOnly_RejectsCustomer(t *testing.T) {
	app := newAppWith(customerSession())
	req := httptest.NewRequest("GET", "/probe", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 403 {
		t.Fatalf("expected 403 for CUSTOMER, got %d", resp.StatusCode)
	}
}

func TestAdminOnly_RejectsMissingSession(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.DefaultErrorHandler(),
	})
	app.Use(middleware.AdminOnly())
	app.Get("/probe", func(c fiber.Ctx) error { return c.SendString("ok") })

	req := httptest.NewRequest("GET", "/probe", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected 401 for missing session, got %d", resp.StatusCode)
	}
}

// Quick check: the rejection payload matches the platform's standard
// envelope so the FE's ApiError handler can render the message.
func TestAdminOnly_RejectionPayloadIsValidJSON(t *testing.T) {
	app := newAppWith(orgAdminSession())
	req := httptest.NewRequest("GET", "/probe", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("expected JSON payload, got %q (err=%v)", string(body), err)
	}
}

// Ensure the middleware is safely constructible without config — confirms
// the call signature matches the existing middleware pattern in the
// router.
func TestAdminOnly_Constructor(t *testing.T) {
	h := middleware.AdminOnly()
	if h == nil {
		t.Fatalf("AdminOnly() returned nil")
	}
	_ = context.Background()
}
