/* Apache v2 license
*  Copyright (C) <2019> Intel Corporation
*
*  SPDX-License-Identifier: Apache-2.0
 */

package mongo

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	parser "github.com/cjlapao/common-go/odata"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrInvalidInput Client errors
var ErrInvalidInput = errors.New("odata syntax error")

// Making FilterObj global for now, so inlinecount can reuse it.
// TODO: Refactor odata library API
var filterObj bson.M

// ODataQuery creates a mgo query based on odata parameters
//nolint :gocyclo
func ODataQuery(query url.Values, object interface{}, collection *mongo.Collection) error {

	// Parse url values
	queryMap, err := parser.ParseURLValues(query)
	if err != nil {
		return err
	}

	limit, _ := queryMap[parser.Top].(int)
	skip, _ := queryMap[parser.Skip].(int)

	filterObj = make(bson.M)
	if queryMap[parser.Filter] != nil {
		filterQuery, _ := queryMap[parser.Filter].(*parser.ParseNode)
		var err error
		filterObj, err = applyFilter(filterQuery)
		if err != nil {
			return ErrInvalidInput
		}
	}

	// Prepare Select
	selectMap := make(bson.M)

	if queryMap["$select"] != nil {
		selectSlice := reflect.ValueOf(queryMap["$select"])
		if selectSlice.Len() > 1 && selectSlice.Index(0).Interface().(string) != "*" {
			for i := 0; i < selectSlice.Len(); i++ {
				fieldName := selectSlice.Index(i).Interface().(string)
				selectMap[fieldName] = 1
			}
		}
	}

	// Sort
	var sortFields []string
	if queryMap[parser.OrderBy] != nil {
		orderBySlice := queryMap[parser.OrderBy].([]parser.OrderItem)
		for _, item := range orderBySlice {
			if item.Order == "desc" {
				item.Field = "-" + item.Field
			}
			sortFields = append(sortFields, item.Field)
		}
	}

	// Query
	// odataFunc := collection.Find(filterObj).Select(selectMap).Limit(limit).Skip(skip).Sort(sortFields...).All(object)

	fmt.Printf("%v %v", limit, skip)
	return nil
}

// ODataCount runs a collection.Count() function based on $count odata parameter
// func ODataCount(collection *mgo.Collection) (int, error) {
// 	return collection.Count()
// }

// ODataInlineCount retrieves the total count from a filtered data
// func ODataInlineCount(collection *mgo.Collection) (int, error) {

// 	return collection.Find(filterObj).Count()
// }

//nolint :gocyclo
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
			//nolint: vet
			value := primitive.Regex{"^" + node.Children[1].Token.Value.(string), "gi"}
			filter[node.Children[0].Token.Value.(string)] = value

		case "endswith":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			//nolint: vet
			value := primitive.Regex{node.Children[1].Token.Value.(string) + "$", "gi"}
			filter[node.Children[0].Token.Value.(string)] = value

		case "contains":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			//nolint: vet
			value := primitive.Regex{node.Children[1].Token.Value.(string), "gi"}
			filter[node.Children[0].Token.Value.(string)] = value

		}
	}
	return filter, nil
}
