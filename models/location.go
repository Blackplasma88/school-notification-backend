package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt    string             `json:"created_at" bson:"created_at"`
	UpdatedAt    string             `json:"updated_at" bson:"updated_at"`
	LocationId   string             `json:"location_id" bson:"location_id"`
	BuildingName string             `json:"building_name" bson:"building_name"`
	Floor        string             `json:"floor" bson:"floor"`
	Room         string             `json:"room" bson:"room"`
	Status       bool               `json:"status" bson:"status"`
	Slot         []Slot             `json:"slot" bson:"slot"`
}

type Slot struct {
	Day      string     `json:"day" bson:"day"`
	TimeSlot []TimeSlot `json:"time_slot" bson:"time_slot"`
}

type TimeSlot struct {
	Time      string `json:"time" bson:"time"`
	Status    bool   `json:"status" bson:"status"`
	CoursesId string `json:"courses_id" bson:"courses_id"`
}

type LocationRequest struct {
	// Event        string `json:"event" bson:"event"`
	LocationId   string `json:"location_id" bson:"location_id"`
	BuildingName string `json:"building_name" bson:"building_name"`
	Floor        string `json:"floor" bson:"floor"`
	Room         string `json:"room" bson:"room"`
	Status       *bool  `json:"status" bson:"status"`
	// Day          string `json:"day" bson:"day"`
	// TimeStart    string `json:"time_start" bson:"time_start"`
	// TimeEnd      string `json:"time_end" bson:"time_end"`
	// CoursesId    string `json:"courses_id" bson:"courses_id"`
}
