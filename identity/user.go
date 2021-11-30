package identity

import "github.com/google/uuid"

// User entity
type User struct {
	ID        string      `json:"id" bson:"_id"`
	Email     string      `json:"email" bson:"email"`
	Username  string      `json:"userName" bson:"userName"`
	FirstName string      `json:"firstName" bson:"firstName"`
	LastName  string      `json:"lastName" bson:"lastName"`
	Password  string      `json:"password" bson:"password"`
	Roles     []UserRole  `json:"roles" bson:"roles"`
	Claims    []UserClaim `json:"claims" bson:"claims"`
}

func NewUser() *User {
	user := User{
		ID: uuid.NewString(),
	}

	return &user
}

type UserClaim struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"claimName" bson:"claimName"`
}
