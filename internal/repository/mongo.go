package repository

import (
	"context"
	"time"

	"github.com/ReLium/crud/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DBName = "DB"
const CollectionName = "cats"

type MongoDBRepo struct {
	client     *mongodb.MongoClient
	collection *mongo.Collection
}

type mongocat struct {
	Name       string `bson:"name"`
	Gender     string `bson:"gender"`
	Color      string `bson:"color"`
	Vaccinated bool   `bson:"vaccinated"`
}

func NewMongoDBRepo(client *mongodb.MongoClient) *MongoDBRepo {
	return &MongoDBRepo{
		client:     client,
		collection: InitDB(client),
	}
}

func (r *MongoDBRepo) Get(name string) (*Cat, error) {
	var mongocat mongocat
	err := r.collection.FindOne(context.TODO(), bson.M{"name": name}).Decode(&mongocat)
	if err != nil {
		return &Cat{}, err
	}
	repoCat := Cat(mongocat)
	return &repoCat, nil
}

func (r *MongoDBRepo) Delete(name string) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"name": name})
	return err
}

func (r *MongoDBRepo) Insert(cat *Cat) error {
	_, err := r.collection.InsertOne(context.TODO(), mongocat(*cat))
	return err
}

func (r *MongoDBRepo) Update(cat *CatUpdate) error {

	update := bson.M{}
	if cat.Color != "" {
		update["color"] = cat.Color
	}
	if cat.Gender != "" {
		update["gender"] = cat.Gender
	}
	if cat.Vaccinated != nil {
		update["vaccinated"] = *cat.Vaccinated
	}

	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"name": cat.Name}, bson.M{"$set": update})
	return err
}

func (r *MongoDBRepo) Find(query *Query) ([]*Cat, error) {

	var result []*Cat

	mongoQuery := bson.D{}
	if query.Color != "" {
		mongoQuery = append(mongoQuery, bson.E{Key: "color", Value: query.Color})
	}
	if query.Gender != "" {
		mongoQuery = append(mongoQuery, bson.E{Key: "gender", Value: query.Gender})
	}
	if query.Vaccinated != nil {
		mongoQuery = append(mongoQuery, bson.E{Key: "vaccinated", Value: *query.Vaccinated})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := r.collection.Find(ctx, mongoQuery)
	if err != nil {
		return result, err
	}
	defer cur.Close(ctx)
	for cur.Next(ctx) {
		var mongoCat mongocat
		err := cur.Decode(&mongoCat)
		if err != nil {
			return result, err
		}
		cat := Cat(mongoCat)
		result = append(result, &cat)
	}
	return result, cur.Err()
}

func (r *MongoDBRepo) Destroy() error {
	return r.client.Disconnect()
}

func InitDB(client *mongodb.MongoClient) *mongo.Collection {
	col := client.Client.Database(DBName).Collection(CollectionName)
	mod := mongo.IndexModel{
		Keys: bson.M{
			"name": 1,
		}, Options: options.Index().SetUnique(true),
	}
	col.Indexes().CreateOne(context.TODO(), mod)
	return col
}
