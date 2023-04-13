package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// กำหนดคะแนนตามหลักสูตรไว้แล้ว เต็ม 100
type Course struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	UpdatedAt string             `json:"updated_at" bson:"updated_at"`
	// CourseId        string                `json:"course_id" bson:"course_id"`
	Status          string              `json:"status" bson:"status"`
	SubjectId       string              `json:"subject_id" bson:"subject_id"`
	InstructorId    string              `json:"instructor_id" bson:"instructor_id"`
	Name            string              `json:"name" bson:"name"`
	Credit          int                 `json:"credit" bson:"credit"`
	Year            string              `json:"year" bson:"year"`
	Term            string              `json:"term" bson:"term"`
	NumberOfStudent int                 `json:"number_of_student" bson:"number_of_student"`
	StudentIdList   []string            `json:"student_id_list" bson:"student_id_list"`
	LocationId      *primitive.ObjectID `json:"location_id" bson:"location_id"`
	DateTime        []DateTime          `json:"date_time" bson:"date_time"`
	ClassId         *primitive.ObjectID `json:"class_id" bson:"class_id"`
	ClassYear       string              `json:"class_year" bson:"class_year"`
	ClassRoom       string              `json:"class_room" bson:"class_room"`
	// ScoreWorkFull   float64
	// ScoreMidFull    float64
	// ScoreFinalFull  float64
}

type DateTime struct {
	Day  string   `json:"day" bson:"day"`
	Time []string `json:"time" bson:"time"`
}

type CourseRequest struct {
	Id           string     `json:"id"`
	Event        string     `json:"event"`
	SubjectId    string     `json:"subject_id"`
	InstructorId string     `json:"instructor_id"`
	LocationId   string     `json:"location_id"`
	DateTime     []DateTime `json:"date_time"`
	ClassId      string     `json:"class_id"`
}
