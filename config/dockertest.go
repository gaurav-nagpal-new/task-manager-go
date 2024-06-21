package config

import (
	"context"
	"fmt"

	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func InitializeMongoContainer() (*dockertest.Pool, *dockertest.Resource, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		zap.L().Error("error creating docker pool", zap.Error(err))
		return nil, nil, err
	}

	envVariables := []string{
		"MONGO_INITDB_ROOT_USERNAME=root",
		"MONGO_INITDB_ROOT_PASSWORD=password",
	}

	resource, err := pool.Run("mongo", "5.0", envVariables)
	if err != nil {
		zap.L().Error("error creating mongo container", zap.Error(err))
		return nil, nil, err
	}

	if err := pool.Retry(func() error {
		var err error
		MongoClient, err = mongo.Connect(context.Background(), options.Client().ApplyURI(
			fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp"))))
		if err != nil {
			return err
		}

		return MongoClient.Ping(context.Background(), nil)
	}); err != nil {
		zap.L().Error("error connecting to mongodb", zap.Error(err))
		return nil, nil, err
	}

	return pool, resource, nil
}

func TearDownContainer(pool *dockertest.Pool, resource *dockertest.Resource) error {
	if err := pool.Purge(resource); err != nil {
		zap.L().Error("could not purge resource", zap.Error(err))
		return err
	}
	return nil
}
