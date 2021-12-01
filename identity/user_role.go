package identity

type UserRole struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"roleName" bson:"roleName"`
}

var AdminRole = UserRole{
	ID:   "_admin",
	Name: "Administrator",
}

var RegularUserRole = UserRole{
	ID:   "_user",
	Name: "User",
}
