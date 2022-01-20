package mongodb

//TODO: Refactor implementation
import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RemoveOneBuilder struct {
	Filter []BuilderElement
}

func NewRemoveOneBuilder() *RemoveOneBuilder {
	return &RemoveOneBuilder{
		Filter: make([]BuilderElement, 0),
	}
}

func (c *RemoveOneBuilder) FilterBy(key string, value interface{}) *RemoveOneBuilder {
	element := BuilderElement{key, value}
	c.Filter = append(c.Filter, element)
	return c
}

func (c *RemoveOneBuilder) Build() mongo.DeleteOneModel {
	model := mongo.DeleteOneModel{}

	if len(c.Filter) == 0 {
		model.Filter = bson.D{{}}
	} else {
		filterPrimitives := make([]primitive.E, 0)

		for _, filterElement := range c.Filter {
			bsonElement := primitive.E{
				Key:   filterElement.Key,
				Value: filterElement.Value,
			}

			filterPrimitives = append(filterPrimitives, bsonElement)
		}

		model.Filter = filterPrimitives
	}

	return model
}
