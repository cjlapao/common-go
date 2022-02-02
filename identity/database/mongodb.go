package database

import (
	"errors"

	"github.com/cjlapao/common-go/database/mongodb"
	identity_constants "github.com/cjlapao/common-go/identity/constants"
	"github.com/cjlapao/common-go/identity/models"
)

var currentDatabase string
var ErrUserNotValid = errors.New("user model is not valid")

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

func (u MongoDBContextAdapter) GetUserByUsername(username string) *models.User {
	var result models.User
	repo := getMongoDBTenantRepository()
	dbUsers := repo.FindOne("username", username)
	dbUsers.Decode(&result)
	return &result
}

func (u MongoDBContextAdapter) UpsertUser(user models.User) error {
	if user.IsValid() {
		repo := getMongoDBTenantRepository()
		logger.Info("Upserting user %v into database %v", currentDatabase)
		builder, err := mongodb.NewUpdateOneModelBuilder().FilterBy("_id", mongodb.Equal, user.ID).Encode(user).Build()
		if err != nil {
			return err
		}
		result := repo.UpsertOne(builder)
		if result.MatchedCount <= 0 {
			logger.Error("There was an error upserting user %v", user.Email)
		}
	} else {
		return ErrUserNotValid
	}

	return nil
}

func (u MongoDBContextAdapter) RemoveUser(id string) bool {
	if id == "" {
		return false
	}

	repo := getMongoDBTenantRepository()
	logger.Info("Removing user %v from database %v", currentDatabase)
	builder, err := mongodb.NewDeleteOneBuilder().FilterBy("_id", mongodb.Equal, id).Build()
	if err != nil {
		return false
	}
	result := repo.DeleteOne(builder)
	if result.DeletedCount <= 0 {
		logger.Error("There was an error removing userid %v", id)
	}

	return true
}

func (u MongoDBContextAdapter) GetUserRefreshToken(id string) *string {
	user := u.GetUserById(id)
	if user != nil {
		return &user.RefreshToken
	}

	return nil
}

func (u MongoDBContextAdapter) UpdateUserRefreshToken(id string, token string) bool {
	user := u.GetUserById(id)
	if user != nil {
		user.RefreshToken = token
		repo := getMongoDBTenantRepository()
		builder, err := mongodb.NewUpdateOneModelBuilder().FilterBy("_id", mongodb.Equal, id).Encode(user).Build()
		if err != nil {
			return false
		}

		result := repo.UpsertOne(builder)
		if result.MatchedCount == 0 && result.ModifiedCount == 0 && result.UpsertedCount == 0 {
			logger.Error("There was an error updating the refresh token for user with id %v", id)
			return false
		}
		return true
	}

	return false
}

func (u MongoDBContextAdapter) GetUserEmailVerifyToken(id string) *string {
	user := u.GetUserById(id)
	if user != nil {
		return &user.EmailVerifyToken
	}

	return nil
}

func (u MongoDBContextAdapter) UpdateUserEmailVerifyToken(id string, token string) bool {
	user := u.GetUserById(id)
	if user != nil {
		user.EmailVerifyToken = token
		repo := getMongoDBTenantRepository()
		builder, _ := mongodb.NewUpdateOneModelBuilder().FilterBy("_id", mongodb.Equal, id).Encode(user).Build()
		result := repo.UpsertOne(builder)
		if result.MatchedCount <= 0 {
			logger.Error("There was an error updating the verify email token for user with id %v", id)
			return false
		}
		return true
	}

	return false
}

func getMongoDBTenantRepository() mongodb.MongoRepository {
	mongodbSvc := mongodb.Get()
	userRepo := mongodbSvc.TenantDatabase().NewRepository(identity_constants.IdentityUsersCollection)
	return userRepo
}
