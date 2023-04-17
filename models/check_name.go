package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// attend , absent , leave , late
type CheckName struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt     string             `json:"created_at" bson:"created_at"`
	UpdatedAt     string             `json:"updated_at" bson:"updated_at"`
	CourseId      string             `json:"course_id" bson:"course_id"`
	Date          string             `json:"date" bson:"date"`
	Status        string             `json:"status" bson:"status"`
	TimeLate      string             `json:"time_late" bson:"time_late"`
	CheckNameData []CheckNameData    `json:"check_name_data,omitempty" bson:"check_name_data,omitempty"`
}

type CheckNameData struct {
	StudentId string  `json:"student_id" bson:"student_id"`
	UpdatedAt string  `json:"updated_at" bson:"updated_at"`
	Time      string  `json:"time" bson:"time"`
	Status    string  `json:"status" bson:"status"`
	CheckBy   string  `json:"check_by" bson:"check_by"`
	Note      *string `json:"note" bson:"note"`
}

type CheckNameRequest struct {
	// Event     string  `json:"event"`
	// StudentId string  `json:"student_id"`
	Date     string `json:"date"`
	TimeLate *int   `json:"time_late"`
	// Time      string  `json:"time"`
	CourseId  string `json:"course_id"`
	StudentId string `json:"student_id" bson:"student_id"`
	CheckBy   string `json:"check_by" bson:"check_by"`
	// Status    string  `json:"status"`
	// Note      *string `json:"note"`

	// TimeEnd string `json:"time_end"`
}

type CheckNameStudentRes struct {
	Date      string `json:"date"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
	Time      string `json:"time" bson:"time"`
	Status    string `json:"status" bson:"status"`
	CheckBy   string `json:"check_by" bson:"check_by"`
}
