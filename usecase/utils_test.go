package usecase

import (
	"os"
	"task-manager/config"
	"testing"

	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	zap.L().Debug("starting main test function")
	pool, resource, err := config.InitializeMongoContainer()
	if err != nil {
		zap.L().Error("error initializing mongo containers", zap.Error(err))
	}

	// start the test cases
	exitCode := m.Run()

	if err := config.TearDownContainer(pool, resource); err != nil {
		zap.L().Error("error tearing down the containers", zap.Error(err))
	}

	os.Exit(exitCode)
}
