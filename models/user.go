package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	ProfileId string             `json:"profile_id" bson:"profile_id"`
	Role      string             `json:"role" bson:"role"`
}

type UserRequest struct {
	Username  string `json:"username" bson:"username"`
	Password  string `json:"password" bson:"password"`
	ProfileId string `json:"profile_id" bson:"profile_id"`
	Role      string `json:"role" bson:"role"`
}
