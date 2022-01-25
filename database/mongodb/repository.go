package mongodb

import (
	"context"
	"time"

	"github.com/cjlapao/common-go/guard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository interface {
	Pipeline() *PipelineBuilder
	Find(filter interface{}) (*mongoCursor, error)
	FindBy(fieldName string, value interface{}) (*mongoCursor, error)
	FindOne(fieldName string, value string) *mongoSingleResult
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

func (repository *MongoDefaultRepository) Pipeline() *PipelineBuilder {
	return NewPipelineBuilder(repository.Collection)
}

func (repository *MongoDefaultRepository) Find(filter interface{}) (*mongoCursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var filterToApply interface{}

	if stringFilter, ok := filter.(string); ok {
		if filter == "" {
			filter = bson.D{{}}
		} else {
			processedFilter, err := NewFilterParser(stringFilter).Parse()
			if err != nil {
				logger.Error("There was an error applying the filter, %v", err.Error())
				filterToApply = bson.D{{}}
			} else {
				filterToApply = processedFilter
			}
		}
	} else {
		filterToApply = filter
	}

	defer cancel()

	cur, err := repository.Collection.coll.Find(ctx, filterToApply)

	return &mongoCursor{cursor: cur}, err
}

func (r *MongoDefaultRepository) FindBy(fieldName string, value interface{}) (*mongoCursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var filter interface{}

	if literalValue, ok := value.(string); ok {
		if value == "" {
			filter = bson.D{{}}
		} else {
			processedFilter, err := NewFilterParser(fieldName + " " + literalValue).Parse()
			if err != nil {
				filter = bson.D{
					{
						Key:   fieldName,
						Value: value,
					},
				}
			} else {
				filter = processedFilter
			}
		}
	} else {
		filter = bson.D{
			{
				Key:   fieldName,
				Value: value,
			},
		}
	}

	defer cancel()

	cur, err := r.Collection.coll.Find(ctx, filter)

	return &mongoCursor{cursor: cur}, err
}

func (r *MongoDefaultRepository) FindOne(fieldName string, value string) *mongoSingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	filter := bson.D{
		{
			Key:   fieldName,
			Value: value,
		},
	}

	result := r.Collection.coll.FindOne(ctx, filter)

	return &mongoSingleResult{sr: result}
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
