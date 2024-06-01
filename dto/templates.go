package dto

import "task-manager/model"

type EmailCronTemplateData struct {
	UserName              string
	UserEmail             string
	Data                  []*model.Task
	PerformancePercentage float64
}
