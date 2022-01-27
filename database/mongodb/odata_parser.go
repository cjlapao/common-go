package mongodb

//TODO: Refactor implementation
import (
	"context"
	"errors"
	"net/url"
	"reflect"

	"github.com/cjlapao/common-go/models"
	"github.com/cjlapao/common-go/odata"
	"github.com/cjlapao/common-go/parser"
	"go.mongodb.org/mongo-driver/bson"
)

// ODataParser Structure element
type ODataParser struct {
	Collection *mongoCollection
}

// ErrInvalidInput OData syntax error definition
var ErrInvalidInput = errors.New("odata syntax error")

// ErrInvalidDestination Invalid destination error
var ErrInvalidDestination = errors.New("destination must be a pointer type")

// EmptyODataParser Creates an empty odata parser for a specific collection
func EmptyODataParser(collection *mongoCollection) *ODataParser {
	result := ODataParser{
		Collection: collection,
	}

	return &result
}

// TODO: Implement inner count
// GetODataResponse Creates a odata response from an odata url query including count
func (odataParser *ODataParser) GetODataResponse(query url.Values) (*models.ODataResponse, error) {
	// ctx := context.Background()
	response := models.ODataResponse{}

	queryMap, err := odata.ParseURLValues(query)
	if err != nil {
		return nil, err
	}

	// Checks if the count flag is true and count the collection records
	if count, ok := queryMap[odata.Count].(bool); ok {
		if count {
			builder := NewEmptyPipeline(odataParser.Collection)
			countField := builder.CountCollection()
			response.Count = countField
		}
	}

	// execute and decode the odata query
	var element []map[string]interface{}
	err = odataParser.Decode(query, &element)
	if err != nil {
		return nil, err
	}

	response.Value = element

	return &response, nil
}

// Decode Decodes a odata query url to a destination object, this needs to be a pointer
// This will return the object but not an odata response
// returns error if there was an error in the query
func (odataParser *ODataParser) Decode(query url.Values, destination interface{}) error {
	var destType = reflect.TypeOf(destination)
	if destType.Kind() != reflect.Ptr {
		return ErrInvalidDestination
	}

	ctx := context.Background()

	cursor, err := odataParser.Query(query)
	if err != nil {
		return err
	}

	err = cursor.cursor.All(ctx, destination)
	if err != nil {
		return err
	}

	return nil
}

// Query creates a mongo query based on odata parameters
// returns a cursor ready to be iterated or an error if something goes wrong
func (odataParser *ODataParser) Query(query url.Values) (*mongoCursor, error) {
	builder := NewEmptyPipeline(odataParser.Collection)

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
		filterObj, err = ApplyFilter(filterQuery)
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

	cursor, err := builder.Aggregate()

	return cursor, err
}
