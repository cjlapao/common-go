package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository interface {
	OData() *ODataParser
	Pipeline() *PipelineBuilder
	Find(filter interface{}) (*mongoCursor, error)
	FindBy(filter string) (*mongoCursor, error)
	FindFieldBy(fieldName string, operation filterOperation, value interface{}) (*mongoCursor, error)

	FindOne(fieldName string, value string) *mongoSingleResult
	InsertOne(element interface{}) *mongo.InsertOneResult
	InsertMany(elements []interface{}) *mongo.InsertManyResult
	UpsertOne(model *MongoUpdateOneModel) *mongo.UpdateResult
	UpdateMany(models ...*MongoUpdateOneModel) (*mongo.BulkWriteResult, error)
	UpsertMany(models ...*MongoUpdateOneModel) (*mongo.BulkWriteResult, error)
	DeleteOne(model *MongoDeleteOneModel) *mongo.DeleteResult
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

// Pipeline Creates an empty pipeline for querying mongodb
func (repository *MongoDefaultRepository) Pipeline() *PipelineBuilder {
	return NewEmptyPipeline(repository.Collection)
}

// OData Creates an OData parser to return data
func (repository *MongoDefaultRepository) OData() *ODataParser {
	return EmptyODataParser(repository.Collection)
}

// Find finds records with a filter and returns a cursor to iterate trough them
//
// Example:
//		repository.Find(bson.M{"userId", "someId"})
//
// You can also use a string query similar to odata
//
// for example
//		repository.Find("userId eq 'someId'")
//
// or a function like startswith
// 		repository.Find("startswith(userId, 'someId'")
func (repository *MongoDefaultRepository) Find(filter interface{}) (*mongoCursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filterToApply interface{}

	if stringFilter, ok := filter.(string); ok {
		if filter == "" {
			filter = bson.D{}
		} else {
			processedFilter, err := NewFilterParser(stringFilter).Parse()
			if err != nil {
				logger.Error("There was an error applying the filter, %v", err.Error())
				return nil, err
			} else {
				filterToApply = processedFilter
			}
		}
	} else {
		filterToApply = filter
	}

	cur, err := repository.Collection.coll.Find(ctx, filterToApply)

	return &mongoCursor{cursor: cur}, err
}

// Find finds records with a filter and returns a cursor to iterate trough them
// You can also use a string query similar to odata, for example
//		repository.Find("userId eq 'someId'")
//
// or a function like startswith
// 		repository.Find("startswith(userId, 'someId'")
func (r *MongoDefaultRepository) FindBy(filter string) (*mongoCursor, error) {
	return r.Find(filter)
}

func (r *MongoDefaultRepository) FindFieldBy(fieldName string, operation filterOperation, value interface{}) (*mongoCursor, error) {
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

func (r *MongoDefaultRepository) UpdateOne(model *MongoUpdateOneModel) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	options := options.Update().SetUpsert(*model.model.Upsert)
	updateOneResult, err := r.Collection.coll.UpdateOne(ctx, model.Filter, model.Update, options)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return updateOneResult
}

func (r *MongoDefaultRepository) UpsertOne(model *MongoUpdateOneModel) *mongo.UpdateResult {
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

func (r *MongoDefaultRepository) UpdateMany(models ...*MongoUpdateOneModel) (*mongo.BulkWriteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(models) == 0 {
		return nil, ErrNoElements
	}

	writeModels := make([]mongo.WriteModel, 0)

	for _, model := range models {
		writeModels = append(writeModels, model.model)
	}

	result, err := r.Collection.coll.BulkWrite(ctx, writeModels)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return result, nil
}

func (r *MongoDefaultRepository) UpsertMany(models ...*MongoUpdateOneModel) (*mongo.BulkWriteResult, error) {
	if len(models) == 0 {
		return nil, ErrNoElements
	}

	for _, model := range models {
		model.model.SetUpsert(true)
	}

	return r.UpdateMany(models...)
}

func (r *MongoDefaultRepository) DeleteOne(model *MongoDeleteOneModel) *mongo.DeleteResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteOptions := options.Delete()
	deleteOptions.Hint = model.Hint

	deleteOneResult, err := r.Collection.coll.DeleteOne(ctx, model.Filter, deleteOptions)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return deleteOneResult
}
