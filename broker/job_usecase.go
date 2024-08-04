package broker

import "github.com/s4mukka/justinject/domain"

type JobUseCase struct {
	JobRepository domain.IJobRepository
	K8sRepository domain.IK8sRepository
}

func (uc *JobUseCase) CreateJob(request domain.CreateJobRequest) (*domain.Job, error) {
	extractor, err := uc.JobRepository.GetExtractorById(request.ExtractorId)
	if err != nil {
		return nil, err
	}

	job := domain.Job{
		CreateJobRequest: domain.CreateJobRequest{
			ExtractorId:   extractor.Id,
			Query:         request.Query,
			UpperBound:    request.UpperBound,
			LowerBound:    request.LowerBound,
			NumPartitions: request.NumPartitions,
		},
	}

	if err := uc.JobRepository.CreateJob(&job); err != nil {
		return nil, err
	}

	if err := uc.K8sRepository.CreateJob(&job); err != nil {
		return nil, err
	}

	return &job, nil
}
