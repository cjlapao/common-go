package mongodb

//TODO: Refactor implementation
import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/cjlapao/common-go/guard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateOneModelBuilder struct {
	LiteralFilter string
	Elements      []BuilderElement
}

// NewUpdateOneModelBuilder Creates a new builder for an updateone model
func NewUpdateOneModelBuilder() *UpdateOneModelBuilder {
	return &UpdateOneModelBuilder{
		Elements: make([]BuilderElement, 0),
	}
}

// Set Sets a field to be updated, this will add the property if it does not exist in the model
func (c *UpdateOneModelBuilder) Set(field string, value interface{}) *UpdateOneModelBuilder {
	guard.FatalEmptyOrNil(field)

	element := BuilderElement{
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

	element := BuilderElement{
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

	element := BuilderElement{
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
	element := BuilderElement{
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
//		builder.Encode()
func (c *UpdateOneModelBuilder) Encode(element interface{}, ignoredFields ...string) *UpdateOneModelBuilder {
	var mapped map[string]interface{}
	customRegistry := createCustomRegistry().Build()

	//converting the document to a bson
	marshalled, err := bson.MarshalExtJSONWithRegistry(customRegistry, element, false, false)
	if err != nil {
		return nil
	}

	// Converting the json marshalled element to a map
	json.Unmarshal(marshalled, &mapped)

	// Removing any ignored field
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

func (c *UpdateOneModelBuilder) Build() mongo.UpdateOneModel {
	model := mongo.UpdateOneModel{}

	if len(c.Elements) == 0 {
		panic(errors.New("no elements to update"))
	}

	setOperations := c.getElements(SetOperation)
	setOnInsertOperations := c.getElements(SetOnInsertOperation)
	filterOperations := c.getElements(FilterOperation)
	unsetOperations := c.getElements(UnsetOperation)

	setElementsPrimitives := make([]primitive.E, 0)
	setOnInsertPrimitives := make([]primitive.E, 0)
	unsetPrimitives := make([]primitive.E, 0)

	for _, updateElement := range setOperations {
		bsonElement := primitive.E{
			Key:   updateElement.key,
			Value: updateElement.value,
		}

		setElementsPrimitives = append(setElementsPrimitives, bsonElement)
	}

	for _, updateElement := range setOnInsertOperations {
		bsonElement := primitive.E{
			Key:   updateElement.key,
			Value: updateElement.value,
		}

		setOnInsertPrimitives = append(setElementsPrimitives, bsonElement)
	}

	for _, unsetElement := range unsetOperations {
		bsonElement := primitive.E{
			Key:   unsetElement.key,
			Value: unsetElement.value,
		}

		unsetPrimitives = append(unsetPrimitives, bsonElement)
	}

	update := bson.M{}
	if len(setElementsPrimitives) > 0 {
		update["$set"] = setElementsPrimitives
	}

	if len(setOnInsertPrimitives) > 0 {
		update["$setOnInsert"] = setOnInsertPrimitives
	}

	if len(unsetOperations) > 0 {
		update["$unset"] = unsetPrimitives
	}
	model.Update = update

	if len(filterOperations) > 0 {
		filterPrimitives := bson.M{}
		for _, filterElement := range filterOperations {
			stringFilter := getOperationString(filterElement.key, filterElement.filterOperation, filterElement.value)

			filterParser := NewFilterParser(stringFilter)
			parsedFilter, err := filterParser.Parse()
			if err == nil {
				filterPrimitives[filterElement.key] = parsedFilter.(primitive.M)[filterElement.key]
			}
		}

		model.Filter = filterPrimitives
	} else if len(c.LiteralFilter) > 0 {
		filterParser := NewFilterParser(c.LiteralFilter)
		parsedFilter, err := filterParser.Parse()
		if err == nil {
			model.Filter = parsedFilter
		}
	} else {
		model.Filter = bson.D{}
	}

	return model
}

func (c *UpdateOneModelBuilder) getElements(operation elementBuilderOperation) []BuilderElement {
	result := make([]BuilderElement, 0)
	for _, element := range c.Elements {
		if element.operation == operation {
			result = append(result, element)
		}
	}

	return result
}

func (c *UpdateOneModelBuilder) hasElement(key string) (bool, int) {
	for idx, element := range c.Elements {
		if strings.EqualFold(key, element.key) {
			return true, idx
		}
	}

	return false, -1
}
