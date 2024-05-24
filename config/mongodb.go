package config

import (
	"context"
	"fmt"
	"os"
	"task-manager/constants"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Initialize mongoDb connection here

var MongoClient *mongo.Client

func InitMongoConnection() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv(constants.MongoDBConnectionString)).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)
	MongoClient = client
	if err != nil {
		//TODO : change to Debug log using zap
		fmt.Println("error in initializing mongodb connection")
		return err
	}

	return nil
}
