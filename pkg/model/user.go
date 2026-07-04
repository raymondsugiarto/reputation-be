package model

import (
	concern "github.com/raymondsugiarto/reputation-be/pkg/model/common"
)

type UserType string

const (
	// ADMIN is the org-owner admin user. Created via /api/auth/sign-up
	// and tied to a single organization. Can manage customers under
	// that organization only.
	ADMIN UserType = "ADMIN"
	// CUSTOMER is the regular end-user (perorangan / perusahaan) of the
	// platform. Authenticated via /api/customer/sign-up.
	CUSTOMER UserType = "CUSTOMER"
	// INTERNAL_ADMIN is a platform-internal employee account. Seeded
	// directly into the DB. Has cross-organization visibility and is
	// the only role allowed to access /api/admin/* endpoints.
	INTERNAL_ADMIN UserType = "INTERNAL_ADMIN"
)

type User struct {
	concern.CommonWithIDs
	OrganizationID  string
	UserType        UserType
	UserCredentials []UserCredential
}
