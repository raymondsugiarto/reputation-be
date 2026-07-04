package model

import (
	concern "github.com/raymondsugiarto/reputation-be/pkg/model/common"
)

type AccountType string

const (
	AccountCash    AccountType = "CASH"
	AccountBank    AccountType = "BANK"
	AccountEWallet AccountType = "E_WALLET"
)

type Account struct {
	concern.CommonWithIDs
	OrganizationID string
	CustomerID     string
	// Customer       *Customer `gorm:"foreignKey:CustomerID;references:ID"`
	AccountName string
	AccountType AccountType
	Balance     float64
	Currency    string // IDR
}
