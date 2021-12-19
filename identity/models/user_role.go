package models

type UserRole struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"roleName" bson:"roleName"`
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
