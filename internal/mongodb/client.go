package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	Client              *mongo.Client
	timeoutMilliseconds int
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

	err = client.Ping(ctx, nil)
	if err != nil {
		return &MongoClient{}, err
	}

	return &MongoClient{
		Client:              client,
		timeoutMilliseconds: timeoutMilliseconds,
	}, nil
}

func (c *MongoClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.timeoutMilliseconds)*time.Millisecond)
	defer cancel()
	return c.Client.Disconnect(ctx)
}
