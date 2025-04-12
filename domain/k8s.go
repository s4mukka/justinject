package domain

type IK8sRepository interface {
	CreateJob(job IJob) error
}
