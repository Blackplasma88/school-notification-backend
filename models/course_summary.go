package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CourseSummary struct {
	Id          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   string             `json:"created_at" bson:"created_at"`
	UpdatedAt   string             `json:"updated_at" bson:"updated_at"`
	CourseId    string             `json:"course_id" bson:"course_id"`
	StudentData []StudentData      `json:"student_data" bson:"student_data"`
}

type StudentData struct {
	StudentId     string  `json:"student_id" bson:"student_id"`
	ScoreWorkGet  float64 `json:"score_work_get" bson:"score_work_get"`
	ScoreWorkFull float64 `json:"score_work_full" bson:"score_work_full"`
	ScoreMidGet   float64 `json:"score_mid_get" bson:"score_mid_get"`
	ScoreMidFull  float64 `json:"score_mid_full" bson:"score_mid_full"`
	ScoreFinalGet float64 `json:"score_final_get" bson:"score_final_get"`
	ScoreFinaFull float64 `json:"score_final_full" bson:"score_final_full"`
	Grade         float64 `json:"grade" bson:"grade"`

	AllDateCount         int `json:"all_date_count" bson:"all_date_count"`
	CheckNameAttendCount int `json:"check_name_attend_count" bson:"check_name_attend_count"`
	CheckNameAbsentCount int `json:"check_name_absent_count" bson:"check_name_absent_count"`
	// CheckNameLeaveCount  int
	CheckNameLateCount int `json:"check_name_late_count" bson:"check_name_late_count"`
	// CheckNameStatus      bool
	// CheckNameStatusMsg   string
}

type StudentDataRes struct {
	CourseName    string  `json:"course_name" bson:"course_name"`
	ScoreWorkGet  float64 `json:"score_work_get" bson:"score_work_get"`
	ScoreWorkFull float64 `json:"score_work_full" bson:"score_work_full"`
	ScoreMidGet   float64 `json:"score_mid_get" bson:"score_mid_get"`
	ScoreMidFull  float64 `json:"score_mid_full" bson:"score_mid_full"`
	ScoreFinalGet float64 `json:"score_final_get" bson:"score_final_get"`
	ScoreFinaFull float64 `json:"score_final_full" bson:"score_final_full"`
	Grade         float64 `json:"grade" bson:"grade"`

	AllDateCount         int `json:"all_date_count" bson:"all_date_count"`
	CheckNameAttendCount int `json:"check_name_attend_count" bson:"check_name_attend_count"`
	CheckNameAbsentCount int `json:"check_name_absent_count" bson:"check_name_absent_count"`
	CheckNameLateCount   int `json:"check_name_late_count" bson:"check_name_late_count"`
}

type CourseSummaryRequest struct {
	CourseId string `json:"course_id"`
}
