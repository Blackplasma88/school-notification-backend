package repository

import (
	"context"
	"school-notification-backend/db"
	"school-notification-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const messageCollection = "messages"

type MessageRepository interface {
	Insert(message *models.Message) (*mongo.InsertOneResult, error)
	Update(message *models.Message) (*mongo.UpdateResult, error)
	GetConversationAllByFilter(filter interface{}) (messages []*models.Message, err error)
}

type messageRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewMessageRepository(conn db.Connection) MessageRepository {
	return &messageRepository{c: conn.DB().Collection(messageCollection), ctx: context.TODO()}
}

func (m *messageRepository) Insert(message *models.Message) (*mongo.InsertOneResult, error) {
	return m.c.InsertOne(m.ctx, message)
}

func (m *messageRepository) Update(message *models.Message) (*mongo.UpdateResult, error) {
	return m.c.UpdateByID(m.ctx, message.Id, bson.M{"$set": message})
}

func (m *messageRepository) GetConversationAllByFilter(filter interface{}) (messages []*models.Message, err error) {

	cur, err := m.c.Find(m.ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(m.ctx) {
		var b *models.Message
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		messages = append(messages, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(m.ctx)

	if len(messages) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return messages, nil
}
