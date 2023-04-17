package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Score struct {
	Id               primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt        string             `json:"created_at" bson:"created_at"`
	UpdatedAt        string             `json:"updated_at" bson:"updated_at"`
	CourseId         string             `json:"course_id" bson:"course_id"`
	Type             string             `json:"type" bson:"type"`
	Name             string             `json:"name" bson:"name"`
	ScoreFull        float64            `json:"score_full" bson:"score_full"`
	ScoreInformation []ScoreInformation `json:"score_information,omitempty" bson:"score_information,omitempty"`
}

type ScoreInformation struct {
	StudentId string   `json:"student_id" bson:"student_id"`
	UpdatedAt string   `json:"updated_at" bson:"updated_at"`
	ScoreGet  *float64 `json:"score_get,omitempty" bson:"score_get,omitempty"`
	Status    string   `json:"status" bson:"status"`
	Note      *string  `json:"note,omitempty" bson:"note,omitempty"`
}

type ScoreRequest struct {
	CourseId  string   `json:"course_id" bson:"course_id"`
	Name      string   `json:"name" bson:"name"`
	ScoreFull *float64 `json:"score_full" bson:"score_full"`
	ScoreGet  *float64 `json:"score_get" bson:"score_get"`
	Status    string   `json:"status" bson:"status"`
	StudentId string   `json:"student_id" bson:"student_id"`
	Type      string   `json:"type" bson:"type"`
}

type ScoreStudentRes struct {
	Name      string   `json:"name" bson:"name"`
	UpdatedAt string   `json:"updated_at" bson:"updated_at"`
	ScoreFull float64  `json:"score_full" bson:"score_full"`
	ScoreGet  *float64 `json:"score_get" bson:"score_get"`
	Status    string   `json:"status" bson:"status"`
}
