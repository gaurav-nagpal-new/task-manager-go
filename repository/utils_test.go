package repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var db *mongo.Client

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		zap.L().Error("error creating docker pool")
		return
	}

	envVariables := []string{
		"MONGODB_INIT_USERNAME=root",
		"MONGODB_INIT_PASSWORD=password",
	}

	resource, err := pool.Run("mongo", "5.0", envVariables)
	if err != nil {
		zap.L().Error("error creating mongo container", zap.Error(err))
		return
	}

	if err := pool.Retry(func() error {
		var err error
		db, err = mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp"))))
		if err != nil {
			return err
		}

		return db.Ping(context.Background(), nil)
	}); err != nil {
		zap.L().Error("error connecting to mongodb", zap.Error(err))
	}

	// start the test cases
	exitCode := m.Run()

	if err := pool.Purge(resource); err != nil {
		zap.L().Error("could not purge resource", zap.Error(err))
	}

	os.Exit(exitCode)
}
