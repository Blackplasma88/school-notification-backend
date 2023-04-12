package repository

import (
	"context"
	"school-notification-backend/db"
	"school-notification-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const checkNameCollection = "check-name"

type CheckNameRepository interface {
	Insert(checkName *models.CheckName) (*mongo.InsertOneResult, error)
	Update(checkName *models.CheckName) (*mongo.UpdateResult, error)
	GetByFilter(filter interface{}) (checkName *models.CheckName, err error)
	GetByFilterAll(filter interface{}) (checkNameList []*models.CheckName, err error)
}

type checkNameRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewCheckNameRepository(conn db.Connection) CheckNameRepository {
	return &checkNameRepository{c: conn.DB().Collection(checkNameCollection), ctx: context.TODO()}
}

func (c *checkNameRepository) Insert(checkName *models.CheckName) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(c.ctx, checkName)
}

func (c *checkNameRepository) Update(checkName *models.CheckName) (*mongo.UpdateResult, error) {
	return c.c.UpdateByID(c.ctx, checkName.Id, bson.M{"$set": checkName})
}

func (c *checkNameRepository) GetByFilter(filter interface{}) (checkName *models.CheckName, err error) {

	result := c.c.FindOne(c.ctx, filter)

	err = result.Decode(&checkName)
	if err != nil {
		return nil, err
	}

	return checkName, result.Err()
}

func (c *checkNameRepository) GetByFilterAll(filter interface{}) (checkNameList []*models.CheckName, err error) {

	cur, err := c.c.Find(c.ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		var b *models.CheckName
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		checkNameList = append(checkNameList, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(c.ctx)

	if len(checkNameList) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return checkNameList, nil
}
