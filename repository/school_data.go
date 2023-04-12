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

const schoolDataCollection = "school-data"

type SchoolDataRepository interface {
	GetAll() (schoolDataList []*models.SchoolData, err error)
	Insert(schoolData interface{}) (*mongo.InsertOneResult, error)
	GetById(id string) (schoolData *models.SchoolData, err error)
	Update(schoolData *models.SchoolData) (*mongo.UpdateResult, error)
	GetByFilter(filter interface{}) (schoolData *models.SchoolData, err error)
	GetByFilterAll(filter interface{}) (schoolDataList []*models.SchoolData, err error)
}

type schoolDataRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewSchoolDataRepository(conn db.Connection) SchoolDataRepository {
	return &schoolDataRepository{c: conn.DB().Collection(schoolDataCollection), ctx: context.TODO()}
}

func (s *schoolDataRepository) Insert(schoolData interface{}) (*mongo.InsertOneResult, error) {
	return s.c.InsertOne(s.ctx, schoolData)
}

func (s *schoolDataRepository) Update(schoolData *models.SchoolData) (*mongo.UpdateResult, error) {
	return s.c.UpdateByID(s.ctx, schoolData.Id, bson.M{"$set": schoolData})
}

func (s *schoolDataRepository) GetById(id string) (schoolData *models.SchoolData, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := s.c.FindOne(s.ctx, bson.M{"_id": oID})

	err = result.Decode(&schoolData)
	if err != nil {
		return nil, err
	}

	return schoolData, result.Err()
}

func (s *schoolDataRepository) GetAll() (schoolDataList []*models.SchoolData, err error) {

	cur, err := s.c.Find(s.ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(s.ctx) {
		var b *models.SchoolData
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		schoolDataList = append(schoolDataList, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(s.ctx)

	if len(schoolDataList) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return schoolDataList, nil
}

func (s *schoolDataRepository) GetByFilter(filter interface{}) (schoolData *models.SchoolData, err error) {

	result := s.c.FindOne(s.ctx, filter)

	err = result.Decode(&schoolData)
	if err != nil {
		return nil, err
	}

	return schoolData, result.Err()
}

func (s *schoolDataRepository) GetByFilterAll(filter interface{}) (schoolDataList []*models.SchoolData, err error) {

	cur, err := s.c.Find(s.ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(s.ctx) {
		var b *models.SchoolData
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		schoolDataList = append(schoolDataList, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(s.ctx)

	if len(schoolDataList) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return schoolDataList, nil
}
