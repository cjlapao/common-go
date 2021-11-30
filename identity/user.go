package identity

// User entity
type User struct {
	ID        string      `json:"id" bson:"_id"`
	Email     string      `json:"email" bson:"email"`
	Username  string      `json:"userName" bson:"userName"`
	FirstName string      `json:"firstName" bson:"firstName"`
	LastName  string      `json:"lastName" bson:"lastName"`
	Password  string      `json:"password" bson:"password"`
	Claims    []UserClaim `json:"claims" bson:"claims"`
}

type UserClaim struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}
