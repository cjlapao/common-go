package models

import "github.com/cjlapao/common-go/cryptorand"

type UserRole struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"roleName" bson:"roleName"`
}

func NewUserRole(name string) UserRole {
	return UserRole{
		ID:   cryptorand.GenerateNumericRandomString(45),
		Name: name,
	}
}

func (ur UserRole) IsValid() bool {
	return ur.ID != "" && ur.Name != ""
}

func (ur UserRole) IsSuperUser() bool {
	return ur.ID == "_su"
}

var SuRole = UserRole{
	ID:   "_su",
	Name: "Super Administrator",
}

var AdminRole = UserRole{
	ID:   "_admin",
	Name: "Administrator",
}

var RegularUserRole = UserRole{
	ID:   "_user",
	Name: "User",
}
