package repository

import (
	"context"
	"school-notification-backend/db"
	"school-notification-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const courseSummaryCollection = "course-summary"

type CourseSummaryRepository interface {
	Insert(courseSummary *models.CourseSummary) (*mongo.InsertOneResult, error)
	Update(courseSummary *models.CourseSummary) (*mongo.UpdateResult, error)
	GetByFilter(filter interface{}) (courseSummary *models.CourseSummary, err error)
	GetAll() (courseSummaryList []*models.CourseSummary, err error)
	GetByFilterAll(filter interface{}) (courseSummaryList []*models.CourseSummary, err error)
}

type courseSummaryRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewCourseSummaryRepository(conn db.Connection) CourseSummaryRepository {
	return &courseSummaryRepository{c: conn.DB().Collection(courseSummaryCollection), ctx: context.TODO()}
}

func (c *courseSummaryRepository) Insert(courseSummary *models.CourseSummary) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(c.ctx, courseSummary)
}

func (c *courseSummaryRepository) Update(courseSummary *models.CourseSummary) (*mongo.UpdateResult, error) {
	return c.c.UpdateByID(c.ctx, courseSummary.Id, bson.M{"$set": courseSummary})
}

func (c *courseSummaryRepository) GetByFilter(filter interface{}) (courseSummary *models.CourseSummary, err error) {

	result := c.c.FindOne(c.ctx, filter)

	err = result.Decode(&courseSummary)
	if err != nil {
		return nil, err
	}

	return courseSummary, result.Err()
}

func (c *courseSummaryRepository) GetByFilterAll(filter interface{}) (courseSummaryList []*models.CourseSummary, err error) {

	cur, err := c.c.Find(c.ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		var b *models.CourseSummary
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		courseSummaryList = append(courseSummaryList, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(c.ctx)

	if len(courseSummaryList) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return courseSummaryList, nil
}

func (c *courseSummaryRepository) GetAll() (courseSummaryList []*models.CourseSummary, err error) {

	cur, err := c.c.Find(c.ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		var b *models.CourseSummary
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		courseSummaryList = append(courseSummaryList, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(c.ctx)

	if len(courseSummaryList) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return courseSummaryList, nil
}
