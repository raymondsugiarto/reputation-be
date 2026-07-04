package model

import concern "github.com/raymondsugiarto/reputation-be/pkg/model/common"

// Accounts : table accounts
type Organization struct {
	concern.CommonWithIDs
	Code   string
	Name   string
	Origin string
}
