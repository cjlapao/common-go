package mongodb

//TODO: Refactor implementation
import (
	"context"
	"encoding/hex"
	"errors"
	"net/url"
	"strings"

	"github.com/cjlapao/common-go/models"
	"github.com/cjlapao/common-go/odata"
	"github.com/cjlapao/common-go/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ODataParser struct {
	Collection *mongoCollection
}

var ErrInvalidInput = errors.New("odata syntax error")

func EmptyODataParser(collection *mongoCollection) *ODataParser {
	result := ODataParser{
		Collection: collection,
	}

	return &result
}

func (odataParser *ODataParser) GetODataResponse(query url.Values) (*models.ODataResponse, error) {
	ctx := context.Background()
	response := models.ODataResponse{}

	queryMap, err := odata.ParseURLValues(query)
	if err != nil {
		return nil, err
	}

	// Checks if the count flag is true and count the collection records
	if count, ok := queryMap[odata.Count].(bool); ok {
		if count {
			builder := NewPipelineBuilder()
			countField := builder.CountCollection(odataParser.Collection.coll)
			response.Count = countField
		}
	}

	cursor, err := odataParser.Query(query)
	if err != nil {
		return nil, err
	}
	var element []map[string]interface{}
	err = cursor.All(ctx, &element)
	if err != nil {
		return nil, err
	}

	response.Value = element

	return &response, nil
}

func (odataParser *ODataParser) Parse(query url.Values, destination interface{}) error {
	return nil
}

// Query creates a mgo query based on odata parameters
func (odataParser *ODataParser) Query(query url.Values) (*mongo.Cursor, error) {
	ctx := context.Background()
	builder := NewPipelineBuilder()

	// Parse url values
	queryMap, err := odata.ParseURLValues(query)
	if err != nil {
		return nil, err
	}

	// Prepares the limit pipeline
	if limit, ok := queryMap[odata.Top].(int); ok {
		builder.Limit(limit)
	}

	// Prepares the skip pipeline
	if skip, ok := queryMap[odata.Skip].(int); ok {
		builder.Skip(skip)
	}

	// Prepares the filter object and build the match pipeline
	filterObj := make(bson.M)
	if queryMap[odata.Filter] != nil {
		filterQuery, _ := queryMap[odata.Filter].(*parser.ParseNode)
		var err error
		filterObj, err = applyFilter(filterQuery)
		if err != nil {
			return nil, ErrInvalidInput
		}

		// Creates the match pipeline for the filter
		builder.Match(filterObj)
	}

	// Prepare Select object and build the project pipeline
	selectMap := make(bson.M)
	if queryMap["$select"] != nil {
		if selectSlice, ok := queryMap["$select"].([]string); ok {
			for i := 0; i < len(selectSlice); i++ {
				fieldName := selectSlice[i]
				selectMap[fieldName] = 1
			}

			// Creates the project pipeline for the select argument
			builder.Project(selectMap)
		}
	}

	// Prepare the sort object and build the sort pipeline
	sortMap := make(bson.M)
	if queryMap[odata.OrderBy] != nil {
		orderBySlice := queryMap[odata.OrderBy].([]odata.OrderItem)
		for _, item := range orderBySlice {
			if item.Order == "desc" {
				sortMap[item.Field] = -1
			} else {
				sortMap[item.Field] = 1
			}
		}

		// Create the sort pipeline for the sorting object
		builder.Sort(sortMap)
	}

	return builder.Aggregate(ctx, odataParser.Collection.coll)
}

// ODataCount runs a collection.Count() function based on $count odata parameter
// func ODataCount(collection *mgo.Collection) (int, error) {
// 	return collection.Count()
// }

// ODataInlineCount retrieves the total count from a filtered data
// func ODataInlineCount(collection *mgo.Collection) (int, error) {

// 	return collection.Find(filterObj).Count()
// }

func applyFilter(node *parser.ParseNode) (bson.M, error) {

	filter := make(bson.M)

	if _, ok := node.Token.Value.(string); ok {
		switch node.Token.Value {

		case "eq":
			// Escape single quotes in the case of strings
			if _, valueOk := node.Children[1].Token.Value.(string); valueOk {
				node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			}
			value := bson.M{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, keyOk := node.Children[0].Token.Value.(string); !keyOk {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "ne":
			// Escape single quotes in the case of strings
			if _, valueOk := node.Children[1].Token.Value.(string); valueOk {
				node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			}
			value := bson.M{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, keyOk := node.Children[0].Token.Value.(string); !keyOk {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "gt":
			var keyString string
			if keyString, ok = node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}

			var value bson.M
			if keyString == "_id" {
				var idString string
				if _, ok := node.Children[1].Token.Value.(string); ok {
					idString = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
				}
				decodedString, err := hex.DecodeString(idString)
				if err != nil || len(decodedString) != 12 {
					return nil, ErrInvalidInput
				}
				objectId := primitive.NewObjectID()
				objectId.UnmarshalText(decodedString)
				value = bson.M{"$" + node.Token.Value.(string): objectId}
			} else {
				value = bson.M{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			}
			filter[keyString] = value

		case "ge":
			value := bson.M{"$gte": node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "lt":
			value := bson.M{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "le":
			value := bson.M{"$lte": node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "and":
			leftFilter, err := applyFilter(node.Children[0]) // Left children
			if err != nil {
				return nil, err
			}
			rightFilter, _ := applyFilter(node.Children[1]) // Right children
			if err != nil {
				return nil, err
			}
			filter["$and"] = []bson.M{leftFilter, rightFilter}

		case "or":
			leftFilter, err := applyFilter(node.Children[0]) // Left children
			if err != nil {
				return nil, err
			}
			rightFilter, err := applyFilter(node.Children[1]) // Right children
			if err != nil {
				return nil, err
			}
			filter["$or"] = []bson.M{leftFilter, rightFilter}

		//Functions
		case "startswith":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)

			value := primitive.Regex{
				Pattern: "^" + node.Children[1].Token.Value.(string),
				Options: "gi",
			}

			filter[node.Children[0].Token.Value.(string)] = value

		case "endswith":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)

			value := primitive.Regex{
				Pattern: node.Children[1].Token.Value.(string) + "$",
				Options: "gi",
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "contains":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)

			value := primitive.Regex{
				Pattern: node.Children[1].Token.Value.(string),
				Options: "gi",
			}
			filter[node.Children[0].Token.Value.(string)] = value

		}
	}
	return filter, nil
}
