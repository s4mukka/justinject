package broker

import "github.com/s4mukka/justinject/domain"

type JobUseCase struct {
	ExtractorRepository domain.IExtractorRepository
	JobRepository       domain.IJobRepository
	K8sRepository       domain.IK8sRepository
}

func (uc *JobUseCase) CreateJob(request domain.CreateJobRequest) (domain.IJob, error) {
	extractor, err := uc.ExtractorRepository.GetById(request.ExtractorId)
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
		Extractor: domain.Extractor{
			Driver: extractor.Driver,
		},
	}

	if err := uc.JobRepository.Create(&job); err != nil {
		return nil, err
	}

	if err := uc.K8sRepository.CreateJob(&job); err != nil {
		return nil, err
	}

	return &job, nil
}

type JobUseCaseFactory struct {
	extractorRepositoryFactory domain.IFactory[domain.IExtractorRepository]
	jobRepositoryFactory       domain.IFactory[domain.IJobRepository]
	k8sRepositoryFactory       domain.IFactory[domain.IK8sRepository]
}

func (f *JobUseCaseFactory) Create() (domain.IJobUseCase, error) {
	extractorRepository, err := f.extractorRepositoryFactory.Create()
	if err != nil {
		return nil, err
	}
	jobRepository, err := f.jobRepositoryFactory.Create()
	if err != nil {
		return nil, err
	}
	k8sRepository, err := f.k8sRepositoryFactory.Create()
	if err != nil {
		return nil, err
	}

	return &JobUseCase{
		ExtractorRepository: extractorRepository,
		JobRepository:       jobRepository,
		K8sRepository:       k8sRepository,
	}, nil
}
