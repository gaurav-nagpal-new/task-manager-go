package main

import (
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

	config.InitZapLogger()

	router := mux.NewRouter()

	// use middleware to log endpoint at which request is made
	router.Use(middleware.LogEndPoint)

	// --------------- v1 routes start here -------------------
	router.HandleFunc(routes.GetTasks, usecase.GetTasksHandler).Methods(http.MethodGet)
	router.HandleFunc(routes.UpdateTasks, usecase.UpdateTasksHandler).Methods(http.MethodGet)
	router.HandleFunc(routes.DeleteTasks, usecase.DeleteTasksHandler).Methods(http.MethodGet)
	router.HandleFunc(routes.CreateTasks, usecase.CreateTasksHandler).Methods(http.MethodGet)
	// --------------- v1 routes end here ---------------------

	zap.L().Debug("Starting server")
	if err := http.ListenAndServe(":8080", router); err != nil {
		zap.L().Debug("error starting server")
		return
	}
}
