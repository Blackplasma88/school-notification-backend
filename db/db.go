package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connection interface {
	Close()
	DB() *mongo.Database
}

type conn struct {
	client *mongo.Client
}

func NewConnection() Connection {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("dbURL:", getURLLocal())
	clientOptions := options.Client().ApplyURI(getURLLocal())
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &conn{client: client}
}

func (c *conn) Close() {
	ctx := context.TODO()
	c.client.Disconnect(ctx)
}

func (c *conn) DB() *mongo.Database {
	return c.client.Database("SchoolManagement")
}

func getURLLocal() string {
	return fmt.Sprintf("mongodb://%s:%s", os.Getenv("APP_IP"), os.Getenv("DB_PORT"))
}
