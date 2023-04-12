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

const scoreCollection = "scores"

type ScoreRepository interface {
	// GetAll() (scores []*models.Score, err error)
	Insert(score *models.Score) (*mongo.InsertOneResult, error)
	// GetScoreById(id string) (score *models.Score, err error)
	GetScoreByFilter(filter interface{}) (score *models.Score, err error)
	GetByFilterAll(filter interface{}) (scores []*models.Score, err error)
	Update(score *models.Score) (*mongo.UpdateResult, error)
}

type scoreRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewScoreRepository(conn db.Connection) ScoreRepository {
	return &scoreRepository{c: conn.DB().Collection(scoreCollection), ctx: context.TODO()}
}

func (s *scoreRepository) Insert(score *models.Score) (*mongo.InsertOneResult, error) {
	return s.c.InsertOne(s.ctx, score)
}

func (s *scoreRepository) Update(score *models.Score) (*mongo.UpdateResult, error) {
	return s.c.UpdateByID(s.ctx, score.Id, bson.M{"$set": score})
}

func (s *scoreRepository) GetScoreByFilter(filter interface{}) (score *models.Score, err error) {

	result := s.c.FindOne(s.ctx, filter)

	err = result.Decode(&score)
	if err != nil {
		return nil, err
	}

	return score, result.Err()
}

func (s *scoreRepository) GetByFilterAll(filter interface{}) (scores []*models.Score, err error) {

	cur, err := s.c.Find(s.ctx, filter, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(s.ctx) {
		var b *models.Score
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		scores = append(scores, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(s.ctx)

	if len(scores) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return scores, nil
}

func (s *scoreRepository) GetScoreById(id string) (score *models.Score, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := s.c.FindOne(s.ctx, bson.M{"_id": oID})

	err = result.Decode(&score)
	if err != nil {
		return nil, err
	}

	return score, result.Err()
}

func (s *scoreRepository) GetAll() (scores []*models.Score, err error) {

	cur, err := s.c.Find(s.ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(s.ctx) {
		var b *models.Score
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		scores = append(scores, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(s.ctx)

	if len(scores) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return scores, nil
}
