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

const informationCollection = "informations"

type InformationRepository interface {
	Insert(information *models.Information) (*mongo.InsertOneResult, error)
	Update(information *models.Information) (*mongo.UpdateResult, error)
	GetAll() (informations []*models.Information, err error)
	GetInformationById(id string) (information *models.Information, err error)
}

type informationRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewInformationRepository(conn db.Connection) InformationRepository {
	return &informationRepository{c: conn.DB().Collection(informationCollection), ctx: context.TODO()}
}

func (i *informationRepository) Insert(information *models.Information) (*mongo.InsertOneResult, error) {
	return i.c.InsertOne(i.ctx, information)
}

func (i *informationRepository) Update(information *models.Information) (*mongo.UpdateResult, error) {
	return i.c.UpdateByID(i.ctx, information.Id, bson.M{"$set": information})
}

func (i *informationRepository) GetInformationById(id string) (information *models.Information, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := i.c.FindOne(i.ctx, bson.M{"_id": oID})

	err = result.Decode(&information)
	if err != nil {
		return nil, err
	}

	return information, result.Err()
}

func (i *informationRepository) GetAll() (informations []*models.Information, err error) {

	cur, err := i.c.Find(i.ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(i.ctx) {
		var b *models.Information
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		informations = append(informations, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(i.ctx)

	if len(informations) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return informations, nil
}
