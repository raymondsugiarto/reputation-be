package model

import concern "github.com/raymondsugiarto/reputation-be/pkg/model/common"

type UserCredential struct {
	concern.CommonWithIDs
	OrganizationID string
	UserID         string
	User           *User
	Username       string
	Password       string
}
