package models

import "github.com/google/uuid"

type Tenant struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func NewTenant() *Tenant {
	tenant := Tenant{
		ID: uuid.NewString(),
	}

	return &tenant
}
