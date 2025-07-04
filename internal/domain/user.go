package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password"`
	Role      Role               `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at" `
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at" `
}
