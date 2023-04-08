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

const classCollection = "classes"

type ClassRepository interface {
	Insert(class *models.ClassData) (*mongo.InsertOneResult, error)
	Update(class *models.ClassData) (*mongo.UpdateResult, error)
	GetAll() (classes []*models.ClassData, err error)
	GetClassById(id string) (class *models.ClassData, err error)
	GetClassByFilter(filter interface{}) (class *models.ClassData, err error)
	GetCountOfClassYear(classYear string) (num int, err error)
}

type classRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewClassRepository(conn db.Connection) ClassRepository {
	return &classRepository{c: conn.DB().Collection(classCollection), ctx: context.TODO()}
}

func (c *classRepository) Insert(class *models.ClassData) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(c.ctx, class)
}

func (c *classRepository) Update(class *models.ClassData) (*mongo.UpdateResult, error) {
	return c.c.UpdateByID(c.ctx, class.Id, bson.M{"$set": class})
}

func (c *classRepository) GetClassById(id string) (class *models.ClassData, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := c.c.FindOne(c.ctx, bson.M{"_id": oID})

	err = result.Decode(&class)
	if err != nil {
		return nil, err
	}

	return class, result.Err()
}

func (c *classRepository) GetAll() (classes []*models.ClassData, err error) {

	cur, err := c.c.Find(c.ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(c.ctx) {
		var b *models.ClassData
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		classes = append(classes, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(c.ctx)

	if len(classes) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return classes, nil
}

func (c *classRepository) GetClassByFilter(filter interface{}) (class *models.ClassData, err error) {

	result := c.c.FindOne(c.ctx, filter)

	err = result.Decode(&class)
	if err != nil {
		return nil, err
	}

	return class, result.Err()
}

func (c *classRepository) GetCountOfClassYear(classYear string) (num int, err error) {

	cur, err := c.c.Find(c.ctx, bson.M{"class_year": classYear}, nil)
	if err != nil {
		return num, err
	}

	num = 0
	for cur.Next(c.ctx) {
		num++
	}

	if err := cur.Err(); err != nil {
		return num, err
	}

	cur.Close(c.ctx)

	return num, nil
}
