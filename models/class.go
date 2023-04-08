package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ClassData struct {
	Id              primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt       string             `json:"created_at" bson:"created_at"`
	UpdatedAt       string             `json:"updated_at" bson:"updated_at"`
	ClassYear       string             `json:"class_year" bson:"class_year"`
	ClassRoom       string             `json:"class_room" bson:"class_room"`
	Year            string             `json:"year" bson:"year"`
	Term            string             `json:"term" bson:"term"`
	NumberOfStudent int                `json:"number_of_student" bson:"number_of_student"`
	StudentIdList   []string           `json:"student_id_list" bson:"student_id_list"`
}
