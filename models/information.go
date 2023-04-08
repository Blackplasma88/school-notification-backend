package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Information struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   string             `json:"created_at" bson:"created_at"`
	UpdatedAt   string             `json:"updated_at" bson:"updated_at"`
	Name        string             `json:"name" bson:"name"`
	FilePath    string             `json:"filepath" bson:"filepath"`
	Description string             `json:"description" bson:"description"`
	Content     string             `json:"content" bson:"content"`
	Category    string             `json:"category" bson:"category"`
}
