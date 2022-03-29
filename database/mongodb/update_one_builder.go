package mongodb

//TODO: Refactor implementation
import (
	"encoding/json"
	"strings"

	"github.com/cjlapao/common-go/guard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUpdateOneModel struct {
	model  *mongo.UpdateOneModel
	Filter interface{}
	Hint   interface{}
	Update interface{}
}

// Transforms the model into a json string representation
func (model MongoUpdateOneModel) String() string {
	result, err := json.MarshalIndent(model.model, "", "  ")
	if err != nil {
		return ""
	}

	return string(result)
}

type UpdateOneModelBuilder struct {
	LiteralFilter string
	Elements      []builderElement
}

// NewUpdateOneModelBuilder Creates a new builder for an updateone model
func NewUpdateOneModelBuilder() *UpdateOneModelBuilder {
	return &UpdateOneModelBuilder{
		Elements: make([]builderElement, 0),
	}
}

// Set Sets a field to be updated, this will add the property if it does not exist in the model
func (c *UpdateOneModelBuilder) Set(field string, value interface{}) *UpdateOneModelBuilder {
	guard.FatalEmptyOrNil(field)

	element := builderElement{
		SetOperation,
		field,
		Equal,
		value,
	}

	has, idx := c.hasElement(field)
	if !has {
		c.Elements = append(c.Elements, element)
	} else {
		c.Elements[idx].operation = SetOperation
		c.Elements[idx].value = value
	}
	return c
}

// Unset Unsets a field in the document, this will effectively remove the property from the document
func (c *UpdateOneModelBuilder) Unset(field string) *UpdateOneModelBuilder {
	guard.FatalEmptyOrNil(field)

	element := builderElement{
		UnsetOperation,
		field,
		Equal,
		"1",
	}

	has, idx := c.hasElement(field)
	if !has {
		c.Elements = append(c.Elements, element)
	} else {
		c.Elements[idx].operation = UnsetOperation
		c.Elements[idx].value = "1"
	}
	return c
}

// SetOnInsert, only updates a value on insert, if the operation is a update it will be ignored
func (c *UpdateOneModelBuilder) SetOnInsert(field string, value interface{}) *UpdateOneModelBuilder {
	guard.FatalEmptyOrNil(field)

	element := builderElement{
		SetOnInsertOperation,
		field,
		Equal,
		value,
	}

	has, idx := c.hasElement(field)
	if !has {
		c.Elements = append(c.Elements, element)
	} else {
		c.Elements[idx].operation = SetOnInsertOperation
		c.Elements[idx].value = value
	}

	return c
}

// Filter creates a filter for the model, this will allow to update just a subset of the collection
// if no filter is present it will apply the operation to all documents in the collection
// You can use odata type of query, for example:
//		builder.Filter("userId eq 'some_id'")
func (c *UpdateOneModelBuilder) Filter(query string) *UpdateOneModelBuilder {
	c.LiteralFilter = query
	return c
}

// FilterBy creates a filter for the model using the field, operation and value to filter on
// if you have several filters applied they will always be joined by AND, no OR is permitted
// for example:
//		builder.FilterBy("userId", Equals, "some_id")
func (c *UpdateOneModelBuilder) FilterBy(key string, operation filterOperation, value interface{}) *UpdateOneModelBuilder {
	element := builderElement{
		FilterOperation,
		key,
		operation,
		value,
	}

	c.Elements = append(c.Elements, element)
	return c
}

// Encode creates a bson representation of an already existing interface, this is helpful to pass
// an object and let the encode generate the sets for all properties. you can also ignore fields
// if you have used any of the Set or SetOnInsert methods to build this model, they superseed the
// encode
// for example:
//		var obj := struct {
//			Name string
//			Timestamp time.Time
//		}{
//			Name: "SomeProperty"
//			Timestamp: time.
//		}
//		builder.Encode(obj, "Timestamp")
func (c *UpdateOneModelBuilder) Encode(element interface{}, ignoredFields ...string) *UpdateOneModelBuilder {
	var mapped map[string]interface{}
	// Creating mongo marshaler custom registry for date, time and oid types
	customRegistry := createCustomRegistry().Build()

	//converting the document to a bson
	marshalled, err := bson.MarshalExtJSONWithRegistry(customRegistry, element, false, false)
	if err != nil {
		return nil
	}

	// Converting the json marshalled element to a map
	json.Unmarshal(marshalled, &mapped)

	// Removing any ignored field and fields that have been explicitly set or setOnInsert
	for field, val := range mapped {
		ignored := false
		for _, ignoredField := range ignoredFields {
			if strings.EqualFold(ignoredField, field) {
				ignored = true
				break
			}
		}

		has, _ := c.hasElement(field)
		if !ignored && !has {
			c.Set(field, val)
		}
	}

	return c
}

// Build Builds the model to use during the query
// you can pass options for the build process, for example
//		builder.Build(UpsertBuildOption)
func (c *UpdateOneModelBuilder) Build(options ...BuilderOptions) (*MongoUpdateOneModel, error) {
	model := mongo.UpdateOneModel{}
	model.SetUpsert(false)

	// if there is no instructions to build
	if len(c.Elements) == 0 {
		return nil, ErrNoElements
	}

	// getting all elements filtered by there respective functions
	setOperations := c.getElements(SetOperation)
	setOnInsertOperations := c.getElements(SetOnInsertOperation)
	filterOperations := c.getElements(FilterOperation)
	unsetOperations := c.getElements(UnsetOperation)

	// creating initial arrays for each of the different functions
	setElementsPrimitives := make([]primitive.E, 0)
	setOnInsertPrimitives := make([]primitive.E, 0)
	unsetPrimitives := make([]primitive.E, 0)

	// creating the primitives for the set operations
	for _, updateElement := range setOperations {
		bsonElement := primitive.E{
			Key:   updateElement.key,
			Value: updateElement.value,
		}

		setElementsPrimitives = append(setElementsPrimitives, bsonElement)
	}

	// creating the primitives for the setOnInsert operations
	for _, updateElement := range setOnInsertOperations {
		bsonElement := primitive.E{
			Key:   updateElement.key,
			Value: updateElement.value,
		}

		setOnInsertPrimitives = append(setOnInsertPrimitives, bsonElement)
	}

	// creating the primitives for the unset operations
	for _, unsetElement := range unsetOperations {
		bsonElement := primitive.E{
			Key:   unsetElement.key,
			Value: unsetElement.value,
		}

		unsetPrimitives = append(unsetPrimitives, bsonElement)
	}

	// Creating the root set object if there is any primitive
	update := bson.M{}
	if len(setElementsPrimitives) > 0 {
		update["$set"] = setElementsPrimitives
	}

	// Creating the root setOnInsert object if there is any primitive
	if len(setOnInsertPrimitives) > 0 {
		update["$setOnInsert"] = setOnInsertPrimitives
	}

	// Creating the root unset object if there is any primitive
	if len(unsetOperations) > 0 {
		update["$unset"] = unsetPrimitives
	}

	model.Update = update

	// Processing the filter elements, this can be the literal or just the operations
	if len(filterOperations) > 0 {
		filterPrimitives := bson.M{}
		for _, filterElement := range filterOperations {
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

	if len(options) > 0 {
		for _, option := range options {
			if option == UpsertBuildOption {
				model.SetUpsert(true)
			}
		}
	}

	return &MongoUpdateOneModel{
		model:  &model,
		Filter: model.Filter,
		Hint:   model.Hint,
		Update: model.Update,
	}, nil
}

// getElements Gets the elements filtered by operation
func (c *UpdateOneModelBuilder) getElements(operation elementBuilderOperation) []builderElement {
	result := make([]builderElement, 0)
	for _, element := range c.Elements {
		if element.operation == operation {
			result = append(result, element)
		}
	}

	return result
}

// hasElement Checks if an element exists in the slice
func (c *UpdateOneModelBuilder) hasElement(key string) (bool, int) {
	for idx, element := range c.Elements {
		if strings.EqualFold(key, element.key) {
			return true, idx
		}
	}

	return false, -1
}
