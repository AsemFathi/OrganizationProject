package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          *string            `json"name" validate:"required, min=2, max=100"`
	Email         *string            `json"email validate:"email, required"`
	Password      *string            `json"password validate:"required, min=2, max=100"`
	Token         *string            `json"token" bson:"token"`
	Refresh_token *string            `json"refresh_token" bson:"refresh_token"`
	User_id       string             `json"user_id" bson:"user_id"`
	User_type     *string            `json"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Created_at    time.Time          `json:"created_at" bson:"created_at"`
	Updated_at    time.Time          `json:"updated_at" bson:"updated_at"`
}
