package domain

import (
	"time"

	"gorm.io/gorm"
)

type CreateJobRequest struct {
	ExtractorId   string `json:"extractorId"   binding:"required"`
	Query         string `json:"query"         binding:"required"`
	LowerBound    int64  `json:"lowerBound"`
	UpperBound    int64  `json:"upperBound"`
	NumPartitions int64  `json:"numPartitions"`
}

type Job struct {
	Id        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	CreateJobRequest
	Extractor Extractor
}

type JobTemplate struct {
	JobId         string
	Query         string
	LowerBound    int64
	UpperBound    int64
	NumPartitions int64
	Driver        string
	Host          string
	Port          int
	User          string
	Password      string
	Db            string
}

type IJob interface {
	ParseTemplate() *JobTemplate
}

type IJobService interface {
	CreateJob(ctx IRestContext)
}

type IJobUseCase interface {
	CreateJob(request CreateJobRequest) (IJob, error)
}

type IJobRepository interface {
	Create(job IJob) error
}

func (j Job) ParseTemplate() *JobTemplate {
	return &JobTemplate{
		JobId:         j.Id,
		Query:         j.Query,
		UpperBound:    j.UpperBound,
		LowerBound:    j.LowerBound,
		NumPartitions: j.NumPartitions,
		Driver:        j.Extractor.Driver.Driver,
		Host:          j.Extractor.Driver.Host,
		Port:          j.Extractor.Driver.Port,
		User:          j.Extractor.Driver.User,
		Password:      j.Extractor.Driver.Password,
		Db:            j.Extractor.Driver.Db,
	}
}
