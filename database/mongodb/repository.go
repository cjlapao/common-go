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
	FindOne(filter interface{}) *mongoSingleResult
	InsertOne(element interface{}) (*mongoInsertOneResult, error)
	InsertMany(elements ...interface{}) (*mongoInsertManyResult, error)
	UpdateOne(model *MongoUpdateOneModel) (*mongoUpdateResult, error)
	UpdateMany(models ...*MongoUpdateOneModel) (*mongoBulkWriteResult, error)
	UpsertOne(model *MongoUpdateOneModel) (*mongoUpdateResult, error)
	UpsertMany(models ...*MongoUpdateOneModel) (*mongoBulkWriteResult, error)
	DeleteOne(model *MongoDeleteOneModel) (*mongoDeleteResult, error)
	DeleteMany(filter interface{}) (*mongoDeleteResult, error)
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
			filterToApply = bson.D{}
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

// FindFieldBy finds a record using a specific field and operation, this is a more restricted filtering
// but helps with the strong typed query that can be produced
//
// Example:
//		repository.FindFieldBy("userId" , mongo.Equal, "someId")
func (r *MongoDefaultRepository) FindFieldBy(fieldName string, operation filterOperation, value interface{}) (*mongoCursor, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	stringFilter := getOperationString(fieldName, operation, value)

	filterParser := NewFilterParser(stringFilter)
	parsedFilter, err := filterParser.Parse()

	defer cancel()

	if err != nil {
		return nil, err
	}

	cur, err := r.Collection.coll.Find(ctx, parsedFilter)

	return &mongoCursor{cursor: cur}, err
}

// FindOne Finds a single result based on a filter, this can be a bson document or
// a odata type of query.
//
// Example:
//		repository.FindOne("userId eq 'someId'")
// or
//		repository.FindOne(bson.M{"userId": "someId"})
// The result will be a `mongoSingleResult` that can be decoded to any interface
func (r *MongoDefaultRepository) FindOne(filter interface{}) *mongoSingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var filterToApply interface{}

	defer cancel()
	if stringFilter, ok := filter.(string); ok {
		if filter == "" {
			filterToApply = bson.D{}
		} else {
			processedFilter, err := NewFilterParser(stringFilter).Parse()
			if err != nil {
				logger.Error("There was an error applying the filter, %v", err.Error())
				return nil
			} else {
				filterToApply = processedFilter
			}
		}
	} else {
		filterToApply = filter
	}

	result := r.Collection.coll.FindOne(ctx, filterToApply)

	return &mongoSingleResult{sr: result}
}

// InsertOne Inserts a record int the collection and returns the inserted id
func (r *MongoDefaultRepository) InsertOne(element interface{}) (*mongoInsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := r.Collection.coll.InsertOne(ctx, element)

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	result := mongoInsertOneResult{}
	result.FromMongo(insertResult)
	return &result, err
}

// InsertMany Inserts multiple records in the collection and returns the inserted id's
func (r *MongoDefaultRepository) InsertMany(elements ...interface{}) (*mongoInsertManyResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := r.Collection.coll.InsertMany(ctx, elements)

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	result := mongoInsertManyResult{}
	result.FromMongo(insertResult)
	return &result, nil
}

// UpdateOne updates a document in a collection using a UpdateOneModel, this can be constructed
// using strong typed language when called the UpdateOneModelBuilder
// It will return the number of affected documents
func (r *MongoDefaultRepository) UpdateOne(model *MongoUpdateOneModel) (*mongoUpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	options := options.Update().SetUpsert(*model.model.Upsert)
	updateOneResult, err := r.Collection.coll.UpdateOne(ctx, model.Filter, model.Update, options)

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	result := mongoUpdateResult{}
	result.FromMongo(updateOneResult)
	return &result, nil
}

// UpsertOne similar to the UpdateOne this updates or inserts (if it does not exists) a document
// in a collection using a UpdateOneModel, this can be constructed using strong typed language
// when called the UpdateOneModelBuilder.
// It will return the number of affected documents
func (r *MongoDefaultRepository) UpsertOne(model *MongoUpdateOneModel) (*mongoUpdateResult, error) {
	model.model.SetUpsert(true)
	return r.UpdateOne(model)
}

// UpdateMany updates documents in the collection using a UpdateOneModel, this can be constructed
// using strong typed language when called the UpdateOneModelBuilder.
func (r *MongoDefaultRepository) UpdateMany(models ...*MongoUpdateOneModel) (*mongoBulkWriteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if len(models) == 0 {
		return nil, ErrNoElements
	}

	writeModels := make([]mongo.WriteModel, 0)

	for _, model := range models {
		writeModels = append(writeModels, model.model)
	}

	bulkWriteResult, err := r.Collection.coll.BulkWrite(ctx, writeModels)
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	result := mongoBulkWriteResult{}
	result.FromMongo(bulkWriteResult)
	return &result, nil
}

// UpsertMany Similar to UpdateMany this updates or inserts (if it does not exists) a document
// in a collection using a UpdateOneModel, this can be constructed using strong typed language
// when called the UpdateOneModelBuilder.
func (r *MongoDefaultRepository) UpsertMany(models ...*MongoUpdateOneModel) (*mongoBulkWriteResult, error) {
	if len(models) == 0 {
		return nil, ErrNoElements
	}

	for _, model := range models {
		model.model.SetUpsert(true)
	}

	return r.UpdateMany(models...)
}

// DeleteOne Deletes a document in a collection using a DeleteOneModel, this can be constructed
// using strong typed language using the DeleteOneModelBuilder
func (r *MongoDefaultRepository) DeleteOne(model *MongoDeleteOneModel) (*mongoDeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteOptions := options.Delete()
	deleteOptions.Hint = model.Hint

	deleteOneResult, err := r.Collection.coll.DeleteOne(ctx, model.Filter, deleteOptions)

	if err != nil {
		logger.Exception(err, "There was an error while deleting collection documents")
		return nil, err
	}

	result := mongoDeleteResult{}
	result.FromMongo(deleteOneResult)
	return &result, nil
}

// DeleteMany Deletes documents in a collection based on a filter, the filter can be a valid bson document
// or you can use a odata type query.
//
// Example:
//		repository.DeleteMany(bson.M{"userId": "someId"})
// or
//		repository.DeleteMany("userId eq 'someId'")
func (r *MongoDefaultRepository) DeleteMany(filter interface{}) (*mongoDeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var filterToApply interface{}

	if stringFilter, ok := filter.(string); ok {
		if filter == "" {
			filterToApply = bson.D{}
		} else {
			processedFilter, err := NewFilterParser(stringFilter).Parse()
			if err != nil {
				logger.Exception(err, "There was an error applying the filter")
				return nil, err
			} else {
				filterToApply = processedFilter
			}
		}
	} else {
		filterToApply = filter
	}

	deleteOptions := options.Delete()

	deleteOneResult, err := r.Collection.coll.DeleteMany(ctx, filterToApply, deleteOptions)

	if err != nil {
		logger.Exception(err, "There was an error while deleting collection documents")
		return nil, err
	}

	result := mongoDeleteResult{}
	result.FromMongo(deleteOneResult)
	return &result, nil
}
