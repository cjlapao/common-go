package database

import (
	"github.com/cjlapao/common-go/database/mongodb"
	identity_constants "github.com/cjlapao/common-go/identity/constants"
	"github.com/cjlapao/common-go/identity/models"
)

var currentDatabase string

type MongoDBContextAdapter struct{}

func (u MongoDBContextAdapter) GetUserById(id string) *models.User {
	var result models.User
	repo := getMongoDBTenantRepository()
	dbUsers := repo.FindOne("_id", id)
	dbUsers.Decode(&result)
	return &result
}

func (u MongoDBContextAdapter) GetUserByEmail(email string) *models.User {
	var result models.User
	repo := getMongoDBTenantRepository()
	dbUsers := repo.FindOne("email", email)
	dbUsers.Decode(&result)
	return &result
}

func (u MongoDBContextAdapter) UpsertUser(user models.User) {
	if user.IsValid() {
		logger.Info("Upserting user %v into database %v", currentDatabase)
		repo := getMongoDBTenantRepository()
		builder := mongodb.NewUpdateOneBuilder().FilterBy("_id", user.ID).Encode(user).Build()
		result := repo.UpsertOne(builder)
		if result.MatchedCount <= 0 {
			logger.Error("There was an error upserting user %v", user.Email)
		}
	}
}

func (u MongoDBContextAdapter) RemoveUser(id string) bool {
	return true
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
		repo := getMongoDBTenantRepository()
		builder := mongodb.NewUpdateOneBuilder().FilterBy("_id", id).Encode(user).Build()
		result := repo.UpsertOne(builder)
		if result.MatchedCount <= 0 {
			logger.Error("No user found with id %v", id)
		}
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
		repo := getMongoDBTenantRepository()
		builder := mongodb.NewUpdateOneBuilder().FilterBy("_id", id).Encode(user).Build()
		result := repo.UpsertOne(builder)
		if result.MatchedCount <= 0 {
			logger.Error("No user found with id %v", id)
		}
	}
}

func getMongoDBTenantRepository() mongodb.Repository {
	mongodbSvc := mongodb.Get()
	factory, currentDatabase := mongodbSvc.GetTenantDatabase()
	userRepo := mongodb.NewRepository(factory, currentDatabase, identity_constants.IdentityUsersCollection)
	return userRepo
}
