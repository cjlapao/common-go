package models

import "time"

type ApiKey struct {
	ID          string       `json:"id" bson:"_id"`
	Name        string       `json:"name" bson:"name"`
	KeyValue    string       `json:"keyValue" bson:"keyValue"`
	Description string       `json:"description,omitempty" bson:"description"`
	ValidFrom   time.Time    `json:"validFrom" bson:"validFrom"`
	ValidTo     time.Time    `json:"validTo" bson:"validTo"`
	Roles       []ApiKeyRole `json:"roles" bson:"roles"`
}
