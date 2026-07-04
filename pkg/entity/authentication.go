package entity

type LoginRequestDto struct {
	Username string `json:"username" bson:"username" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type LoginDto struct {
	Token   string `json:"token" bson:"token"`
	Expired string `json:"expired" bson:"exp"`
}

var UserSessionKey = "userSessionKey"

type UserSessionDto struct {
	ID             string             `json:"id"`   // user credential id
	UserID         string             `json:"uid"`  // user id
	OrganizationID string             `json:"oid"`  // organization id
	CustomerID     string             `json:"cid"`  // user id
	AccountID      string             `json:"aid"`  // user id
	UserCredential *UserCredentialDto `json:"user"` // additional data
}

func NewUserSessionDtoFromClaims(claims map[string]interface{}) *UserSessionDto {
	return &UserSessionDto{
		ID:             claims["id"].(string),
		UserID:         claims["uid"].(string),
		OrganizationID: claims["oid"].(string),
	}
}

type SignUpRequestDto struct {
	Name     string `json:"name" bson:"name" validate:"required"`
	Username string `json:"username" bson:"username" validate:"required"`
	Password string `json:"password" bson:"password" validate:"required"`
}

type SignUpResponseDto struct {
	OrganizationID string `json:"organizationId" bson:"organizationId"`
	UserID         string `json:"userId" bson:"userId"`
}
