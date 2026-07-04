package entity

import "time"

// AdminProfileDto is the response of GET /api/admin/me — a trimmed view of
// the internal-admin user that the FE uses to render the sidebar (display
// name) and to verify session freshness on page load.
type AdminProfileDto struct {
	ID             string `json:"id"`
	UserID         string `json:"userId"`
	OrganizationID string `json:"organizationId"`
	Username       string `json:"username"`
	UserType       string `json:"userType"`
}

// AdminStatsDto powers the simple admin dashboard. Numbers are kept flat
// (no nesting) so the FE can render KPI cards without extra plumbing.
//
// All counts are platform-wide — internal admins see every customer, every
// organization, regardless of tenant.
type AdminStatsDto struct {
	TotalCustomers      int64     `json:"totalCustomers"`
	TotalOrganizations  int64     `json:"totalOrganizations"`
	TotalInternalAdmins int64     `json:"totalInternalAdmins"`
	CustomersThisMonth  int64     `json:"customersThisMonth"`
	IndividualCustomers int64     `json:"individualCustomers"`
	CompanyCustomers    int64     `json:"companyCustomers"`
	GeneratedAt         time.Time `json:"generatedAt"`
}
