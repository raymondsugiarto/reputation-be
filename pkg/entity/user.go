package entity

import (
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/database/pagination"
	"github.com/raymondsugiarto/reputation-be/pkg/shared/security"
)

// to be
type UserDto struct {
	ID              string              `json:"id"`
	OrganizationID  string              `json:"organizationId"`
	UserType        model.UserType      `json:"userType"`
	UserCredentials []UserCredentialDto `json:"userCredential"`
}

func NewUserDtoFromModel(m *model.User) *UserDto {
	return &UserDto{
		ID:              m.ID,
		OrganizationID:  m.OrganizationID,
		UserType:        m.UserType,
		UserCredentials: []UserCredentialDto{},
	}
}

func (f *UserDto) FromModel(m *model.User) UserDto {
	return *NewUserDtoFromModel(m)
}

func (f *UserDto) ToModel() *model.User {
	m := &model.User{
		OrganizationID: f.OrganizationID,
	}
	if len(f.UserCredentials) > 0 {
		m.UserCredentials = make([]model.UserCredential, len(f.UserCredentials))
		for i, uc := range f.UserCredentials {
			m.UserCredentials[i] = *uc.ToModel()
		}
	}
	if f.ID != "" {
		m.ID = f.ID
	}
	return m
}

// UserCredentialDto
type UserCredentialDto struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organizationId"`
	UserID         string   `json:"userId"`
	Username       string   `json:"username"`
	Password       string   `json:"-"`
	User           *UserDto `json:"user,omitempty"`
}

func NewUserCredentialDtoFromModel(m *model.UserCredential) *UserCredentialDto {
	d := &UserCredentialDto{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		Username:       m.Username,
	}
	if m.User != nil {
		d.User = NewUserDtoFromModel(m.User)
	}
	return d
}

func (f *UserCredentialDto) FromModel(m *model.UserCredential) UserCredentialDto {
	return *NewUserCredentialDtoFromModel(m)
}

func (f *UserCredentialDto) ToModel() *model.UserCredential {
	encryptedPass, _ := security.HashPassword(f.Password)
	m := &model.UserCredential{
		OrganizationID: f.OrganizationID,
		UserID:         f.UserID,
		Username:       f.Username,
		Password:       encryptedPass,
	}
	if f.ID != "" {
		m.ID = f.ID
	}
	return m
}

type UserFilterDto struct {
	pagination.GetListRequest
}
