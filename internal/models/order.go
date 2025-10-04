package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ProductID primitive.ObjectID `json:"productId" bson:"productId"`
	Quantity  int                `json:"quantity" bson:"quantity"`
	Price     float64            `json:"price" bson:"price"`
}

type Order struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID          primitive.ObjectID `json:"userId" bson:"userId"`
	Items           []OrderItem        `json:"items" bson:"items"`
	Total           float64            `json:"total" bson:"total"`
	Status          string             `json:"status" bson:"status"`
	ShippingAddress string             `json:"shippingAddress" bson:"shippingAddress"`
	PaymentStatus   string             `json:"paymentStatus" bson:"paymentStatus"`
	PaymentID       string             `json:"paymentId" bson:"paymentId"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
}
