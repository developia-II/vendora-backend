package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}
type User struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	Email            string             `json:"email" bson:"email" validate:"required,email"`
	Name             string             `json:"name" bson:"name"`
	Address          string             `json:"address" bson:"address"`
	Phone            string             `json:"phone" bson:"phone"`
	Password         string             `json:"-" bson:"password" validate:"required,min=6"`
	Role             string             `json:"role" bson:"role"`
	CreatedAt        time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsVerified       bool               `json:"isverified" bson:"Isverified"`
	ResetToken       string             `json:"-" bson:"resetToken,omitempty"`
	ResetTokenExpiry time.Time          `json:"-" bson:"resetTokenExpiry,omitempty"`
	PasswordResetAt  time.Time          `json:"-" bson:"passwordResetAt,omitempty"`
}
