package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/dkischenko/chat/internal/user"
	"github.com/dkischenko/chat/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongodb struct {
	logger     *logger.Logger
	collection *mongo.Collection
}

func NewStorage(database *mongo.Database, collection string, logger *logger.Logger) user.Repository {
	return &mongodb{
		collection: database.Collection(collection),
		logger:     logger,
	}
}

func (db *mongodb) Create(ctx context.Context, user *user.User) (id string, err error) {
	result, err := db.collection.InsertOne(ctx, user)
	if err != nil {
		db.logger.Entry.Errorf("failed insert data with error: %v", err)
		return "", err
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		db.logger.Entry.Infof("user ID mongo: %v", oid.Hex())
		return oid.Hex(), nil
	}

	db.logger.Entry.Errorf("failed to convert objectid to hex with oid: %s", oid)
	return "", fmt.Errorf("failed to convert objectid to hex with oid: %s", oid)
}

func (db *mongodb) FindOne(ctx context.Context, username string) (u *user.User, err error) {
	filter := bson.M{"username": username}
	result := db.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return u, fmt.Errorf("user with '%s' username not found", username)
		}

		return u, fmt.Errorf("failed find user with error: %s", result.Err())
	}

	if err = result.Decode(&u); err != nil {
		return u, fmt.Errorf("error decode data with error: %s", err)
	}

	return u, nil
}
