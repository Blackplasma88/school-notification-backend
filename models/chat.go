package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Conversation struct {
	Id        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt string             `json:"created_at" bson:"created_at"`
	UpdatedAt string             `json:"updated_at" bson:"updated_at"`
	Members   []string           `json:"members" bson:"members"`
}

type Message struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt      string             `json:"created_at" bson:"created_at"`
	UpdatedAt      string             `json:"updated_at" bson:"updated_at"`
	ConversationId string             `json:"conversation_id" bson:"conversation_id"`
	Sender         string             `json:"sender" bson:"sender"`
	Text           string             `json:"text" bson:"text"`
}

type ConversationRequest struct {
	SenderId   string `json:"sender_id"`
	ReceiverId string `json:"receiver_id"`
}

type MessageRequest struct {
	ConversationId string `json:"conversation_id" `
	SenderId       string `json:"sender_id"`
	Text           string `json:"text"`
}
