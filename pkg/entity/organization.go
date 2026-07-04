package entity

import (
	"github.com/raymondsugiarto/reputation-be/pkg/model"
	concern "github.com/raymondsugiarto/reputation-be/pkg/model/common"
)

const (
	OriginKey       = "x-origin"
	OriginTypeKey   = "x-origin-type" // ADMIN, COMPANY, EMPLOYEE
	OrganizationKey = "organization"
)

type OrganizationDto struct {
	concern.CommonWithIDs
	Code   string
	Name   string
	Origin string
}

func NewOrganizationDtoFromModel(m *model.Organization) *OrganizationDto {
	return &OrganizationDto{
		CommonWithIDs: concern.CommonWithIDs{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
		},
		Code:   m.Code,
		Name:   m.Name,
		Origin: m.Origin,
	}
}
