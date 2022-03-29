package models

type UserRole struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"roleName" bson:"roleName"`
}

func NewUserRole(id string, name string) UserRole {
	return UserRole{
		ID:   id,
		Name: name,
	}
}

func (ur UserRole) IsValid() bool {
	return ur.ID != "" && ur.Name != ""
}

func (ur UserRole) IsSuperUser() bool {
	return ur.ID == "_su"
}
