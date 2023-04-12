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

const courseCollection = "courses"

type CourseRepository interface {
	Insert(course *models.Course) (*mongo.InsertOneResult, error)
	Update(course *models.Course) (*mongo.UpdateResult, error)
	GetCourseById(id string) (course *models.Course, err error)
	GetCourseAllByFilter(filter interface{}) (courses []*models.Course, err error)
}

type courseRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewCoursesRepository(conn db.Connection) CourseRepository {
	return &courseRepository{c: conn.DB().Collection(courseCollection), ctx: context.TODO()}
}

func (c *courseRepository) Insert(course *models.Course) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(c.ctx, course)
}

func (c *courseRepository) Update(course *models.Course) (*mongo.UpdateResult, error) {
	return c.c.UpdateByID(c.ctx, course.Id, bson.M{"$set": course})
}

func (c *courseRepository) GetCourseById(id string) (course *models.Course, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := c.c.FindOne(c.ctx, bson.M{"_id": oID})

	err = result.Decode(&course)
	if err != nil {
		return nil, err
	}

	return course, result.Err()
}

func (c *courseRepository) GetCourseAllByFilter(filter interface{}) (courses []*models.Course, err error) {

	cur, err := c.c.Find(c.ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		var b *models.Course
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		courses = append(courses, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(c.ctx)

	if len(courses) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return courses, nil
}
