package models

import (
	cryptorand "github.com/cjlapao/common-go-cryptorand"
	"github.com/cjlapao/common-go/constants"
	"github.com/cjlapao/common-go/guard"
)

type ApiKeyRole struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func NewApiKeyRole(name string) *ApiKeyRole {
	if err := guard.EmptyOrNil(name); err != nil {
		logger.Exception(err, "There was an error creating the api Key")
	}

	result := ApiKeyRole{
		ID:   cryptorand.GetRandomString(constants.ID_SIZE),
		Name: name,
	}

	return &result
}
