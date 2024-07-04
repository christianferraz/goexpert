package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/christianferraz/goexpert/26-Leilao/configuration/logger"
	"github.com/christianferraz/goexpert/26-Leilao/internal/entity/user_entity"
	"github.com/christianferraz/goexpert/26-Leilao/internal/internal_error"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	Id   string `bson:"_id"`
	Name string `bson:"name"`
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
	}
}

func (r *UserRepository) FindUserById(ctx context.Context, userId string) (*user_entity.User, *internal_error.InternalError) {
	var user UserEntityMongo
	filter := bson.M{"_id": userId}
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			logger.Error(fmt.Sprintf("user not found with id: %s", userId), err)
			return nil, internal_error.NewNotFoundError(fmt.Sprintf("user not found with id: %s", userId))
		}
		logger.Error("Error trying to find user by userId", err)

		return nil, internal_error.NewNotFoundError("Error trying to find user by userId")
	}
	return &user_entity.User{
		Id:   user.Id,
		Name: user.Name,
	}, nil
}
