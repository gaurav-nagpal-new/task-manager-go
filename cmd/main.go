package main

import (
	"context"
	"net/http"
	"path/filepath"
	"task-manager/config"
	"task-manager/middleware"
	"task-manager/routes"
	"task-manager/usecase"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {

	// Load the .env file
	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
		zap.L().Error("error loading .env file", zap.Error(err))
		return
	}

	// Initialize MongoDB connections
	if err := config.InitMongoConnection(); err != nil {
		zap.L().Error("error creating mongodb connection", zap.Error(err))
		return
	}

	// close mongodb connection
	defer func() {
		if err := config.MongoClient.Disconnect(context.Background()); err != nil {
			zap.L().Error("error closing mongodb connection", zap.Error(
				err,
			))
		}
	}()

	config.InitZapLogger()

	router := mux.NewRouter()
	authRouter := router.NewRoute().Subrouter()
	apiRouter := router.NewRoute().Subrouter()

	// Common middleware for logging
	authRouter.Use(middleware.LogEndPoint)
	apiRouter.Use(middleware.LogEndPoint, middleware.VerifyJWTAuth)

	// --------------- v1 task routes start here -------------------
	apiRouter.HandleFunc(routes.GetTasks, usecase.GetTasksHandler).Methods(http.MethodGet)
	apiRouter.HandleFunc(routes.UpdateTasks, usecase.UpdateTasksHandler).Methods(http.MethodPut)
	apiRouter.HandleFunc(routes.DeleteTasks, usecase.DeleteTasksHandler).Methods(http.MethodDelete)
	apiRouter.HandleFunc(routes.CreateTasks, usecase.CreateTasksHandler).Methods(http.MethodPost)
	// --------------- v1 task routes end here ---------------------

	// --------------- Auth routes start here----------------
	authRouter.HandleFunc(routes.SignUp, usecase.SignUpHandler).Methods(http.MethodPost)
	authRouter.HandleFunc(routes.SignIn, usecase.SignInHandler).Methods(http.MethodPost)
	// --------------- Auth routes end here ------------------

	zap.L().Debug("Starting server")
	if err := http.ListenAndServe(":8080", router); err != nil {
		zap.L().Debug("error starting server")
		return
	}
}
