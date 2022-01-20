package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/cjlapao/common-go/guard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository interface {
	Filter(filter interface{}) []*interface{}
	Find(fieldName string, value string) []*interface{}
	FindOne(fieldName string, value string) *mongo.SingleResult
	InsertOne(element interface{}) *mongo.InsertOneResult
	InsertMany(elements []interface{}) *mongo.InsertManyResult
	UpsertOne(model mongo.UpdateOneModel) *mongo.UpdateResult
	UpsertMany(filter interface{}, elements []interface{}) *mongo.UpdateResult
	DeleteOne(model mongo.DeleteOneModel) *mongo.DeleteResult
}

type MongoDefaultRepository struct {
	factory    *MongoFactory
	Database   *mongoDatabase
	Collection *mongoCollection
}

// NewRepository Creates a new repository for a specific collection, this will allow you to perform
// queries and aggregations in the collection
// Returns an implemented interface MongoRepository
func (mongoFactory *MongoFactory) NewRepository(collection string) MongoRepository {
	defaultRepo := MongoDefaultRepository{
		factory: mongoFactory,
	}

	defaultRepo.Database = mongoFactory.GetDatabase(mongoFactory.Database.name)
	defaultRepo.Collection = mongoFactory.GetCollection(collection)

	return &defaultRepo
}

// NewDatabaseRepository Creates a new repository for a specific collection in a database,
// this will allow you to perform queries and aggregations in the collection
// Returns an implemented interface MongoRepository
func (mongoFactory *MongoFactory) NewDatabaseRepository(database string, collection string) MongoRepository {
	defaultRepo := MongoDefaultRepository{}

	defaultRepo.Database = mongoFactory.GetDatabase(database)
	defaultRepo.Collection = mongoFactory.GetCollection(collection)

	return &defaultRepo
}

func (r *MongoDefaultRepository) Filter(filter interface{}) []*interface{} {
	logger.Info("Session %v", fmt.Sprintf("%v", r.factory.Client.cl.NumberSessionsInProgress()))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var filterToApply interface{}

	if stringFilter, ok := filter.(string); ok {
		if filter == "" {
			filter = bson.D{{}}
		} else {
			processedFilter, err := NewFilterParser(stringFilter).Parse()
			if err != nil {
				filterToApply = bson.D{{}}
			} else {
				filterToApply = processedFilter
			}
		}
	} else {
		filterToApply = filter
	}

	defer cancel()

	cur, err := r.Collection.coll.Find(ctx, filterToApply)
	if err != nil {
		logger.LogError(err)
		return nil
	}
	var elements []*interface{}
	for cur.Next(ctx) {
		var element interface{}
		err := cur.Decode(&element)
		if err != nil {
			logger.LogError(err)
			return nil
		}
		elements = append(elements, &element)
	}

	return elements
}

func (r *MongoDefaultRepository) Find(fieldName string, value string) []*interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	filter := bson.D{
		{
			Key:   fieldName,
			Value: value,
		},
	}

	cur, err := r.Collection.coll.Find(ctx, filter)
	if err != nil {
		logger.LogError(err)
		return nil
	}
	var elements []*interface{}
	for cur.Next(ctx) {
		var element *interface{}
		err := cur.Decode(&element)
		if err != nil {
			logger.LogError(err)
			return nil
		}
		elements = append(elements, element)
	}

	return elements
}

func (r *MongoDefaultRepository) FindOne(fieldName string, value string) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	filter := bson.D{
		{
			Key:   fieldName,
			Value: value,
		},
	}

	cur := r.Collection.coll.FindOne(ctx, filter)

	return cur
}

func (r *MongoDefaultRepository) InsertOne(element interface{}) *mongo.InsertOneResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := r.Collection.coll.InsertOne(ctx, element)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return insertResult
}

func (r *MongoDefaultRepository) InsertMany(elements []interface{}) *mongo.InsertManyResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := r.Collection.coll.InsertMany(ctx, elements)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return insertResult
}

func (r *MongoDefaultRepository) UpsertOne(model mongo.UpdateOneModel) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Update().SetUpsert(true)
	updateOneResult, err := r.Collection.coll.UpdateOne(ctx, model.Filter, model.Update, options)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return updateOneResult
}

func (r *MongoDefaultRepository) UpsertMany(filter interface{}, elements []interface{}) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if guard.IsNill(filter) {
		filter = bson.D{{}}
	}

	options := options.Update().SetUpsert(true)

	updateOneResult, err := r.Collection.coll.UpdateMany(ctx, filter, elements, options)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return updateOneResult
}

func (r *MongoDefaultRepository) DeleteOne(model mongo.DeleteOneModel) *mongo.DeleteResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteOptions := options.Delete()
	deleteOptions.Collation = model.Collation
	deleteOptions.Hint = model.Hint

	deleteOneResult, err := r.Collection.coll.DeleteOne(ctx, model.Filter, deleteOptions)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return deleteOneResult
}
