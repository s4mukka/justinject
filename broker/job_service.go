package broker

import (
	"net/http"

	"github.com/s4mukka/justinject/domain"
)

type IRestContext interface {
	JSON(code int, obj any)
	ShouldBindJSON(obj any) error
}

type JobService struct {
	jobUseCase domain.IJobUseCase
}

func (c *JobService) CreateJob(ctx IRestContext) {
	var request domain.CreateJobRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	job, err := c.jobUseCase.CreateJob(request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, *job)
}

type JobServiceFactory struct{}

func (f *JobServiceFactory) MakeJobService(jobUseCase domain.IJobUseCase) *JobService {
	return &JobService{jobUseCase}
}
