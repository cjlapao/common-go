package models

type TenantFeature struct {
	ID    string             `json:"id" bson:"_id"`
	Name  string             `json:"name" bson:"name"`
	State TenantFeatureState `json:"state" bson:"state"`
}
