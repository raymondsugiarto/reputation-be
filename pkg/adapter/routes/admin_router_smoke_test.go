// Package routes_test — smoke test that verifies the admin router mounts
// the expected endpoints without requiring a live database. The test
// manually creates a Fiber app, registers a minimal fake
// admin/customer service pair into app.State(), then asks the router to
// mount and inspects the resulting route stack.
//
// We can't exercise the handlers end-to-end (they need a real DB), but
// this guards against the most common failure mode: forgetting to wire
// the admin service into the DI container, which would cause a panic
// at route-mount time with `fiber.MustGetState`.
package routes_test

import (
	"context"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/raymondsugiarto/reputation-be/pkg/adapter/routes"
	"github.com/raymondsugiarto/reputation-be/pkg/entity"
	"github.com/raymondsugiarto/reputation-be/pkg/module/admin"
	"github.com/raymondsugiarto/reputation-be/pkg/module/customer"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
)

// fakeAdminService is a stub that satisfies admin.Service. Method bodies
// are intentionally empty — the test only validates route-mount, not
// handler behaviour.
type fakeAdminService struct{}

func (f *fakeAdminService) GetProfile(_ context.Context, _ *entity.UserSessionDto) *entity.AdminProfileDto {
	return nil
}
func (f *fakeAdminService) GetStats(_ context.Context) (*entity.AdminStatsDto, error) {
	return nil, nil
}

// fakeCustomerService is a stub that satisfies customer.Service so that
// routes.AdminRouter can resolve it from app.State().
type fakeCustomerService struct{}

func (f *fakeCustomerService) SignUp(_ context.Context, _ *entity.CustomerSignUpRequestDto) (*entity.CustomerSignUpResponseDto, error) {
	return nil, nil
}
func (f *fakeCustomerService) Create(_ context.Context, _ *entity.CustomerDto) (*entity.CustomerDto, error) {
	return nil, nil
}
func (f *fakeCustomerService) FindByID(_ context.Context, _ string) (*entity.CustomerDto, error) {
	return nil, nil
}
func (f *fakeCustomerService) FindByUserID(_ context.Context, _ string) (*entity.CustomerDto, error) {
	return nil, nil
}
func (f *fakeCustomerService) FindAll(_ context.Context, _ pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	return nil, nil
}
func (f *fakeCustomerService) Update(_ context.Context, _ *entity.CustomerDto) (*entity.CustomerDto, error) {
	return nil, nil
}
func (f *fakeCustomerService) Delete(_ context.Context, _ string) error {
	return nil
}
func (f *fakeCustomerService) Approve(_ context.Context, _ string, _ string, _ string) (*customer.ApprovalResultDto, error) {
	return nil, nil
}
func (f *fakeCustomerService) Reject(_ context.Context, _ string, _ string, _ string) (*customer.ApprovalResultDto, error) {
	return nil, nil
}
func (f *fakeCustomerService) FindPendingApprovals(_ context.Context, _ pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	return nil, nil
}
func (f *fakeCustomerService) FindApprovalHistory(_ context.Context, _ pagination.PaginationRequestDto) (*pagination.ResultPagination[entity.CustomerDto], error) {
	return nil, nil
}
func (f *fakeCustomerService) GetApprovalStats(_ context.Context) (*entity.CustomerApprovalStatsDto, error) {
	return nil, nil
}

// Compile-time interface satisfaction guards. If the real interfaces
// change, these will fail to compile — which is exactly what we want.
var (
	_ admin.Service    = (*fakeAdminService)(nil)
	_ customer.Service = (*fakeCustomerService)(nil)
)

func TestAdminRouterMounts(t *testing.T) {
	app := fiber.New()
	app.State().Set(admin.ServiceName, &fakeAdminService{})
	app.State().Set(customer.ServiceName, &fakeCustomerService{})

	api := app.Group("/api")
	adminGroup := api.Group("/admin")
	// If AdminRouter is broken (service lookup panics), this call will
	// fatal-error before the assertion.
	routes.AdminRouter(app, adminGroup)

	// Walk the registered routes and confirm the expected paths.
	expected := map[string]bool{
		"/api/admin/me":                            false,
		"/api/admin/stats":                         false,
		"/api/admin/customers":                     false,
		"/api/admin/approvals/":                    false,
		"/api/admin/approvals/history":             false,
		"/api/admin/approvals/stats":               false,
		"/api/admin/approvals/:customerId/approve": false,
		"/api/admin/approvals/:customerId/reject":  false,
	}
	for _, route := range app.Stack() {
		for _, r := range route {
			if _, ok := expected[r.Path]; ok {
				expected[r.Path] = true
			}
		}
	}
	for path, found := range expected {
		if !found {
			t.Errorf("expected admin route %q to be registered", path)
		}
	}
}
