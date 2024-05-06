package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	FirstName    *string            `json:"firstname" validation:"required,min=2,max=100"`
	LastName     *string            `json:"lastname" validation:"required,min=2,max=100"`
	Password     *string            `json:"password" validation:"required,min=2,max=10"`
	Email        *string            `json:"email" validation:"email,required"`
	Phone        *string            `json:"phone" validation:"required"`
	Token        *string            `json:"token"`
	UserType     *string            `json:"usertype" validation:"required,eq=ADMIN|eq=USER"` //enum
	RefreshToken *string            `json:"refreshtoken"`
	CreatedAt    time.Time          `json:"createdat"`
	UpdatedAt    time.Time          `json:"updatedat"`
	UserID       *string            `json:"userid"`
}
