package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Subject struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt    string             `json:"created_at" bson:"created_at"`
	UpdatedAt    string             `json:"updated_at" bson:"updated_at"`
	SubjectId    string             `json:"subject_id" bson:"subject_id"`
	Name         string             `json:"name" bson:"name"`
	Category     string             `json:"category" bson:"category"`
	Credit       int                `json:"credit" bson:"credit"`
	ClassYear    string             `json:"class_year" bson:"class_year"`
	InstructorId []string           `json:"instructor_id" bson:"instructor_id"`
}

type SubjectRequest struct {
	Event        string   `json:"event"`
	SubjectId    string   `json:"subject_id"`
	Name         string   `json:"name"`
	InstructorId []string `json:"instructor_id"`
	Credit       *int     `json:"credit"`
	Category     string   `json:"category" bson:"category"`
	ClassYear    string   `json:"class_year" bson:"class_year"`
}
