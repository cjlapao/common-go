package mongodb

import (
	"context"
	"time"

	"github.com/cjlapao/common-go/database"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoFactory MongoFactory Entity
type MongoFactory struct {
	Client      *mongo.Client
	Database    *mongo.Database
	Context     *execution_context.Context
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
	factory.Logger.Info("MongoDB Factory initiated successfully.")
	return &factory
}

func (f *MongoFactory) WithDatabase(databaseName string) *MongoFactory {
	f.DatabaseCxt.DatabaseName = databaseName
	f.GetDatabase()
	f.Logger.Info("MongoDB Factory database %v initiated successfully.", databaseName)
	return f
}

// GetContext Gets the Execution context
func (f *MongoFactory) GetContext() *execution_context.Context {
	if f.Context != nil {
		return f.Context
	}

	f.Context = execution_context.Get()

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

func (f *MongoFactory) GetCollection(collectionName string) *mongo.Collection {
	if f.Client == nil {
		f.Client = f.GetClient()
	}
	if f.Database == nil {
		f.Database = f.GetDatabase()
	}

	collection := f.Database.Collection(collectionName)
	if collection == nil {
		f.Logger.Error("There was an error getting the collection " + collectionName)
		return nil
	}

	f.Logger.Debug("Database was retrieved successfully")
	return collection
}
