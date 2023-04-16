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

const faceDetectionCollection = "face-detection"

type FaceDetectionRepository interface {
	Insert(faceDetectData *models.FaceDetectData) (*mongo.InsertOneResult, error)
	Update(faceDetectData *models.FaceDetectData) (*mongo.UpdateResult, error)
	GetById(id string) (faceDetectData *models.FaceDetectData, err error)
	GetAll() (faceDetectDataList []*models.FaceDetectData, err error)
	GetByFilter(filter interface{}) (faceDetectData *models.FaceDetectData, err error)
	GetByFilterAll(filter interface{}) (faceDetectDataList []*models.FaceDetectData, err error)
}

type faceDetectionRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewFaceDetectionRepository(conn db.Connection) FaceDetectionRepository {
	return &faceDetectionRepository{c: conn.DB().Collection(faceDetectionCollection), ctx: context.TODO()}
}

func (c *faceDetectionRepository) Insert(faceDetectData *models.FaceDetectData) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(c.ctx, faceDetectData)
}

func (c *faceDetectionRepository) Update(faceDetectData *models.FaceDetectData) (*mongo.UpdateResult, error) {
	return c.c.UpdateByID(c.ctx, faceDetectData.Id, bson.M{"$set": faceDetectData})
}

func (f *faceDetectionRepository) GetById(id string) (faceDetectData *models.FaceDetectData, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := f.c.FindOne(f.ctx, bson.M{"_id": oID})

	err = result.Decode(&faceDetectData)
	if err != nil {
		return nil, err
	}

	return faceDetectData, result.Err()
}

func (f *faceDetectionRepository) GetAll() (faceDetectDataList []*models.FaceDetectData, err error) {

	cur, err := f.c.Find(f.ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(f.ctx) {
		var b *models.FaceDetectData
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		faceDetectDataList = append(faceDetectDataList, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(f.ctx)

	if len(faceDetectDataList) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return faceDetectDataList, nil
}

func (f *faceDetectionRepository) GetByFilter(filter interface{}) (faceDetectData *models.FaceDetectData, err error) {

	result := f.c.FindOne(f.ctx, filter)

	err = result.Decode(&faceDetectData)
	if err != nil {
		return nil, err
	}

	return faceDetectData, result.Err()
}

func (f *faceDetectionRepository) GetByFilterAll(filter interface{}) (faceDetectDataList []*models.FaceDetectData, err error) {

	cur, err := f.c.Find(f.ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(f.ctx) {
		var b *models.FaceDetectData
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		faceDetectDataList = append(faceDetectDataList, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(f.ctx)

	if len(faceDetectDataList) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return faceDetectDataList, nil
}
