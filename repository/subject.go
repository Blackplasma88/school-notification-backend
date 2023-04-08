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

const subjectCollection = "subjects"

type SubjectRepository interface {
	GetAll() (subjects []*models.Subject, err error)
	Insert(subject *models.Subject) (*mongo.InsertOneResult, error)
	GetSubjectById(id string) (subject *models.Subject, err error)
	GetSubjectByFilter(filter interface{}) (subject *models.Subject, err error)
	Update(subject *models.Subject) (*mongo.UpdateResult, error)
}

type subjectRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewSubjectRepository(conn db.Connection) SubjectRepository {
	return &subjectRepository{c: conn.DB().Collection(subjectCollection), ctx: context.TODO()}
}

func (s *subjectRepository) Insert(subject *models.Subject) (*mongo.InsertOneResult, error) {
	return s.c.InsertOne(s.ctx, subject)
}

func (s *subjectRepository) Update(subject *models.Subject) (*mongo.UpdateResult, error) {
	return s.c.UpdateByID(s.ctx, subject.Id, bson.M{"$set": subject})
}

func (s *subjectRepository) GetSubjectByFilter(filter interface{}) (subject *models.Subject, err error) {

	result := s.c.FindOne(s.ctx, filter)

	err = result.Decode(&subject)
	if err != nil {
		return nil, err
	}

	return subject, result.Err()
}

func (s *subjectRepository) GetSubjectById(id string) (subject *models.Subject, err error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result := s.c.FindOne(s.ctx, bson.M{"_id": oID})

	err = result.Decode(&subject)
	if err != nil {
		return nil, err
	}

	return subject, result.Err()
}

func (s *subjectRepository) GetAll() (subjects []*models.Subject, err error) {

	cur, err := s.c.Find(s.ctx, bson.M{}, nil)
	if err != nil {
		return nil, err
	}

	for cur.Next(s.ctx) {
		var b *models.Subject
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		subjects = append(subjects, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(s.ctx)

	if len(subjects) == 0 {
		return nil, mongo.ErrNoDocuments
	}

	return subjects, nil
}
