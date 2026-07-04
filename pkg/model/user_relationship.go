package model

import concern "github.com/raymondsugiarto/reputation-be/pkg/model/common"

type UserRelationshipType string

const (
	UserRelationshipFollow UserRelationshipType = "FOLLOW"
)

type UserRelationship struct {
	concern.CommonWithIDs
	UserID         string
	UserIDFollower string
	Relationship   UserRelationshipType
}
