package mongodb

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Pipeline struct {
	Type      string
	Primitive primitive.D
}

type Filter struct {
	Field string
	Value interface{}
}

type PipelineBuilder struct {
	Pipelines []Pipeline
	Filters   []Filter
}

type MongoSort int

const (
	Asc  MongoSort = 1
	Desc MongoSort = -1
)

func NewPipelineBuilder() *PipelineBuilder {
	builder := PipelineBuilder{}
	builder.Pipelines = make([]Pipeline, 0)

	return &builder
}

func (pipelineBuilder *PipelineBuilder) Add(pipeline bson.D) *PipelineBuilder {
	pipelineEntry := Pipeline{
		Type:      "USER",
		Primitive: pipeline,
	}

	pipelineBuilder.Pipelines = append(pipelineBuilder.Pipelines, pipelineEntry)
	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) Page(page int, pageSize int) *PipelineBuilder {
	if page == -1 || pageSize <= 0 {
		return pipelineBuilder
	}

	skip := 0
	if page > 0 {
		skip = page * pageSize
	}

	pipelineBuilder.Skip(skip)
	pipelineBuilder.Limit(pageSize)

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) GetCount(collection *mongo.Collection) int {
	ctx := context.Background()
	cursor, err := pipelineBuilder.Count().AggregateUserWithCount(ctx, collection)
	if err != nil {
		return -1
	}
	var element []map[string]interface{}
	err = cursor.All(ctx, &element)
	if err != nil || len(element) == 0 {
		return 0
	}

	return int(element[0]["count"].(int32))
}

func (pipelineBuilder *PipelineBuilder) Count() *PipelineBuilder {
	countPipeline := Pipeline{
		Type: "COUNT",
		Primitive: bson.D{
			{
				Key:   "$count",
				Value: "count",
			},
		},
	}

	has, index := pipelineBuilder.Has("COUNT")
	if !has {
		pipelineBuilder.Pipelines = append(pipelineBuilder.Pipelines, countPipeline)
	} else {
		pipelineBuilder.Pipelines[index] = countPipeline
	}

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) Filter(field string, value interface{}) *PipelineBuilder {
	has, index := pipelineBuilder.hasField(field)
	if !has {
		filter := Filter{
			Field: field,
			Value: value,
		}
		pipelineBuilder.Filters = append(pipelineBuilder.Filters, filter)
	} else {
		pipelineBuilder.Filters[index].Value = value
	}

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) Skip(skip int) *PipelineBuilder {
	skipPipeline := Pipeline{
		Type: "SKIP",
		Primitive: bson.D{
			{
				Key:   "$skip",
				Value: skip,
			},
		},
	}

	has, index := pipelineBuilder.Has("SKIP")
	if !has {
		pipelineBuilder.Pipelines = append(pipelineBuilder.Pipelines, skipPipeline)
	} else {
		pipelineBuilder.Pipelines[index] = skipPipeline
	}

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) Limit(limit int) *PipelineBuilder {
	limitPipeline := Pipeline{
		Type: "LIMIT",
		Primitive: bson.D{
			{
				Key:   "$limit",
				Value: limit,
			},
		},
	}

	has, index := pipelineBuilder.Has("LIMIT")
	if !has {
		pipelineBuilder.Pipelines = append(pipelineBuilder.Pipelines, limitPipeline)
	} else {
		pipelineBuilder.Pipelines[index] = limitPipeline
	}

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) Sort(field string, order MongoSort) *PipelineBuilder {
	sortPipeline := Pipeline{
		Type: "SORT",
		Primitive: bson.D{
			{
				Key: "$sort",
				Value: bson.D{
					{
						Key:   field,
						Value: order,
					},
				},
			},
		},
	}

	has, index := pipelineBuilder.Has("SORT")
	if !has {
		pipelineBuilder.Pipelines = append(pipelineBuilder.Pipelines, sortPipeline)
	} else {
		pipelineBuilder.Pipelines[index] = sortPipeline
	}

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) SortAfter(field string, order MongoSort) *PipelineBuilder {
	sortPipeline := Pipeline{
		Type: "SORT_AFTER",
		Primitive: bson.D{
			{
				Key: "$sort",
				Value: bson.D{
					{
						Key:   field,
						Value: order,
					},
				},
			},
		},
	}

	has, index := pipelineBuilder.Has("SORT_AFTER")
	if !has {
		pipelineBuilder.Pipelines = append(pipelineBuilder.Pipelines, sortPipeline)
	} else {
		pipelineBuilder.Pipelines[index] = sortPipeline
	}

	return pipelineBuilder
}

func (pipelineBuilder *PipelineBuilder) GetUserPipeline() *bson.A {
	pipelines := bson.A{}
	// Appending the user pipelines if they exist
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "USER" {
			pipelines = append(pipelines, pipeline.Primitive)
		}
	}

	return &pipelines
}

func (pipelineBuilder *PipelineBuilder) GetUserPipelineWithCount() *bson.A {
	pipelines := bson.A{}
	// Appending the user pipelines if they exist
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "USER" {
			pipelines = append(pipelines, pipeline.Primitive)
		}
	}

	// Appending Match  if it exists
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "MATCH" {
			pipelines = append(pipelines, pipeline.Primitive)
			break
		}
	}

	// Appending the count pipelines if they exist
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "COUNT" {
			pipelines = append(pipelines, pipeline.Primitive)
			break
		}
	}

	return &pipelines
}

func (pipelineBuilder *PipelineBuilder) Get() *bson.A {
	pipelines := bson.A{}
	pipelineBuilder.getFilterPipeline()

	// Appending first the sort if it exists
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "SORT" {
			pipelines = append(pipelines, pipeline.Primitive)
			break
		}
	}

	// Appending the user pipelines if they exist
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "USER" {
			pipelines = append(pipelines, pipeline.Primitive)
		}
	}

	// Appending first the sort if it exists
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "SORT_AFTER" {
			pipelines = append(pipelines, pipeline.Primitive)
			break
		}
	}

	// Appending Match  if it exists
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "MATCH" {
			pipelines = append(pipelines, pipeline.Primitive)
			break
		}
	}

	// Appending the skip pipelines if they exist
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "SKIP" {
			pipelines = append(pipelines, pipeline.Primitive)
			break
		}
	}

	// Appending the limit pipelines if they exist
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "LIMIT" {
			pipelines = append(pipelines, pipeline.Primitive)
			break
		}
	}

	// Appending the count pipelines if they exist
	for _, pipeline := range pipelineBuilder.Pipelines {
		if pipeline.Type == "COUNT" {
			pipelines = append(pipelines, pipeline.Primitive)
			break
		}
	}

	return &pipelines
}

func (pipelineBuilder *PipelineBuilder) Aggregate(ctx context.Context, collection *mongo.Collection) (*mongo.Cursor, error) {
	return collection.Aggregate(ctx, *pipelineBuilder.Get())
}

func (pipelineBuilder *PipelineBuilder) AggregateUser(ctx context.Context, collection *mongo.Collection) (*mongo.Cursor, error) {
	return collection.Aggregate(ctx, *pipelineBuilder.GetUserPipeline())
}

func (pipelineBuilder *PipelineBuilder) AggregateUserWithCount(ctx context.Context, collection *mongo.Collection) (*mongo.Cursor, error) {
	return collection.Aggregate(ctx, *pipelineBuilder.GetUserPipelineWithCount())
}

func (pipelineBuilder *PipelineBuilder) Has(key string) (bool, int) {
	for index, pipeline := range pipelineBuilder.Pipelines {
		if strings.EqualFold(pipeline.Type, key) {
			return true, index
		}
	}

	return false, -1
}

func (pipelineBuilder *PipelineBuilder) hasField(fieldName string) (bool, int) {
	for index, field := range pipelineBuilder.Filters {
		if strings.EqualFold(field.Field, fieldName) {
			return true, index
		}
	}

	return false, -1
}

func (pipelineBuilder *PipelineBuilder) getFilterPipeline() bool {
	if len(pipelineBuilder.Filters) == 0 {
		return false
	}

	fields := primitive.D{}

	for _, filter := range pipelineBuilder.Filters {
		primitiveField := bson.D{
			{
				Key: filter.Field,
				Value: bson.D{
					{
						Key:   "$regex",
						Value: filter.Value,
					},
					{
						Key:   "$options",
						Value: "i",
					},
				},
			},
		}
		fields = append(fields, primitiveField...)
	}
	matchPipeline := Pipeline{
		Type: "MATCH",
		Primitive: bson.D{
			{
				Key:   "$match",
				Value: fields,
			},
		},
	}

	has, index := pipelineBuilder.Has("MATCH")
	if !has {
		pipelineBuilder.Pipelines = append(pipelineBuilder.Pipelines, matchPipeline)
	} else {
		pipelineBuilder.Pipelines[index] = matchPipeline
	}

	return true
}
