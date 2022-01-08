package database

import (
	"fmt"

	"github.com/cjlapao/common-go/database/mongodb"
	"github.com/cjlapao/common-go/identity"
	"github.com/cjlapao/common-go/identity/models"
)

type MongoDBContextAdapter struct{}

func (u MongoDBContextAdapter) GetUserById(id string) *models.User {
	var result models.User
	repo := GetRepository()
	dbUsers := repo.FindOne("_id", id)
	dbUsers.Decode(&result)
	return &result
}

func (u MongoDBContextAdapter) GetUserByEmail(email string) *models.User {
	var result models.User
	repo := GetRepository()
	dbUsers := repo.FindOne("email", email)
	dbUsers.Decode(&result)
	return &result
}

func (u MongoDBContextAdapter) UpsertUser(user models.User) {

}

func (u MongoDBContextAdapter) GetUserRefreshToken(id string) *string {
	user := u.GetUserById(id)
	if user != nil {
		return &user.RefreshToken
	}

	return nil
}

func (u MongoDBContextAdapter) UpdateUserRefreshToken(id string, token string) {
	user := u.GetUserById(id)
	if user != nil {
		user.RefreshToken = token
		repo := GetRepository()
		builder := mongodb.NewUpdateOneBuilder().FilterBy("_id", id).Encode(user).Build()
		result := repo.UpsertOne(builder)
		fmt.Printf("%v %v %v %v\n", result.MatchedCount, result.ModifiedCount, result.UpsertedCount, result.UpsertedID)
	}
}

func (u MongoDBContextAdapter) GetUserEmailVerifyToken(id string) *string {
	user := u.GetUserById(id)
	if user != nil {
		return &user.EmailVerifyToken
	}

	return nil
}

func (u MongoDBContextAdapter) UpdateUserEmailVerifyToken(id string, token string) {
	user := u.GetUserById(id)
	if user != nil {
		user.EmailVerifyToken = token
		repo := GetRepository()
		builder := mongodb.NewUpdateOneBuilder().FilterBy("_id", id).Encode(user).Build()
		repo.UpsertOne(builder)
	}
}

func GetRepository() mongodb.Repository {
	mongosvc := mongodb.Get()
	factory, database := mongosvc.GetTenantDatabase()
	userRepo := mongodb.NewRepository(factory, database, identity.IdentityUsersCollection)
	return userRepo
}
