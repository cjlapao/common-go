package database

import (
	"errors"
	"fmt"

	"github.com/cjlapao/common-go/database/mongodb"
	identity_constants "github.com/cjlapao/common-go/identity/constants"
	"github.com/cjlapao/common-go/identity/models"
)

var currentDatabase string
var ErrUserNotValid = errors.New("user model is not valid")
var ErrUnknown = errors.New("unknown error occurred")

type MongoDBContextAdapter struct{}

func (u MongoDBContextAdapter) GetUserById(id string) *models.User {
	var result models.User
	repo := getMongoDBTenantRepository()
	dbUsers := repo.FindOne(fmt.Sprintf("_id eq %v", id))
	dbUsers.Decode(&result)
	return &result
}

func (u MongoDBContextAdapter) GetUserByEmail(email string) *models.User {
	var result models.User
	repo := getMongoDBTenantRepository()
	dbUsers := repo.FindOne(fmt.Sprintf("email eq %v", email))
	dbUsers.Decode(&result)
	return &result
}

func (u MongoDBContextAdapter) GetUserByUsername(username string) *models.User {
	var result models.User
	repo := getMongoDBTenantRepository()
	dbUsers := repo.FindOne(fmt.Sprintf("username eq %v", username))
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
		result, err := repo.UpsertOne(builder)
		if err != nil {
			logger.Error("There was an error upserting user %v, %v", user.Email, err.Error())
			return err
		}
		if result.MatchedCount <= 0 {
			logger.Error("There was an error upserting user %v", user.Email)
			return ErrUnknown
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
	result, err := repo.DeleteOne(builder)
	if err != nil {
		logger.Exception(err, "there was an error removing user from collection with id %v", id)
	}
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

		result, err := repo.UpsertOne(builder)
		if err != nil {
			logger.Exception(err, "There was an error while upserting the refresh token with id %v", id)
			return false
		}
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
		result, err := repo.UpsertOne(builder)
		if err != nil {
			logger.Exception(err, "There was an error upserting the email verification token with id %v", id)
			return false
		}
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
