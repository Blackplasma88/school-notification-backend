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

const locationCollection = "locations"

type LocationRepository interface {
	Insert(location *models.Location) (*mongo.InsertOneResult, error)
	Update(location *models.Location) (*mongo.UpdateResult, error)
	GetAll() (locations []*models.Location, err error)
	GetLocationById(id string) (location *models.Location, err error)
	GetLocationByFilter(filter interface{}) (location *models.Location, err error)
}

type locationRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewLocationRepository(conn db.Connection) LocationRepository {
	return &locationRepository{c: conn.DB().Collection(locationCollection), ctx: context.TODO()}
}

func (l *locationRepository) Insert(location *models.Location) (*mongo.InsertOneResult, error) {
	return l.c.InsertOne(l.ctx, location)
}

func (l *locationRepository) Update(location *models.Location) (*mongo.UpdateResult, error) {
	return l.c.UpdateByID(l.ctx, location.Id, bson.M{"$set": location})
}

func (l *locationRepository) GetLocationById(id string) (location *models.Location, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := l.c.FindOne(l.ctx, bson.M{"_id": oID})

	err = result.Decode(&location)
	if err != nil {
		return nil, err
	}

	return location, result.Err()
}

func (l *locationRepository) GetAll() (locations []*models.Location, err error) {

	cur, err := l.c.Find(l.ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(l.ctx) {
		var b *models.Location
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		locations = append(locations, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(l.ctx)

	if len(locations) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return locations, nil
}

func (l *locationRepository) GetLocationByFilter(filter interface{}) (location *models.Location, err error) {

	result := l.c.FindOne(l.ctx, filter)

	err = result.Decode(&location)
	if err != nil {
		return nil, err
	}

	return location, result.Err()
}
