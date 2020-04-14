package infrastructure

import (
	"context"
	"fmt"
	"github.com/pedidosya/minesweeper-API/utils"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type IDataBaseProvider interface {
	ConnectDatabase() (interface{}, error)

	Insert(collectionName string, val interface{}) (interface{}, error)

	Upsert(collectionName string, id interface{}, val interface{}) error

	Update(collectionName string, id interface{}, val interface{}) (bool, error)

	ReplaceById(collectionName string, id interface{}, val interface{}) error

	GetById(collectionName string, id interface{}) (*mongo.SingleResult, error)

	Find(collectionName string, filter interface{}, ops *options.FindOptions) (*mongo.Cursor, error)

	Aggregate(collectionName string, pipeline interface{}, opts *options.AggregateOptions) (*mongo.Cursor, error)
}

type MongoDataBaseProvider struct {
	client *mongo.Database
}

var instanceDatabase *mongo.Database
var onceDatabase sync.Once

func (provider *MongoDataBaseProvider) ConnectDatabase() (interface{}, error) {

	onceDatabase.Do(func() {
		// Get config with viper
		dbHost := viper.GetString("database.host")
		dbName := viper.GetString("database.name")
		dbUser := viper.GetString("database.user")
		dbPassword := viper.GetString("database.password")

		uriDS := fmt.Sprintf("mongodb://%s:%s@%s/%s", dbUser, dbPassword, dbHost, dbName)
		utils.LogInfo("Connecting with MongoDB...")
		ctx, _ := context.WithTimeout(context.Background(), viper.GetDuration("database.timeout")*time.Second)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(uriDS))

		if err != nil {
			utils.LogInfo("couldn't reach database: %v", err)
			panic(err)
		}

		// Check the connection
		err = client.Ping(context.TODO(), nil)

		if err != nil {
			utils.LogError(fmt.Errorf("couldn't reach database: %v", err))
			panic(err)
		}

		utils.LogInfo("Connected to MongoDB!")
		instanceDatabase = client.Database(dbName)
	})
	return instanceDatabase, nil
}

func (provider *MongoDataBaseProvider) Insert(collectionName string, val interface{}) (interface{}, error) {
	collection := provider.client.Collection(collectionName)
	insertResult, err := collection.InsertOne(context.TODO(), val)
	if err != nil {
		return nil, fmt.Errorf("error to insert in collection: %s", collectionName)
	}
	utils.LogInfo("inserted a single document: %s, in collection: %s", insertResult.InsertedID, collectionName)
	return insertResult.InsertedID, nil
}

func (provider *MongoDataBaseProvider) Upsert(collectionName string, id interface{}, val interface{}) error {
	collection := provider.client.Collection(collectionName)
	filter := bson.M{"_id": id}
	updateResult, err := collection.ReplaceOne(context.TODO(), filter, val)
	if err != nil {
		return fmt.Errorf("error to replace in collection: %s, %v", collectionName, err)
	}
	if updateResult.ModifiedCount == 0 && updateResult.MatchedCount == 0 {
		insertResult, err := collection.InsertOne(context.TODO(), val)
		if err != nil {
			return fmt.Errorf("error to insert in collection: %s, %v", collectionName, err)
		}
		utils.LogInfo("inserted a single document: %s, in collection: %s", insertResult.InsertedID, collectionName)
	}
	return nil
}

func (provider *MongoDataBaseProvider) Update(collectionName string, id interface{}, val interface{}) (bool, error) {
	collection := provider.client.Collection(collectionName)
	filter := bson.M{"_id": id}
	updateResult, err := collection.UpdateOne(context.TODO(), filter, val)
	if err != nil {
		return false, fmt.Errorf("error to replace in collection: %s, %v", collectionName, err)
	}
	if updateResult.ModifiedCount == 0 && updateResult.MatchedCount == 0 {
		return false, nil
	}
	utils.LogInfo("updated a single document: %v, in collection: %s", id, collectionName)
	return true, nil
}

func (provider *MongoDataBaseProvider) ReplaceById(collectionName string, id interface{}, val interface{}) error {
	collection := provider.client.Collection(collectionName)

	// Set filters
	filter := bson.M{"_id": id}

	updateResult, err := collection.ReplaceOne(context.TODO(), filter, val)
	if err != nil {
		return fmt.Errorf("not can create object_id %s: %v", collectionName, err)
	}
	utils.LogInfo("modified %d document: ", updateResult.ModifiedCount)
	return nil
}

func (provider *MongoDataBaseProvider) GetById(collectionName string, id interface{}) (*mongo.SingleResult, error) {
	collection := provider.client.Collection(collectionName)

	// Set filters
	filter := bson.M{"_id": id}

	// Find document
	sr := collection.FindOne(context.TODO(), filter)
	return sr, nil

}

func (provider *MongoDataBaseProvider) Find(collectionName string, filter interface{}, options *options.FindOptions) (*mongo.Cursor, error) {
	collection := provider.client.Collection(collectionName)
	cur, err := collection.Find(context.TODO(), filter, options)
	if err != nil {
		return nil, fmt.Errorf("error to find in collection: %s", collectionName)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error from database: %v", err)
	}
	return cur, nil
}

func (provider *MongoDataBaseProvider) Aggregate(collectionName string, pipeline interface{}, opts *options.AggregateOptions) (*mongo.Cursor, error) {
	collection := provider.client.Collection(collectionName)
	cur, err := collection.Aggregate(context.TODO(), pipeline, opts)
	if err != nil {
		return nil, fmt.Errorf("error to find in collection: %s", collectionName)
	}

	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("error from database: %v", err)
	}
	return cur, nil
}

func NewDataBaseClient() IDataBaseProvider {
	provider := &MongoDataBaseProvider{}
	m, err := provider.ConnectDatabase()
	if err != nil {
		panic(err)
	}
	provider.client = m.(*mongo.Database)
	return provider
}
