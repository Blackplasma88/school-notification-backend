package repository

import (
	"context"
	"school-notification-backend/db"
	"school-notification-backend/models"
	"school-notification-backend/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const conversationCollection = "conversations"

type ConversationRepository interface {
	Insert(conversation *models.Conversation) (*mongo.InsertOneResult, error)
	Update(conversation *models.Conversation) (*mongo.UpdateResult, error)
	GetConversationAllByFilter(filter interface{}) (conversations []*models.Conversation, err error)
	GetConversationById(id string) (conversation *models.Conversation, err error)
	GetByFilter(filter interface{}) (conversation *models.Conversation, err error)
}

type conversationRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewConversationRepository(conn db.Connection) ConversationRepository {
	return &conversationRepository{c: conn.DB().Collection(conversationCollection), ctx: context.TODO()}
}

func (c *conversationRepository) Insert(conversation *models.Conversation) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(c.ctx, conversation)
}

func (c *conversationRepository) Update(conversation *models.Conversation) (*mongo.UpdateResult, error) {
	return c.c.UpdateByID(c.ctx, conversation.Id, bson.M{"$set": conversation})
}

func (c *conversationRepository) GetConversationAllByFilter(filter interface{}) (conversations []*models.Conversation, err error) {

	cur, err := c.c.Find(c.ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		var b *models.Conversation
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		conversations = append(conversations, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(c.ctx)

	if len(conversations) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return conversations, nil
}

func (c *conversationRepository) GetByFilter(filter interface{}) (conversation *models.Conversation, err error) {

	result := c.c.FindOne(c.ctx, filter)

	err = result.Decode(&conversation)
	if err != nil {
		return nil, err
	}

	return conversation, result.Err()
}

func (c *conversationRepository) GetConversationById(id string) (conversation *models.Conversation, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := c.c.FindOne(c.ctx, bson.M{"_id": oID})

	err = result.Decode(&conversation)
	if err != nil {
		return nil, err
	}

	return conversation, result.Err()
}
