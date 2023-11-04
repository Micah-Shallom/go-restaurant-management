package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID 				primitive.ObjectID 		`bson:"_id"`
	FirstName		*string					`json:"firstname" validate:"required,min=2,max=100"`
	LastName		*string					`json:"lastname" validate:"required,min=2,max=100"`
	Password		*string					`json:"password" validate:"required,min=6"`
	Email			*string					`json:"email" validate:"email,required"`
	Phone			*string					`json:"phone" validate:"required"`
	Token			*string					`json:"token"`
	UserType		*string					`json:"usertype" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken	*string					`json:"refreshtoken"`
	CreatedAt		time.Time				`json:"created_at"`
	UpdatedAt		time.Time				`json:"updated_at"`
	UserID			string					`json:"userid"`
} 