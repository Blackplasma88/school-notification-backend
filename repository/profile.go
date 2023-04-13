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
	Insert(profile interface{}) (*mongo.InsertOneResult, error)
	Update(id primitive.ObjectID, filter interface{}) (*mongo.UpdateResult, error)
	GetProfileByFilterForCheckExists(filter interface{}) (err error)
	GetProfileById(filter interface{}, role string) (profile interface{}, err error)
	GetAll(role string) (profiles []interface{}, err error)
	GetProfileByFilterAll(filter interface{}, role string) (profiles []interface{}, err error)
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
	} else if role == "student" {
		for cur.Next(p.ctx) {
			var b *models.ProfileStudent
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

func (p *profileRepository) GetProfileById(filter interface{}, role string) (profile interface{}, err error) {

	result := p.c.FindOne(p.ctx, filter)

	if role == "teacher" {
		p := models.ProfileTeacher{}
		err = result.Decode(&p)
		if err != nil {
			return nil, err
		}
		profile = p
	} else if role == "student" {
		p := models.ProfileStudent{}
		err = result.Decode(&p)
		if err != nil {
			return nil, err
		}
		profile = p
	}

	return profile, result.Err()
}

func (p *profileRepository) GetProfileByFilterAll(filter interface{}, role string) (profiles []interface{}, err error) {

	cur, err := p.c.Find(p.ctx, filter, nil)
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
	} else if role == "student" {
		for cur.Next(p.ctx) {
			var b *models.ProfileStudent
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
