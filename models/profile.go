package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ProfileAdmin struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	UpdatedAt string             `json:"updated_at" bson:"updated_at"`
	ProfileId string             `json:"profile_id" bson:"profile_id"`
	Name      string             `json:"name" bson:"name"`
	Role      string             `json:"role" bson:"role"`
}

type ProfileTeacher struct {
	Id                 primitive.ObjectID   `json:"id" bson:"_id"`
	CreatedAt          string               `json:"created_at" bson:"created_at"`
	UpdatedAt          string               `json:"updated_at" bson:"updated_at"`
	ProfileId          string               `json:"profile_id" bson:"profile_id"`
	Name               string               `json:"name" bson:"name"`
	Role               string               `json:"role" bson:"role"`
	SubjectId          string               `json:"subject_id" bson:"subject_id"`
	ClassInCounseling  string               `json:"class_in_counseling" bson:"class_in_counseling"`
	CoursesTeachesList []CoursesTeachesList `json:"courses_teaches_list" bson:"courses_teaches_list"`
}

type CoursesTeachesList struct {
	CoursesIdList []string `json:"courses_id_list" bson:"courses_id_list"`
	Year          int      `json:"year" bson:"year"`
	Term          int      `json:"term" bson:"term"`
}
