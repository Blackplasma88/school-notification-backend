package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// admin
type ProfileAdmin struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	UpdatedAt string             `json:"updated_at" bson:"updated_at"`
	ProfileId string             `json:"profile_id" bson:"profile_id"`
	Name      string             `json:"name" bson:"name"`
	Role      string             `json:"role" bson:"role"`
}

// teacher
type ProfileTeacher struct {
	Id                primitive.ObjectID  `json:"id" bson:"_id"`
	CreatedAt         string              `json:"created_at" bson:"created_at"`
	UpdatedAt         string              `json:"updated_at" bson:"updated_at"`
	ProfileId         string              `json:"profile_id" bson:"profile_id"`
	Name              string              `json:"name" bson:"name"`
	Role              string              `json:"role" bson:"role"`
	Category          string              `json:"category" bson:"category"`
	SubjectId         string              `json:"subject_id" bson:"subject_id"`
	ClassInCounseling string              `json:"class_in_counseling" bson:"class_in_counseling"`
	CourseTeachesList []CourseTeachesList `json:"course_teaches_list" bson:"course_teaches_list"`
	Slot              []Slot              `json:"slot" bson:"slot"`
}

type CourseTeachesList struct {
	CourseIdList []primitive.ObjectID `json:"course_id_list" bson:"course_id_list"`
	Year         string               `json:"year" bson:"year"`
	Term         string               `json:"term" bson:"term"`
}

// student
type ProfileStudent struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	UpdatedAt string             `json:"updated_at" bson:"updated_at"`
	ProfileId string             `json:"profile_id" bson:"profile_id"`
	Name      string             `json:"name" bson:"name"`
	Role      string             `json:"role" bson:"role"`
	ClassId   string             `json:"class_id" bson:"class_id"`
	ParentId  string             `json:"parent_id" bson:"parent_id"`
	GPA       float64            `json:"gpa" bson:"gpa"`
	AllCredit int                `json:"all_credit" bson:"all_credit"`
	TermScore []TermScore        `json:"term_score" bson:"term_score"`
}

type TermScore struct {
	Year       string       `json:"year" bson:"year"`
	Term       string       `json:"term" bson:"term"`
	GPA        float64      `json:"gpa" bson:"gpa"`
	TermCredit int          `json:"term_credit" bson:"term_credit"`
	CourseList []CourseList `json:"course_list" bson:"course_list"`
}

type CourseList struct {
	// CreatedAt string             `json:"created_at" bson:"created_at"`
	Id            primitive.ObjectID `json:"id" bson:"id"`
	Grade         float64            `json:"grade" bson:"grade"`
	ScoreWorkGet  float64            `json:"score_work_get" bson:"score_work_get"`
	ScoreWorkFull float64            `json:"score_work_full" bson:"score_work_full"`
	ScoreMidGet   float64            `json:"score_mid_get" bson:"score_mid_get"`
	ScoreMidFull  float64            `json:"score_mid_full" bson:"score_mid_full"`
	ScoreFinalGet float64            `json:"score_final_get" bson:"score_final_get"`
	ScoreFinaFull float64            `json:"score_final_full" bson:"score_final_full"`
	Credit        int                `json:"credit" bson:"credit"`

	AllDateCount         int `json:"all_date_count" bson:"all_date_count"`
	CheckNameAttendCount int `json:"check_name_attend_count" bson:"check_name_attend_count"`
	CheckNameAbsentCount int `json:"check_name_absent_count" bson:"check_name_absent_count"`
	CheckNameLateCount   int `json:"check_name_late_count" bson:"check_name_late_count"`
}

// request
type ProfileRequest struct {
	Event     string `json:"event"`
	Role      string `json:"role"`
	ProfileId string `json:"profile_id" bson:"profile_id"`
	Name      string `json:"name"`

	ProfileIdNew      string `json:"profile_id_new" bson:"profile_id_new"`
	AdvisorId         string `json:"advisor_id" bson:"advisor_id"`
	SubjectId         string `json:"subject_id"`
	ClassId           string `json:"class_id"`
	ClassInCounseling string `json:"class_in_counseling" bson:"class_in_counseling"`
	ParentId          string `json:"parent_id"`
	StudentId         string `json:"student_id"`
	CoursesId         string `json:"courses_id" bson:"courses_id"`
	Category          string `json:"category" bson:"category"`
}

type ProfileForChat struct {
	Id   primitive.ObjectID `json:"id" bson:"_id"`
	Name string             `json:"name" bson:"name"`
	Role string             `json:"role" bson:"role"`
}
