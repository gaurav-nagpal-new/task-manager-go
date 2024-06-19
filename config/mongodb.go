package config

import (
	"context"
	"os"
	"task-manager/constants"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// Initialize mongoDb connection here

var MongoClient *mongo.Client

func InitMongoConnection() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv(constants.MongoDBConnectionString)).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	MongoClient = client
	if err != nil {
		zap.L().Error("error in initializing mongodb connection", zap.Error(err))
		return err
	}

	return nil
}
