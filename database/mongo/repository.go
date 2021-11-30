package database

import (
	"context"
	"time"

	"github.com/cjlapao/common-go/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var logger = log.Get()

type Repository interface {
	Find(fieldName string, value string) []*primitive.M
	FindOne(fieldName string, value string) *mongo.SingleResult
	InsertOne(element interface{}) *mongo.InsertOneResult
	InsertMany(elements []interface{}) *mongo.InsertManyResult
	UpsertOne(filterField string, filterValue string, element interface{}) *mongo.UpdateResult
	UpsertMany(filterField string, filterValue string, elements []interface{}) *mongo.UpdateResult
}

type DefaultRepository struct {
	Client         *mongo.Client
	Database       *mongo.Database
	Collection     *mongo.Collection
	DatabaseName   string
	CollectionName string
}

func NewRepository(factory *MongoFactory, database string, collection string) Repository {
	defaultRepo := DefaultRepository{
		DatabaseName:   database,
		CollectionName: collection,
	}
	defaultRepo.Client = factory.Client
	defaultRepo.Database = factory.GetDatabase()
	defaultRepo.Collection = factory.GetCollection(collection)

	return &defaultRepo
}

func (r *DefaultRepository) Find(fieldName string, value string) []*primitive.M {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	filter := bson.D{{fieldName, value}}
	cur, err := r.Collection.Find(ctx, filter)
	if err != nil {
		logger.LogError(err)
		return nil
	}
	var elements []*bson.M
	for cur.Next(ctx) {
		var element bson.M
		err := cur.Decode(&element)
		if err != nil {
			logger.LogError(err)
			return nil
		}
		elements = append(elements, &element)
	}

	return elements
}

func (r *DefaultRepository) FindOne(fieldName string, value string) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	filter := bson.D{{fieldName, value}}

	cur := r.Collection.FindOne(ctx, filter)

	return cur
}

func (r *DefaultRepository) InsertOne(element interface{}) *mongo.InsertOneResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := r.Collection.InsertOne(ctx, element)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return insertResult
}

func (r *DefaultRepository) InsertMany(elements []interface{}) *mongo.InsertManyResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	insertResult, err := r.Collection.InsertMany(ctx, elements)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return insertResult
}

func (r *DefaultRepository) UpsertOne(filterField string, filterValue string, element interface{}) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{filterField, filterValue}}
	options := options.Update().SetUpsert(true)
	updatePipeline := bson.D{{"$set", element}}
	updateOneResult, err := r.Collection.UpdateOne(ctx, filter, updatePipeline, options)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return updateOneResult
}

func (r *DefaultRepository) UpsertMany(filterField string, filterValue string, elements []interface{}) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{{filterField, filterValue}}
	options := options.Update().SetUpsert(true)

	updateOneResult, err := r.Collection.UpdateMany(ctx, filter, elements, options)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return updateOneResult
}
