package mongodb

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/cjlapao/common-go/database"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient interface {
	Database(string) MongoDatabaseClient
	Connect() error
	StartSession() (mongoSession, error)
}

type MongoDatabaseClient interface {
	Collection(name string)
}

type MongoCollectionClient interface {
	Find(context.Context, interface{}) (mongoCursor, error)
	FindOne(context.Context, interface{}) mongoSingleResult
	InsertOne(context.Context, interface{}) (interface{}, error)
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	name string
	db   *mongo.Database
}

type mongoCollection struct {
	name string
	coll *mongo.Collection
}

type mongoSession struct {
	mongo.Session
}

type mongoCursor struct {
	cursor *mongo.Cursor
}

func (cursor mongoCursor) Decode(destination interface{}) error {
	ctx := context.Background()
	var destType = reflect.TypeOf(destination)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer type")
	}

	return cursor.cursor.All(ctx, destination)
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

func (cursor mongoSingleResult) Decode(destination interface{}) error {
	var destType = reflect.TypeOf(destination)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer type")
	}

	return cursor.sr.Decode(destination)
}

// MongoFactory MongoFactory Entity
type MongoFactory struct {
	Context         *execution_context.Context
	Client          *mongoClient
	Database        *mongoDatabase
	DatabaseContext *database.DatabaseContext
	Logger          *log.Logger
}

// NewFactory Creates a brand new factory for a specific connection string
// this will create and attach a mongo client that it will use for all connections
// returns a pointer to a MongoFactory object
func NewFactory(connectionString string) *MongoFactory {
	factory := MongoFactory{}
	factory.DatabaseContext = &database.DatabaseContext{
		ConnectionString: connectionString,
	}

	factory.Logger = log.Get()
	factory.Context = execution_context.Get()

	factory.GetClient()
	factory.Logger.Info("MongoDB Factory initiated successfully.")
	return &factory
}

func (f *MongoFactory) WithDatabase(databaseName string) *MongoFactory {
	f.GetDatabase(databaseName)
	f.Logger.Info("MongoDB Factory database %v initiated successfully.", databaseName)
	return f
}

// GetClient This will either return an already initiated client or the current
// active client in the factory, this will avoid having unclosed clients
// if you need a brand new client please use the NewFactory method to create a brand
// new factory.
// returns a mongoClient object
func (f *MongoFactory) GetClient() *mongoClient {

	if f.Client != nil {
		return f.Client
	}

	connectionString := f.DatabaseContext.ConnectionString
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

	f.Client = &mongoClient{cl: client}

	f.Logger.Debug("Client connection created successfully")
	return f.Client
}

// GetDatabase Get a database from the current cluster and sets it in the database context
// returns a mongoDatabase object
func (f *MongoFactory) GetDatabase(databaseName string) *mongoDatabase {
	if f.Client == nil {
		f.Client = f.GetClient()
	}

	database := f.Client.cl.Database(databaseName)
	if database == nil {
		f.Logger.Error("There was an error getting the database %v", databaseName)
		return nil
	}

	f.DatabaseContext.CurrentDatabaseName = databaseName
	f.Database = &mongoDatabase{db: database, name: databaseName}

	f.Logger.Debug("Database was retrieved successfully")
	return f.Database
}

// GetCollection Get a collection from the current database
// returns a mongoCollection object
func (f *MongoFactory) GetCollection(collectionName string) *mongoCollection {
	if f.Client == nil {
		f.Client = f.GetClient()
	}

	if f.Database == nil {
		f.Database = f.GetDatabase(f.DatabaseContext.CurrentDatabaseName)
	}

	collection := f.Database.db.Collection(collectionName)
	if collection == nil {
		f.Logger.Error("There was an error getting the collection %v", collectionName)
		return nil
	}

	f.DatabaseContext.CurrentCollection = collectionName
	f.Logger.Debug("Collection was retrieved successfully")
	return &mongoCollection{coll: collection, name: collectionName}
}

// StartSession Starts a session in the mongodb client
func (f *MongoFactory) StartSession() (mongo.Session, error) {
	session, err := f.Client.cl.StartSession()
	return &mongoSession{Session: session}, err
}
