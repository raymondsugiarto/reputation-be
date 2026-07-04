package entity

import (
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
)

type AccountRequest struct {
	CustomerID  string  `json:"customerId"`
	AccountName string  `json:"accountName"`
	AccountType string  `json:"accountType"` // CASH, BANK, E_WALLET
	Balance     float64 `json:"balance,omitempty"`
	Currency    string  `json:"currency,omitempty"`
}

func (r *AccountRequest) ToDto() *AccountDto {
	balance := r.Balance
	if balance == 0 {
		balance = 0.0
	}
	return &AccountDto{
		CustomerID:  r.CustomerID,
		AccountName: r.AccountName,
		AccountType: model.AccountType(r.AccountType),
		Balance:     balance,
		Currency:    r.Currency,
	}
}

type AccountDto struct {
	ID             string            `json:"id"`
	OrganizationID string            `json:"organizationId"`
	CustomerID     string            `json:"customerId"`
	AccountName    string            `json:"accountName"`
	AccountType    model.AccountType `json:"accountType"`
	Balance        float64           `json:"balance"`
	Currency       string            `json:"currency,omitempty"`
}

func NewAccountDtoFromModel(m *model.Account) *AccountDto {
	if m == nil {
		return nil
	}

	return &AccountDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		CustomerID:     m.CustomerID,
		AccountName:    m.AccountName,
		AccountType:    m.AccountType,
		Balance:        m.Balance,
		Currency:       m.Currency,
	}
}

func (a *AccountDto) ToModel() *model.Account {
	balance := a.Balance
	if balance == 0 {
		balance = 0.0
	}
	account := &model.Account{
		OrganizationID: a.OrganizationID,
		AccountName:    a.AccountName,
		AccountType:    a.AccountType,
		Balance:        balance,
		Currency:       a.Currency,
	}
	if a.ID != "" {
		account.ID = a.ID
	}
	// Convert empty string to nil for nullable fields
	if a.CustomerID != "" {
		account.CustomerID = a.CustomerID
	}
	return account
}

// AccountFilterDto for query filtering
type AccountFilterDto struct {
	pagination.GetListRequest
	CustomerID  string `query:"customerId"`
	AccountName string `query:"accountName"`
	AccountType string `query:"accountType"`
}
