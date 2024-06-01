package usecase

import (
	"net/http"
	"os"

	"task-manager/constants"
	"task-manager/dto"
	"task-manager/utils"
	"time"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

func StartEmailCronHandler(w http.ResponseWriter, r *http.Request) {
	zap.L().Debug("EmailCron started")

	// make scheduler
	scheduler := gocron.NewScheduler(time.UTC)

	// defining time expression
	scheduler.Cron(os.Getenv(constants.EmailCronExpression)).Do(GenerateSummaryAndSendEmail)

	utils.Response(w, &dto.GenericResponse{
		Code:    http.StatusOK,
		Message: "Cron started",
	}, http.StatusOK)

	scheduler.StartAsync()

}
