package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name        string              `json:"name" bson:"name" validate:"required"`
	Description string              `json:"description,omitempty" bson:"description"`
	ParentID    *primitive.ObjectID `json:"parentId,omitempty" bson:"parentId"` // For subcategories
}
