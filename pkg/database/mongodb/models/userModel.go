package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          *string            `json"name" validate:"required, min=2, max=100"`
	Email         *string            `json"email validate:"required, min=2, max=100"`
	Password      *string            `json"password validate:"email, required"`
	Token         *string            `json"token"`
	Refresh_token *string            `json"refresh_token"`
	User_id       string             `json"user_id"`
	User_type     *string            `json"user_type" validate:"required, eq=ADMIN|eq=USER"`
	Created_at    time.Time          `json"created_at"`
	Updated_at    time.Time          `json"updated_at"`
}
