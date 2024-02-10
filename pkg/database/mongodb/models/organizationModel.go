package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	ID          primitive.ObjectID `bson:"_id, omitempty" json:"_id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Members     []Member           `bson:"members" json:"members"`
}

type Member struct {
	UserID string `bson:"user_id" json:"user_id"`
	Role   string `bson:"role" json:"role"`
}
type InviteUserRequest struct {
	UserEmail string `json:"user_email" binding:"required"`
	UserID    string `bson:"user_id" json:"user_id"`
	Role      string `bson:"role" json:"role"`
}
