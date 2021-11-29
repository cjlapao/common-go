package database

import (
	"context"
	"time"

	"github.com/cjlapao/common-go/database"
	"github.com/cjlapao/common-go/executionctx"
	"github.com/cjlapao/common-go/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoFactory MongoFactory Entity
type MongoFactory struct {
	Client      *mongo.Client
	Database    *mongo.Database
	Context     *executionctx.Context
	DatabaseCxt *database.DatabaseCtx
	Logger      *log.Logger
}

// NewFactory Instantiate a new MongoDb Factory
func NewFactory(connectionString string) *MongoFactory {
	factory := MongoFactory{}
	factory.DatabaseCxt = &database.DatabaseCtx{
		ConnectionString: connectionString,
	}
	factory.Logger = log.Get()
	factory.GetClient()
	factory.GetContext()
	return &factory
}

func (f *MongoFactory) WithDatabase(databaseName string) *MongoFactory {
	f.DatabaseCxt.DatabaseName = databaseName
	f.GetDatabase()
	return f
}

// GetContext Gets the Execution context
func (f *MongoFactory) GetContext() *executionctx.Context {
	if f.Context != nil {
		return f.Context
	}

	f.Context = executionctx.GetContext()

	return f.Context
}

// GetClient returns mongodb client
func (f *MongoFactory) GetClient() *mongo.Client {

	if f.Client != nil {
		return f.Client
	}

	connectionString := f.DatabaseCxt.ConnectionString
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		f.Logger.LogError(err)
		return nil
	}

	// Check the connection
	err = client.Ping(ctx, nil)

	if err != nil {
		f.Logger.LogError(err)
		return nil
	}

	f.Client = client

	f.Logger.Debug("Client connection created successfully")
	return client
}

// GetDatabase returns MongoDb database
func (f *MongoFactory) GetDatabase() *mongo.Database {
	if f.Client == nil {
		f.Client = f.GetClient()
	}

	database := f.Client.Database(f.DatabaseCxt.DatabaseName)

	if database == nil {
		f.Logger.Error("There was an error getting the database " + f.DatabaseCxt.DatabaseName)
		return nil
	}

	f.Database = database

	f.Logger.Debug("Database was retrieved successfully")
	return database
}

// Find Finds documents in the database
func (f *MongoFactory) Find(collectionName string, filter bson.D) []*bson.M {
	if f.Database == nil {
		f.GetDatabase()
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	collection := f.Database.Collection(collectionName)

	if collection == nil {
		return nil
	}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		f.Logger.LogError(err)
		return nil
	}
	var elements []*bson.M
	for cur.Next(ctx) {
		var element bson.M
		err := cur.Decode(&element)
		if err != nil {
			f.Logger.LogError(err)
			return nil
		}
		elements = append(elements, &element)
	}

	return elements
}

// InsertOne Inserts one document into the selected collection
func (f *MongoFactory) InsertOne(collectionName string, element interface{}) *mongo.InsertOneResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := f.Database.Collection(collectionName)

	insertResult, err := collection.InsertOne(ctx, element)

	if err != nil {
		f.Logger.LogError(err)
		return nil
	}

	return insertResult
}

// InsertMany Inserts one document into the selected collection
func (f *MongoFactory) InsertMany(collectionName string, elements []interface{}) *mongo.InsertManyResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := f.Database.Collection(collectionName)

	insertResult, err := collection.InsertMany(ctx, elements)

	if err != nil {
		f.Logger.LogError(err)
		return nil
	}

	return insertResult
}
