package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SchoolData struct {
	Id                  primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt           string             `json:"created_at" bson:"created_at"`
	UpdatedAt           string             `json:"updated_at" bson:"updated_at"`
	Type                string             `json:"type" bson:"type"`
	Year                *string            `json:"year,omitempty" bson:"year,omitempty"`
	Term                *string            `json:"term,omitempty" bson:"term,omitempty"`
	Status              *bool              `json:"status,omitempty" bson:"status,omitempty"`
	SubjectCategory     *string            `json:"subject_category,omitempty" bson:"subject_category,omitempty"`
	InformationCategory *string            `json:"information_category,omitempty" bson:"information_category,omitempty"`
}

type SchoolDataRequest struct {
	Id       string `json:"id"`
	Year     string `json:"year"`
	Term     string `json:"term"`
	Status   *bool  `json:"status"`
	Category string `json:"category" bson:"category"`
}

type YearAndTerm struct {
}
