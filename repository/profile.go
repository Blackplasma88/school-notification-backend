package repository

import (
	"context"
	"school-notification-backend/db"
	"school-notification-backend/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const profileCollection = "profiles"

type ProfileRepository interface {
}

type profileRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewProfileRepository(conn db.Connection) ProfileRepository {
	return &profileRepository{c: conn.DB().Collection(profileCollection), ctx: context.TODO()}
}

func (p *profileRepository) Insert(profile interface{}) (*mongo.InsertOneResult, error) {
	return p.c.InsertOne(p.ctx, profile)
}

func (p *profileRepository) Update(id primitive.ObjectID, filter interface{}) (*mongo.UpdateResult, error) {
	return p.c.UpdateByID(p.ctx, id, bson.M{"$set": filter})
}

func (p *profileRepository) GetProfileByFilterForCheckExists(filter interface{}) (err error) {

	result := p.c.FindOne(p.ctx, filter)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

// func (p *profileRepository) GetAll(role string) (profiles []*models.Profile, err error) {
func (p *profileRepository) GetAll(role string) (profiles []interface{}, err error) {

	cur, err := p.c.Find(p.ctx, bson.M{"role": role}, nil)
	if err != nil {
		return nil, err
	}

	if role == "teacher" {
		for cur.Next(p.ctx) {
			var b *models.ProfileTeacher
			err := cur.Decode(&b)
			if err != nil {
				return nil, err
			}

			profiles = append(profiles, b)
		}
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(p.ctx)

	if len(profiles) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return profiles, nil
}
