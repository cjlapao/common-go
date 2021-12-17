package mongodb

import (
	"encoding/json"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UpdateOneBuilder struct {
	Filter   []BuilderElement
	Elements []BuilderElement
}

func NewUpdateOneBuilder() *UpdateOneBuilder {
	return &UpdateOneBuilder{
		Filter:   make([]BuilderElement, 0),
		Elements: make([]BuilderElement, 0),
	}
}

func (c *UpdateOneBuilder) SetElement(key string, value interface{}) *UpdateOneBuilder {
	// guard.FatalEmptyOrNil(key)
	// guard.FatalEmptyOrNil(value)

	element := BuilderElement{key, value}
	c.Elements = append(c.Elements, element)
	return c
}

func (c *UpdateOneBuilder) Encode(element interface{}, ignoredFields ...string) *UpdateOneBuilder {
	var mapped map[string]interface{}
	//converting the document to a bson
	marshalled, err := bson.MarshalExtJSON(element, false, true)
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
		if !ignored {
			c.SetElement(field, val)
		}
	}

	return c
}

func (c *UpdateOneBuilder) FilterBy(key string, value interface{}) *UpdateOneBuilder {
	element := BuilderElement{key, value}
	c.Filter = append(c.Filter, element)
	return c
}

func (c *UpdateOneBuilder) Build() mongo.UpdateOneModel {
	model := mongo.UpdateOneModel{}

	if len(c.Elements) == 0 {
		panic(errors.New("no elements to update"))
	}

	updatePrimitives := make([]primitive.E, 0)

	for _, updateElement := range c.Elements {
		bsonElement := primitive.E{
			Key:   updateElement.Key,
			Value: updateElement.value,
		}

		updatePrimitives = append(updatePrimitives, bsonElement)
	}

	update := bson.D{{Key: "$set", Value: updatePrimitives}}
	model.Update = update

	if len(c.Filter) == 0 {
		model.Filter = bson.D{{}}
	} else {
		filterPrimitives := make([]primitive.E, 0)

		for _, filterElement := range c.Filter {
			bsonElement := primitive.E{
				Key:   filterElement.Key,
				Value: filterElement.value,
			}

			filterPrimitives = append(filterPrimitives, bsonElement)
		}

		model.Filter = filterPrimitives
	}

	return model
}
