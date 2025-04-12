package broker

import (
	"net/http"

	"github.com/s4mukka/justinject/domain"
)

type JobService struct {
	jobUseCase domain.IJobUseCase
}

func (c *JobService) CreateJob(ctx domain.IRestContext) {
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

	ctx.JSON(http.StatusOK, job)
}

type JobServiceFactory struct {
	jobUseCaseFactory domain.IFactory[domain.IJobUseCase]
}

func (f JobServiceFactory) Create() (domain.IJobService, error) {
	jobUseCase, err := f.jobUseCaseFactory.Create()
	if err != nil {
		return nil, err
	}

	return &JobService{jobUseCase}, nil
}
