package mongodb

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strings"

	parser "github.com/cjlapao/common-go/odata"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidInput = errors.New("odata syntax error")

func Parser(query url.Values) {
	queryMap, err := parser.ParseURLValues(query)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	filterObj := make(bson.M)
	if queryMap[parser.Filter] != nil {
		filterQuery, _ := queryMap[parser.Filter].(*parser.ParseNode)
		var err error
		filterObj, err = applyFilter(filterQuery)
		if err != nil {
			fmt.Errorf("%v", err.Error())
		}
	}

	fmt.Printf("%v", queryMap["$top"])
	fmt.Printf("%v", filterObj)
}

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
				objId := primitive.NewObjectID()
				objId.UnmarshalText(decodedString)
				value = bson.M{"$" + node.Token.Value.(string): objId}
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
			//nolint: vet
			value := primitive.Regex{Pattern: "^" + node.Children[1].Token.Value.(string), Options: "gi"}
			filter[node.Children[0].Token.Value.(string)] = value

		case "endswith":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			//nolint: vet
			value := primitive.Regex{Pattern: node.Children[1].Token.Value.(string) + "$", Options: "gi"}
			filter[node.Children[0].Token.Value.(string)] = value

		case "contains":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			//nolint: vet
			value := primitive.Regex{Pattern: node.Children[1].Token.Value.(string), Options: "gi"}
			filter[node.Children[0].Token.Value.(string)] = value

		}
	}
	return filter, nil
}
