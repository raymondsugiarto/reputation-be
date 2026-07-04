package shared

import (
	"context"

	"github.com/raymondsugiarto/reputation-be/pkg/entity"
)

var UserContextKey = "user"
var UserCredentialDataKey = "userCredentialData"

type UserCredentialData struct {
	ID             string `json:"id"`  // user credential id
	UserID         string `json:"uid"` // user id
	OrganizationID string `json:"oid"` // organization id
	CustomerID     string `json:"cid"` // user id
	AccountID      string `json:"aid"` // user id
}

type OrganizationData struct {
	ID string `json:"id"`
}

func GetOrigin(ctx context.Context) string {
	return ctx.Value(entity.OriginKey).(string)
}

func GetOriginTypeKey(ctx context.Context) string {
	return ctx.Value(entity.OriginTypeKey).(string)
}

func GetOrganization(ctx context.Context) *OrganizationData {
	return ctx.Value(entity.OrganizationKey).(*OrganizationData)
}

func GetUserCredential(ctx context.Context) *UserCredentialData {
	return ctx.Value(UserCredentialDataKey).(*UserCredentialData)
}
