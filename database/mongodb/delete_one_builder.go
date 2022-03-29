package mongodb

//TODO: Refactor implementation
import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDeleteOneModel struct {
	model  *mongo.DeleteOneModel
	Filter interface{}
	Hint   interface{}
}

// Transforms the model into a json string representation
func (model MongoDeleteOneModel) String() string {
	result, err := json.MarshalIndent(model.model, "", "  ")
	if err != nil {
		return ""
	}

	return string(result)
}

type DeleteOneBuilder struct {
	LiteralFilter string
	Filters       []builderElement
}

func NewDeleteOneBuilder() *DeleteOneBuilder {
	return &DeleteOneBuilder{
		Filters: make([]builderElement, 0),
	}
}

// Filter creates a filter for the model, this will allow to update just a subset of the collection
// if no filter is present it will apply the operation to all documents in the collection
// You can use odata type of query, for example:
//		builder.Filter("userId eq 'some_id'")
func (c *DeleteOneBuilder) Filter(query string) *DeleteOneBuilder {
	c.LiteralFilter = query
	return c
}

// FilterBy creates a filter for the model using the field, operation and value to filter on
// if you have several filters applied they will always be joined by AND, no OR is permitted
// for example:
//		builder.FilterBy("userId", Equals, "some_id")
func (c *DeleteOneBuilder) FilterBy(key string, operation filterOperation, value interface{}) *DeleteOneBuilder {
	element := builderElement{
		FilterOperation,
		key,
		operation,
		value,
	}

	c.Filters = append(c.Filters, element)
	return c
}

// Build Builds the model to use during the query
func (c *DeleteOneBuilder) Build() (*MongoDeleteOneModel, error) {
	model := mongo.DeleteOneModel{}

	// Processing the filter elements, this can be the literal or just the operations
	if len(c.Filters) > 0 {
		filterPrimitives := bson.M{}
		for _, filterElement := range c.Filters {
			stringFilter := getOperationString(filterElement.key, filterElement.filterOperation, filterElement.value)

			filterParser := NewFilterParser(stringFilter)
			parsedFilter, err := filterParser.Parse()
			if err != nil {
				return nil, err
			}
			filterPrimitives[filterElement.key] = parsedFilter.(primitive.M)[filterElement.key]
		}

		model.Filter = filterPrimitives
	} else if len(c.LiteralFilter) > 0 {
		filterParser := NewFilterParser(c.LiteralFilter)
		parsedFilter, err := filterParser.Parse()
		if err != nil {
			return nil, err
		}
		model.Filter = parsedFilter
	} else {
		model.Filter = bson.D{}
	}

	return &MongoDeleteOneModel{
		model:  &model,
		Filter: model.Filter,
		Hint:   model.Hint,
	}, nil
}
