package repositories

import (
	"context"
	"structure-golang/core/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type userRepo struct {
	db         *mongo.Database
	collection string
}

func NewUserRepository(db *mongo.Database, collection string) UserRepository {
	return userRepo{
		db:         db,
		collection: collection,
	}
}

func (r userRepo) Get(filter models.RepoFilterUserModel) (result models.RepoUserModel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = r.db.Collection(r.collection).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, err
	}

	return result, err
}
