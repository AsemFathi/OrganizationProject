package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	ID          primitive.ObjectID  `bson:"_id, omitempty" json:"_id,omitempty"`
	Name        string              `bson:"name" json:"name"`
	Description string              `bson:"description" json:"description"`
	Org_ID      string              `bson:"organization_id" json:"organization_id"`
	Created_By  string              `bson:"created_by" json:"created_by"`
	Members     []InviteUserRequest `bson:"organization_members" json:"organization_members"`
}

type InviteUserRequest struct {
	UserEmail string `bson:"user_email" json:"user_email" binding:"required"`
	UserID    string `bson:"user_id" json:"user_id"`
	Role      string `bson:"access_level" json:"access_level"`
}
