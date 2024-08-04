package domain

type IK8sRepository interface {
	CreateJob(job *Job) error
}
