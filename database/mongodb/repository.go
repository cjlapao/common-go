package mongodb

import (
	"context"
	"time"

	"github.com/cjlapao/common-go/guard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository interface {
	Filter(filter interface{}) []*interface{}
	Find(fieldName string, value string) []*interface{}
	FindOne(fieldName string, value string) *mongo.SingleResult
	InsertOne(element interface{}) *mongo.InsertOneResult
	InsertMany(elements []interface{}) *mongo.InsertManyResult
	UpsertOne(model mongo.UpdateOneModel) *mongo.UpdateResult
	UpsertMany(filter interface{}, elements []interface{}) *mongo.UpdateResult
	DeleteOne(model mongo.DeleteOneModel) *mongo.DeleteResult
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

func (r *DefaultRepository) Filter(filter interface{}) []*interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	if filter == "" {
		filter = bson.D{{}}
	}
	defer cancel()

	cur, err := r.Collection.Find(ctx, filter)
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

func (r *DefaultRepository) Find(fieldName string, value string) []*interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	filter := bson.D{
		{
			Key:   fieldName,
			Value: value,
		},
	}

	cur, err := r.Collection.Find(ctx, filter)
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

func (r *DefaultRepository) FindOne(fieldName string, value string) *mongo.SingleResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()
	filter := bson.D{
		{
			Key:   fieldName,
			Value: value,
		},
	}

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

func (r *DefaultRepository) UpsertOne(model mongo.UpdateOneModel) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Update().SetUpsert(true)
	updateOneResult, err := r.Collection.UpdateOne(ctx, model.Filter, model.Update, options)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return updateOneResult
}

func (r *DefaultRepository) UpsertMany(filter interface{}, elements []interface{}) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if guard.IsNill(filter) {
		filter = bson.D{{}}
	}

	options := options.Update().SetUpsert(true)

	updateOneResult, err := r.Collection.UpdateMany(ctx, filter, elements, options)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return updateOneResult
}

func (r *DefaultRepository) DeleteOne(model mongo.DeleteOneModel) *mongo.DeleteResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	deleteOptions := options.Delete()
	deleteOptions.Collation = model.Collation
	deleteOptions.Hint = model.Hint

	deleteOneResult, err := r.Collection.DeleteOne(ctx, model.Filter, deleteOptions)

	if err != nil {
		logger.LogError(err)
		return nil
	}

	return deleteOneResult
}
