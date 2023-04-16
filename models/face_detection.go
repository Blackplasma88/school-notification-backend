package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type FaceDetectData struct {
	Id                   primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt            string             `json:"created_at" bson:"created_at"`
	UpdatedAt            string             `json:"updated_at" bson:"updated_at"`
	Status               string             `json:"status" bson:"status"`
	Name                 string             `json:"name" bson:"name"`
	ClassId              string             `json:"class_id" bson:"class_id"`
	NumberOfImage        int                `json:"number_of_image" bson:"number_of_image"`
	NumberOfStudent      int                `json:"number_of_student" bson:"number_of_student"`
	StudentIdList        []string           `json:"student_id_list" bson:"student_id_list"`
	ImageStudentPathList [][]string         `json:"image_student_path_list" bson:"image_student_path_list"`
}

type FaceDetectDataRequest struct {
	ClassId       string   `json:"class_id"`
	Id            string   `json:"id" bson:"_id"`
	StudentId     string   `json:"student_id"`
	ImagePathList []string `json:"image_path_list" bson:"image_path_list"`
}
