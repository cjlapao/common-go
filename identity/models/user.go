package models

import "github.com/google/uuid"

// User entity
type User struct {
	ID           string      `json:"id" bson:"_id"`
	Email        string      `json:"email" bson:"email"`
	Username     string      `json:"userName" bson:"userName"`
	FirstName    string      `json:"firstName" bson:"firstName"`
	LastName     string      `json:"lastName" bson:"lastName"`
	DisplayName  string      `json:"displayName" bson:"displayName"`
	Password     string      `json:"password" bson:"password"`
	Token        string      `json:"-" bson:"-"`
	RefreshToken string      `json:"-" bson:"-"`
	Roles        []UserRole  `json:"roles" bson:"roles"`
	Claims       []UserClaim `json:"claims" bson:"claims"`
}

func NewUser() *User {
	user := User{
		ID: uuid.NewString(),
	}

	return &user
}

func (u User) IsValid() bool {
	_, err := uuid.Parse(u.ID)
	if err != nil {
		return false
	}
	if u.Username == "" {
		return false
	}
	if u.Email == "" {
		return false
	}

	return true
}
