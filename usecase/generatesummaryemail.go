package usecase

import (
	"context"
	"math"
	"task-manager/config"
	"task-manager/constants"
	"task-manager/dto"
	"task-manager/email"
	"task-manager/model"
	"task-manager/repository"
	"time"

	"go.uber.org/zap"
)

func GenerateSummaryAndSendEmail() {
	// use summary.html template and fill it with data calculated here
	// and then convert that .html file to pdf
	// send that pdf to user's email fetched from DB

	// calculate the data needed in html template
	/*
		Get all users from DB
		Iterate through all users
		Get all the tasks with their count (use facet) for that particular user using task_collection from user object
		Data will be like:
		Todo : 10
		InProgress : 5
		Done : 15
		Performace Percentage : ( Done/( Todo + InProgress)) * 100
	*/

	mongoRepo := repository.MongoRepository{
		Client:     config.MongoClient,
		Collection: constants.UserCollectionName,
	}

	now := time.Now()
	startDate := now.AddDate(0, 0, -7)
	endDate := now
	ctx := context.Background()
	usersData, err := mongoRepo.FetchAllUsersData(ctx)
	if err != nil {
		zap.L().Error("error fetching users data", zap.Error(err))
		return
	}

	for _, user := range usersData {
		mongoRepo.Collection = user.TaskCollection
		tasks, err := mongoRepo.FetchAllTasksDetails(ctx, startDate, endDate)
		if err != nil {
			zap.L().Error("error fetching tasks data for the user", zap.Error(err))
		}

		emailData := &dto.EmailCronTemplateData{
			UserName:              user.Name,
			UserEmail:             user.Email,
			Data:                  tasks,
			PerformancePercentage: calculatePerformance(tasks),
		}

		if err := email.SendSummaryEmailToUser(constants.PerformanceSummaryTemplate, emailData); err != nil {
			zap.L().Error("error creating html template and sending email", zap.Error(err))
		}
	}

	zap.L().Debug("email cron executed successfully")
}

func calculatePerformance(tasks []*model.Task) float64 {
	var todo, inProgress, done float64

	for _, t := range tasks {
		switch t.Status {
		case model.Todo:
			todo++
		case model.InProgress:
			inProgress++
		case model.Done:
			done++
		}
	}

	if todo+inProgress == 0 {
		return 100
	}
	if todo+inProgress+done == 0 {
		return 0
	}

	return math.Round((done/(inProgress+todo+done))*100*100) / 100
}
