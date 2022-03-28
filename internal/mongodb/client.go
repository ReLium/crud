package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client *mongo.Client
	ctx    context.Context
}

func NewClient(host string, timeoutMilliseconds int) (*MongoClient, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(host))
	if err != nil {
		return &MongoClient{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutMilliseconds)*time.Millisecond)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return &MongoClient{}, err
	}

	client.Ping(context.TODO(), nil)
	if err != nil {
		return &MongoClient{}, err
	}

	return &MongoClient{
		Client: client,
		ctx:    ctx,
	}, nil
}

func (c *MongoClient) Disconnect() error {
	return c.Client.Disconnect(c.ctx)
}
