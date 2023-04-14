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

const usersCollection = "users"

type UsersRepository interface {
	InsertUser(user *models.User) (*mongo.InsertOneResult, error)
	Update(user *models.User) (*mongo.UpdateResult, error)
	GetById(id string) (user *models.User, err error)
	GetByUsername(username string) (user *models.User, err error)
	GetAll() (users []*models.User, err error)
	Delete(id string) (*mongo.DeleteResult, error)
}

type usersRepository struct {
	c   *mongo.Collection
	ctx context.Context
}

func NewUsersRepository(conn db.Connection) UsersRepository {
	return &usersRepository{c: conn.DB().Collection(usersCollection), ctx: context.TODO()}
}

func (u *usersRepository) InsertUser(user *models.User) (*mongo.InsertOneResult, error) {
	return u.c.InsertOne(u.ctx, user)
}

func (u *usersRepository) Update(user *models.User) (*mongo.UpdateResult, error) {
	return u.c.UpdateByID(u.ctx, user.Id, bson.M{"$set": user})
}

func (u *usersRepository) GetById(id string) (user *models.User, err error) {

	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = u.c.FindOne(u.ctx, bson.M{"_id": oID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *usersRepository) GetByUsername(username string) (user *models.User, err error) {

	err = u.c.FindOne(u.ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *usersRepository) GetAll() (users []*models.User, err error) {
	cur, err := u.c.Find(u.ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	for cur.Next(u.ctx) {
		var b *models.User
		err := cur.Decode(&b)
		if err != nil {
			return nil, err
		}

		users = append(users, b)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(u.ctx)

	if len(users) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return users, err
}

func (u *usersRepository) Delete(id string) (*mongo.DeleteResult, error) {
	if ok := primitive.IsValidObjectID(id); ok == false {
		return nil, util.ErrIdIsNotPrimitiveObjectID
	}

	oID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	result, err := u.c.DeleteOne(u.ctx, bson.M{"_id": oID})
	if err != nil {
		return nil, err
	}

	return result, nil
}
