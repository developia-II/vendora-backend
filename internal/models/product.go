package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `json:"name" bson:"name" validate:"required"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price" validate:"required,gt=0"`
	Stock       int                `json:"stock" bson:"stock" validate:"gte=0"`
	CategoryID  primitive.ObjectID `json:"categoryId" bson:"categoryId"` // Reference
	Images      []string           `json:"images" bson:"images"`         // URLs
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}
